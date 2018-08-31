package entities

import (
	"time"
)

// Resource is a entity which holds information about the active resources
type Resource struct {
	Name        string
	Description string
	Pin         int
	Kind        ResourceKind
	Status      ResourceStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ResourceRequest is a struct that resembles a request performed to edit or create a resource instance
type ResourceRequest struct {
	Name        string
	Description string
	Pin         int
	Kind        ResourceKind
}
