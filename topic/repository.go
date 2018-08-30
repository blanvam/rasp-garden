package topic

import (
	"context"

	broker "github.com/blanvam/rasp-garden/broker"
	entity "github.com/blanvam/rasp-garden/entities"
)

// Repository interface definition for a respository
type Repository interface {
	IsConnected(ctx context.Context) bool
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Publish(ctx context.Context, topic string, qos uint8, r *entity.Message) error
	Subscribe(ctx context.Context, topic string, qos uint8, callback broker.CallbackHandler) error
	Unsubscribe(ctx context.Context, topic string) error
}
