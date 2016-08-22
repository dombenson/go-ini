package ini


type File struct {
	sections map[string]*section

}


// Returns a named Section. A Section will be created if one does not already exist for the given name.
func (f File) section(name string) *section {
	theSection := f.sections[name]
	if theSection == nil {
		theSection = &section{stringValues: make(map[string]string), arrayValues: make(map[string][]string)}
		f.sections[name] = theSection
	}
	return theSection
}

func (f File) Get(section, key string) (value string, ok bool) {
	return f.section(section).Get(key)
}

func (f File) Set(section, key string, value string) (ok bool) {
	return f.section(section).Set(key, value)
}

func (f File) SetInt(section, key string, value int) (ok bool) {
	return f.section(section).SetInt(key, value)
}

func (f File) SetBool(section, key string, value bool) (ok bool) {
	return f.section(section).SetBool(key, value)
}

func (f File) GetInt(section, key string) (value int, ok bool) {
	return f.section(section).GetInt(key)
}


func (f File) GetBool(section, key string) (value bool, ok bool) {
	return f.section(section).GetBool(key)
}

func (f File) GetArr(section, key string) (value []string, ok bool) {
	return f.section(section).GetArr(key)
}
