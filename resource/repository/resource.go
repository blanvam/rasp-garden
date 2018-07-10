package repository

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	entity "github.com/blanvam/rasp-garden/entities"
	"github.com/blanvam/rasp-garden/resource"
)

type resourceRepository struct {
	database resource.Database
	minPin   int
	maxPin   int
}

// NewResourceRepository aaa
func NewResourceRepository(bd resource.Database, min int, max int) resource.Repository {
	return &resourceRepository{
		database: bd,
		minPin:   min,
		maxPin:   max,
	}
}

func (r *resourceRepository) save(ctx context.Context, re *entity.Resource) (int, error) {
	re.UpdatedAt = time.Now()
	b := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(b)
	err := encoder.Encode(re)
	if err != nil {
		return re.Pin, err
	}

	c := make(chan bool)
	var result bool

	go r.database.Write(c, re.Pin, b)

	select {
	case <-ctx.Done():
		log.Println("Context is done, why? I dont now")
		return re.Pin, ctx.Err()
	case result = <-c:
		log.Println("Went to db successfuly :)")
		if result == false {
			return re.Pin, entity.ErrStore
		}
	}

	return re.Pin, nil
}

func (r *resourceRepository) GetByID(ctx context.Context, id int) (*entity.Resource, error) {
	c := make(chan []byte)
	var result []byte

	go r.database.Read(c, id)

	select {
	case <-ctx.Done():
		log.Println("Context is done, why? I dont now")
		return nil, ctx.Err()
	case result = <-c:
		resource := &entity.Resource{}
		decoder := gob.NewDecoder(bytes.NewBuffer(result))
		err := decoder.Decode(resource)
		if err != nil {
			return nil, entity.ErrNotFound
		}
		log.Println(fmt.Sprintf("Resource %s in pin %d successfuly got from bd", resource.Name, resource.Pin))
		return resource, nil
	}
}

func (r *resourceRepository) All(ctx context.Context) ([]*entity.Resource, error) {
	listResources := []*entity.Resource{}
	for i := r.minPin; i <= r.maxPin; i++ {
		existedResource, _ := r.GetByID(ctx, i)
		if existedResource != nil {
			listResources = append(listResources, existedResource)
		}
	}

	return listResources, nil
}

func (r *resourceRepository) Update(ctx context.Context, re *entity.Resource) (int, error) {
	return r.save(ctx, re)
}

func (r *resourceRepository) Store(ctx context.Context, re *entity.Resource) (int, error) {
	if r.minPin > re.Pin || re.Pin > r.maxPin {
		return re.Pin, fmt.Errorf("Resource PIN must be between %d and %d", r.minPin, r.maxPin)
	}
	return r.save(ctx, re)
}

func (r *resourceRepository) Delete(ctx context.Context, id int) (bool, error) {
	c := make(chan error)
	go r.database.Delete(c, id)

	select {
	case <-ctx.Done():
		log.Println("Context is done, why? I dont now")
		return false, ctx.Err()
	case result := <-c:
		log.Println("Went to db successfuly :)")
		return result == nil, result
	}
}
