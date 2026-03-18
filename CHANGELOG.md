# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<!--
## [Unreleased]

## `x.y.z` (YYYY-MM-DD) [CURRENT | YANKED]

### Added (for new features)
### Changed (for changes in existing functionality)
### Deprecated (for soon-to-be removed features)
### Removed (for now removed features)
### Fixed (for any bug fixes)
### Security
-->

## [Unreleased]

### Added

- Add `CLAUDE.md`
- Bump to Go `1.26`
- Bump to testify `v1.11.1`
- Add `TestErr_Is_Nil`, `TestErr_ValueEq_Nil`, `TestErr_Eq_DifferentChainLengths`, `TestErr_Eq_Nil`, `TestErr_Clone_Nil` tests
- Add `TestErr_ImplementsError`, `TestErr_ErrorsIs_Compatibility`, `TestErr_NilIsNilError` tests for standard `error` interface compatibility

### Changed

- [BREAKING] All constructors (`New`, `NewSimple`, `FromError`) and `Wrap` now return `*Err` instead of `Err` — `nil` means no error, compatible with standard `if err != nil` checks
- [BREAKING] `Empty()` now returns `nil` (`*Err`)
- [BREAKING] `JSON()` now returns `([]byte, error)` instead of `([]byte, Err)`
- [BREAKING] `ValueEq()` and `Eq()` now accept `*Err` instead of `Err`
- `ToError()` now returns `e` directly (since `*Err` implements `error`) instead of wrapping the message in a new `errors.New`
- `JSON()` and `JSONOrEmpty()` now operate on a clone to avoid mutating the receiver's `StackTrace` field
- All methods are nil-safe (nil pointer receiver handled explicitly)
- `Eq()`: fix potential nil pointer dereference when chains have different lengths

### Fixed

- `Eq()`: no longer panics when comparing chains of different lengths

## `0.6.0` (2025-07-07) [CURRENT]

### Added

- Add `JSONOrEmpty()` method to avoid returning an error like the `JSON()` method.

## `0.5.0` (2025-06-27)

### Changed

- [Breaking] Apply good practice: Struct Err has methods on both value and pointer receivers. All methods now require a pointer receiver.

## `0.4.0` (2025-06-27)

### Added

- Add a new optional parameter in `New`, `NewSimple` and `Wrap` to change the value of `runtime.Caller()`

## `0.3.0` (2025-05-22)

### Added

- Add `Code` field to `Err`
- Add `Unwrap()` and `ToError()` method to `Err`

### Change

- [BREAKING] Change `NewErr` function to add the error code to `New`

## `0.2.0` (2025-05-06)

### Added

- Add methods Clone, ValueEq and Eq to Err

## `0.1.0` (2025-05-06)

### Added

- First version
