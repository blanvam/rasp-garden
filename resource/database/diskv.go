package database

import (
	"io"
	"strconv"

	"github.com/peterbourgon/diskv"
)

func getdb() *diskv.Diskv {
	flatTransform := func(s string) []string { return []string{} }
	db := diskv.New(diskv.Options{
		BasePath:     "my-data-dir",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})

	return db
}

// Read function to read from diskv database
func Read(c chan []byte, id int) {
	bd := getdb()
	result, err := bd.Read(strconv.Itoa(id))
	if err == nil {
		c <- result
	}
	c <- nil
}

// Write function to write data with key into diskv database
func Write(c chan bool, id int, r io.Reader) {
	bd := getdb()
	err := bd.WriteStream(strconv.Itoa(id), r, false)
	c <- err != nil
}

// Delete the id on diskv database
func Delete(c chan bool, id int) {
	bd := getdb()
	err := bd.Erase(strconv.Itoa(id))
	c <- err != nil
}
