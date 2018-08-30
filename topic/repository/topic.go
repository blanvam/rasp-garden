package repository

import (
	"bytes"
	"context"
	"encoding/gob"
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
		return result == true
	}
}

// Connect return error if the client cannot connect to the server.
func (t *topicRepository) Connect(ctx context.Context) error {
	if t.IsConnected(ctx) {
		return nil
	}

	c := make(chan error)

	go t.client.Connect(c)

	return t.waitForError(ctx, c)
}

// Disconnect will disconnect the client
func (t *topicRepository) Disconnect(ctx context.Context) error {
	if !t.IsConnected(ctx) {
		return nil
	}

	c := make(chan error)

	go t.client.Disconnect(c, 1000)

	return t.waitForError(ctx, c)
}

// Publish will publish the given payload
func (t *topicRepository) Publish(ctx context.Context, topic string, qos uint8, msg *entity.Message) error {
	if !t.IsConnected(ctx) {
		t.Connect(ctx)
	}

	msgByte := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(msgByte)
	err := encoder.Encode(msg)
	if err != nil {
		return err
	}

	c := make(chan error)

	go t.client.Publish(c, topic, qos, msgByte.Bytes())

	return t.waitForError(ctx, c)
}

// Subscribe will subscribe to the given topic
func (t *topicRepository) Subscribe(ctx context.Context, topic string, qos uint8, callback topic.CallbackHandler) error {
	if !t.IsConnected(ctx) {
		t.Connect(ctx)
	}

	c := make(chan error)

	go t.client.Subscribe(c, topic, qos, callback)

	return t.waitForError(ctx, c)
}

// Unsubscribe will unsubscribe from the given topic
func (t *topicRepository) Unsubscribe(ctx context.Context, topic string) error {
	if !t.IsConnected(ctx) {
		t.Connect(ctx)
	}

	c := make(chan error)

	go t.client.Unsubscribe(c, topic)

	return t.waitForError(ctx, c)
}

func (t *topicRepository) waitForError(ctx context.Context, c chan error) error {
	var result error
	select {
	case <-ctx.Done():
		log.Println("Context is done")
		return entity.ErrCtxDone
	case result = <-c:
		return result
	}
}
