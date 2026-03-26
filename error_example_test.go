package xerr

import (
	"errors"
	"fmt"
)

func ExampleNew() {
	err := New(errors.New("not found"), "user lookup failed", map[string]int{"user_id": 42}, 404, nil)
	fmt.Println(err.Value)
	fmt.Println(err.Code)
	fmt.Println(err.Msg)

	// Output:
	// not found
	// 404
	// user lookup failed
}

func ExampleNew_nil() {
	err := New(nil, "msg", nil, 0, nil)
	fmt.Println(err == nil)

	// Output: true
}

func ExampleNewSimple() {
	err := NewSimple(errors.New("timeout"), "service unavailable", nil)
	fmt.Println(err.Value)
	fmt.Println(err.Msg)

	// Output:
	// timeout
	// service unavailable
}

func ExampleErr_Wrap() {
	inner := NewSimple(errors.New("connection refused"), "db error", nil)
	outer := inner.Wrap(errors.New("query failed"), "cannot fetch user", nil, 500)

	fmt.Println(outer.Value)
	fmt.Println(outer.Prev.Value)

	// Output:
	// query failed
	// connection refused
}

func ExampleEmpty() {
	err := Empty()
	fmt.Println(err == nil)

	// Output: true
}

func ExampleErr_IsEmpty() {
	var err *Err
	fmt.Println(err.IsEmpty())

	err = New(errors.New("fail"), "", nil, 0, nil)
	fmt.Println(err.IsEmpty())

	// Output:
	// true
	// false
}

func ExampleErr_IsError() {
	err := New(errors.New("fail"), "", nil, 0, nil)
	fmt.Println(err.IsError())

	var nilErr *Err
	fmt.Println(nilErr.IsError())

	// Output:
	// true
	// false
}

func ExampleErr_Is() {
	sentinel := errors.New("not found")
	inner := New(sentinel, "inner", nil, 0, nil)
	outer := New(errors.New("wrapper"), "outer", nil, 0, inner)

	fmt.Println(outer.Is(sentinel))
	fmt.Println(outer.Is(errors.New("other")))

	// Output:
	// true
	// false
}

func ExampleErr_Unwrap() {
	inner := NewSimple(errors.New("root cause"), "inner", nil)
	outer := inner.Wrap(errors.New("wrapper"), "outer", nil, 0)

	prev := outer.Unwrap()
	fmt.Println(prev != nil)

	// Output: true
}

func ExampleErr_Clone() {
	original := New(errors.New("fail"), "msg", nil, 0, nil)
	clone := original.Clone()
	fmt.Println(clone.Value)
	fmt.Println(original.ValueEq(clone))

	// Output:
	// fail
	// true
}

func ExampleErr_ValueEq() {
	sentinel := errors.New("fail")
	err1 := New(sentinel, "first", nil, 0, nil)
	err2 := New(sentinel, "second", nil, 0, nil)
	err3 := New(errors.New("other"), "third", nil, 0, nil)

	fmt.Println(err1.ValueEq(err2))
	fmt.Println(err1.ValueEq(err3))

	// Output:
	// true
	// false
}

func ExampleErr_Eq() {
	sentinel := errors.New("fail")
	prev := New(errors.New("root"), "", nil, 0, nil)
	err1 := New(sentinel, "a", nil, 0, prev)
	err2 := New(sentinel, "b", nil, 0, prev)

	fmt.Println(err1.Eq(err2))

	// Output: true
}

func ExampleErr_ToError() {
	err := New(errors.New("fail"), "msg", nil, 0, nil)
	var e error = err.ToError()
	fmt.Println(e != nil)

	var nilErr *Err
	fmt.Println(nilErr.ToError() == nil)

	// Output:
	// true
	// true
}

func ExampleFromError() {
	stdErr := errors.New("standard error")
	err := FromError(stdErr)
	fmt.Println(err.Value)
	fmt.Println(FromError(nil) == nil)

	// Output:
	// standard error
	// true
}

func ExampleErr_JSON() {
	err := New(errors.New("fail"), "msg", nil, 0, nil)
	b, _ := err.JSON()
	fmt.Println(len(b) > 0)

	// Output: true
}

func ExampleErr_JSONOrEmpty() {
	var err *Err
	b := err.JSONOrEmpty()
	fmt.Println(len(b))

	// Output: 0
}
