# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased](https://github.com/rasa/compat/compare/v0.4.1...HEAD)

### Added

- Change `UID()` and `GID()` to return POSIX values on Windows.
- Add `Rename()` function (works atomically on Windows).
- Add `WriteFileAtomic()` function.
- Add `WriteReaderAtomic()` function.
- Add deprecated `ReplaceFile()` (to mimic github.com/natefinch/atomic).

### Fixed

- Fix `Stat()` to always return nil, err with err != nil on Windows.

### Changed

- Change `Chmod()` to set correct ACLs on Windows.
- Change `IsAdmin()` to `IsRoot()` to be more unix-centric.
- Deprecate `IsAdmin()`.
- Move `cmd/main.go` to `cmd/demo/main.go`.
- Move debug logic to only run with `debug` tag is set.
- Various code refactorings.
- Improve code comments.

## [0.4.1](https://github.com/rasa/compat/releases/tag/v0.4.1)

### Added

- Add `Is$GOOS` and `Is$GOARCH` constants (`IsAIX`, `Is386`, etc.).
- Add `BuildBits()` function.
- Add `CPUBits()` function.
- Add `BTime()` support for `Stat()`/`Lstat()` on Linux and Android.
- Add `BTime()` support for `Stat()`/`Lstat()` on FreeBSD and NetBSD.
- Add `CTime()` support for `Stat()`/`Lstat()` on Windows.
- Add `Mode()` support for `Stat()`/`Lstat()` on Windows.
- Add `Chmod()` function.
- Add `Create()`, `CreateEx()`, `CreateTemp()` and `CreateTempEx()` functions.
- Add `Mkdir()`, `MkdirAll()` and `MkdirTemp()` functions.
- Add `OpenFile()`, `WriteFile()` and `WriteFileEx()` functions.
- Add `Umask()` and `GetUmask()` functions.
- Add running tests on all BSD variants, Illumos, Solaris, and JS/Wasm.

## [0.4.0](https://github.com/rasa/compat/releases/tag/v0.4.0)

### Changed

- Rename Device* functions to Partition* to be more user friendly.

## [0.3.0](https://github.com/rasa/compat/releases/tag/v0.3.0)

### Added

- Add `IsAdmin()` and `IsWSL()` functions.

## [0.2.0](https://github.com/rasa/compat/releases/tag/v0.2.0)

### Added

- Add `Nice()` and `Renice()` functions.

## [0.1.0](https://github.com/rasa/compat/releases/tag/v0.1.0)

### Added

- Initial release.

<!-- markdownlint-configure-file
MD024:
  # Only check sibling headings
  siblings_only: true
-->
