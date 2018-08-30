package usecase

import (
	"context"
	"time"

	"github.com/blanvam/rasp-garden/broker"

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

func (t *topicUsecase) Publish(c context.Context, topic string, r *entity.Resource) error {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()

	err := t.repository.Publish(ctx, topic, t.qos, r)
	if err != nil {
		return err
	}

	return nil
}

func (t *topicUsecase) Subscribe(c context.Context, topic string, callback broker.CallbackHandler) error {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()

	err := t.repository.Subscribe(ctx, topic, t.qos, callback)
	return err
}

func (t *topicUsecase) Unsubscribe(c context.Context, topic string) error {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()

	err := t.repository.Unsubscribe(ctx, topic)

	return err
}
