package entities

// ResourceRequest is a struct that resembles a request performed to edit or create a resource instance
type ResourceRequest struct {
	Name        string
	Description string
	Pin         int
	Kind        resourceKind
}
