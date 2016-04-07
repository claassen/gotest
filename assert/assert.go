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

func isZeroValue(v reflect.Value) bool {
    switch v.Kind() {
    case reflect.Func, reflect.Map, reflect.Slice, reflect.Chan, reflect.Interface, reflect.Ptr:
        return v.IsNil()
    }
    
    return false
}

func areEqualValues(x, y interface{}) bool {

	if areEqualReferences(x, y) {
		fmt.Println("equal refs")
		return true
	}

	if reflect.DeepEqual(x, y) {
		fmt.Println("deep equal")
		return true
	}

	//nil interface value will have nil type
	tx := reflect.TypeOf(x)
	ty := reflect.TypeOf(y)

	vx := reflect.ValueOf(x)
	vy := reflect.ValueOf(y)

	//Check for zero values for types which can be assigned nil and consider those to equal nil
	if (x == nil && isZeroValue(vy)) || (y == nil && isZeroValue(vx)) {
		return true
	}

	if tx != nil && ty != nil && tx.ConvertibleTo(ty) {
		return reflect.DeepEqual(vx.Convert(ty).Interface(), y)
	}

	return false
}

func areEqualReferences(x, y interface{}) bool {

	if x == nil || y == nil {
		return x == y
	}

	return reflect.Indirect(reflect.ValueOf(x)) == reflect.Indirect(reflect.ValueOf(y))
}

func (v AssertValue) IsNil() {
	if !areEqualValues(v.value, nil) {
		message := fmt.Sprintf("Expected %#v to be nil.", v.value)
		fail(message)
	}
}

func (v AssertValue) IsNotNil() {
	if areEqualValues(v.value, nil) {
		message := fmt.Sprintf("Expected %#v to not be nil.", v.value)
		fail(message)
	}
}

func (e AssertValue) IsEqualTo(expected interface{}) {
	if !areEqualValues(e.value, expected) {
		message := fmt.Sprintf("Expected %#v to be equal to %#v.", expected, e.value)
		fail(message)
	}
}

func (e AssertValue) IsNotEqualTo(expected interface{}) {
	if areEqualValues(e.value, expected) {
		message := fmt.Sprintf("Expected %#v to not be equal to %#v.", expected, e.value)
		fail(message)
	}
}

func (e AssertValue) Is(expected interface{}) {
	if !areEqualReferences(e.value, expected) {
		message := fmt.Sprintf("Expected %#v to be the same object as %#v", expected, e.value)
		fail(message)
	}
}

func (e AssertValue) IsNot(expected interface{}) {
	if areEqualReferences(e.value, expected) {
		message := fmt.Sprintf("Expected %#v to not be the same object as %#v", expected, e.value)
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
