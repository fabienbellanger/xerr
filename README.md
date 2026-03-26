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

func divide(a, b int) (int, *xerr.Err) {
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
	_, err = divide(10, 0)
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

func divide(a, b int) (int, *xerr.Err) {
	if b == 0 {
		return 0, xerr.New(errors.New("cannot divide by 0"), "Cannot divide by 0", nil, 20, nil)
	}
	return a / b, xerr.Empty()
}

func main() {
	_, xe := divide(10, 0)
	if xe.IsError() {
		log.Printf("Error: %v\n",
			xe.Wrap(errors.New("error in main()"), "Error in main()", [2]int{10, 0}, 10))
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
BenchmarkErr_New-8                 254626       4621 ns/op      1416 B/op       5 allocs/op
BenchmarkErr_NewSimple-8           184974       6640 ns/op      1664 B/op       7 allocs/op
BenchmarkErr_FromError-8           236146       5148 ns/op      1400 B/op       4 allocs/op
BenchmarkErr_Error-8              1668969        721 ns/op       600 B/op      14 allocs/op
BenchmarkErr_Wrap-8                176904       6825 ns/op      1792 B/op       8 allocs/op
BenchmarkErr_JSON_Simple-8         573890       2063 ns/op       752 B/op       5 allocs/op
BenchmarkErr_JSON_WithDetails-8    460956       2603 ns/op       784 B/op       7 allocs/op
BenchmarkErr_JSON_NestedErrors-8    21578      55670 ns/op     15920 B/op      20 allocs/op
BenchmarkErr_Is_Simple-8        152763013       7.86 ns/op         0 B/op       0 allocs/op
BenchmarkErr_Is_NestedErrors-8   70081856      17.11 ns/op         0 B/op       0 allocs/op
BenchmarkErr_JSONOrEmpty-8         559858       2082 ns/op       752 B/op       5 allocs/op
BenchmarkErr_Eq-8               100000000      10.16 ns/op         0 B/op       0 allocs/op
BenchmarkErr_Clone_4-8            9479798        128 ns/op       512 B/op       4 allocs/op
BenchmarkErr_Clone_8-8            4535299        258 ns/op      1024 B/op       8 allocs/op
BenchmarkErr_Clone_16-8           2370424        509 ns/op      2048 B/op      16 allocs/op
```