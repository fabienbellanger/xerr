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
		Prev:      next,
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

	return fmt.Sprintf("value=%v, msg=%s, details=%v, file=%s, line=%d, timestamp=%s, prev={%s}",
		e.Value,
		e.Msg,
		e.Details,
		e.File,
		e.Line,
		time.UnixMicro(e.Timestamp).Format(time.RFC3339Nano),
		e.Prev,
	)
}

// Is checks if the error in the Err struct is of a specific type or value.
// It uses the errors.Is function to check if the error in the Err struct
// matches the provided error value.
func (e Err) Is(err error) bool {
	return errors.Is(e.Value, err)
}

// As checks if the error in the Err struct can be cast to a specific type.
//
// As panics if target is not a non-nil pointer to either a type that implements error, or to any interface type.
func (e Err) As(target any) bool {
	return errors.As(e.Value, &target)
}

// JSON converts the Err struct into a JSON representation.
func (e Err) JSON() ([]byte, Err) {
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
// The Value field is converted to a string using the Error() method if it is not nil.
// The function returns the JSON representation of the Err struct.
// If the Value field is nil, it returns an empty string for the Value field in the JSON output.
func (e Err) MarshalJSON() ([]byte, error) {
	type Alias Err // Use an alias to avoid infinite recursion

	return json.Marshal(&struct {
		Value   string `json:"value"`
		Details any    `json:"details"`
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
		Alias: (Alias)(e),
	})
}
