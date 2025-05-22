package xerr

import (
	"errors"
	"testing"
)

func BenchmarkErr_Error(b *testing.B) {
	for b.Loop() {
		e := New(errors.New("test"), "My error message", nil, 0, nil)
		_ = e
	}
}

func BenchmarkErr_JSON_Simple(b *testing.B) {
	e := New(errors.New("test"), "My error message", nil, 0, nil)

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
	e := New(errors.New("test"), "My error message", details, 0, nil)

	for b.Loop() {
		r, _ := e.JSON()
		_ = r
	}
}

func BenchmarkErr_JSON_NestedErrors(b *testing.B) {
	e3 := New(errors.New("test 2"), "My error message 2", nil, 503, nil)
	e2 := New(errors.New("test 2"), "My error message 2", nil, 502, &e3)
	e1 := New(errors.New("test 1"), "My error message 1", nil, 501, &e2)
	e := New(errors.New("test"), "My error message", nil, 500, &e1)

	for b.Loop() {
		r, _ := e.JSON()
		_ = r
	}
}

func BenchmarkErr_Is_Simple(b *testing.B) {
	myErr := errors.New("my error")
	e := New(myErr, "My error message", nil, 0, nil)

	for b.Loop() {
		ok := e.Is(myErr)
		_ = ok
	}
}

func BenchmarkErr_Is_NestedErrors(b *testing.B) {
	myErr := errors.New("my error")
	e1 := New(myErr, "My error message 1", nil, 401, nil)
	e := New(errors.New("test"), "My error message", nil, 400, &e1)

	for b.Loop() {
		ok := e.Is(myErr)
		_ = ok
	}
}

func BenchmarkErr_Clone_4(b *testing.B) {
	myErr1 := errors.New("my error 1")
	myErr2 := errors.New("my error 2")
	myErr3 := errors.New("my error 3")
	e3 := New(myErr3, "My error message 3", nil, 0, nil)
	e2 := New(myErr2, "My error message 2", nil, 0, &e3)
	e1 := New(myErr1, "My error message 1", nil, 0, &e2)
	e := New(errors.New("test"), "My error message", nil, 0, &e1)

	for b.Loop() {
		err := e.Clone()
		_ = err
	}
}

func BenchmarkErr_Clone_8(b *testing.B) {
	myErr1 := errors.New("my error 1")
	myErr2 := errors.New("my error 2")
	myErr3 := errors.New("my error 3")
	myErr4 := errors.New("my error 4")
	myErr5 := errors.New("my error 5")
	myErr6 := errors.New("my error 6")
	myErr7 := errors.New("my error 7")
	e7 := New(myErr7, "My error message 7", nil, 0, nil)
	e6 := New(myErr6, "My error message 6", nil, 0, &e7)
	e5 := New(myErr5, "My error message 5", nil, 0, &e6)
	e4 := New(myErr4, "My error message 4", nil, 0, &e5)
	e3 := New(myErr3, "My error message 3", nil, 0, &e4)
	e2 := New(myErr2, "My error message 2", nil, 0, &e3)
	e1 := New(myErr1, "My error message 1", nil, 0, &e2)
	e := New(errors.New("test"), "My error message", nil, 0, &e1)

	for b.Loop() {
		err := e.Clone()
		_ = err
	}
}

func BenchmarkErr_Clone_16(b *testing.B) {
	myErr1 := errors.New("my error 1")
	myErr2 := errors.New("my error 2")
	myErr3 := errors.New("my error 3")
	myErr4 := errors.New("my error 4")
	myErr5 := errors.New("my error 5")
	myErr6 := errors.New("my error 6")
	myErr7 := errors.New("my error 7")
	myErr8 := errors.New("my error 8")
	myErr9 := errors.New("my error 9")
	myErr10 := errors.New("my error 10")
	myErr11 := errors.New("my error 11")
	myErr12 := errors.New("my error 12")
	myErr13 := errors.New("my error 13")
	myErr14 := errors.New("my error 14")
	myErr15 := errors.New("my error 15")
	e15 := New(myErr15, "My error message 15", nil, 0, nil)
	e14 := New(myErr14, "My error message 14", nil, 0, &e15)
	e13 := New(myErr13, "My error message 13", nil, 0, &e14)
	e12 := New(myErr12, "My error message 12", nil, 0, &e13)
	e11 := New(myErr11, "My error message 11", nil, 0, &e12)
	e10 := New(myErr10, "My error message 10", nil, 0, &e11)
	e9 := New(myErr9, "My error message 9", nil, 0, &e10)
	e8 := New(myErr8, "My error message 8", nil, 0, &e9)
	e7 := New(myErr7, "My error message 7", nil, 0, &e8)
	e6 := New(myErr6, "My error message 6", nil, 0, &e7)
	e5 := New(myErr5, "My error message 5", nil, 0, &e6)
	e4 := New(myErr4, "My error message 4", nil, 0, &e5)
	e3 := New(myErr3, "My error message 3", nil, 0, &e4)
	e2 := New(myErr2, "My error message 2", nil, 0, &e3)
	e1 := New(myErr1, "My error message 1", nil, 0, &e2)
	e := New(errors.New("test"), "My error message", nil, 0, &e1)

	for b.Loop() {
		err := e.Clone()
		_ = err
	}
}
