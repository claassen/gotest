package assert

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

//Stolen from the go testing package: https://golang.org/src/testing/testing.go
func decorate(message string) string {
	_, file, line, ok := runtime.Caller(3)

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

	//gotest renames _test.go files to _testx.go files, we want to show the original file name
	file = strings.Replace(file, "x.go", ".go", 1)

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

type AssertValue struct {
	value interface{}
}

func AssertThat(val interface{}) AssertValue {
	return AssertValue{value: val}
}

func fail(message string) {
	panic(decorate(message))
}

func Fail(message string) {
	fail(message)
}

func (v AssertValue) IsNil() {
	if v.value != nil {
		message := fmt.Sprintf("Expected %#v to be nil.")
		fail(message)
	}
}

func (v AssertValue) IsNotNil() {
	if v.value == nil {
		message := fmt.Sprintf("Expected %#v to not be nil.")
		fail(message)
	}
}

func (e AssertValue) IsEqualTo(expected interface{}) {
	if e.value != expected {
		message := fmt.Sprintf("Expected %#v to be equal to %#v.", expected, e.value)
		fail(message)
	}
}

func (e AssertValue) IsNotEqualTo(expected interface{}) {
	if e.value == expected {
		message := fmt.Sprintf("Expected %#v to not be equal to %#v.", expected, e.value)
		fail(message)
	}
}

func panics(f func()) bool {

	didPanic := false

	func() {
		defer func() {
			if message := recover(); message != nil {
				didPanic = true
			}
		}()
		f()
	}()

	return didPanic
}

func (e AssertValue) Panics() {
	v := reflect.ValueOf(e.value)
	t := v.Type()
	if t.Kind() == reflect.Func {
		if t.NumIn() != 0 {
			fail("Cannot assert that function accepting arguments panics.")
		}

		f, _ := e.value.(func())

		if !panics(f) {
			fail("Expected function to panic but it did not.")
		}
	} else {
		fail("Cannot assert that non function object panics.")
	}
}

func (e AssertValue) DoesNotPanic() {
	v := reflect.ValueOf(e.value)
	t := v.Type()
	if t.Kind() == reflect.Func {
		if t.NumIn() != 0 {
			fail("Cannot assert that function accepting arguments does not panic.")
		}

		f, _ := e.value.(func())

		if panics(f) {
			fail("Expected function not to panic but it did.")
		}
	} else {
		fail("Cannot assert that non function object does not panic.")
	}
}
