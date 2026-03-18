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

`Err` is a value type (not a pointer) that wraps `error` with:
- `Value` — the underlying Go error
- `Code` — optional integer error code
- `Msg` — human-readable message
- `Details` — any serializable context (struct, map, etc.)
- `File`, `Line` — call site captured via `runtime.Caller`
- `Timestamp` — microseconds since epoch
- `Prev *Err` — linked list of previous errors (error chain)
- `StackTrace []byte` — captured via `debug.Stack()`

### Design patterns

- Functions return `xerr.Err` by value, not `error`. Callers check `err.IsEmpty()` / `err.IsError()` instead of `err != nil`.
- `Empty()` serves as the "no error" sentinel — returned when there's nothing to report.
- `New()` accepts a variadic `skip` parameter to control `runtime.Caller` depth, allowing wrapper functions to report the caller's file/line rather than their own. `NewSimple` and `Wrap` do the same.
- `Clone()` deep-copies the entire `Prev` chain; used internally by `Wrap` to preserve the chain without mutation.
- `JSON()` returns `([]byte, Err)` — the second return is itself an `Err`. `JSONOrEmpty()` silently drops marshaling errors and returns `[]byte{}` on failure.
- `MarshalJSON` uses a local `Alias` type to avoid infinite recursion while customizing JSON output (converts `Value error` → string, `Timestamp int64` → `time.Time`).
- `Is()` walks the `Prev` chain, not just `e.Value`, so `errors.Is`-style checks work across the full chain.
