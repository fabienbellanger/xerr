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
	"runtime/debug"
	"time"
)

// Err is a custom error type that wraps an error value with additional context.
//
// It includes fields for a message, details, file name, line number,
// timestamp, and a pointer to a previous error.
type Err struct {
	Value      error  `json:"value"`
	Code       int    `json:"code,omitzero"`
	Msg        string `json:"msg"`
	Details    any    `json:"details"`
	File       string `json:"file"`
	Line       int    `json:"line"`
	Timestamp  int64  `json:"timestamp"`
	Prev       *Err   `json:"prev"`
	StackTrace []byte `json:"stack_trace,omitempty"`
}

// New creates a new *Err with the provided error value, message, details, code,
// and a pointer to a previous Err.
//
// Returns nil if value is nil, making it compatible with standard nil checks.
//
// Example:
//
//	var myError = errors.New("my error")
//	type Person struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//	details := Person{Name: "John", Age: 30}
//	err := New(myError, "My error message", details, 0, nil)
func New(value error, msg string, details any, code int, prev *Err, skip ...int) *Err {
	if value == nil {
		return nil
	}

	callerSkip := 1
	if len(skip) == 1 {
		callerSkip = skip[0]
	}
	_, file, line, _ := runtime.Caller(callerSkip)
	stack := debug.Stack()

	return &Err{
		Value:      value,
		Code:       code,
		Msg:        msg,
		Details:    details,
		File:       file,
		Line:       line,
		Timestamp:  time.Now().UnixMicro(),
		Prev:       prev.Clone(),
		StackTrace: stack,
	}
}

// NewSimple creates a new *Err with the provided error value and message.
//
// It sets the details to nil, the code to 0, and the previous error to nil.
func NewSimple(value error, msg string, prev *Err, skip ...int) *Err {
	e := New(value, msg, nil, 0, prev)
	if e == nil {
		return nil
	}

	callerSkip := 1
	if len(skip) == 1 {
		callerSkip = skip[0]
	}
	_, file, line, _ := runtime.Caller(callerSkip)

	e.File = file
	e.Line = line

	return e
}

// Wrap creates a new *Err that wraps an existing error with additional context,
// chaining the current error as Prev.
func (e *Err) Wrap(value error, msg string, details any, code int, skip ...int) *Err {
	err := New(value, msg, details, code, e)

	callerSkip := 1
	if len(skip) == 1 {
		callerSkip = skip[0]
	}
	_, file, line, _ := runtime.Caller(callerSkip)

	err.File = file
	err.Line = line

	return err
}

// Clone creates a deep copy of the Err struct.
//
// It recursively clones the Prev field to ensure that the entire error chain
// is duplicated. Returns nil if called on a nil pointer.
func (e *Err) Clone() *Err {
	if e == nil {
		return nil
	}

	var clonedPrev *Err
	if e.Prev != nil {
		clonedPrev = e.Prev.Clone()
	}

	return &Err{
		Value:      e.Value,
		Code:       e.Code,
		Msg:        e.Msg,
		Details:    e.Details,
		File:       e.File,
		Line:       e.Line,
		Timestamp:  e.Timestamp,
		Prev:       clonedPrev,
		StackTrace: e.StackTrace,
	}
}

// Empty returns nil, representing the absence of an error.
// Callers may also use nil directly.
func Empty() *Err {
	return nil
}

// IsEmpty checks if the Err is nil or has no error value.
func (e *Err) IsEmpty() bool {
	return e == nil || e.Value == nil
}

// IsError checks if the Err contains an error value.
func (e *Err) IsError() bool {
	return e != nil && e.Value != nil
}

// Error implements the error interface.
func (e *Err) Error() string {
	if e.IsEmpty() {
		return ""
	}

	result := fmt.Sprintf("value=%v", e.Value)

	if e.Code != 0 {
		result += fmt.Sprintf(", code=%d", e.Code)
	}

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

// Is checks if any error in the chain matches the target, using errors.Is semantics.
//
// Example:
//
//	var myError = errors.New("my error")
//	err := New(myError, "My error message", nil, 0, nil)
//	if err.Is(myError) {
//		fmt.Println("The error matches myError")
//	}
func (e *Err) Is(err error) bool {
	if e == nil {
		return false
	}
	if errors.Is(e.Value, err) {
		return true
	}

	prev := e.Prev
	for prev != nil {
		if errors.Is(prev.Value, err) {
			return true
		}
		prev = prev.Prev
	}
	return false
}

// Unwrap returns the previous error in the chain, implementing the errors.Unwrap interface.
func (e *Err) Unwrap() error {
	if e == nil || e.Prev == nil {
		return nil
	}
	return e.Prev
}

// FromError creates a new *Err from an existing error.
// Returns nil if err is nil.
func FromError(err error) *Err {
	if err == nil {
		return nil
	}
	return New(err, "", nil, 0, nil)
}

// JSON converts the Err into a JSON representation.
// Pass true to include the stack trace in the output.
func (e *Err) JSON(stackTrace ...bool) ([]byte, error) {
	if e.IsEmpty() {
		return []byte{}, nil
	}

	clone := e.Clone()
	if len(stackTrace) == 0 || !stackTrace[0] {
		clone.StackTrace = nil
	}

	s, err := json.Marshal(clone)
	if err != nil {
		return []byte{}, err
	}

	return s, nil
}

// JSONOrEmpty converts the Err into a JSON representation.
// Returns an empty byte slice if the Err is empty or marshaling fails.
func (e *Err) JSONOrEmpty(stackTrace ...bool) []byte {
	if e.IsEmpty() {
		return []byte{}
	}

	clone := e.Clone()
	if len(stackTrace) == 0 || !stackTrace[0] {
		clone.StackTrace = nil
	}

	s, err := json.Marshal(clone)
	if err != nil {
		return []byte{}
	}

	return s
}

// MarshalJSON implements the json.Marshaler interface for the Err type.
//
// It customizes the JSON representation of the Err struct to include the error message
// and other fields in a specific format.
func (e *Err) MarshalJSON() ([]byte, error) {
	type Alias Err // Use an alias to avoid infinite recursion

	return json.Marshal(&struct {
		Value      string    `json:"value"`
		Details    any       `json:"details"`
		Timestamp  time.Time `json:"timestamp"`
		StackTrace string    `json:"stack_trace,omitempty"`
		Alias
	}{
		Value: func() string {
			if e.Value != nil {
				return e.Value.Error()
			}
			return ""
		}(),
		Details: func() any {
			if e.Details == nil {
				return nil
			}
			if _, err := json.Marshal(e.Details); err != nil {
				return nil
			}
			return e.Details
		}(),
		Timestamp:  time.UnixMicro(e.Timestamp),
		StackTrace: string(e.StackTrace),
		Alias:      (Alias)(*e),
	})
}

// ValueEq checks if the Value field of two Err structs are equal.
//
// Example:
//
//	var myError = errors.New("my error")
//	err1 := New(myErr, "My error message 1", nil, 0, nil)
//	err2 := New(myErr, "My error message 2", nil, 0, nil)
//	println(err1.ValueEq(err2)) // true
func (e *Err) ValueEq(other *Err) bool {
	if e == nil || other == nil {
		return e == other
	}
	return errors.Is(e.Value, other.Value)
}

// Eq checks if two Err structs are equal by comparing Value and the Prev chain.
func (e *Err) Eq(other *Err) bool {
	if e == nil || other == nil {
		return e == other
	}
	if !errors.Is(e.Value, other.Value) {
		return false
	}

	ep, op := e.Prev, other.Prev
	for ep != nil && op != nil {
		if !errors.Is(ep.Value, op.Value) {
			return false
		}
		ep, op = ep.Prev, op.Prev
	}
	return ep == nil && op == nil
}

// ToError converts the Err to a standard error interface value.
// Returns nil if the Err is nil or empty.
func (e *Err) ToError() error {
	if e == nil || e.IsEmpty() {
		return nil
	}
	return e
}
