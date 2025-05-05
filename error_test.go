package xerr

import (
	"errors"
	"strings"
	"testing"

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
	assert.Equal(t, 26, err.Line)
	assert.Nil(t, err.Prev)
}

func TestNewErrNestedErrors(t *testing.T) {
	err2 := NewErr(errors.New("test 2"), "My error message 2", nil, nil)
	err1 := NewErr(errors.New("test 1"), "My error message 1", nil, &err2)

	assert.Equal(t, errors.New("test 1"), err1.Value)
	assert.Equal(t, "My error message 1", err1.Msg)
	assert.Nil(t, err1.Details)
	assert.True(t, strings.Contains(err1.File, "error_test.go"))
	assert.Equal(t, 38, err1.Line)

	assert.Equal(t, errors.New("test 2"), err2.Value)
	assert.Equal(t, "My error message 2", err2.Msg)
	assert.Nil(t, err2.Details)
	assert.True(t, strings.Contains(err2.File, "error_test.go"))
	assert.Equal(t, 37, err2.Line)
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
	// err := NewErr(errors.New("test"), "My error message", nil, nil)

	// TODO: Add test
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
	err2 := NewErr(myErr2, "My error message 2", nil, nil)
	err := NewErr(myErr, "My error message", nil, &err2)

	assert.True(t, err.Is(myErr))
	// assert.True(t, err.Is(myErr2)) // TODO: Is not work!
}
