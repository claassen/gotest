package assert

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

//Source: https://golang.org/src/testing/testing.go
func decorate(message string) string {
	_, file, line, ok := runtime.Caller(2)

	if ok {
		// Truncate file name at last file name separator
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
		line = 1
	}

	buf := new(bytes.Buffer)
	// Every line is indented at least one tab.
	buf.WriteByte('\t')
	fmt.Fprintf(buf, "%s:%d: ", file, line)
	lines := strings.Split(message, "\n")
	if l := len(lines); l > 1 && lines[l-1] == "" {
		lines = lines[:l-1]
	}
	for i, line := range lines {
		if i > 0 {
			// Second and subsequent lines are indented an extra tab.
			buf.WriteString("\n\t\t")
		}
		buf.WriteString(line)
	}
	buf.WriteByte('\n')

	return buf.String()
}

func Fail(message string) {
	panic(decorate(message))
}
