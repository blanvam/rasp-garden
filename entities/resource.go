package entities

import "time"

// Resource is a entity which holds information about the active resources
type Resource struct {
	Name        string
	Description string
	Pin         int
	Kind        resourceKind
	Status      resourceStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
