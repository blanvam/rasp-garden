package entities

// resourceKind is a string that determines the resource' kind
type resourceKind string

// resourceStatus is a string that determines the resource' state
type resourceStatus string

const (
	resourceKindIn       = resourceKind("in")
	resourceKindOut      = resourceKind("out")
	resourceStatusOpen   = resourceStatus("open")
	resourceStatusClosed = resourceStatus("closed")
)

// ResourceList is a list of Resource instances
type ResourceList []Resource
