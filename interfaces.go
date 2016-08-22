package ini

import "io"


// An Ini instance represents a parsed INI file.
type Ini interface {
	// Looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
	Get(section, key string) (value string, ok bool)
	// Looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
	// The `ok` boolean will be false in the event that the value could not be parsed as an int
	GetInt(section, key string) (value int, ok bool)
	// Looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
	// The `ok` boolean will be false in the event that the value could not be parsed as a bool
	GetBool(section, key string) (value bool, ok bool)
	// Looks up a value for an array key in a section and returns that value, along with a boolean result similar to a map lookup.
	GetArr(section, key string) (value []string, ok bool)
	// Set the value for a key in a section, along with a boolean result similar to a map lookup.
	Set(section, key, value string) bool
	// Set a key in a section to an integer value
	SetInt(section, key string, value int) bool
	// Set a key in a section to a boolean value
	SetBool(section, key string, value bool) bool
}

type readerInterface interface {
	io.WriterTo
	io.Writer
}
type writerInterface interface {
	io.ReaderFrom
//	io.Reader
}

type Reader interface {
	Ini
	readerInterface
}


type Writer interface {
	Ini
	writerInterface
}

type ReadWriter interface {
	Reader
	writerInterface
}
