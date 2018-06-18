package usecase

import (
	"context"
	"time"

	entity "github.com/blanvam/rasp-garden/entities"
	"github.com/blanvam/rasp-garden/resource"
)

type resourceUsecase struct {
	repository     resource.Repository
	contextTimeout time.Duration
}

// NewResourceUsecase interface
func NewResourceUsecase(r resource.Repository, timeout time.Duration) resource.Usecase {
	return &resourceUsecase{
		repository:     r,
		contextTimeout: timeout,
	}
}

func (r *resourceUsecase) All(c context.Context) ([]*entity.Resource, error) {

	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	listResources, err := r.repository.All(ctx)
	if err != nil {
		return nil, err
	}

	return listResources, nil
}

func (r *resourceUsecase) GetByID(c context.Context, id int) (*entity.Resource, error) {

	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	res, err := r.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *resourceUsecase) Update(c context.Context, re *entity.Resource) (*entity.Resource, error) {

	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	re.UpdatedAt = time.Now()
	return r.repository.Update(ctx, re)
}

func (r *resourceUsecase) Store(c context.Context, re *entity.Resource) (*entity.Resource, error) {

	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	existedResource, _ := r.GetByID(ctx, re.Pin)
	if existedResource != nil {
		return nil, entity.ErrConflict
	}

	_, err := r.repository.Store(ctx, re)
	if err != nil {
		return nil, err
	}

	return re, nil
}

func (r *resourceUsecase) Delete(c context.Context, id int) (bool, error) {

	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	existedResource, _ := r.GetByID(ctx, id)
	if existedResource != nil {
		return false, entity.ErrNotFound
	}

	return r.repository.Delete(ctx, id)
}