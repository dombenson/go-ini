package ini

import "io"

// This implements the full ini.StreamReadWriter interface
type File struct {
	sections map[string]*section
	reader io.Reader
}


// Returns a named Section. A Section will be created if one does not already exist for the given name.
func (f *File) section(name string) *section {
	theSection := f.sections[name]
	if theSection == nil {
		theSection = &section{stringValues: make(map[string]string), arrayValues: make(map[string][]string)}
		f.sections[name] = theSection
	}
	return theSection
}

// Looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
func (f *File) Get(section, key string) (value string, ok bool) {
	return f.section(section).Get(key)
}

// Set the value for a key in a section, along with a boolean result similar to a map lookup.
func (f *File) Set(section, key string, value string) (ok bool) {
	return f.section(section).Set(key, value)
}

// Set a key in a section to an integer value
func (f *File) SetInt(section, key string, value int) (ok bool) {
	return f.section(section).SetInt(key, value)
}

// Set a key in a section to a boolean value
func (f *File) SetBool(section, key string, value bool) (ok bool) {
	return f.section(section).SetBool(key, value)
}

// Looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
// The `ok` boolean will be false in the event that the value could not be parsed as an int
func (f *File) GetInt(section, key string) (value int, ok bool) {
	return f.section(section).GetInt(key)
}


// Looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
// The `ok` boolean will be false in the event that the value could not be parsed as a bool
func (f *File) GetBool(section, key string) (value bool, ok bool) {
	return f.section(section).GetBool(key)
}

// Looks up a value for an array key in a section and returns that value, along with a boolean result similar to a map lookup.
func (f *File) GetArr(section, key string) (value []string, ok bool) {
	return f.section(section).GetArr(key)
}

func(f *File) Remove(section, key string) {
	f.section(section).Remove(key)
}


func(f *File) RemoveSection(section string) {
	_, found := f.sections[section]
	if(found) {
		delete(f.sections, section)
	}
}
