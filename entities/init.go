package entities

// ResourceKind is a string that determines the resource' kind
type ResourceKind string

// ResourceStatus is a string that determines the resource' state
type ResourceStatus string

const (
	// ResourceKindIn refer to input pin
	ResourceKindIn = ResourceKind("in")

	// ResourceKindOut refer to output pin
	ResourceKindOut = ResourceKind("out")

	// ResourceStatusOpen refer to open pin status
	ResourceStatusOpen = ResourceStatus("open")

	// ResourceStatusClosed refer to closed pin status
	ResourceStatusClosed = ResourceStatus("closed")
)

// ResourceList is a list of Resource instances
type ResourceList []Resource
