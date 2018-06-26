package database

import (
	"io"
	"strconv"

	"github.com/blanvam/rasp-garden/resource"
	"github.com/peterbourgon/diskv"
)

type diskvDatabase struct {
	connection *diskv.Diskv
}

// NewDiskvDatabase aaa
func NewDiskvDatabase(Conn *diskv.Diskv) resource.Database {
	return &diskvDatabase{
		connection: Conn,
	}
}

// Read function to read from diskv database
func (d *diskvDatabase) Read(c chan []byte, id int) {
	result, err := d.connection.Read(strconv.Itoa(id))
	if err == nil {
		c <- result
	}
	c <- nil
}

// Write function to write data with key into diskv database
func (d *diskvDatabase) Write(c chan bool, id int, r io.Reader) {
	err := d.connection.WriteStream(strconv.Itoa(id), r, false)
	c <- err == nil
}

// Delete the id on diskv database
func (d *diskvDatabase) Delete(c chan error, id int) {
	err := d.connection.Erase(strconv.Itoa(id))
	c <- err
}
