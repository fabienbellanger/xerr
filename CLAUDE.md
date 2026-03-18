# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make test           # Run tests with coverage
make test-verbose   # Run tests with verbose output
make lint           # Format + vet
make bench          # Run benchmarks
make doc            # Launch godoc on port 9898
go test -run TestName ./...  # Run a single test
```

## Architecture

This is a single-file Go library (`error.go`) with no subpackages. The entire public API lives in `package xerr`.

### Core type: `Err`

`Err` is a struct that wraps `error`. All constructors return `*Err`, so `nil` means no error — fully compatible with the standard `error` interface.

Fields:
- `Value` — the underlying Go error
- `Code` — optional integer error code
- `Msg` — human-readable message
- `Details` — any serializable context (struct, map, etc.)
- `File`, `Line` — call site captured via `runtime.Caller`
- `Timestamp` — microseconds since epoch
- `Prev *Err` — linked list of previous errors (error chain)
- `StackTrace []byte` — captured via `debug.Stack()`

### Design patterns

- Constructors (`New`, `NewSimple`, `FromError`) and `Wrap` return `*Err`. Callers use `if err != nil` as usual.
- `Empty()` returns `nil`; kept for readability but `nil` is equivalent.
- All methods have nil-safe pointer receivers: `IsEmpty()`, `IsError()`, `Error()`, `Is()`, etc.
- `New()` accepts a variadic `skip` parameter to control `runtime.Caller` depth, allowing wrapper functions to report the caller's file/line rather than their own. `NewSimple` and `Wrap` do the same.
- `Clone()` deep-copies the entire `Prev` chain; used internally by `Wrap` and `JSON()` to avoid mutating the receiver.
- `JSON()` returns `([]byte, error)`. `JSONOrEmpty()` silently drops marshaling errors and returns `[]byte{}` on failure. Both work on a clone to avoid side effects.
- `ToError()` returns `e` directly (since `*Err` implements `error`), or `nil` if empty.
- `MarshalJSON` uses a local `Alias` type to avoid infinite recursion while customizing JSON output (converts `Value error` → string, `Timestamp int64` → `time.Time`).
- `Is()` walks the `Prev` chain, not just `e.Value`, so error matching works across the full chain.

### Testing line numbers

Tests that verify `err.File` and `err.Line` use `runtime.Caller(0)` to capture the expected line dynamically:

```go
_, _, wantLine, _ := runtime.Caller(0); wantLine++
err := New(errors.New("test"), ...)
assert.Equal(t, wantLine, err.Line)
```

Tests for the `skip` parameter only assert the filename (`error.go` vs `error_test.go`), not the exact line.
