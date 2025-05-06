// Package xerr is a simple error wrapper that provides additional context and
// functionality for error handling in Go applications.
//
// It includes features such as JSON serialization, nested error handling,
// and custom error messages. The package is designed to be straightforward to use
// and integrate into existing Go codebases, making it a valuable tool.
package xerr

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"time"
)

// Err is a custom error type that wraps an error value with additional context.
//
// It includes fields for a message, details, file name, line number,
// timestamp, and a pointer to a previous error.
type Err struct {
	Value     error  `json:"value"`
	Msg       string `json:"msg"`
	Details   any    `json:"details"`
	File      string `json:"file"`
	Line      int    `json:"line"`
	Timestamp int64  `json:"timestamp"`
	Prev      *Err   `json:"prev"`
}

// NewErr creates a new Err struct with the provided error value, message,
// details, and a pointer to a previous Err struct.
//
// The timestamp is set to the current time in microseconds since the epoch.
//
// Example:
//
//	var myError = errors.New("my error")
//	type Person struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//	details := Person{Name: "John", Age: 30}
//	err := NewErr(myError, "My error message", details, nil)
func NewErr(value error, msg string, details any, next *Err) Err {
	if value == nil {
		return EmptyErr()
	}

	_, file, line, _ := runtime.Caller(1)

	return Err{
		Value:     value,
		Msg:       msg,
		Details:   details,
		File:      file,
		Line:      line,
		Timestamp: time.Now().UnixMicro(),
		Prev:      next.Clone(),
	}
}

// Clone creates a deep copy of the Err struct.
//
// It recursively clones the Prev field to ensure that the entire error chain
// is duplicated.
func (e *Err) Clone() *Err {
	if e == nil {
		return nil
	}

	// Reservely clone the Prev field to avoid modifying the original
	var clonedPrev *Err
	if e.Prev != nil {
		clonedPrev = e.Prev.Clone()
	}

	return &Err{
		Value:     e.Value,
		Msg:       e.Msg,
		Details:   e.Details,
		File:      e.File,
		Line:      e.Line,
		Timestamp: e.Timestamp,
		Prev:      clonedPrev,
	}
}

// EmptyErr returns an empty Err struct.
func EmptyErr() Err {
	return Err{}
}

// IsEmpty checks if the Err struct is empty, meaning it has no error value.
func (e Err) IsEmpty() bool {
	return e.Value == nil
}

// IsError checks if the Err struct contains an error value.
func (e Err) IsError() bool {
	return e.Value != nil
}

func (e Err) Error() string {
	if e.IsEmpty() {
		return ""
	}

	result := fmt.Sprintf("value=%v", e.Value)

	if e.Msg != "" {
		result += fmt.Sprintf(", msg=%+v", e.Msg)
	}

	if e.Details != nil {
		result += fmt.Sprintf(", details=%+v", e.Details)
	}

	if e.File != "" {
		result += fmt.Sprintf(", source=%s:%d", e.File, e.Line)
	}

	if e.Timestamp != 0 {
		result += fmt.Sprintf(", timestamp=%s", time.UnixMicro(e.Timestamp).Format(time.RFC3339Nano))
	}

	if e.Prev != nil {
		result += fmt.Sprintf(", prev={%s}", e.Prev.Error())
	}

	return result
}

// Is checks if the error in the Err struct is of a specific type or value.
//
// It uses the errors.Is function to check if the error in the Err struct
// matches the provided error value.
//
// Example:
//
//	var myError = errors.New("my error")
//	details := Person{Name: "John", Age: 30}
//	err := NewErr(myError, "My error message", details, nil)
//	if err.Is(myError) {
//		fmt.Println("The error matches myError")
//	}
func (e Err) Is(err error) bool {
	if errors.Is(e.Value, err) {
		return true
	}

	prev := e.Prev
	for !(prev == nil) {
		if errors.Is(prev.Value, err) {
			return true
		}

		prev = prev.Prev
	}
	return false
}

// JSON converts the Err struct into a JSON representation.
func (e Err) JSON() ([]byte, Err) {
	if e.IsEmpty() {
		return []byte{}, EmptyErr()
	}

	s, err := json.Marshal(e)
	if err != nil {
		return []byte{}, NewErr(err, "Error when converting Err into JSON", nil, nil)
	}

	return s, EmptyErr()
}

// MarshalJSON implements the json.Marshaler interface for the Err type.
//
// It customizes the JSON representation of the Err struct to include the error message
// and other fields in a specific format.
//
// The Value field is converted to a string using the Error() method if it is not nil.
//
// The function returns the JSON representation of the Err struct.
//
// If the Value field is nil, it returns an empty string for the Value field in the JSON output.
func (e Err) MarshalJSON() ([]byte, error) {
	type Alias Err // Use an alias to avoid infinite recursion

	return json.Marshal(&struct {
		Value     string    `json:"value"`
		Details   any       `json:"details"`
		Timestamp time.Time `json:"timestamp"`
		Alias
	}{
		Value: func() string {
			if e.Value != nil {
				return e.Value.Error()
			}

			return ""
		}(),
		Details: func() any {
			// If Details is nil, return nil
			if e.Details == nil {
				return nil
			}

			// Check if Details is serializable
			if _, err := json.Marshal(e.Details); err != nil {
				return nil
			}

			return e.Details
		}(),
		Timestamp: time.UnixMicro(e.Timestamp),
		Alias:     (Alias)(e),
	})
}

// ValueEq checks if the Value field of the Err struct is equal to
// the Value field of another Err struct.
//
// Example:
//
//	var myError = errors.New("my error")
//	err1 := NewErr(myErr, "My error message 1", nil, nil)
//	err2 := NewErr(myErr, "My error message 2", nil, nil)
//	println(err1.ValueEq(err2)) // true
func (e *Err) ValueEq(other Err) bool {
	return e.Value == other.Value
}

// Eq checks if the Err struct is equal to another Err struct only by comparing
// the Value field and the previous errors in the chain.
func (e *Err) Eq(other Err) bool {
	if e.Value != other.Value {
		return false
	}

	prev := e.Prev
	for !(prev == nil) {
		if prev.Value != other.Prev.Value {
			return false
		}

		prev = prev.Prev
	}
	return true
}
