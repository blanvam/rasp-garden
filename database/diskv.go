package database

import (
	"io"
	"strconv"

	"github.com/peterbourgon/diskv"
)

type diskvDatabase struct {
	connection *diskv.Diskv
}

// NewDiskvDatabase aaa
func NewDiskvDatabase(bdPath string) Database {
	conn := getdb(bdPath)
	return &diskvDatabase{
		connection: conn,
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

func getdb(bdPath string) *diskv.Diskv {
	flatTransform := func(s string) []string { return []string{} }
	db := diskv.New(diskv.Options{
		BasePath:     bdPath,
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})
	return db
}
