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
		return 0, xerr.NewErr(errors.New("cannot divide by 0"), "Cannot divide by 0", nil, 0, nil)
	}
	return a / b, xerr.EmptyErr()
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
		return 0, xerr.NewErr(errors.New("cannot divide by 0"), "Cannot divide by 0", nil, 20, nil)
	}
	return a / b, xerr.EmptyErr()
}

func main() {
    _, err := divide(10, 0)
	if err.IsError() {
		log.Printf("Error: %v\n", 
            xerr.NewErr(errors.New("error in main()"), "Error in main()", [2]int{10, 0}, 10, &err))
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
goarch: arm64
pkg: github.com/fabienbellanger/xerr
cpu: Apple M3
BenchmarkErr_Error-8                     2832218               422.7 ns/op           264 B/op          3 allocs/op
BenchmarkErr_JSON_Simple-8                549771              2137 ns/op             720 B/op          5 allocs/op
BenchmarkErr_JSON_WithDetails-8           448258              2639 ns/op             752 B/op          7 allocs/op
BenchmarkErr_JSON_NestedErrors-8           75326             15828 ns/op            3989 B/op         14 allocs/op
BenchmarkErr_Is_Simple-8                100000000               11.65 ns/op            0 B/op          0 allocs/op
BenchmarkErr_Is_NestedErrors-8          55479355                22.90 ns/op            0 B/op          0 allocs/op
BenchmarkErr_Clone-8                    11963326               103.4 ns/op           288 B/op          3 allocs/op
```