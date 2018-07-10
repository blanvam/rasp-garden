package topic

import (
	"context"

	entity "github.com/blanvam/rasp-garden/entities"
)

// Usecase interface definition for case of use
type Usecase interface {
	Read(ctx context.Context) ([]*entity.Topic, error)
	Write(ctx context.Context, id int) (*entity.Topic, error)
	Do(ctx context.Context, r *entity.Topic) (bool, error)
}
