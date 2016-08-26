package ini

import "io"

// This implements the full ini.StreamReadWriter interface
type file struct {
	sections map[string]*section
	reader io.Reader
}


// Returns a named Section. A Section will be created if one does not already exist for the given name.
func (f *file) section(name string) *section {
	theSection := f.sections[name]
	if theSection == nil {
		theSection = &section{stringValues: make(map[string]string), arrayValues: make(map[string][]string)}
		f.sections[name] = theSection
	}
	return theSection
}

// Looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
func (f *file) Get(section, key string) (value string, ok bool) {
	return f.section(section).Get(key)
}

// Set the value for a key in a section, along with a boolean result similar to a map lookup.
func (f *file) Set(section, key string, value string) (ok bool) {
	return f.section(section).Set(key, value)
}

// Set a key in a section to an integer value
func (f *file) SetInt(section, key string, value int) (ok bool) {
	return f.section(section).SetInt(key, value)
}

// Set a key in a section to a boolean value
func (f *file) SetBool(section, key string, value bool) (ok bool) {
	return f.section(section).SetBool(key, value)
}
// Set a key in a section to an array
func (f *file) SetArr(section, key string, value []string) (ok bool) {
	return f.section(section).SetArr(key, value)
}

// Looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
// The `ok` boolean will be false in the event that the value could not be parsed as an int
func (f *file) GetInt(section, key string) (value int, ok bool) {
	return f.section(section).GetInt(key)
}


// Looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
// The `ok` boolean will be false in the event that the value could not be parsed as a bool
func (f *file) GetBool(section, key string) (value bool, ok bool) {
	return f.section(section).GetBool(key)
}

// Looks up a value for an array key in a section and returns that value, along with a boolean result similar to a map lookup.
func (f *file) GetArr(section, key string) (value []string, ok bool) {
	return f.section(section).GetArr(key)
}

func(f *file) Remove(section, key string) {
	f.section(section).Remove(key)
}


func(f *file) RemoveSection(section string) {
	_, found := f.sections[section]
	if(found) {
		delete(f.sections, section)
	}
}

func (f *file) Copy(w Setter) {
	for secName, sec := range f.sections {
		for keyName, val := range sec.stringValues {
			w.Set(secName, keyName, val)
		}
		for keyName, arVal := range sec.arrayValues {
			w.SetArr(secName, keyName, arVal)
		}
	}
}