# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased](https://github.com/rasa/compat/compare/v0.5.3...HEAD)

### Added

- Add `Fchmod()` function.
- Add `Link()` function.
- Add `SupportsNice()` function.
- Add `WithRetrySeconds()` function.
- Add optional `Option` param to `Rename()`.
- Add optional `Option` param to `RemoveAll()`.

### Fixed

- Increase code coverage 29% to 83%.
- Implement `Nice()` on Plan9 OS.

### Changed

- Rename `O_DELETE` constant to `O_FILE_FLAG_DELETE_ON_CLOSE` to align with constant in golang's standard library.
- Rename `O_NOROATTR` constant to `O_FILE_FLAG_NO_RO_ATTR`.
- Deprecate `O_DELETE` and `O_NOROATTR` constants.
- Freshen cloned code from golang's standard library.
- Always return `os.PathError` or `os.LinkError`, as appropriate.
- Rework `Nice()` to return `0` and no error, when not supported by the OS.

## [v0.5.3](https://github.com/rasa/compat/compare/v0.5.2...v0.5.3)

### Added

- Add support for running tests on APFS, ExFAT, FAT32, HFS+, HFS+J, HFSX, JHFS+,
  JHFS+X, and UDF filesystems on macOS (as root).
- Add support for running tests on ext2, ext3, ext4, exFAT, FAT32, F2FS, NTFS,
  ReiserFS and XFS filesystems on Linux (as root).
- Add support for running tests on exFAT, FAT32, NTFS, and ReFS filesystems on
  Windows. Not all systems support ReFS.
- Add `Geteuid()` and `Getegid()` functions.
- Add `IsBSDLike` constant (`IsBSD` or `IsApple`).
- Add `UnderlyingGoVersion()` to report go version under Tinygo compiler.
- Add `cmd/updater` to build binaries.

### Fixed

- Fix `IsRoot()` to check EUID == 0 on non-Windows systems.
- Fix `SupportsSymlinks()` to return false on the Plan 9 OS.
- Add -buildvcs=false to fix build on DragonflyBSD.

### Changed

- (Temporarily) removed support for mode 0o000 on Windows. ***BREAKING CHANGE***
- Change `Links()` from a uint64 to a uint (to align with APIs). ***BREAKING CHANGE***
- Bump Tinygo to 0.39.0 (uses go1.25).
- Bump Wasitime to 36.0.1.
- Bump Android NDK to r28c.
- Reorder field order of `Stat()`'s `String()` function.

## [v0.5.2](https://github.com/rasa/compat/compare/v0.5.1...v0.5.2)

### Added

- Add optional `Option` param to `OpenFile()`.
- Add optional `Option` param to `Chmod()`.
- Add `WithReadOnlyMode()` function.
- Add `ReadOnlyModeIgnore`, `ReadOnlyModeSet` and `ReadOnlyModeReset` constants.
- Add optional `Option` param to `CreateTemp()`.
- Add optional `Option` param to `MkdirTemp()`.
- Add `Symlink()` function.
- Add `Remove()` and `RemoveAll()` functions.
- Add `WalkDir()` function (mimics `io/fs`' version).
- Add `ReadDir()` and `FormatDirEntry()` functions.
- Add cmd/updater program to keep upstream code up to date.

### Fixed

- Fix ability to set perms to u+w after setting o-w on Windows.
- Enable `TestLstatUser` test on Windows.
- Enable `TestLstatGroup` test on Windows.

### Changed

- Rename `DefaultFileMode()` to  `WithDefaultFileMode()`.
- Rename `Flag()` to  `WithFlags()`.
- Rename `KeepFileMode()` to  `WithKeepFileMode()`.
- Renamed `FileMode()` to  `WithFileMode()`. ***BREAKING CHANGE***
- Rename `FileOptions` to `Options`.
- Deprecate `DefaultFileMode()`, `Flag()`, `KeepFileMode()`, and `UseFileMode()`
  functions.
- Deprecate `CreateEx()`, `CreateTempEx()` and `WriteFileEx()` functions.

## [0.5.1](https://github.com/rasa/compat/compare/v0.5.0...v0.5.1)

### Added

- Add `SupportsLinks()`, `SupportsATime()`, `SupportsBTime()`, `SupportsCTime()`
  and `SupportsSymlinks()` functions.
- Add `String()` function to `FileInfo` interface.
- Add missing `Lstat()` tests.
- Add `UserIDSourceIsSID` return value for `UserIDSource()` function under Windows.

### Fixed

- Fix to only follow symlinks via `Stat()` not `Lstat()`.

### Changed

- Deprecate `Supported()` function.
- Deprecate `Links`, `ATime`, `BTime`, `CTime`, `UID` and `GID` constants.
- Simplify `Stat()` call example in cmd/demo.

## [0.5.0](https://github.com/rasa/compat/compare/v0.4.4...v0.5.0)

### Added

- Add `User()` function to `Stat()`s `FileInfo` interface, to return the file's
  user's name.
- Add `Group()` function to `Stat()`s `FileInfo` interface to return the file's
  group name.
- Add `Error()` function to `Stat()`s `FileInfo` interface to return the last
  error received when calling `BTime()`, `CTime()`, `UID()`, `GID()`, `User()`,
  or `Group()` if an additional system call is executed.
- Add `UserIDSource()` function to return if the `UID()` function or the
  `User()` function is the user's actual ID in the OS, or if the OS doesn't use
  user IDs.
- Add `PartitionType()` function.
- Add `Getuid()` and `Getgid()` functions.
- Add `IsX86CPU`, `IsArmCPU`, `IsMipsCPU`, and `IsPpcCPU` constants.

### Fixed

### Changed

- Change `Stat()`'s `UID()` and `GID()` to return int values, to be more
  aligned with linux/unix. ***BREAKING CHANGE***
- Speed up `Stat()` calls by deferring additional system calls to execute the
  first time `BTime()`, `CTime()`, `UID()`, `GID()`, `User()`, or `Group()` is
  called.

## [0.4.4](https://github.com/rasa/compat/compare/v0.4.3...v0.4.4)

### Added

- Set perms when `WriteFileAtomic()` and `WriteReaderAtomic()` create files, not after closing.
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
