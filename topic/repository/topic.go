package repository

import (
	"context"
	"log"

	entity "github.com/blanvam/rasp-garden/entities"
	"github.com/blanvam/rasp-garden/topic"
	paho "github.com/eclipse/paho.mqtt.golang"
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
func (t *topicRepository) IsConnected() bool {
	if t.client == nil {
		return false
	}
	return t.client.IsConnected()
}

// Disconnect will disconnect from the given MQTT server and clean up all client resources
func (t *topicRepository) Disconnect(ctx context.Context) error {
	if t.IsConnected() {
		t.client.Disconnect(1000)
		t.client = nil
	}
	return nil
}

// Publish will publish the given payload to the given topic with the given quality of service level
func (t *topicRepository) Publish(ctx context.Context, topic string, qos uint8, payload interface{}) error {
	if !t.IsConnected() {
		return entity.ErrNotConnected
	}
	return t.client.Publish(topic, qos, true, payload)
}

// Subscribe will subscribe to the given topic with the given quality of service level and message handler
func (t *topicRepository) Subscribe(c chan []byte, ctx context.Context, topic string, qos uint8) ([]byte, error) {
	if !t.IsConnected() {
		return nil, entity.ErrNotConnected
	}
	handler := func(i paho.Client, message paho.Message) {
		log.Printf("RECEIVED - Topic: %s, Message Length: %d bytes", message.Topic(), len(message.Payload()))
		if c != nil {
			c <- message.Payload()
		}
	}
	return t.client.Subscribe(topic, qos, handler)
}

// Unsubscribe will unsubscribe from the given topic
func (t *topicRepository) Unsubscribe(ctx context.Context, topic string) error {
	if !t.IsConnected() {
		return entity.ErrNotConnected
	}
	return t.client.Unsubscribe(topic)
}
