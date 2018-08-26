package api

import (
	"log"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/blanvam/rasp-garden/middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

const (
	apiversion string = "v1"
	apiRoute   string = "/api"
	idRoute    string = "/{id:[0-9a-f]+}"
	port       string = "8000"
	jwtSecret  string = "MySecret"
)

var baseRoute string
var apiCommonMiddleware *negroni.Negroni

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

func Api(routePath string, controller Controller, middleware Middleware) {
	baseRoute = apiRoute + "/" + apiversion
	apiCommonMiddleware = getAPICommonMiddleware()

	// ROUTES
	log.Println("Setting up routes")
	middlewareRouter := mux.NewRouter()
	router := mux.NewRouter() //two routers are neccesary due to negroni

	// API routes
	apiRouter := router.PathPrefix(baseRoute).Subrouter()
	// Root routes
	rootRouter := apiRouter.PathPrefix(routePath).Subrouter()
	rootRouter.HandleFunc("", controller.ListHandler).Methods("GET")
	rootRouter.HandleFunc("", controller.CreateHandler).Methods("POST")
	// Detail routes
	detailRouter := rootRouter.PathPrefix(idRoute).Subrouter()
	detailRouter.HandleFunc("", controller.DetailHandler).Methods("GET")
	detailRouter.HandleFunc("", controller.DeleteHandler).Methods("DELETE")

	// Middlewares
	// Order matters, we have to go from most to least specific routes
	middlewareRouter.PathPrefix(baseRoute + routePath + idRoute).Handler(apiCommonMiddleware.With(
		middleware,
		negroni.Wrap(detailRouter),
	))
	middlewareRouter.PathPrefix(baseRoute).Handler(apiCommonMiddleware.With(
		negroni.Wrap(apiRouter),
	))

	log.Println("Server starting at port", port)
	log.Panic(http.ListenAndServe(":"+port, middlewareRouter))
}
