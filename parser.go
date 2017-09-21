package ini

import (
	"bufio"
	"regexp"
	"strings"
)

var (
	sectionRegex   = regexp.MustCompile(`^\[(.*)\]$`)
	assignArrRegex = regexp.MustCompile(`^([^=\[\]]+)\[\][^=]*=(.*)$`)
	assignRegex    = regexp.MustCompile(`^([^=]+)=(.*)$`)
	quotesRegex    = regexp.MustCompile(`^(['"])(.*)(['"])$`)
)

func trimWithQuotes(inputVal string) (ret string) {
	ret = strings.TrimSpace(inputVal)
	groups := quotesRegex.FindStringSubmatch(ret)
	if groups != nil {
		if groups[1] == groups[3] {
			ret = groups[2]
		}
	}
	return
}

func parseFile(in *bufio.Scanner, file *file) (bytes int64, err error) {
	section := ""
	lineNum := 0
	bytes = -1
	readLine := true
	for readLine = in.Scan(); readLine; readLine = in.Scan() {
		line := in.Text()
		bytes++
		bytes += int64(len(line))
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
			curVal, ok := file.section(section).arrayValues[key]
			if ok {
				file.section(section).arrayValues[key] = append(curVal, val)
			} else {
				file.section(section).arrayValues[key] = make([]string, 1, 4)
				file.section(section).arrayValues[key][0] = val
			}
		} else if groups := assignRegex.FindStringSubmatch(line); groups != nil {
			key, val := groups[1], groups[2]
			key, val = strings.TrimSpace(key), trimWithQuotes(val)
			file.section(section).stringValues[key] = val
		} else if groups := sectionRegex.FindStringSubmatch(line); groups != nil {
			name := strings.TrimSpace(groups[1])
			section = name
			// Create the section if it does not exist
			file.section(section)
		} else {
			err = ErrSyntax{lineNum, line}
			return
		}

	}
	if bytes < 0 {
		bytes = 0
	}
	err = in.Err()
	return
}
