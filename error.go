// Package xerr provides a structured error type with call-site capture,
// error chaining, stack traces, and JSON serialization.
//
// All constructors return *Err, so a nil return means no error and is fully
// compatible with standard Go nil checks.
package xerr

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

// Err wraps an error with structured context: an optional code, human-readable
// message, arbitrary details, call-site location (File, Line), timestamp, stack
// trace, and a Prev pointer that forms a linked chain of errors.
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
// and a pointer to a previous Err. Returns nil if value is nil.
//
// The optional skip parameter controls the depth passed to [runtime.Caller] for
// capturing the call site. It defaults to 1 (the caller of New). Wrapper
// functions should pass a higher value so that File and Line reflect their own
// caller rather than the wrapper itself.
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

// NewSimple creates a new *Err with only a value, message, and optional prev
// chain. Details is nil, Code is 0.
//
// The optional skip parameter works the same as in [New].
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

// Wrap creates a new *Err with value, msg, details, and code, chaining the
// receiver as Prev. The receiver is cloned to avoid mutation.
//
// The optional skip parameter works the same as in [New].
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

// Error implements the error interface, returning a human-readable string with
// all non-zero fields formatted as key=value pairs (e.g. "value=…, code=…").
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

// Is reports whether any Value in the Prev chain matches target using [errors.Is].
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

// Unwrap returns the previous error in the chain, satisfying the interface
// expected by [errors.Unwrap].
func (e *Err) Unwrap() error {
	if e == nil || e.Prev == nil {
		return nil
	}
	return e.Prev
}

// FromError creates a new *Err from a plain error, capturing the caller's
// file and line. Returns nil if err is nil.
func FromError(err error) *Err {
	if err == nil {
		return nil
	}
	return New(err, "", nil, 0, nil)
}

// JSON returns the JSON encoding of the Err. The stack trace is omitted
// unless stackTrace is true. It operates on a clone to avoid mutating the
// receiver.
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

// JSONOrEmpty is like [Err.JSON] but silently returns an empty byte slice on
// error or if the Err is empty.
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

// MarshalJSON implements [json.Marshaler]. It converts Value to its string
// representation, Timestamp to a [time.Time], StackTrace to a string, and
// drops non-serializable Details. An internal Alias type prevents infinite
// recursion.
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

// ValueEq reports whether e and other have the same Value (compared with
// [errors.Is]). Both nil returns true; one nil returns false.
//
// Example:
//
//	var myError = errors.New("my error")
//	err1 := New(myError, "My error message 1", nil, 0, nil)
//	err2 := New(myError, "My error message 2", nil, 0, nil)
//	println(err1.ValueEq(err2)) // true
func (e *Err) ValueEq(other *Err) bool {
	if e == nil || other == nil {
		return e == other
	}
	return errors.Is(e.Value, other.Value)
}

// Eq reports whether e and other have the same Value and an identical Prev chain
// (each pair compared with [errors.Is]).
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

// ToError returns the receiver as an error interface value, or nil if the Err
// is nil or empty. This avoids the typed-nil pitfall when returning *Err from
// a function with an error return type.
func (e *Err) ToError() error {
	if e == nil || e.IsEmpty() {
		return nil
	}
	return e
}
