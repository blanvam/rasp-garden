package resource

import "io"

// Database interface definition for a resource
type Database interface {
	Read(c chan []byte, id int)
	Write(c chan bool, id int, r io.Reader)
	Delete(c chan error, id int)
}
