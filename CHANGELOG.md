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

### Changed

- [Breaking] Apply good practice: Struct Err has methods on both value and pointer receivers. All methods now require a pointer receiver.

## `0.4.0` (2025-06-27) [CURRENT]

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
