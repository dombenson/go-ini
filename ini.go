// Package ini provides functions for parsing INI configuration files.
package ini

import (
	"io"
	"os"
)


// Newfile will create and initialise a File object
func NewFile() *File {
	file := File{}
	file.sections = make(map[string]*section)
	return &file
}

// Load creates a File and populates it with data from a reader.
func Load(in io.Reader) (*File, error) {
	file := NewFile()
	_, err := file.ReadFrom(in)
	return file, err
}

// LoadFile creates a File and populates it with data from a file on disk
// This is a convenience helper since it is a very common use case
func LoadFile(filename string) (file *File, err error) {
	file = nil
	fh, err := os.Open(filename)
	if(err != nil) {
		return
	}
	defer fh.Close()
	return Load(fh)
}

// Create a file and populate with data from an existing ini.Reader
func Copy(in Copier) (*File) {
	file := NewFile()
	in.Copy(file)
	return file
}
