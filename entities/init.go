package entities

// resourceKind is a string that determines the resource' kind
type resourceKind string

// resourceStatus is a string that determines the resource' state
type resourceStatus string

const (
	// ResourceKindIn refer to input pin
	ResourceKindIn = resourceKind("in")

	// ResourceKindOut refer to output pin
	ResourceKindOut = resourceKind("out")

	// ResourceStatusOpen refer to open pin status
	ResourceStatusOpen = resourceStatus("open")

	// ResourceStatusClosed refer to closed pin status
	ResourceStatusClosed = resourceStatus("closed")
)

// ResourceList is a list of Resource instances
type ResourceList []Resource
