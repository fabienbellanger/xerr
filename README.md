# xerr: An extended error library for Go

[![Build status](https://github.com/fabienbellanger/xerr/actions/workflows/CI.yml/badge.svg?branch=main)](https://github.com/fabienbellanger/xerr/actions/workflows/CI.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/fabienbellanger/xerr)](https://goreportcard.com/report/github.com/fabienbellanger/xerr)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=square)](https://pkg.go.dev/github.com/fabienbellanger/xerr)

> xerr is a Go package that provides extended error handling capabilities, including error wrapping, stack traces, error codes, and additional context information to make debugging and error management easier in Go applications.

## Installation

In your project:
```bash
go get github.com/fabienbellanger/xerr
```

## Examples

### Simple error
```go
package main

import (
	"errors"
	"log"

	"github.com/fabienbellanger/xerr"
)

func divide(a, b int) (int, xerr.Err) {
	if b == 0 {
		return 0, xerr.New(errors.New("cannot divide by 0"), "Cannot divide by 0", nil, 0, nil)
	}
	return a / b, xerr.Empty()
}

func main() {
    // Without error
	d, err := divide(10, 10)
	if err.IsEmpty() {
		log.Printf("No error: %d\n", d)
	}

    // With error
	_, err := divide(10, 0)
	if err.IsError() {
		log.Printf("Error: %v\n", err)
	}
}
```

### Nested error
```go
package main

import (
	"errors"
	"log"

	"github.com/fabienbellanger/xerr"
)

func divide(a, b int) (int, xerr.Err) {
	if b == 0 {
		return 0, xerr.New(errors.New("cannot divide by 0"), "Cannot divide by 0", nil, 20, nil)
	}
	return a / b, xerr.Empty()
}

func main() {
    _, err := divide(10, 0)
	if err.IsError() {
		log.Printf("Error: %v\n", 
            xerr.New(errors.New("error in main()"), "Error in main()", [2]int{10, 0}, 10, &err))
	}
}
```

## Benchmarks

Run:
```bash
make bench
```

Results:
```
goos: darwin
goarch: amd64
pkg: github.com/fabienbellanger/xerr
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkErr_Error-12                    2479927               506.7 ns/op           264 B/op          3 allocs/op
BenchmarkErr_JSON_Simple-12               559191              2186 ns/op             688 B/op          5 allocs/op
BenchmarkErr_JSON_WithDetails-12          465482              2469 ns/op             720 B/op          7 allocs/op
BenchmarkErr_JSON_NestedErrors-12          83658             14154 ns/op            3989 B/op         14 allocs/op
BenchmarkErr_Is_Simple-12               100000000               10.34 ns/op            0 B/op          0 allocs/op
BenchmarkErr_Is_NestedErrors-12         57608085                20.82 ns/op            0 B/op          0 allocs/op
BenchmarkErr_Clone_4-12                  8457126               153.7 ns/op           384 B/op          4 allocs/op
BenchmarkErr_Clone_8-12                  4106058               285.3 ns/op           768 B/op          8 allocs/op
BenchmarkErr_Clone_16-12                 2079298               575.0 ns/op          1536 B/op         16 allocs/op
```