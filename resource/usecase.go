package resource

import (
	"context"

	entity "github.com/blanvam/rasp-garden/entities"
)

// Usecase interface definition for case of use
type Usecase interface {
	All(ctx context.Context) ([]*entity.Resource, error)
	GetByID(ctx context.Context, id int) (*entity.Resource, error)
	Update(ctx context.Context, r *entity.Resource) (*entity.Resource, error)
	Store(ctx context.Context, r *entity.Resource) (*entity.Resource, error)
	Delete(ctx context.Context, id int) (bool, error)
}
