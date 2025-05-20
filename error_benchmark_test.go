package xerr

import (
	"errors"
	"testing"
)

func BenchmarkErr_Error(b *testing.B) {
	for b.Loop() {
		e := NewErr(errors.New("test"), "My error message", nil, 0, nil)
		_ = e
	}
}

func BenchmarkErr_JSON_Simple(b *testing.B) {
	e := NewErr(errors.New("test"), "My error message", nil, 0, nil)

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
	e := NewErr(errors.New("test"), "My error message", details, 0, nil)

	for b.Loop() {
		r, _ := e.JSON()
		_ = r
	}
}

func BenchmarkErr_JSON_NestedErrors(b *testing.B) {
	e3 := NewErr(errors.New("test 2"), "My error message 2", nil, 503, nil)
	e2 := NewErr(errors.New("test 2"), "My error message 2", nil, 502, &e3)
	e1 := NewErr(errors.New("test 1"), "My error message 1", nil, 501, &e2)
	e := NewErr(errors.New("test"), "My error message", nil, 500, &e1)

	for b.Loop() {
		r, _ := e.JSON()
		_ = r
	}
}

func BenchmarkErr_Is_Simple(b *testing.B) {
	var myErr = errors.New("my error")
	e := NewErr(myErr, "My error message", nil, 0, nil)

	for b.Loop() {
		ok := e.Is(myErr)
		_ = ok
	}
}

func BenchmarkErr_Is_NestedErrors(b *testing.B) {
	var myErr = errors.New("my error")
	e1 := NewErr(myErr, "My error message 1", nil, 401, nil)
	e := NewErr(errors.New("test"), "My error message", nil, 400, &e1)

	for b.Loop() {
		ok := e.Is(myErr)
		_ = ok
	}
}

func BenchmarkErr_Clone(b *testing.B) {
	var myErr1 = errors.New("my error 1")
	var myErr2 = errors.New("my error 2")
	e2 := NewErr(myErr2, "My error message 2", nil, 0, nil)
	e1 := NewErr(myErr1, "My error message 1", nil, 0, &e2)
	e := NewErr(errors.New("test"), "My error message", nil, 0, &e1)

	for b.Loop() {
		err := e.Clone()
		_ = err
	}
}
