package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/blanvam/rasp-garden/api"
	"github.com/blanvam/rasp-garden/broker"
	"github.com/blanvam/rasp-garden/database"
	_resourceController "github.com/blanvam/rasp-garden/resource/controller"
	_resourceMiddleware "github.com/blanvam/rasp-garden/resource/middleware"
	_resourceRepo "github.com/blanvam/rasp-garden/resource/repository"
	_resourceUsecase "github.com/blanvam/rasp-garden/resource/usecase"

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
	bdPath := os.Getenv("BD_PATH")
	database := database.NewDiskvDatabase(bdPath)

	t := time.Duration(1) * time.Second
	cid := "start"
	u := "username"
	p := "password"
	s := []string{"0.0.0.0:1883"}
	brokerClient := broker.NewPahoClient(t, cid, u, p, s)

	resourceRepo := _resourceRepo.NewResourceRepository(database, minpin, maxpin)
	resourceUsecase := _resourceUsecase.NewResourceUsecase(resourceRepo, time.Duration(timeout)*time.Second)
	resourceController := _resourceController.NewResourceHTTPpHandler(resourceUsecase)
	resourceMiddleware := _resourceMiddleware.NewRequireResourceMiddleware(resourceUsecase)

	log.Println("Despues api")

	topicRepo := _topicRepo.NewTopicRepository(brokerClient)
	var qoS uint8
	qoS = 2 // At most once (0) // At least once (1) //Exactly once (2).
	topicUsecase := _topicUsecase.NewTopicUsecase(topicRepo, qoS, time.Duration(timeout)*time.Second)

	c := context.Background()

	topic := "pin12"
	topicUsecase.Subscribe(c, topic)
	msgt := time.Now()
	msg := entity.Resource{"E1", "Prueba", 12, entity.ResourceKindOut, entity.ResourceStatusClosed, msgt, msgt}

	err := topicUsecase.Publish(c, topic, &msg)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	api.Api(resourceRoute, resourceController, resourceMiddleware)
}
