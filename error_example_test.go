package xerr

import (
	"errors"
	"fmt"
)

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
