package xerr

import (
	"errors"
	"testing"
)

func BenchmarkErr_Error(b *testing.B) {
	for b.Loop() {
		e := NewErr(errors.New("test"), "My error message", nil, nil)
		_ = e
	}
}

func BenchmarkErr_JSON_Simple(b *testing.B) {
	e := NewErr(errors.New("test"), "My error message", nil, nil)

	for b.Loop() {
		r, _ := e.JSON()
		_ = r
	}
}

func BenchmarkErr_JSON_WithDetails(b *testing.B) {
	details := struct {
		Channel chan int
	}{
		Channel: make(chan int),
	}
	e := NewErr(errors.New("test"), "My error message", details, nil)

	for b.Loop() {
		r, _ := e.JSON()
		_ = r
	}
}

func BenchmarkErr_JSON_NestedErrors(b *testing.B) {
	e3 := NewErr(errors.New("test 2"), "My error message 2", nil, nil)
	e2 := NewErr(errors.New("test 2"), "My error message 2", nil, &e3)
	e1 := NewErr(errors.New("test 1"), "My error message 1", nil, &e2)
	e := NewErr(errors.New("test"), "My error message", nil, &e1)

	for b.Loop() {
		r, _ := e.JSON()
		_ = r
	}
}

func BenchmarkErr_Is_Simple(b *testing.B) {
	var myErr = errors.New("my error")
	e := NewErr(myErr, "My error message", nil, nil)

	for b.Loop() {
		ok := e.Is(myErr)
		_ = ok
	}
}

func BenchmarkErr_Is_NestedErrors(b *testing.B) {
	var myErr = errors.New("my error")
	e1 := NewErr(myErr, "My error message 1", nil, nil)
	e := NewErr(errors.New("test"), "My error message", nil, &e1)

	for b.Loop() {
		ok := e.Is(myErr)
		_ = ok
	}
}
