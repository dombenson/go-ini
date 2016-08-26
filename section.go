package ini

import (
	"strconv"
	"strings"
)


// A Section represents a single section of an INI file.
type section struct {
	stringValues stringSection
	arrayValues  arraySection
}

// All ini settings for a section except arrays are stored in this
// Helper methods like GetInt parse entries in this map
type stringSection map[string]string

// Used for storing array values for a section
type arraySection map[string][]string

func makeSection(values stringSection) *section {
	return &section{stringValues: values, arrayValues: map[string][]string{}}
}

// Looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
func (s *section) Get(key string) (value string, ok bool) {
	value, ok = s.stringValues[key]
	return
}

// Looks up a value for a key in this section and attempts to parse that value as a boolean, along with a boolean result similar to a map lookup.
// The `ok` boolean will be false in the event that the value could not be parsed as a bool
func (s *section) GetBool(key string) (value bool, ok bool) {
	rawValue, ok := s.Get(key)
	if !ok {
		return
	}
	ok = true
	lowerCase := strings.ToLower(rawValue)
	switch lowerCase {
	case "", "0", "false", "no":
		value = false
	case "1", "true", "yes":
		value = true
	default:
		ok = false
	}
	return
}

// Looks up a value for a key in this section and attempts to parse that value as an integer, along with a boolean result similar to a map lookup.
// The `ok` boolean will be false in the event that the value could not be parsed as an int
func (s *section) GetInt(key string) (value int, ok bool) {
	rawValue, ok := s.Get(key)
	if !ok {
		return
	}
	ok = false
	value, err := strconv.Atoi(rawValue)
	if err != nil {
		return
	}
	ok = true
	return
}

// Looks up a value for an array key in a section and returns that value, along with a boolean result similar to a map lookup.
func (s *section) GetArr(key string) (value []string, ok bool) {
	value, ok = s.arrayValues[key]
	return
}

func (s *section) Set(key string, value string) (ok bool) {
	s.stringValues[key] = value
	return true
}

func (s *section) SetArr(key string, value []string) (ok bool) {
	s.arrayValues[key] = value
	return true
}

func (s *section) SetInt(key string, value int) (ok bool) {
	ok = s.Set(key, strconv.Itoa(value))
	return
}

func (s *section) SetBool(key string, value bool) (ok bool) {
	var useVal string
	if(value) {
		useVal = "true"
	} else {
		useVal = "false"
	}
	return s.Set(key, useVal)
}

func (s *section) Remove(key string) {
	_, found := s.stringValues[key]
	if(found) {
		delete(s.stringValues, key)
	}
	_, found = s.arrayValues[key]
	if(found) {
		delete(s.arrayValues, key)
	}
}