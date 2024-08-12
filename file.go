package ini

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

// This implements the full ini.StreamReadWriter interface
type file struct {
	sections map[string]*section
	reader   io.Reader
}

func (f *file) ParseEnvironmentVariables() error {
	environ := os.Environ()
	environMap := make(map[string]string, len(environ))

	for _, envvar := range environ {
		pair := strings.SplitN(envvar, "=", 2)
		environMap[pair[0]] = pair[1]
	}

	templateValues := struct {
		Env map[string]string
	}{
		Env: environMap,
	}

	for sectionName, currentSection := range f.sections {
		newStringSection := make(stringSection, len(currentSection.stringValues))
		for k, v := range currentSection.stringValues {
			tmpl, err := template.New(fmt.Sprintf("[%s]%s", sectionName, k)).Parse(v)
			if err != nil {
				return err
			}

			var res bytes.Buffer
			err = tmpl.Execute(&res, templateValues)
			if err != nil {
				return err
			}
			newStringSection[k] = res.String()
		}
		f.sections[sectionName].stringValues = newStringSection

		newArraySection := make(arraySection, len(currentSection.arrayValues))

		for k, valueSlice := range currentSection.arrayValues {
			newArraySection[k] = make([]string, len(valueSlice))
			for i, v := range valueSlice {
				tmpl, err := template.New(fmt.Sprintf("[%s][%d]%s", sectionName, i, k)).Parse(v)
				if err != nil {
					return err
				}

				var res bytes.Buffer
				err = tmpl.Execute(&res, templateValues)
				if err != nil {
					return err
				}
				newArraySection[k][i] = res.String()
			}
		}
		f.sections[sectionName].arrayValues = newArraySection
	}

	return nil
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

func (f *file) Sections() (value []string) {
	value = make([]string, 0, len(f.sections))
	for sect, _ := range f.sections {
		value = append(value, sect)
	}
	return
}

func (f *file) Values(section string) (value map[string]string) {
	value = make(map[string]string)
	sect := f.section(section)
	if sect != nil {
		for k, v := range sect.stringValues {
			value[k] = v
		}
	}
	return
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

func (f *file) Remove(section, key string) {
	f.section(section).Remove(key)
}

func (f *file) RemoveSection(section string) {
	_, found := f.sections[section]
	if found {
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
