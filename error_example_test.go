package xerr

import (
	"errors"
	"fmt"
)

func ExampleErr_JSON() {
	err := Err{
		Value:     errors.New("test"),
		Msg:       "My error message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: 1691234567890,
		Prev:      nil,
	}
	result, _ := err.JSON()

	fmt.Println(string(result))

	// Output: {"value":"test","details":null,"timestamp":"1970-01-20T14:47:14.56789+01:00","msg":"My error message","file":"error_test.go","line":26,"prev":null}
}

func ExampleErr_Is() {
	var myErr = errors.New("my error")
	err := Err{
		Value:     myErr,
		Msg:       "My error message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: 1691234567890,
		Prev:      nil,
	}

	fmt.Println(err.Is(myErr))

	// Output: true
}
