package repository

import (
	"context"
	"log"

	entity "github.com/blanvam/rasp-garden/entities"
	"github.com/blanvam/rasp-garden/topic"
)

type topicRepository struct {
	client topic.Client
}

// NewTopicRepository creates an topicRepository instance.
func NewTopicRepository(cli topic.Client) topic.Repository {
	return &topicRepository{
		client: cli,
	}
}

// IsConnected return true if the client is connected to the server.
func (t *topicRepository) IsConnected(ctx context.Context) bool {
	if t.client == nil {
		return false
	}
	c := make(chan bool)
	var result bool

	go t.client.IsConnected(c)

	select {
	case <-ctx.Done():
		log.Println("Context is done (topic IsConnected)")
		return false
	case result = <-c:
		log.Println("Went to client successfuly :)")
		return result == true
	}
}

func (t *topicRepository) Connect(ctx context.Context) error {
	c := make(chan error)

	go t.client.Connect(c)

	return t.waitForError(ctx, c)
}

// Disconnect will disconnect from the given MQTT server and clean up all client resources
func (t *topicRepository) Disconnect(ctx context.Context) error {
	if !t.IsConnected(ctx) {
		return entity.ErrNotConnected
	}

	c := make(chan error)

	go t.client.Disconnect(c, 1000)

	return t.waitForError(ctx, c)
}

// Publish will publish the given payload to the given topic with the given quality of service level
func (t *topicRepository) Publish(ctx context.Context, topic string, qos uint8, payload interface{}) error {
	if !t.IsConnected(ctx) {
		return entity.ErrConnected
	}

	c := make(chan error)

	go t.client.Publish(c, topic, qos, payload)

	return t.waitForError(ctx, c)
}

// Subscribe will subscribe to the given topic with the given quality of service level and message handler
func (t *topicRepository) Subscribe(ctx context.Context, topic string, qos uint8, callback topic.CallbackHandler) error {
	if !t.IsConnected(ctx) {
		return entity.ErrConnected
	}

	c := make(chan error)

	go t.client.Subscribe(c, topic, qos, callback)

	return t.waitForError(ctx, c)
}

// Unsubscribe will unsubscribe from the given topic
func (t *topicRepository) Unsubscribe(ctx context.Context, topic string) error {
	if !t.IsConnected(ctx) {
		return entity.ErrNotConnected
	}

	c := make(chan error)

	go t.client.Unsubscribe(c, topic)

	return t.waitForError(ctx, c)
}

func (t *topicRepository) waitForError(ctx context.Context, c chan error) error {
	var result error
	select {
	case <-ctx.Done():
		log.Println("Context is done (topic Publish)")
		return entity.ErrCtxDone
	case result = <-c:
		log.Println("Went to client successfuly :)")
		return result
	}
}
