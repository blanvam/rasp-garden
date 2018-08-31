package usecase

import (
	"context"
	"encoding/json"
	"time"

	entity "github.com/blanvam/rasp-garden/entities"
	"github.com/blanvam/rasp-garden/resource"
	"github.com/blanvam/rasp-garden/utils"
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

func (r *resourceUsecase) Bind(ctx context.Context, request *entity.ResourceRequest) (*entity.Resource, error) {
	resource := &entity.Resource{
		Name:        request.Name,
		Description: request.Description,
		Pin:         request.Pin,
		Kind:        request.Kind,
		Status:      entity.ResourceStatusClosed,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return resource, nil
}

func (r *resourceUsecase) BindBytes(c context.Context, payload []byte) (*entity.Resource, error) {

	resource := &entity.Resource{}
	if err := json.Unmarshal(payload, &resource); err != nil {
		return nil, entity.ErrBrokerReceived
	}
	return resource, nil
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

func (r *resourceUsecase) Update(c context.Context, re *entity.Resource) error {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	existedResource, _ := r.GetByID(ctx, re.Pin)
	if existedResource == nil {
		return entity.ErrNotFound
	}
	re.CreatedAt = existedResource.CreatedAt
	_, err := r.repository.Update(ctx, re)
	if existedResource.Status != re.Status {
		r.pinChange(re, existedResource.Status)
	}
	return err
}

func (r *resourceUsecase) Store(c context.Context, re *entity.Resource) error {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	existedResource, _ := r.GetByID(ctx, re.Pin)
	if existedResource != nil {
		return entity.ErrConflict
	}

	_, err := r.repository.Store(ctx, re)
	r.pinChange(re, re.Status)
	return err
}

func (r *resourceUsecase) Delete(c context.Context, id int) (bool, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	existedResource, _ := r.GetByID(ctx, id)
	if existedResource == nil {
		return false, entity.ErrNotFound
	}
	r.pinChange(existedResource, entity.ResourceStatusOpen)
	return r.repository.Delete(ctx, id)
}

func (r *resourceUsecase) pinChange(re *entity.Resource, a entity.ResourceStatus) error {
	if re.Kind == entity.ResourceKindIn {
		return utils.GPIOInChange(re.Pin, re.Status == entity.ResourceStatusClosed)
	}
	return utils.GPIOOutChange(re.Pin, re.Status == entity.ResourceStatusClosed)
}
