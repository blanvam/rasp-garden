package topic

import (
	"context"

	broker "github.com/blanvam/rasp-garden/broker"
)

// Usecase interface definition for case of use
type Usecase interface {
	Publish(ctx context.Context, topic string, i interface{}) error
	Subscribe(ctx context.Context, topic string, callback broker.CallbackHandler) error
	Unsubscribe(ctx context.Context, topic string) error
}
