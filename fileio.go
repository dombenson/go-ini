package ini

import (
	"io"
	"bufio"
	"os"
	"sort"
	"fmt"
	"strings"
	"errors"
)



// Loads INI data from a reader and stores the data in the File.
func (f File) ReadFrom(in io.Reader) (n int64, err error) {
	n = 0
	scanner := bufio.NewScanner(in)
	n, err = parseFile(scanner, f)
	return
}

// Loads INI data from a named file and stores the data in the File.
func (f File) LoadFile(file string) (err error) {
	in, err := os.Open(file)
	if err != nil {
		return
	}
	defer in.Close()
	_, err = f.ReadFrom(in)
	return
}


// Write out an INI File representing the current state to a writer.
func (f File) WriteTo(out io.Writer) (n int64, err error) {
	orderedSections := make([]string, len(f.sections))
	counter := 0
	n = 0
	thisWrite := 0
	for section, _ := range f.sections {
		orderedSections[counter] = section
		counter++
	}
	sort.Strings(orderedSections)
	for _, section := range orderedSections {
		options := f.sections[section]
		thisWrite, err = fmt.Fprintln(out, "["+section+"]")
		n += int64(thisWrite)
		if(err) != nil {
			return
		}
		orderedStringKeys := make([]string, len(options.stringValues))
		counter = 0
		for key, _ := range options.stringValues {
			orderedStringKeys[counter] = key
			counter++
		}
		sort.Strings(orderedStringKeys)
		for _, key := range orderedStringKeys {
			thisWrite, err = fmt.Fprintln(out, key, "=", options.stringValues[key])
			n += int64(thisWrite)
			if(err) != nil {
				return
			}
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
				thisWrite, err = fmt.Fprintln(out, key, "[]=", value)
				n += int64(thisWrite)
				if(err) != nil {
					return
				}
			}
		}
		thisWrite, err = fmt.Fprintln(out)
		n += int64(thisWrite)
		if(err) != nil {
			return
		}
	}
	return
}

func (f File) Write(p []byte) (n int, err error) {
	reader := strings.NewReader(string(p))
	var m int64 = 0
	m, err = f.ReadFrom(reader)
	n = int(m)
	if(n != len(p) && err == nil) {
		err = errors.New("Internal error: failed to write")
	}
	return
}

