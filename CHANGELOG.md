# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased](https://github.com/rasa/compat/compare/v0.4.0...HEAD)

## Added

- Add `Chmod()`, `Create()`, `CreateTemp()`, `OpenFile()` and `WriteFile()` functions.
- Add `Mkdir()`, `MkdirAll()` and `MkdirTemp()` functions.
- Add `CreateEx(name string, perm os.FileMode, flags int)` function.
- Add `WriteFileEx(name string, data []byte, perm os.FileMode, flags int)` function.
- Add `Umask()` function.

## [0.4.0](https://github.com/rasa/compat/releases/tag/v0.4.0)

## Changed

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
