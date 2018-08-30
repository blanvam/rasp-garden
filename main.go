package main

import (
	"context"
	"fmt"
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

	_topicClient "github.com/blanvam/rasp-garden/topic/client"
	_topicRepo "github.com/blanvam/rasp-garden/topic/repository"
	_topicUsecase "github.com/blanvam/rasp-garden/topic/usecase"

	entity "github.com/blanvam/rasp-garden/entities"
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

	go api.Api(resourceRoute, resourceController, resourceMiddleware)

	log.Println("Despues api")

	t := time.Duration(1) * time.Second
	cid := "start"
	u := "username"
	p := "password"
	s := []string{"0.0.0.0:1883"}

	topicClient := _topicClient.NewPahoClient(t, cid, u, p, s)
	topicRepo := _topicRepo.NewTopicRepository(topicClient)
	topicUsecase := _topicUsecase.NewTopicUsecase(topicRepo, 1, time.Duration(timeout)*time.Second)

	c := context.Background()

	topic := "pin112"
	topicUsecase.Subscribe(c, topic)
	msgt := time.Now()
	msg := entity.Message{2, "Hello2", msgt, msgt}

	err := topicUsecase.Publish(c, topic, &msg)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
