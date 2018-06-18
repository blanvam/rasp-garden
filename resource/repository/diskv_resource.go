package repository

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	entity "github.com/blanvam/rasp-garden/entities"
	"github.com/blanvam/rasp-garden/resource"
)

// Pins able to control on board
const Pins int = 26

type diskvRepository struct {
	database resource.Database
}

// NewDiskvResourceRepository aaa
func NewDiskvResourceRepository(bd resource.Database) resource.Repository {

	return &diskvRepository{
		database: bd,
	}
}

func (r *diskvRepository) save(ctx context.Context, re *entity.Resource) (int, error) {
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
		fmt.Println("Context is done, why? I dont now")
		return re.Pin, ctx.Err()
	case result = <-c:
		fmt.Println("Went to db successfuly :)")
		if result == false {
			return re.Pin, entity.ErrStore
		}
	}

	return re.Pin, nil
}

func (r *diskvRepository) GetByID(ctx context.Context, id int) (*entity.Resource, error) {
	c := make(chan []byte)
	var result []byte
	var resource entity.Resource
	go r.database.Read(c, id)

	select {
	case <-ctx.Done():
		fmt.Println("Context is done, why? I dont now")
		return nil, ctx.Err()
	case result = <-c:
		fmt.Println("Went to db successfuly :)")
		decoder := gob.NewDecoder(bytes.NewBuffer(result))
		resource := &entity.Resource{}
		if err := decoder.Decode(resource); err != nil {
			return nil, fmt.Errorf("decoding: %v", err)
		}
	}

	return &resource, nil
}

func (r *diskvRepository) All(ctx context.Context) ([]*entity.Resource, error) {
	var listResources []*entity.Resource
	for i := 0; i < Pins; i++ {
		existedResource, _ := r.GetByID(ctx, i)
		listResources = append(listResources, existedResource)
	}

	return listResources, nil
}

func (r *diskvRepository) Update(ctx context.Context, re *entity.Resource) (int, error) {
	existedResource, _ := r.GetByID(ctx, re.Pin)
	if existedResource == nil {
		return re.Pin, entity.ErrNotFound
	}
	return r.save(ctx, re)
}

func (r *diskvRepository) Store(ctx context.Context, re *entity.Resource) (int, error) {
	existedResource, _ := r.GetByID(ctx, re.Pin)
	if existedResource != nil {
		return re.Pin, entity.ErrConflict
	}
	return r.save(ctx, re)
}

func (r *diskvRepository) Delete(ctx context.Context, id int) (bool, error) {
	existedResource, _ := r.GetByID(ctx, id)
	if existedResource != nil {
		return false, entity.ErrNotFound
	}

	c := make(chan bool)
	result := true
	var resource entity.Resource
	go r.database.Delete(c, id)

	select {
	case <-ctx.Done():
		fmt.Println("Context is done, why? I dont now")
		return false, ctx.Err()
	case result = <-c:
		fmt.Println("Went to db successfuly :)")
		if result == false {
			return false, entity.ErrInternalServer
		}
	}
	return result, nil
}
