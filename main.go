package main

import (
	"log"
	"os"
	"time"

	"github.com/blanvam/rasp-garden/api"
	_resourceController "github.com/blanvam/rasp-garden/resource/controller"
	_resourceDB "github.com/blanvam/rasp-garden/resource/database"
	_resourceMiddleware "github.com/blanvam/rasp-garden/resource/middleware"
	_resourceRepo "github.com/blanvam/rasp-garden/resource/repository"
	_resourceUsecase "github.com/blanvam/rasp-garden/resource/usecase"
	"github.com/peterbourgon/diskv"
)

const (
	resourceRoute string = "/resource"
	timeout       int    = 3
	jwtSecret     string = "MySecret"
	minpin        int    = 1
	maxpin        int    = 26
)

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
	checkOrSetEnv("BD_PATH", "/Users/asfarto/go/src/github.com/blanvam/rasp-garden/diskv_db")
	checkOrSetEnv("JWT_SECRET", "MySecret")
}

func main() {
	log.Println("Setting up resources")
	dbConn := getdb()
	resourceDB := _resourceDB.NewDiskvDatabase(dbConn)
	resourceRepo := _resourceRepo.NewResourceRepository(resourceDB, minpin, maxpin)
	resourceUsecase := _resourceUsecase.NewResourceUsecase(resourceRepo, time.Duration(timeout)*time.Second)
	resourceController := _resourceController.NewResourceHTTPpHandler(resourceUsecase)
	resourceMiddleware := _resourceMiddleware.NewRequireResourceMiddleware(resourceUsecase)

	api.Api(resourceRoute, resourceController, resourceMiddleware)
}
