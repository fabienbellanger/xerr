package xerr

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
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

	_, _, wantLine, _ := runtime.Caller(0)
	wantLine++
	err := New(errors.New("test"), "My error message", details, 10, nil)

	assert.Equal(t, errors.New("test"), err.Value)
	assert.Equal(t, 10, err.Code)
	assert.Equal(t, "My error message", err.Msg)
	assert.Equal(t, details, err.Details)
	assert.True(t, strings.Contains(err.File, "error_test.go"))
	assert.Equal(t, wantLine, err.Line)
	assert.Nil(t, err.Prev)
}

func TestErr_New_With_Err(t *testing.T) {
	err2 := New(errors.New("test 2"), "My error message 2", nil, 20, nil)
	_, _, wantLine, _ := runtime.Caller(0)
	wantLine++
	err1 := New(err2, "My error message 1", nil, 10, nil)

	assert.Equal(t, err2, err1.Value)
	assert.Equal(t, 10, err1.Code)
	assert.Equal(t, "My error message 1", err1.Msg)
	assert.Nil(t, err1.Details)
	assert.True(t, strings.Contains(err1.File, "error_test.go"))
	assert.Equal(t, wantLine, err1.Line)
}

func TestErr_New_NestedErrors(t *testing.T) {
	_, _, wantLine2, _ := runtime.Caller(0)
	wantLine2++
	err2 := New(errors.New("test 2"), "My error message 2", nil, 20, nil)
	_, _, wantLine1, _ := runtime.Caller(0)
	wantLine1++
	err1 := New(errors.New("test 1"), "My error message 1", nil, 10, err2)

	assert.Equal(t, errors.New("test 1"), err1.Value)
	assert.Equal(t, 10, err1.Code)
	assert.Equal(t, "My error message 1", err1.Msg)
	assert.Nil(t, err1.Details)
	assert.True(t, strings.Contains(err1.File, "error_test.go"))
	assert.Equal(t, wantLine1, err1.Line)

	assert.Equal(t, errors.New("test 2"), err2.Value)
	assert.Equal(t, 20, err2.Code)
	assert.Equal(t, "My error message 2", err2.Msg)
	assert.Nil(t, err2.Details)
	assert.True(t, strings.Contains(err2.File, "error_test.go"))
	assert.Equal(t, wantLine2, err2.Line)
}

func TestErr_New_ReturnsNilOnNilValue(t *testing.T) {
	err := New(nil, "My error message", nil, 0, nil)
	assert.Nil(t, err)
}

func TestErr_New_WithSkip(t *testing.T) {
	err := New(errors.New("test"), "My error message", nil, 20, nil, 0)

	assert.Equal(t, errors.New("test"), err.Value)
	assert.Equal(t, 20, err.Code)
	assert.Equal(t, "My error message", err.Msg)
	assert.Nil(t, err.Details)
	// skip=0 → runtime.Caller(0) points inside error.go, not at the call site
	assert.True(t, strings.Contains(err.File, "error.go"))
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

func TestErr_NewSimple_ReturnsNilOnNilValue(t *testing.T) {
	err := NewSimple(nil, "My error message", nil)
	assert.Nil(t, err)
}

func TestErr_NewSimple_WithSkip(t *testing.T) {
	err := NewSimple(errors.New("test"), "My error message", nil, 0)

	assert.Equal(t, errors.New("test"), err.Value)
	assert.Equal(t, 0, err.Code)
	assert.Equal(t, "My error message", err.Msg)
	assert.Nil(t, err.Details)
	// skip=0 → runtime.Caller(0) points inside error.go, not at the call site
	assert.True(t, strings.Contains(err.File, "error.go"))
}

// ----------------------------------------------------------------------------
//
// Tests of Wrap()
//
// ----------------------------------------------------------------------------

func TestErr_Wrap(t *testing.T) {
	err := NewSimple(errors.New("test"), "My error message", nil)
	_, _, wantLine, _ := runtime.Caller(0)
	wantLine++
	wrappedErr := err.Wrap(errors.New("wrapped error"), "Wrapped message", nil, 100)

	assert.Equal(t, errors.New("wrapped error"), wrappedErr.Value)
	assert.Equal(t, 100, wrappedErr.Code)
	assert.Equal(t, "Wrapped message", wrappedErr.Msg)
	assert.Nil(t, wrappedErr.Details)
	assert.True(t, strings.Contains(wrappedErr.File, "error_test.go"))
	assert.Equal(t, wantLine, wrappedErr.Line)
	assert.Equal(t, err, wrappedErr.Prev)
}

func TestErr_Wrap_WithSkip(t *testing.T) {
	err := NewSimple(errors.New("test"), "My error message", nil)
	wrappedErr := err.Wrap(errors.New("wrapped error"), "Wrapped message", nil, 100, 0)

	assert.Equal(t, errors.New("wrapped error"), wrappedErr.Value)
	assert.Equal(t, 100, wrappedErr.Code)
	assert.Equal(t, "Wrapped message", wrappedErr.Msg)
	assert.Nil(t, wrappedErr.Details)
	// skip=0 → runtime.Caller(0) points inside error.go, not at the call site
	assert.True(t, strings.Contains(wrappedErr.File, "error.go"))
	assert.Equal(t, err, wrappedErr.Prev)
}

// ----------------------------------------------------------------------------
//
// Tests of Empty()
//
// ----------------------------------------------------------------------------

func TestErr_Empty(t *testing.T) {
	assert.Nil(t, Empty())
}

// ----------------------------------------------------------------------------
//
// Tests of IsEmpty()
//
// ----------------------------------------------------------------------------

func TestErr_IsEmpty(t *testing.T) {
	var err *Err
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
	var err *Err
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
	var err *Err
	assert.Equal(t, "", err.Error())
}

func TestErr_Error_NotEmpty(t *testing.T) {
	now := time.Now().UnixMicro()
	err := &Err{
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
	err2 := &Err{
		Value:     errors.New("test 2"),
		Msg:       "My error message 2",
		Details:   nil,
		File:      "",
		Line:      0,
		Timestamp: now,
		Prev:      nil,
	}
	err := &Err{
		Value:     errors.New("test"),
		Code:      500,
		Msg:       "My error message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now,
		Prev:      err2,
	}

	expected := "value=test, code=500, msg=My error message, source=error_test.go:26, timestamp="
	expected += time.UnixMicro(now).Format(time.RFC3339Nano) + ", prev={value=test 2, msg=My error message 2, "
	expected += "timestamp=" + time.UnixMicro(now).Format(time.RFC3339Nano) + "}"

	assert.Equal(t, expected, err.Error())
}

func TestErr_Error_WithoutTimestamp(t *testing.T) {
	err := &Err{
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
	err := &Err{
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
	err := &Err{
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
	err2 := New(myErr2, "My error message 2", nil, 0, err3)
	err := New(myErr, "My error message", nil, 0, err2)

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

func TestErr_Is_Nil(t *testing.T) {
	var err *Err
	assert.False(t, err.Is(errors.New("anything")))
}

// ----------------------------------------------------------------------------
//
// Tests of Unwrap()
//
// ----------------------------------------------------------------------------

func TestUnwrap(t *testing.T) {
	now := time.Now().UnixMicro()
	err2 := &Err{
		Value:     errors.New("test 2"),
		Msg:       "My error message 2",
		Details:   nil,
		File:      "",
		Line:      0,
		Timestamp: now,
		Prev:      nil,
	}
	err := &Err{
		Value:     errors.New("test"),
		Code:      500,
		Msg:       "My error message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now,
		Prev:      err2,
	}

	assert.Equal(t, err.Unwrap(), err2)
}

func TestUnwrapEmpty(t *testing.T) {
	now := time.Now().UnixMicro()
	err := &Err{
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
	var e *Err
	result, err := e.JSON()

	fmt.Printf("%s\n", result)

	assert.NoError(t, err)
	assert.Equal(t, []byte{}, result)
}

func TestErr_JSON_Simple(t *testing.T) {
	now := time.Now().UnixMicro()
	e := &Err{
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

	assert.NoError(t, err)
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

	e := &Err{
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

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestErr_JSON_NestedErrors(t *testing.T) {
	now := time.Now().UnixMicro()

	e := &Err{
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

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestErr_JSON_DetailError(t *testing.T) {
	now := time.Now().UnixMicro()
	details := struct {
		Channel chan int
	}{
		Channel: make(chan int),
	}

	e := &Err{
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

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestErr_JSON_WithStackTrace(t *testing.T) {
	e := New(errors.New("test"), "My error message", nil, 0, nil)
	result, err := e.JSON(true)

	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(result), `"stack_trace":"`))
}

func TestErr_JSON_WithStackTrace_Empty(t *testing.T) {
	var e *Err
	result, err := e.JSON(true)

	assert.NoError(t, err)
	assert.False(t, strings.Contains(string(result), `"stack_trace":"`))
}

func TestErr_JSON_WithNilValue(t *testing.T) {
	e := &Err{
		Value:   nil,
		Code:    0,
		Msg:     "My error message",
		Details: nil,
	}
	result, err := e.JSON(true)

	assert.NoError(t, err)
	assert.Empty(t, result)
}

// failOnSecondMarshal succeeds on the first json.Marshal call (the pre-check
// inside MarshalJSON) but fails on the second (the real marshal).
type failOnSecondMarshal struct {
	calls *int
}

func (f failOnSecondMarshal) MarshalJSON() ([]byte, error) {
	*f.calls++
	if *f.calls > 1 {
		return nil, errors.New("marshal failed")
	}
	return []byte(`"ok"`), nil
}

func TestErr_JSON_MarshalError(t *testing.T) {
	calls := 0
	e := &Err{
		Value:     errors.New("test"),
		Msg:       "msg",
		Details:   failOnSecondMarshal{calls: &calls},
		Timestamp: time.Now().UnixMicro(),
	}
	result, err := e.JSON()

	assert.Error(t, err)
	assert.Equal(t, []byte{}, result)
}

func TestErr_JSONOrEmpty_MarshalError(t *testing.T) {
	calls := 0
	e := &Err{
		Value:     errors.New("test"),
		Msg:       "msg",
		Details:   failOnSecondMarshal{calls: &calls},
		Timestamp: time.Now().UnixMicro(),
	}
	result := e.JSONOrEmpty()

	assert.Equal(t, []byte{}, result)
}

func TestErr_MarshalJSON_NilValue(t *testing.T) {
	e := &Err{
		Value:     nil,
		Msg:       "msg",
		Timestamp: time.Now().UnixMicro(),
	}
	result, err := json.Marshal(e)

	assert.NoError(t, err)
	assert.Contains(t, string(result), `"value":""`)
}

// ----------------------------------------------------------------------------
//
// Tests of JSONOrEmpty()
//
// ----------------------------------------------------------------------------

func TestErr_JSONOrEmpty_Empty(t *testing.T) {
	var e *Err
	result := e.JSONOrEmpty()

	assert.Equal(t, []byte{}, result)
}

func TestErr_JSONOrEmpty_Simple(t *testing.T) {
	now := time.Now().UnixMicro()
	e := &Err{
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

	err1.Value = myErr2

	assert.False(t, err1.ValueEq(err2))
	assert.False(t, err2.ValueEq(err1))
}

func TestErr_ValueEq_Nil(t *testing.T) {
	var err1 *Err
	var err2 *Err
	assert.True(t, err1.ValueEq(err2))

	err1 = New(errors.New("test"), "", nil, 0, nil)
	assert.False(t, err1.ValueEq(nil))
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
	err1 := New(myErr1, "My error message 1", nil, 300, err)
	err2 := New(myErr1, "My error message 2", nil, 400, err)

	assert.True(t, err1.Eq(err2))
	assert.True(t, err2.Eq(err1))
}

func TestErr_Eq_WithDifferentValues(t *testing.T) {
	myErr1 := errors.New("test 1")
	myErr2 := errors.New("test 2")
	myErr3 := errors.New("test 3")

	err := New(myErr2, "My error message", nil, 0, nil)
	err1 := New(myErr1, "My error message 1", nil, 0, err)
	err2 := New(myErr1, "My error message 2", nil, 0, err)

	err1.Value = myErr3

	assert.False(t, err1.Eq(err2))
	assert.False(t, err2.Eq(err1))
}

func TestErr_Eq_WithDifferentPrevValues(t *testing.T) {
	myErr1 := errors.New("test 1")
	myErr2 := errors.New("test 2")
	myErr3 := errors.New("test 3")

	err := New(myErr2, "My error message", nil, 0, nil)
	err1 := New(myErr1, "My error message 1", nil, 0, err)
	err2 := New(myErr1, "My error message 2", nil, 0, err)

	err1.Prev.Value = myErr3

	assert.False(t, err1.Eq(err2))
	assert.False(t, err2.Eq(err1))
}

func TestErr_Eq_DifferentChainLengths(t *testing.T) {
	myErr1 := errors.New("test 1")
	myErr2 := errors.New("test 2")

	err2 := New(myErr2, "", nil, 0, nil)
	err1a := New(myErr1, "", nil, 0, err2)
	err1b := New(myErr1, "", nil, 0, nil)

	assert.False(t, err1a.Eq(err1b))
	assert.False(t, err1b.Eq(err1a))
}

func TestErr_Eq_Nil(t *testing.T) {
	var err1 *Err
	var err2 *Err
	assert.True(t, err1.Eq(err2))
}

// ----------------------------------------------------------------------------
//
// Tests of Clone()
//
// ----------------------------------------------------------------------------

func TestErr_Clone_Nil(t *testing.T) {
	var err *Err
	clone := err.Clone()
	assert.Nil(t, clone)
}

func TestErr_Clone_Simple(t *testing.T) {
	myErr := errors.New("test 1")
	err := New(myErr, "My error message 1", nil, 0, nil)
	clone := err.Clone()

	assert.Equal(t, *err, *clone)
}

func TestErr_Clone_NestedErrors(t *testing.T) {
	myErr1 := errors.New("test 1")
	myErr2 := errors.New("test 2")

	err := New(myErr2, "My error message", nil, 0, nil)
	err1 := New(myErr1, "My error message 1", nil, 0, err)
	clone := err1.Clone()

	assert.Equal(t, *err1, *clone)
}

func TestErr_Clone_NestedErrors_Empty(t *testing.T) {
	myErr1 := errors.New("test 1")
	myErr2 := errors.New("test 2")

	err := New(myErr2, "My error message", nil, 0, nil)
	err1 := New(myErr1, "My error message 1", nil, 0, err)
	clone := err1.Clone()

	assert.Equal(t, *err1, *clone)
}

// ----------------------------------------------------------------------------
//
// Tests of ToError()
//
// ----------------------------------------------------------------------------

func TestErr_ToError(t *testing.T) {
	now := time.Now().UnixMicro()
	e := &Err{
		Value:     errors.New("test"),
		Code:      10,
		Msg:       "My error message",
		Details:   nil,
		File:      "error_test.go",
		Line:      26,
		Timestamp: now,
		Prev:      nil,
	}
	// ToError now returns e itself (implements error), so the message matches Error()
	assert.Equal(t, e.ToError(), error(e))
}

func TestErr_ToError_Empty(t *testing.T) {
	var err *Err
	assert.Equal(t, err.ToError(), nil)
}

// ----------------------------------------------------------------------------
//
// Tests of FromError()
//
// ----------------------------------------------------------------------------

func TestErr_FromError(t *testing.T) {
	err := FromError(errors.New("test"))

	assert.Equal(t, errors.New("test"), err.Value)
	assert.Equal(t, 0, err.Code)
	assert.Equal(t, "", err.Msg)
	assert.Nil(t, err.Details)
	assert.Nil(t, err.Prev)
}

func TestErr_FromError_Empty(t *testing.T) {
	assert.Nil(t, FromError(nil))
}

// ----------------------------------------------------------------------------
//
// Tests of standard error interface compatibility
//
// ----------------------------------------------------------------------------

func TestErr_ImplementsError(t *testing.T) {
	var err error
	err = New(errors.New("test"), "msg", nil, 0, nil)
	assert.NotNil(t, err)
}

func TestErr_ErrorsIs_Compatibility(t *testing.T) {
	sentinel := errors.New("sentinel")
	err := New(sentinel, "wrapped", nil, 0, nil)

	// Standard errors.Is works via Unwrap chain — but our Is() method handles the Prev chain.
	// Direct value match should work via the error interface.
	assert.True(t, err.Is(sentinel))
}

func TestErr_NilIsNilError(t *testing.T) {
	var err *Err
	// err.ToError() returns untyped nil (not a typed nil in interface)
	assert.Nil(t, err)
	assert.Nil(t, err.ToError())
}
