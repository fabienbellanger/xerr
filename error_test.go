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

func TestErr_NewErr_SimpleError(t *testing.T) {
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

func TestErr_NewErr_NestedErrors(t *testing.T) {
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

func TestErr_NewErr_EmptyError(t *testing.T) {
	err := NewErr(nil, "My error message", nil, nil)

	assert.Equal(t, Err{}, err)
}

// ----------------------------------------------------------------------------
//
// Tests of EmptyErr()
//
// ----------------------------------------------------------------------------

func TestErr_EmptyErr(t *testing.T) {
	assert.Equal(t, Err{}, EmptyErr())
}

// ----------------------------------------------------------------------------
//
// Tests of IsEmpty()
//
// ----------------------------------------------------------------------------

func TestErr_IsEmpty(t *testing.T) {
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

func TestErr_IsError(t *testing.T) {
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

func TestErr_Error_Empty(t *testing.T) {
	err := EmptyErr()

	assert.Equal(t, "", err.Error())
}

func TestErr_Error_NotEmpty(t *testing.T) {
	now := time.Now().UnixMicro()
	err := Err{
		Value:     errors.New("test"),
		Msg:       "My error message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now,
		Prev:      nil,
	}

	expected := "value=test, msg=My error message, source=error_test.go:26, timestamp=" +
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
		Msg:       "My error message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now,
		Prev:      &err2,
	}

	expected := "value=test, msg=My error message, source=error_test.go:26, timestamp="
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
	err := NewErr(myErr, "My error message", nil, nil)

	assert.True(t, err.Is(myErr))
}

func TestErr_Is_NestedErrors(t *testing.T) {
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

func TestErr_Is_False(t *testing.T) {
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

func TestErr_JSON_Empty(t *testing.T) {
	e := EmptyErr()
	expected := []byte("")
	result, err := e.JSON()

	fmt.Printf("%s\n", result)

	assert.Equal(t, EmptyErr(), err)
	assert.Equal(t, expected, result)
}

func TestErr_JSON_Simple(t *testing.T) {
	now := time.Now().UnixMicro()
	e := Err{
		Value:     errors.New("test"),
		Msg:       "My error message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now,
		Prev:      nil,
	}

	expected := []byte(`{"value":"test","details":null,"timestamp":"` + time.UnixMicro(now).Format(time.RFC3339Nano) +
		`","msg":"My error message","file":"error_test.go","line":26,"prev":null}`)
	result, err := e.JSON()

	assert.Equal(t, EmptyErr(), err)
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

	assert.Equal(t, EmptyErr(), err)
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

	assert.Equal(t, EmptyErr(), err)
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
		Msg:       "My error message",
		Details:   details,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now,
		Prev:      nil,
	}

	expected := []byte(`{"value":"test","details":null,"timestamp":"` + time.UnixMicro(now).Format(time.RFC3339Nano) +
		`","msg":"My error message","file":"error_test.go","line":26,"prev":null}`)
	result, err := e.JSON()

	assert.Equal(t, EmptyErr(), err)
	assert.Equal(t, expected, result)
}

// ----------------------------------------------------------------------------
//
// Tests of ValueEq()
//
// ----------------------------------------------------------------------------

func TestErr_ValueEq(t *testing.T) {
	var myErr = errors.New("test 1")

	err1 := NewErr(myErr, "My error message 1", nil, nil)
	err2 := NewErr(myErr, "My error message 2", nil, nil)

	assert.True(t, err1.ValueEq(err2))
	assert.True(t, err2.ValueEq(err1))
}

func TestErr_ValueEq_WithDifferentValues(t *testing.T) {
	var myErr1 = errors.New("test 1")
	var myErr2 = errors.New("test 2")

	err1 := NewErr(myErr1, "My error message 1", nil, nil)
	err2 := NewErr(myErr1, "My error message 2", nil, nil)

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
	var myErr1 = errors.New("test 1")
	var myErr2 = errors.New("test 2")

	err := NewErr(myErr2, "My error message", nil, nil)
	err1 := NewErr(myErr1, "My error message 1", nil, &err)
	err2 := NewErr(myErr1, "My error message 2", nil, &err)

	assert.True(t, err1.Eq(err2))
	assert.True(t, err2.Eq(err1))
}

func TestErr_Eq_WithDifferentValues(t *testing.T) {
	var myErr1 = errors.New("test 1")
	var myErr2 = errors.New("test 2")
	var myErr3 = errors.New("test 3")

	err := NewErr(myErr2, "My error message", nil, nil)
	err1 := NewErr(myErr1, "My error message 1", nil, &err)
	err2 := NewErr(myErr1, "My error message 2", nil, &err)

	// With different values for value
	err1.Value = myErr3

	assert.False(t, err1.Eq(err2))
	assert.False(t, err2.Eq(err1))
}

func TestErr_Eq_WithDifferentPrevValues(t *testing.T) {
	var myErr1 = errors.New("test 1")
	var myErr2 = errors.New("test 2")
	var myErr3 = errors.New("test 3")

	err := NewErr(myErr2, "My error message", nil, nil)
	err1 := NewErr(myErr1, "My error message 1", nil, &err)
	err2 := NewErr(myErr1, "My error message 2", nil, &err)

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
	err := EmptyErr()
	clone := err.Clone()

	assert.Equal(t, err, *clone)
}

func TestErr_Clone_Simple(t *testing.T) {
	var myErr = errors.New("test 1")
	err := NewErr(myErr, "My error message 1", nil, nil)
	clone := err.Clone()

	assert.Equal(t, err, *clone)
}

func TestErr_Clone_NestedErrors(t *testing.T) {
	var myErr1 = errors.New("test 1")
	var myErr2 = errors.New("test 2")

	err := NewErr(myErr2, "My error message", nil, nil)
	err1 := NewErr(myErr1, "My error message 1", nil, &err)
	clone := err1.Clone()

	assert.Equal(t, err1, *clone)
}
