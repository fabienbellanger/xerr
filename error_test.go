package xerr

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ----------------------------------------------------------------------------
//
// Tests of NewErr()
//
// ----------------------------------------------------------------------------

func TestNewErrSimpleError(t *testing.T) {
	details := struct {
		Name string
		Age  int
	}{
		Name: "John Doe",
		Age:  23,
	}

	err := NewErr(errors.New("test"), "My error message", details, nil)

	assert.Equal(t, errors.New("test"), err.Value)
	assert.Equal(t, "My error message", err.Msg)
	assert.Equal(t, details, err.Details)
	assert.True(t, strings.Contains(err.File, "error_test.go"))
	assert.Equal(t, 28, err.Line)
	assert.Nil(t, err.Prev)
}

func TestNewErrNestedErrors(t *testing.T) {
	err2 := NewErr(errors.New("test 2"), "My error message 2", nil, nil)
	err1 := NewErr(errors.New("test 1"), "My error message 1", nil, &err2)

	assert.Equal(t, errors.New("test 1"), err1.Value)
	assert.Equal(t, "My error message 1", err1.Msg)
	assert.Nil(t, err1.Details)
	assert.True(t, strings.Contains(err1.File, "error_test.go"))
	assert.Equal(t, 40, err1.Line)

	assert.Equal(t, errors.New("test 2"), err2.Value)
	assert.Equal(t, "My error message 2", err2.Msg)
	assert.Nil(t, err2.Details)
	assert.True(t, strings.Contains(err2.File, "error_test.go"))
	assert.Equal(t, 39, err2.Line)
}

func TestNewErrEmptyError(t *testing.T) {
	err := NewErr(nil, "My error message", nil, nil)

	assert.Equal(t, Err{}, err)
}

// ----------------------------------------------------------------------------
//
// Tests of EmptyErr()
//
// ----------------------------------------------------------------------------

func TestEmptyErr(t *testing.T) {
	assert.Equal(t, Err{}, EmptyErr())
}

// ----------------------------------------------------------------------------
//
// Tests of IsEmpty()
//
// ----------------------------------------------------------------------------

func TestIsEmpty(t *testing.T) {
	err := EmptyErr()
	assert.True(t, err.IsEmpty())

	err = NewErr(errors.New("test"), "My error message", nil, nil)
	assert.False(t, err.IsEmpty())
}

// ----------------------------------------------------------------------------
//
// Tests of IsError()
//
// ----------------------------------------------------------------------------

func TestIsError(t *testing.T) {
	err := EmptyErr()
	assert.False(t, err.IsError())

	err = NewErr(errors.New("test"), "My error message", nil, nil)
	assert.True(t, err.IsError())
}

// ----------------------------------------------------------------------------
//
// Tests of Error()
//
// ----------------------------------------------------------------------------

func TestErrorEmpty(t *testing.T) {
	err := EmptyErr()

	assert.Equal(t, "", err.Error())
}

func TestErrorNotEmpty(t *testing.T) {
	now := time.Now()
	err := Err{
		Value:     errors.New("test"),
		Msg:       "My error message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now.UnixMicro(),
		Prev:      nil,
	}

	expected := "value=test, msg=My error message, source=error_test.go:26, timestamp=" + now.Format(time.RFC3339Nano)

	assert.Equal(t, expected, err.Error())
}

func TestErrorNestedErrors(t *testing.T) {
	now := time.Now()
	err2 := Err{
		Value:     errors.New("test 2"),
		Msg:       "My error message 2",
		Details:   nil,
		File:      "",
		Line:      0,
		Timestamp: now.UnixMicro(),
		Prev:      nil,
	}
	err := Err{
		Value:     errors.New("test"),
		Msg:       "My error message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now.UnixMicro(),
		Prev:      &err2,
	}

	expected := "value=test, msg=My error message, source=error_test.go:26, timestamp="
	expected += now.Format(time.RFC3339Nano) + ", prev={value=test 2, msg=My error message 2, "
	expected += "timestamp=" + now.Format(time.RFC3339Nano) + "}"

	assert.Equal(t, expected, err.Error())
}

func TestErrorWithoutTimestamp(t *testing.T) {
	err := Err{
		Value:     errors.New("test"),
		Msg:       "My error message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: 0,
		Prev:      nil,
	}

	expected := "value=test, msg=My error message, source=error_test.go:26"

	assert.Equal(t, expected, err.Error())
}

func TestErrorWithDetails(t *testing.T) {
	details := struct {
		Name string
		Age  int
	}{
		Name: "John Doe",
		Age:  23,
	}

	now := time.Now()
	err := Err{
		Value:     errors.New("test"),
		Msg:       "My error message",
		Details:   details,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now.UnixMicro(),
		Prev:      nil,
	}

	expected := "value=test, msg=My error message, details={Name:John Doe Age:23}, source=error_test.go:26, timestamp=" + now.Format(time.RFC3339Nano)

	assert.Equal(t, expected, err.Error())
}

// ----------------------------------------------------------------------------
//
// Tests of Is()
//
// ----------------------------------------------------------------------------

func TestIs(t *testing.T) {
	myErr := errors.New("my error")
	err := NewErr(myErr, "My error message", nil, nil)

	assert.True(t, err.Is(myErr))
}

func TestIsNested(t *testing.T) {
	myErr := errors.New("my error")
	myErr2 := errors.New("my error 2")
	myErr3 := errors.New("my error 3")
	err3 := NewErr(myErr3, "My error message 3", nil, nil)
	err2 := NewErr(myErr2, "My error message 2", nil, &err3)
	err := NewErr(myErr, "My error message", nil, &err2)

	assert.True(t, err.Is(myErr))
	assert.True(t, err.Is(myErr2))
	assert.True(t, err.Is(myErr3))
}

func TestIsFalse(t *testing.T) {
	myErr := errors.New("my error")
	myErr2 := errors.New("my error 2")
	err := NewErr(myErr, "My error message", nil, nil)

	assert.True(t, err.Is(myErr))
	assert.False(t, err.Is(myErr2))
}

// ----------------------------------------------------------------------------
//
// Tests of JSON()
//
// ----------------------------------------------------------------------------

func TestJSONEmpty(t *testing.T) {
	e := EmptyErr()
	expected := []byte("")
	result, err := e.JSON()

	fmt.Printf("%s\n", result)

	assert.Equal(t, EmptyErr(), err)
	assert.Equal(t, expected, result)
}

func TestJSONSimple(t *testing.T) {
	now := time.Now()
	e := Err{
		Value:     errors.New("test"),
		Msg:       "My error message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now.UnixMicro(),
		Prev:      nil,
	}

	expected := []byte(`{"value":"test","details":null,"timestamp":"` + now.Format(time.RFC3339Nano) + `","msg":"My error message","file":"error_test.go","line":26,"prev":null}`)
	result, err := e.JSON()

	assert.Equal(t, EmptyErr(), err)
	assert.Equal(t, expected, result)
}

func TestJSONDetail(t *testing.T) {
	now := time.Now()
	details := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "John Doe",
		Age:  23,
	}

	e := Err{
		Value:     errors.New("test"),
		Msg:       "My error message",
		Details:   details,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now.UnixMicro(),
		Prev:      nil,
	}

	expected := []byte(`{"value":"test","details":{"name":"John Doe","age":23},"timestamp":"` +
		now.Format(time.RFC3339Nano) +
		`","msg":"My error message","file":"error_test.go","line":26,"prev":null}`)
	result, err := e.JSON()

	assert.Equal(t, EmptyErr(), err)
	assert.Equal(t, expected, result)
}

func TestJSONNestedErrors(t *testing.T) {
	now := time.Now()

	e := Err{
		Value:     errors.New("test"),
		Msg:       "My message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now.UnixMicro(),
		Prev: &Err{
			Value:     errors.New("test 2"),
			Msg:       "My message 2",
			Details:   nil,
			File:      "error_test.go",
			Line:      87,
			Timestamp: now.UnixMicro(),
			Prev:      nil,
		},
	}

	expected := []byte(`{"value":"test","details":null,"timestamp":"` +
		now.Format(time.RFC3339Nano) +
		`","msg":"My message","file":"error_test.go","line":26,"prev":{"value":"test 2","details":null,"timestamp":"` +
		now.Format(time.RFC3339Nano) + `","msg":"My message 2","file":"error_test.go","line":87,"prev":null}}`)
	result, err := e.JSON()

	assert.Equal(t, EmptyErr(), err)
	assert.Equal(t, expected, result)
}

func TestJSONDetailError(t *testing.T) {
	now := time.Now()
	details := struct {
		Channel chan int
	}{
		Channel: make(chan int),
	}

	e := Err{
		Value:     errors.New("test"),
		Msg:       "My error message",
		Details:   details,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now.UnixMicro(),
		Prev:      nil,
	}

	expected := []byte(`{"value":"test","details":null,"timestamp":"` + now.Format(time.RFC3339Nano) + `","msg":"My error message","file":"error_test.go","line":26,"prev":null}`)
	result, err := e.JSON()

	assert.Equal(t, EmptyErr(), err)
	assert.Equal(t, expected, result)
}
