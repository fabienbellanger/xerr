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
// Tests of New()
//
// ----------------------------------------------------------------------------

func TestErr_New_SimpleError(t *testing.T) {
	details := struct {
		Name string
		Age  int
	}{
		Name: "John Doe",
		Age:  23,
	}

	err := New(errors.New("test"), "My error message", details, 10, nil)

	assert.Equal(t, errors.New("test"), err.Value)
	assert.Equal(t, 10, err.Code)
	assert.Equal(t, "My error message", err.Msg)
	assert.Equal(t, details, err.Details)
	assert.True(t, strings.Contains(err.File, "error_test.go"))
	assert.Equal(t, 28, err.Line)
	assert.Nil(t, err.Prev)
}

func TestErr_New_With_Err(t *testing.T) {
	err2 := New(errors.New("test 2"), "My error message 2", nil, 20, nil)
	err1 := New(&err2, "My error message 1", nil, 10, nil)

	assert.Equal(t, &err2, err1.Value)
	assert.Equal(t, 10, err1.Code)
	assert.Equal(t, "My error message 1", err1.Msg)
	assert.Nil(t, err1.Details)
	assert.True(t, strings.Contains(err1.File, "error_test.go"))
	assert.Equal(t, 41, err1.Line)
}

func TestErr_New_NestedErrors(t *testing.T) {
	err2 := New(errors.New("test 2"), "My error message 2", nil, 20, nil)
	err1 := New(errors.New("test 1"), "My error message 1", nil, 10, &err2)

	assert.Equal(t, errors.New("test 1"), err1.Value)
	assert.Equal(t, 10, err1.Code)
	assert.Equal(t, "My error message 1", err1.Msg)
	assert.Nil(t, err1.Details)
	assert.True(t, strings.Contains(err1.File, "error_test.go"))
	assert.Equal(t, 53, err1.Line)

	assert.Equal(t, errors.New("test 2"), err2.Value)
	assert.Equal(t, 20, err2.Code)
	assert.Equal(t, "My error message 2", err2.Msg)
	assert.Nil(t, err2.Details)
	assert.True(t, strings.Contains(err2.File, "error_test.go"))
	assert.Equal(t, 52, err2.Line)
}

func TestErr_New_Emptyor(t *testing.T) {
	err := New(nil, "My error message", nil, 0, nil)

	assert.Equal(t, Err{}, err)
}

func TestErr_New_WithSkip(t *testing.T) {
	err := New(errors.New("test"), "My error message", nil, 20, nil, 0)

	assert.Equal(t, errors.New("test"), err.Value)
	assert.Equal(t, 20, err.Code)
	assert.Equal(t, "My error message", err.Msg)
	assert.Nil(t, err.Details)
	assert.True(t, strings.Contains(err.File, "error.go"))
	assert.Equal(t, 57, err.Line)
}

// ----------------------------------------------------------------------------
//
// Tests of NewSimple()
//
// ----------------------------------------------------------------------------

func TestErr_NewSimple(t *testing.T) {
	err := NewSimple(errors.New("test"), "My error message", nil)
	assert.Equal(t, errors.New("test"), err.Value)
	assert.Equal(t, 0, err.Code)
	assert.Equal(t, "My error message", err.Msg)
	assert.Nil(t, err.Details)
	assert.Nil(t, err.Prev)
}

func TestErr_NewSimple_WithSkip(t *testing.T) {
	err := NewSimple(errors.New("test"), "My error message", nil, 0)

	assert.Equal(t, errors.New("test"), err.Value)
	assert.Equal(t, 0, err.Code)
	assert.Equal(t, "My error message", err.Msg)
	assert.Nil(t, err.Details)
	assert.True(t, strings.Contains(err.File, "error.go"))
	assert.Equal(t, 83, err.Line)
}

// ----------------------------------------------------------------------------
//
// Tests of Wrap()
//
// ----------------------------------------------------------------------------

func TestErr_Wrap(t *testing.T) {
	err := NewSimple(errors.New("test"), "My error message", nil)
	wrappedErr := err.Wrap(errors.New("wrapped error"), "Wrapped message", nil, 100)
	expected := Err{
		Value:   errors.New("wrapped error"),
		Code:    100,
		Msg:     "Wrapped message",
		Details: nil,
		Prev:    &err,
	}

	assert.Equal(t, expected.Value, wrappedErr.Value)
	assert.Equal(t, expected.Code, wrappedErr.Code)
	assert.Equal(t, expected.Msg, wrappedErr.Msg)
	assert.Equal(t, expected.Details, wrappedErr.Details)
	assert.True(t, strings.Contains(wrappedErr.File, "error_test.go"))
	assert.Equal(t, 121, wrappedErr.Line)
	assert.Equal(t, expected.Prev, wrappedErr.Prev)
}

func TestErr_Wrap_WithSkip(t *testing.T) {
	err := NewSimple(errors.New("test"), "My error message", nil)
	wrappedErr := err.Wrap(errors.New("wrapped error"), "Wrapped message", nil, 100, 0)
	expected := Err{
		Value:   errors.New("wrapped error"),
		Code:    100,
		Msg:     "Wrapped message",
		Details: nil,
		Prev:    &err,
	}

	assert.Equal(t, expected.Value, wrappedErr.Value)
	assert.Equal(t, expected.Code, wrappedErr.Code)
	assert.Equal(t, expected.Msg, wrappedErr.Msg)
	assert.Equal(t, expected.Details, wrappedErr.Details)
	assert.True(t, strings.Contains(wrappedErr.File, "error.go"))
	assert.Equal(t, 102, wrappedErr.Line)
	assert.Equal(t, expected.Prev, wrappedErr.Prev)
}

// ----------------------------------------------------------------------------
//
// Tests of Empty()
//
// ----------------------------------------------------------------------------

func TestErr_Empty(t *testing.T) {
	assert.Equal(t, Err{}, Empty())
}

// ----------------------------------------------------------------------------
//
// Tests of IsEmpty()
//
// ----------------------------------------------------------------------------

func TestErr_IsEmpty(t *testing.T) {
	err := Empty()
	assert.True(t, err.IsEmpty())

	err = New(errors.New("test"), "My error message", nil, 0, nil)
	assert.False(t, err.IsEmpty())
}

// ----------------------------------------------------------------------------
//
// Tests of IsError()
//
// ----------------------------------------------------------------------------

func TestErr_IsError(t *testing.T) {
	err := Empty()
	assert.False(t, err.IsError())

	err = New(errors.New("test"), "My error message", nil, 0, nil)
	assert.True(t, err.IsError())
}

// ----------------------------------------------------------------------------
//
// Tests of Error()
//
// ----------------------------------------------------------------------------

func TestErr_Error_Empty(t *testing.T) {
	err := Empty()

	assert.Equal(t, "", err.Error())
}

func TestErr_Error_NotEmpty(t *testing.T) {
	now := time.Now().UnixMicro()
	err := Err{
		Value:     errors.New("test"),
		Code:      100,
		Msg:       "My error message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now,
		Prev:      nil,
	}

	expected := "value=test, code=100, msg=My error message, source=error_test.go:26, timestamp=" +
		time.UnixMicro(now).Format(time.RFC3339Nano)

	assert.Equal(t, expected, err.Error())
}

func TestErr_Error_NestedErrors(t *testing.T) {
	now := time.Now().UnixMicro()
	err2 := Err{
		Value:     errors.New("test 2"),
		Msg:       "My error message 2",
		Details:   nil,
		File:      "",
		Line:      0,
		Timestamp: now,
		Prev:      nil,
	}
	err := Err{
		Value:     errors.New("test"),
		Code:      500,
		Msg:       "My error message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now,
		Prev:      &err2,
	}

	expected := "value=test, code=500, msg=My error message, source=error_test.go:26, timestamp="
	expected += time.UnixMicro(now).Format(time.RFC3339Nano) + ", prev={value=test 2, msg=My error message 2, "
	expected += "timestamp=" + time.UnixMicro(now).Format(time.RFC3339Nano) + "}"

	assert.Equal(t, expected, err.Error())
}

func TestErr_Error_WithoutTimestamp(t *testing.T) {
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

func TestErr_Error_WithDetails(t *testing.T) {
	details := struct {
		Name string
		Age  int
	}{
		Name: "John Doe",
		Age:  23,
	}

	now := time.Now().UnixMicro()
	err := Err{
		Value:     errors.New("test"),
		Msg:       "My message",
		Details:   details,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now,
		Prev:      nil,
	}

	expected := "value=test, msg=My message, details={Name:John Doe Age:23}, source=error_test.go:26, timestamp=" +
		time.UnixMicro(now).Format(time.RFC3339Nano)

	assert.Equal(t, expected, err.Error())
}

func TestErr_Error_WithMsg(t *testing.T) {
	now := time.Now().UnixMicro()
	err := Err{
		Value:     errors.New("test"),
		Msg:       "",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now,
		Prev:      nil,
	}

	expected := "value=test, source=error_test.go:26, timestamp=" + time.UnixMicro(now).Format(time.RFC3339Nano)

	assert.Equal(t, expected, err.Error())
}

// ----------------------------------------------------------------------------
//
// Tests of Is()
//
// ----------------------------------------------------------------------------

func TestErr_Is(t *testing.T) {
	myErr := errors.New("my error")
	err := New(myErr, "My error message", nil, 0, nil)

	assert.True(t, err.Is(myErr))
}

func TestErr_Is_NestedErrors(t *testing.T) {
	myErr := errors.New("my error")
	myErr2 := errors.New("my error 2")
	myErr3 := errors.New("my error 3")
	err3 := New(myErr3, "My error message 3", nil, 0, nil)
	err2 := New(myErr2, "My error message 2", nil, 0, &err3)
	err := New(myErr, "My error message", nil, 0, &err2)

	assert.True(t, err.Is(myErr))
	assert.True(t, err.Is(myErr2))
	assert.True(t, err.Is(myErr3))
}

func TestErr_Is_False(t *testing.T) {
	myErr := errors.New("my error")
	myErr2 := errors.New("my error 2")
	err := New(myErr, "My error message", nil, 0, nil)

	assert.True(t, err.Is(myErr))
	assert.False(t, err.Is(myErr2))
}

// ----------------------------------------------------------------------------
//
// Tests of Unwrap()
//
// ----------------------------------------------------------------------------

func TestUnwrap(t *testing.T) {
	now := time.Now().UnixMicro()
	err2 := Err{
		Value:     errors.New("test 2"),
		Msg:       "My error message 2",
		Details:   nil,
		File:      "",
		Line:      0,
		Timestamp: now,
		Prev:      nil,
	}
	err := Err{
		Value:     errors.New("test"),
		Code:      500,
		Msg:       "My error message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now,
		Prev:      &err2,
	}

	assert.Equal(t, err.Unwrap(), &err2)
}

func TestUnwrapEmpty(t *testing.T) {
	now := time.Now().UnixMicro()
	err := Err{
		Value:     errors.New("test"),
		Code:      500,
		Msg:       "My error message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now,
		Prev:      nil,
	}

	assert.Equal(t, err.Unwrap(), nil)
}

// ----------------------------------------------------------------------------
//
// Tests of JSON()
//
// ----------------------------------------------------------------------------

func TestErr_JSON_Empty(t *testing.T) {
	e := Empty()
	expected := []byte("")
	result, err := e.JSON()

	fmt.Printf("%s\n", result)

	assert.Equal(t, Empty(), err)
	assert.Equal(t, expected, result)
}

func TestErr_JSON_Simple(t *testing.T) {
	now := time.Now().UnixMicro()
	e := Err{
		Value:     errors.New("test"),
		Code:      404,
		Msg:       "My error message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now,
		Prev:      nil,
	}

	expected := []byte(`{"value":"test","details":null,"timestamp":"` + time.UnixMicro(now).Format(time.RFC3339Nano) +
		`","code":404,"msg":"My error message","file":"error_test.go","line":26,"prev":null}`)
	result, err := e.JSON()

	assert.Equal(t, Empty(), err)
	assert.Equal(t, expected, result)
}

func TestErr_JSON_Detail(t *testing.T) {
	now := time.Now().UnixMicro()
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
		Timestamp: now,
		Prev:      nil,
	}

	expected := []byte(`{"value":"test","details":{"name":"John Doe","age":23},"timestamp":"` +
		time.UnixMicro(now).Format(time.RFC3339Nano) +
		`","msg":"My error message","file":"error_test.go","line":26,"prev":null}`)
	result, err := e.JSON()

	assert.Equal(t, Empty(), err)
	assert.Equal(t, expected, result)
}

func TestErr_JSON_NestedErrors(t *testing.T) {
	now := time.Now().UnixMicro()

	e := Err{
		Value:     errors.New("test"),
		Msg:       "My message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now,
		Prev: &Err{
			Value:     errors.New("test 2"),
			Msg:       "My message 2",
			Details:   nil,
			File:      "error_test.go",
			Line:      87,
			Timestamp: now,
			Prev:      nil,
		},
	}

	expected := []byte(`{"value":"test","details":null,"timestamp":"` +
		time.UnixMicro(now).Format(time.RFC3339Nano) +
		`","msg":"My message","file":"error_test.go","line":26,"prev":{"value":"test 2","details":null,"timestamp":"` +
		time.UnixMicro(now).Format(time.RFC3339Nano) +
		`","msg":"My message 2","file":"error_test.go","line":87,"prev":null}}`)
	result, err := e.JSON()

	assert.Equal(t, Empty(), err)
	assert.Equal(t, expected, result)
}

func TestErr_JSON_DetailError(t *testing.T) {
	now := time.Now().UnixMicro()
	details := struct {
		Channel chan int
	}{
		Channel: make(chan int),
	}

	e := Err{
		Value:     errors.New("test"),
		Code:      10,
		Msg:       "My error message",
		Details:   details,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now,
		Prev:      nil,
	}

	expected := []byte(`{"value":"test","details":null,"timestamp":"` + time.UnixMicro(now).Format(time.RFC3339Nano) +
		`","code":10,"msg":"My error message","file":"error_test.go","line":26,"prev":null}`)
	result, err := e.JSON()

	assert.Equal(t, Empty(), err)
	assert.Equal(t, expected, result)
}

func TestErr_JSON_WithStackTrace(t *testing.T) {
	e := New(errors.New("test"), "My error message", nil, 0, nil)
	result, err := e.JSON(true)

	assert.False(t, err.IsError())
	assert.True(t, strings.Contains(string(result), `"stack_trace":"`))
}

func TestErr_JSON_WithStackTrace_Empty(t *testing.T) {
	e := Empty()
	result, err := e.JSON(true)

	assert.False(t, err.IsError())
	assert.False(t, strings.Contains(string(result), `"stack_trace":"`))
}

func TestErr_JSON_WithNilValue(t *testing.T) {
	e := Err{
		Value:   nil,
		Code:    0,
		Msg:     "My error message",
		Details: nil,
	}
	result, err := e.JSON(true)

	assert.False(t, err.IsError())
	assert.Empty(t, result)
}

// ----------------------------------------------------------------------------
//
// Tests of JSONOrEmpty()
//
// ----------------------------------------------------------------------------

func TestErr_JSONOrEmpty_Empty(t *testing.T) {
	e := Empty()
	expected := []byte("")
	result := e.JSONOrEmpty()

	assert.Equal(t, expected, result)
}

func TestErr_JSONOrEmpty_Simple(t *testing.T) {
	now := time.Now().UnixMicro()
	e := Err{
		Value:     errors.New("test"),
		Code:      404,
		Msg:       "My error message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now,
		Prev:      nil,
	}

	expected := []byte(`{"value":"test","details":null,"timestamp":"` + time.UnixMicro(now).Format(time.RFC3339Nano) +
		`","code":404,"msg":"My error message","file":"error_test.go","line":26,"prev":null}`)
	result := e.JSONOrEmpty()

	assert.Equal(t, expected, result)
}

// ----------------------------------------------------------------------------
//
// Tests of ValueEq()
//
// ----------------------------------------------------------------------------

func TestErr_ValueEq(t *testing.T) {
	myErr := errors.New("test 1")

	err1 := New(myErr, "My error message 1", nil, 0, nil)
	err2 := New(myErr, "My error message 2", nil, 0, nil)

	assert.True(t, err1.ValueEq(err2))
	assert.True(t, err2.ValueEq(err1))
}

func TestErr_ValueEq_WithDifferentValues(t *testing.T) {
	myErr1 := errors.New("test 1")
	myErr2 := errors.New("test 2")

	err1 := New(myErr1, "My error message 1", nil, 0, nil)
	err2 := New(myErr1, "My error message 2", nil, 0, nil)

	// With different values
	err1.Value = myErr2

	assert.False(t, err1.ValueEq(err2))
	assert.False(t, err2.ValueEq(err1))
}

// ----------------------------------------------------------------------------
//
// Tests of Eq()
//
// ----------------------------------------------------------------------------

func TestErr_Eq_Simple(t *testing.T) {
	myErr1 := errors.New("test 1")
	myErr2 := errors.New("test 2")

	err := New(myErr2, "My error message", nil, 200, nil)
	err1 := New(myErr1, "My error message 1", nil, 300, &err)
	err2 := New(myErr1, "My error message 2", nil, 400, &err)

	assert.True(t, err1.Eq(err2))
	assert.True(t, err2.Eq(err1))
}

func TestErr_Eq_WithDifferentValues(t *testing.T) {
	myErr1 := errors.New("test 1")
	myErr2 := errors.New("test 2")
	myErr3 := errors.New("test 3")

	err := New(myErr2, "My error message", nil, 0, nil)
	err1 := New(myErr1, "My error message 1", nil, 0, &err)
	err2 := New(myErr1, "My error message 2", nil, 0, &err)

	// With different values for value
	err1.Value = myErr3

	assert.False(t, err1.Eq(err2))
	assert.False(t, err2.Eq(err1))
}

func TestErr_Eq_WithDifferentPrevValues(t *testing.T) {
	myErr1 := errors.New("test 1")
	myErr2 := errors.New("test 2")
	myErr3 := errors.New("test 3")

	err := New(myErr2, "My error message", nil, 0, nil)
	err1 := New(myErr1, "My error message 1", nil, 0, &err)
	err2 := New(myErr1, "My error message 2", nil, 0, &err)

	// With different values for prev
	err1.Prev.Value = myErr3

	assert.False(t, err1.Eq(err2))
	assert.False(t, err2.Eq(err1))
}

// ----------------------------------------------------------------------------
//
// Tests of Clone()
//
// ----------------------------------------------------------------------------

func TestErr_Clone_Empty(t *testing.T) {
	err := Empty()
	clone := err.Clone()

	assert.Equal(t, err, *clone)
}

func TestErr_Clone_Simple(t *testing.T) {
	myErr := errors.New("test 1")
	err := New(myErr, "My error message 1", nil, 0, nil)
	clone := err.Clone()

	assert.Equal(t, err, *clone)
}

func TestErr_Clone_NestedErrors(t *testing.T) {
	myErr1 := errors.New("test 1")
	myErr2 := errors.New("test 2")

	err := New(myErr2, "My error message", nil, 0, nil)
	err1 := New(myErr1, "My error message 1", nil, 0, &err)
	clone := err1.Clone()

	assert.Equal(t, err1, *clone)
}

func TestErr_Clone_NestedErrors_Empty(t *testing.T) {
	myErr1 := errors.New("test 1")
	myErr2 := errors.New("test 2")
	empty := Empty()

	err := New(myErr2, "My error message", nil, 0, &empty)
	err1 := New(myErr1, "My error message 1", nil, 0, &err)
	clone := err1.Clone()

	assert.Equal(t, err1, *clone)
}

// ----------------------------------------------------------------------------
//
// Tests of ToError()
//
// ----------------------------------------------------------------------------

func TestErr_ToError(t *testing.T) {
	now := time.Now().UnixMicro()
	e := Err{
		Value:     errors.New("test"),
		Code:      10,
		Msg:       "My error message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now,
		Prev:      nil,
	}
	assert.Equal(t, e.ToError(), errors.New("value=test, code=10, msg=My error message, source=error_test.go:26, timestamp="+time.UnixMicro(now).Format(time.RFC3339Nano)))
}

func TestErr_ToError_Empty(t *testing.T) {
	err := Empty()

	assert.Equal(t, err.ToError(), nil)
}

// ----------------------------------------------------------------------------
//
// Tests of Wrap()
//
// ----------------------------------------------------------------------------

// func TestErr_Wrap(t *testing.T) {
// }

// ----------------------------------------------------------------------------
//
// Tests of FromError()
//
// ----------------------------------------------------------------------------

func TestErr_FromError(t *testing.T) {
	expected := Err{
		Value:   errors.New("test"),
		Code:    0,
		Msg:     "",
		Details: nil,
		Prev:    nil,
	}

	err := FromError(errors.New("test"))

	assert.Equal(t, err.Value, expected.Value)
	assert.Equal(t, err.Code, expected.Code)
	assert.Equal(t, err.Msg, expected.Msg)
	assert.Equal(t, err.Details, expected.Details)
	assert.Equal(t, err.Prev, expected.Prev)
}

func TestErr_FromError_Empty(t *testing.T) {
	assert.Equal(t, FromError(nil), Empty())
}
