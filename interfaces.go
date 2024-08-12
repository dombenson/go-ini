package ini

import "io"

type Getter interface {
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
	// Lists the sections in the file
	Sections() (value []string)
	// Lists the values in a section the file
	Values(section string) (value map[string]string)

	// ParseEnvironmentVariables runs the ini file's values through Go Template with the environment variables available
	// in the format '{{ .Env.ENV_VAR_NAME }}'.
	ParseEnvironmentVariables() error
}

type Copier interface {
	// Copy loaded data to a writer
	Copy(Setter)
}

type Setter interface {
	// Set the value for a key in a section, along with a boolean result similar to a map lookup.
	Set(section, key, value string) bool
	// Set a key in a section to an integer value
	SetInt(section, key string, value int) bool
	// Set a key in a section to a boolean value
	SetBool(section, key string, value bool) bool
	// Set a key in a section to a string slice
	SetArr(section, key string, value []string) bool
}

// A Reader is able to load and extract data from an io.Reader
type Reader interface {
	io.ReaderFrom
	Getter
	Copier
}

// A Writer is able to set and write data to an io.Writer
type Writer interface {
	io.WriterTo
	Setter
	// Remove a key from a section (OK if it does not exist)
	Remove(section, key string)
	// RemoveSection removes a whole section from an ini file (OK if it does not exist)
	RemoveSection(section string)
}

// A ReadWriter is able to load, get, modify and save data
type ReadWriter interface {
	Reader
	Writer
}

// A StreamReader can additionally accept data by being used as an io.Writer
type StreamReader interface {
	Reader
	io.Writer
}

// A StreamWriter can additionally be passed as an io.Reader;
type StreamWriter interface {
	Writer
	io.ReadCloser
}

// A StreamReadWriter can get/set data and be treated as an io.ReadWriteCloser
type StreamReadWriter interface {
	StreamReader
	StreamWriter
}

type File interface {
	StreamReadWriter
}
