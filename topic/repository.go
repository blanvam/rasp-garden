package topic

import (
	"context"

	broker "github.com/blanvam/rasp-garden/broker"
)

// Repository interface definition for a respository
type Repository interface {
	IsConnected(ctx context.Context) bool
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Publish(ctx context.Context, topic string, qos uint8, bytes []byte) error
	Subscribe(ctx context.Context, topic string, qos uint8, callback broker.CallbackHandler) error
	Unsubscribe(ctx context.Context, topic string) error
}
