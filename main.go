package main

import (
	"log"
	"net/http"
	"os"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/blanvam/rasp-garden/middleware"
	_resourceController "github.com/blanvam/rasp-garden/resource/controller"
	_resourceDB "github.com/blanvam/rasp-garden/resource/database"
	_resourceMiddleware "github.com/blanvam/rasp-garden/resource/middleware"
	_resourceRepo "github.com/blanvam/rasp-garden/resource/repository"
	_resourceUsecase "github.com/blanvam/rasp-garden/resource/usecase"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/peterbourgon/diskv"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

const (
	apiversion    string = "v1"
	apiRoute      string = "/api"
	resourceRoute string = "/resource"
	idRoute       string = "/{id:[0-9a-f]+}"
	port          string = "8000"
	timeout       int    = 3
	jwtSecret     string = "MySecret"
	minpin        int    = 1
	maxpin        int    = 26
)

var apiCommonMiddleware *negroni.Negroni
var baseRoute string

func getAPICommonMiddleware() *negroni.Negroni {
	optionsMiddleware := cors.AllowAll()
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})
	return negroni.New(
		negroni.NewLogger(),
		optionsMiddleware,
		middleware.NewRequireJSONMiddleware(),
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
	)
}

func getdb() *diskv.Diskv {
	bdPath := os.Getenv("BD_PATH")
	flatTransform := func(s string) []string { return []string{} }
	db := diskv.New(diskv.Options{
		BasePath:     bdPath,
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})

	return db
}

func checkOrSetEnv(key string, value string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, value)
	}
}
func init() {
	checkOrSetEnv("JWT_SECRET", "MySecret")
	checkOrSetEnv("BD_PATH", "/home/vmacias/go/src/github.com/blanvam/rasp-garden/diskv_db")
	baseRoute = apiRoute + "/" + apiversion
	apiCommonMiddleware = getAPICommonMiddleware()
}

func main() {
	log.Println("Setting up resources")
	dbConn := getdb()
	resourceDB := _resourceDB.NewDiskvDatabase(dbConn)
	resourceRepo := _resourceRepo.NewResourceRepository(resourceDB, minpin, maxpin)
	resourceUsecase := _resourceUsecase.NewResourceUsecase(resourceRepo, time.Duration(timeout)*time.Second)
	resourceController := _resourceController.NewResourceHTTPpHandler(resourceUsecase)
	resourceMiddleware := _resourceMiddleware.NewRequireResourceMiddleware(resourceUsecase)

	// ROUTES
	log.Println("Setting up routes")
	middlewareRouter := mux.NewRouter()
	router := mux.NewRouter() //two routers are neccesary due to negroni

	// API routes
	apiRouter := router.PathPrefix(baseRoute).Subrouter()
	// Resource router
	resourceRouter := apiRouter.PathPrefix(resourceRoute).Subrouter()
	resourceRouter.HandleFunc("", resourceController.ListResourcesHandler).Methods("GET")
	resourceRouter.HandleFunc("", resourceController.CreateResourceHandler).Methods("POST")
	// Resource detail routes
	resourceDetailRouter := resourceRouter.PathPrefix(idRoute).Subrouter()
	resourceDetailRouter.HandleFunc("", resourceController.ResourceDetailHandler).Methods("GET")
	resourceDetailRouter.HandleFunc("", resourceController.ResourceDeleteHandler).Methods("DELETE")

	// Middlewares
	// Order matters, we have to go from most to least specific routes
	middlewareRouter.PathPrefix(baseRoute + resourceRoute + idRoute).Handler(apiCommonMiddleware.With(
		resourceMiddleware,
		negroni.Wrap(resourceDetailRouter),
	))
	middlewareRouter.PathPrefix(baseRoute).Handler(apiCommonMiddleware.With(
		negroni.Wrap(apiRouter),
	))

	log.Println("Server starting at port", port)
	log.Panic(http.ListenAndServe(":"+port, middlewareRouter))
}
