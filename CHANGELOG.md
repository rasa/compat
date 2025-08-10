# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased](https://github.com/rasa/compat/compare/v0.4.4...HEAD)

### Added

### Fixed

### Changed

## [0.4.4](https://github.com/rasa/compat/compare/v0.4.3...v0.4.4)

### Added

- Set perms when `WriteFileAtomic()` and `WriteReaderAtomic()` create file, not after closing.
- Add `Option` param to `Create()` and `CreateTemp()`.
- Add `Flag()` to `FileOptions` functions.
- Add Windows example `Stat()` call to readme.
- Add `icacls` and powershell's `Get-Acl` output to readme.

### Fixed

- Fix ACL rights by using GetTokenInformation(.., TokenUser, ..) on Windows.
- Fix `Create()` failure due to missing flag.
- If a mode of `0` is passed to a function, use the function's default mode.
- Fix tests by using go's standard test framework on js/wasm.

### Changed

- Run more tests against many `FileMode` values, not just one.
- Rework demo to create hello.txt, instead of using .exe.
- Rework test framework's setting the expected `Mode()` result for selected OSes.

## [0.4.3](https://github.com/rasa/compat/compare/v0.4.2...v0.4.3)

### Added

- Run tests on tinygo.
- Run tests for wasip1 on both regular go, and tinygo.
- Add `IsAct`, `IsApple`, `IsBSD`, `IsSolaria`, `IsTinygo`, and `IsUnix` constants.
- Modularized Github actions.

### Fixed

- Fix all tests on wasip1.
- Fix `CPUBits()` to return the value returned by `BuildBits()` if we can't determine the CPU's bits.
- Enhance and fix typos in comments.
- Fix tests when running under `act`.
- Fix intermittent test failures on Windows 2025.
- Fix Stat() by re-adding `CTime()` support on wasip1 (using regular go).
- Fix Stat() by removing `ATime()` and `CTime()` support on wasip1 (using tinygo).

### Changed

- Change `Stat()`'s `Mode()` to return `0o600` for files and `0o700` for directories, on wasip1.
- Change `Stat()`'s `UID()` to return the value `os.Getuid()` returns (1), on wasip1.
- Change `Stat()`'s `GID()` to return the value `os.Getgid()` returns (1), on wasip1.
- Bump tinygo to 0.38.0.

## [0.4.2](https://github.com/rasa/compat/compare/v0.4.1...v0.4.2)

### Added

- Change `UID()` and `GID()` to return POSIX values on Windows.
- Add `Rename()` function (works atomically on Windows).
- Add `WriteFileAtomic()` function.
- Add `WriteReaderAtomic()` function.
- Add deprecated `ReplaceFile()` (to mimic github.com/natefinch/atomic).

### Fixed

- Fix `Stat()` to always return nil when err != nil on Windows.
- Fix `Stat()` by removing `CTime()` support on wasip1.

### Changed

- Change `Chmod()` to set correct ACLs on Windows.
- Change `IsAdmin()` to `IsRoot()` to be more unix-centric.
- Deprecate `IsAdmin()`.
- Move `cmd/main.go` to `cmd/demo/main.go`.
- Move debug logic to only run with `debug` tag is set.
- Various code refactorings.
- Improve code comments.

## [0.4.1](https://github.com/rasa/compat/compare/v0.4.0...v0.4.1)

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

## [0.4.0](https://github.com/rasa/compat/compare/v0.3.0...v0.4.0)

### Changed

- Rename Device* functions to Partition* to be more user friendly.

## [0.3.0](https://github.com/rasa/compat/compare/v0.2.0...v0.3.0)

### Added

- Add `IsAdmin()` and `IsWSL()` functions.

## [0.2.0](https://github.com/rasa/compat/compare/v0.1.0...v0.2.0)

### Added

- Add `Nice()` and `Renice()` functions.

## [0.1.0](https://github.com/rasa/compat/compare/bcf970117c696f70992faaa061148a206a3c4b8e...v0.1.0)

### Added

- Initial release.

<!-- markdownlint-configure-file
MD024:
  # Only check sibling headings
  siblings_only: true
-->
