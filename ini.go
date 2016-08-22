// Package ini provides functions for parsing INI configuration files.
package ini

import (
	"io"
	"os"
)



func NewFile() *File {
	file := File{}
	file.sections = make(map[string]*section)
	return &file
}

// Loads and returns a File from a reader.
func Load(in io.Reader) (*File, error) {
	file := NewFile()
	_, err := file.ReadFrom(in)
	return file, err
}

// Loads and returns a File from a reader.
func LoadFile(filename string) (file *File, err error) {
	file = nil
	fh, err := os.Open(filename)
	if(err != nil) {
		return
	}
	defer fh.Close()
	return Load(fh)
}
