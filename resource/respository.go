package resource

import (
	"context"

	entity "github.com/blanvam/rasp-garden/entities"
)

// Repository interface definition for a respository
type Repository interface {
	All(ctx context.Context) ([]*entity.Resource, error)
	GetByID(ctx context.Context, id int) (*entity.Resource, error)
	Update(ctx context.Context, r *entity.Resource) (int, error)
	Store(ctx context.Context, r *entity.Resource) (int, error)
	Delete(ctx context.Context, id int) (bool, error)
}
