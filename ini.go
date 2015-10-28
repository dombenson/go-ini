// Package ini provides functions for parsing INI configuration files.
package ini

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	sectionRegex   = regexp.MustCompile(`^\[(.*)\]$`)
	assignArrRegex = regexp.MustCompile(`^([^=\[\]]+)\[\][^=]*=(.*)$`)
	assignRegex    = regexp.MustCompile(`^([^=]+)=(.*)$`)
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
	StringValues StringSection
	ArrayValues  ArraySection
}

// All ini settings for a section except arrays are stored in this
// Helper methods like GetInt parse entries in this map
type StringSection map[string]string

// Used for storing array values for a section
type ArraySection map[string][]string

func makeSection(values StringSection) *Section {
	return &Section{StringValues: values, ArrayValues: map[string][]string{}}
}

// Returns a named Section. A Section will be created if one does not already exist for the given name.
func (f File) Section(name string) *Section {
	section := f[name]
	if section == nil {
		section = &Section{StringValues: make(map[string]string), ArrayValues: make(map[string][]string)}
		f[name] = section
	}
	return section
}

// Looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
func (f File) Get(section, key string) (value string, ok bool) {
	return f.Section(section).Get(key)
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

// Looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
func (s *Section) Get(key string) (value string, ok bool) {
	value, ok = s.StringValues[key]
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
	value, ok = s.ArrayValues[key]
	return
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
			key, val = strings.TrimSpace(key), strings.TrimSpace(val)
			curVal, ok := file.Section(section).ArrayValues[key]
			if ok {
				file.Section(section).ArrayValues[key] = append(curVal, val)
			} else {
				file.Section(section).ArrayValues[key] = make([]string, 1, 4)
				file.Section(section).ArrayValues[key][0] = val
			}
		} else if groups := assignRegex.FindStringSubmatch(line); groups != nil {
			key, val := groups[1], groups[2]
			key, val = strings.TrimSpace(key), strings.TrimSpace(val)
			file.Section(section).StringValues[key] = val
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

// Write writes INI File into a writer.
func Write(out io.Writer, file File) {
	for section, options := range file {
		fmt.Fprintln(out, "["+section+"]")
		for key, value := range options {
			fmt.Fprintln(out, key, "=", value)
		}
		fmt.Fprintln(out)
	}
}

// WriteFile writes INI File into a file.
func WriteFile(filename string, file File) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	Write(f, file)
	return nil
}
