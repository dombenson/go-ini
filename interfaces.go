package ini

import "io"


// An instance is able to load and extract data from a reader
type Reader interface {
	io.ReaderFrom
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
}

// An instance is able to set and write data to a writer
type Writer interface {
	io.WriterTo
	// Set the value for a key in a section, along with a boolean result similar to a map lookup.
	Set(section, key, value string) bool
	// Set a key in a section to an integer value
	SetInt(section, key string, value int) bool
	// Set a key in a section to a boolean value
	SetBool(section, key string, value bool) bool
}

// An instance is able to load, get, modify and save data
type ReadWriter interface {
	Reader
	Writer
}


// An instance can additionally accept data by being used as an io.Writer
type StreamReader interface {
	Reader
	io.Writer

}

// An instance can additionally be passed as an io.Reader;
type StreamWriter interface {
	Writer
	io.ReadCloser
}


type StreamReadWriter interface {
	StreamReader
	StreamWriter
}
