package ini

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Loads INI data from a reader and stores the data in the File.
func (f *file) ReadFrom(in io.Reader) (n int64, err error) {
	n = 0
	scanner := bufio.NewScanner(in)
	n, err = parseFile(scanner, f)
	return
}

// Loads INI data from a named file and stores the data in the File.
func (f *file) LoadFile(file string) (err error) {
	in, err := os.Open(file)
	if err != nil {
		return
	}
	defer in.Close()
	_, err = f.ReadFrom(in)
	return
}

// Write out an INI File representing the current state to a writer.
func (f *file) WriteTo(out io.Writer) (n int64, err error) {
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
		if (err) != nil {
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
			if (err) != nil {
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
				if (err) != nil {
					return
				}
			}
		}
		thisWrite, err = fmt.Fprintln(out)
		n += int64(thisWrite)
		if (err) != nil {
			return
		}
	}
	return
}

// Load ini data from the bytestream provided
// This is provided so that data can be loaded by treating File as an io.Writer
func (f *file) Write(p []byte) (n int, err error) {
	reader := strings.NewReader(string(p))
	var m int64 = 0
	m, err = f.ReadFrom(reader)
	n = int(m)
	if n != len(p) && err == nil {
		err = errors.New("Internal error: failed to write")
	}
	return
}

// Write out ini data to the bytestream provided
// This is provided so that data can be saved by treating File as an io.Reader
// The returned stream is a consistent representation of the ini file when the read started
// Call close to start a new read from a freshly serialized file
func (f *file) Read(p []byte) (n int, err error) {
	n = 0
	if f.reader == nil {
		buf := new(bytes.Buffer)
		_, err = f.WriteTo(buf)
		if err != nil {
			return
		}
		f.reader = buf
	}
	n, err = f.reader.Read(p)
	return
}

// Close a read stream
// This is semantically a ReaderCloser, in that it does not affect Write
// After a call to Close() a subsequent Read() will start from the beginning (and reflect any new changes)
func (f *file) Close() (err error) {
	err = nil
	f.reader = nil
	return
}
