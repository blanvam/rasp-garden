package usecase

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	entity "github.com/blanvam/rasp-garden/entities"
	"github.com/blanvam/rasp-garden/topic"
)

type topicUsecase struct {
	repository     topic.Repository
	qos            uint8
	contextTimeout time.Duration
}

// NewTopicUsecase interface
func NewTopicUsecase(r topic.Repository, qos uint8, timeout time.Duration) topic.Usecase {
	topicUsecase := &topicUsecase{
		repository:     r,
		qos:            qos,
		contextTimeout: timeout,
	}

	ctx, cancel := context.WithTimeout(context.Background(), topicUsecase.contextTimeout)
	defer cancel()

	topicUsecase.repository.Connect(ctx)

	return topicUsecase
}

func (t *topicUsecase) Publish(c context.Context, topic string, msg *entity.Message) error {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()

	err := t.repository.Publish(ctx, topic, t.qos, msg)
	if err != nil {
		return err
	}

	return nil
}

func (t *topicUsecase) Subscribe(c context.Context, topic string) error {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()

	callback := func(topic string, id string, payload []byte) {
		message := &entity.Message{}
		decoder := gob.NewDecoder(bytes.NewBuffer(payload))
		err := decoder.Decode(message)
		if err != nil {
			log.Println("Error decoding message payload from broker")
		}
		log.Println(fmt.Sprintf("Message '%d' with content '%s' successfuly received from broker", message.ID, message.Content))
	}

	err := t.repository.Subscribe(ctx, topic, t.qos, callback)
	return err
}

func (t *topicUsecase) Unsubscribe(c context.Context, topic string) error {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()

	err := t.repository.Unsubscribe(ctx, topic)

	return err
}
