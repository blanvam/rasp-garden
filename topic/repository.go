package topic

import (
	"context"
)

// Repository interface definition for a respository
type Repository interface {
	IsConnected(ctx context.Context) bool
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Publish(ctx context.Context, topic string, qos uint8, payload interface{}) error
	Subscribe(ctx context.Context, topic string, qos uint8, callback CallbackHandler) error
	Unsubscribe(ctx context.Context, topic string) error
}
