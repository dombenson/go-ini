// Package ini provides functions for parsing INI configuration files.
package ini

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	sectionRegex   = regexp.MustCompile(`^\[(.*)\]$`)
	assignArrRegex = regexp.MustCompile(`^([^=\[\]]+)\[\][^=]*=(.*)$`)
	assignRegex    = regexp.MustCompile(`^([^=]+)=(.*)$`)
	quotesRegex    = regexp.MustCompile(`^(['"])(.*)(['"])$`)
)

// ErrSyntax is returned when there is a syntax error in an INI file.
type ErrSyntax struct {
	Line   int
	Source string // The contents of the erroneous line, without leading or trailing whitespace
}

func (e ErrSyntax) Error() string {
	return fmt.Sprintf("invalid INI syntax on line %d: %s", e.Line, e.Source)
}

// A File represents a parsed INI file.
type File map[string]*Section

// A Section represents a single section of an INI file.
type Section struct {
	stringValues StringSection
	arrayValues  ArraySection
}

// All ini settings for a section except arrays are stored in this
// Helper methods like GetInt parse entries in this map
type StringSection map[string]string

// Used for storing array values for a section
type ArraySection map[string][]string

func makeSection(values StringSection) *Section {
	return &Section{stringValues: values, arrayValues: map[string][]string{}}
}

// Returns a named Section. A Section will be created if one does not already exist for the given name.
func (f File) Section(name string) *Section {
	section := f[name]
	if section == nil {
		section = &Section{stringValues: make(map[string]string), arrayValues: make(map[string][]string)}
		f[name] = section
	}
	return section
}

// Looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
func (f File) Get(section, key string) (value string, ok bool) {
	return f.Section(section).Get(key)
}

// Set the value for a key in a section, along with a boolean result similar to a map lookup.
func (f File) Set(section, key string, value string) (ok bool) {
        return f.Section(section).Set(key, value)
}

// Set a key in a section to an integer value
func (f File) SetInt(section, key string, value int) (ok bool) {
	return f.Section(section).SetInt(key, value)
}

// Set a key in a section to a boolean value
func (f File) SetBool(section, key string, value bool) (ok bool) {
        return f.Section(section).SetBool(key, value)
}

// Looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
// The `ok` boolean will be false in the event that the value could not be parsed as an int
func (f File) GetInt(section, key string) (value int, ok bool) {
        return f.Section(section).GetInt(key)
}


// Looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
// The `ok` boolean will be false in the event that the value could not be parsed as a bool
func (f File) GetBool(section, key string) (value bool, ok bool) {
	return f.Section(section).GetBool(key)
}

// Looks up a value for an array key in a section and returns that value, along with a boolean result similar to a map lookup.
func (f File) GetArr(section, key string) (value []string, ok bool) {
	return f.Section(section).GetArr(key)
}

func (s *Section) StringValues() (map[string]string) {
	return s.stringValues
}

func (s *Section) ArrayValues() (map[string][]string) {
        return s.arrayValues
}

// Looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
func (s *Section) Get(key string) (value string, ok bool) {
	value, ok = s.stringValues[key]
	return
}

// Looks up a value for a key in this section and attempts to parse that value as a boolean, along with a boolean result similar to a map lookup.
// The `ok` boolean will be false in the event that the value could not be parsed as a bool
func (s *Section) GetBool(key string) (value bool, ok bool) {
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
func (s *Section) GetInt(key string) (value int, ok bool) {
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
func (s *Section) GetArr(key string) (value []string, ok bool) {
	value, ok = s.arrayValues[key]
	return
}

func (s *Section) Set(key string, value string) (ok bool) {
	s.stringValues[key] = value
	return true
}

func (s *Section) SetInt(key string, value int) (ok bool) {
	ok = s.Set(key, strconv.Itoa(value))
	return
}

func (s *Section) SetBool(key string, value bool) (ok bool) {
	var useVal string
	if(value) {
		useVal = "true"
	} else {
		useVal = "false"
	}
	return s.Set(key, useVal)
}

// Loads INI data from a reader and stores the data in the File.
func (f File) Load(in io.Reader) (err error) {
	bufin, ok := in.(*bufio.Reader)
	if !ok {
		bufin = bufio.NewReader(in)
	}
	return parseFile(bufin, f)
}

// Loads INI data from a named file and stores the data in the File.
func (f File) LoadFile(file string) (err error) {
	in, err := os.Open(file)
	if err != nil {
		return
	}
	defer in.Close()
	return f.Load(in)
}

func trimWithQuotes(inputVal string) (ret string) {
	ret = strings.TrimSpace(inputVal)
	groups := quotesRegex.FindStringSubmatch(ret)
	if groups != nil {
		if (groups[1] == groups[3]) {
			ret = groups[2]
		}
	}
	return
}

func parseFile(in *bufio.Reader, file File) (err error) {
	section := ""
	lineNum := 0
	for done := false; !done; {
		var line string
		if line, err = in.ReadString('\n'); err != nil {
			if err == io.EOF {
				done = true
			} else {
				return
			}
		}
		lineNum++
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			// Skip blank lines
			continue
		}
		if line[0] == ';' || line[0] == '#' {
			// Skip comments
			continue
		}

		if groups := assignArrRegex.FindStringSubmatch(line); groups != nil {
			key, val := groups[1], groups[2]
			key, val = strings.TrimSpace(key), trimWithQuotes(val)
			curVal, ok := file.Section(section).arrayValues[key]
			if ok {
				file.Section(section).arrayValues[key] = append(curVal, val)
			} else {
				file.Section(section).arrayValues[key] = make([]string, 1, 4)
				file.Section(section).arrayValues[key][0] = val
			}
		} else if groups := assignRegex.FindStringSubmatch(line); groups != nil {
			key, val := groups[1], groups[2]
			key, val = strings.TrimSpace(key), trimWithQuotes(val)
			file.Section(section).stringValues[key] = val
		} else if groups := sectionRegex.FindStringSubmatch(line); groups != nil {
			name := strings.TrimSpace(groups[1])
			section = name
			// Create the section if it does not exist
			file.Section(section)
		} else {
			return ErrSyntax{lineNum, line}
		}

	}
	return nil
}

// Loads and returns a File from a reader.
func Load(in io.Reader) (File, error) {
	file := make(File)
	err := file.Load(in)
	return file, err
}

// Loads and returns an INI File from a file on disk.
func LoadFile(filename string) (File, error) {
	file := make(File)
	err := file.LoadFile(filename)
	return file, err
}

// Write out an INI File representing the current state to a writer.
func (f File) Write(out io.Writer) {
	orderedSections := make([]string, len(f))
	counter := 0
	for section, _ := range f {
		orderedSections[counter] = section
		counter++
	}
	sort.Strings(orderedSections)
	for _, section := range orderedSections {
		options := f[section]
		fmt.Fprintln(out, "["+section+"]")
		orderedStringKeys := make([]string, len(options.stringValues))
		counter = 0
		for key, _ := range options.stringValues {
			orderedStringKeys[counter] = key
			counter++
		}
		sort.Strings(orderedStringKeys)
		for _, key := range orderedStringKeys {
			fmt.Fprintln(out, key, "=", options.stringValues[key])
		}
		orderedArrayKeys := make([]string, len(options.arrayValues))
		counter = 0
		for key, _ := range options.arrayValues {
			orderedArrayKeys[counter] = key
			counter++
		}
		sort.Strings(orderedArrayKeys)
		for _, key := range orderedArrayKeys {
			for _, value := range options.arrayValues[key] {
				fmt.Fprintln(out, key, "[]=", value)
			}
		}
		fmt.Fprintln(out)
	}
}

// Write out an INI File representing the current state to a file.
func (f File) WriteFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	f.Write(file)
	return nil
}
