package topic

import "context"

// Repository interface definition for a respository
type Repository interface {
	IsConnected() bool
	Connect(ctx context.Context, servers ...string) error
	Disconnect(ctx context.Context) error
	Publish(ctx context.Context, topic string, qos uint8, payload interface{}) error
	Subscribe(c chan []byte, ctx context.Context, topic string, qos uint8) error
	Unsubscribe(ctx context.Context, topic string) error
}
