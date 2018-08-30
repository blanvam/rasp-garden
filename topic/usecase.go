package topic

import (
	"context"

	entity "github.com/blanvam/rasp-garden/entities"
)

// Usecase interface definition for case of use
type Usecase interface {
	Publish(ctx context.Context, topic string, msg *entity.Message) error
	Subscribe(ctx context.Context, topic string) error
	Unsubscribe(ctx context.Context, topic string) error
}
