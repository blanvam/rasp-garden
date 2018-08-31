package resource

import (
	"context"

	entity "github.com/blanvam/rasp-garden/entities"
)

// Usecase interface definition for case of use
type Usecase interface {
	Bind(ctx context.Context, rq *entity.ResourceRequest) (*entity.Resource, error)
	BindBytes(ctx context.Context, payload []byte) (*entity.Resource, error)
	All(ctx context.Context) ([]*entity.Resource, error)
	GetByID(ctx context.Context, id int) (*entity.Resource, error)
	Update(ctx context.Context, r *entity.Resource) (bool, error)
	Store(ctx context.Context, r *entity.Resource) (bool, error)
	Delete(ctx context.Context, id int) (bool, error)
}
