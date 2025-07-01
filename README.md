# compat

[![Build](https://github.com/rasa/compat/actions/workflows/build.yml/badge.svg)](https://github.com/rasa/compat/actions/workflows/build.yml)
[![Last Commit](https://img.shields.io/github/last-commit/rasa/compat.svg)](https://github.com/rasa/compat/commits)
[![Codecov](https://codecov.io/gh/rasa/compat/branch/main/graph/badge.svg)](https://codecov.io/gh/rasa/compat)
[![Release](https://img.shields.io/github/v/release/rasa/compat.svg?style=flat)](https://github.com/rasa/compat/releases)
[![go.mod](https://img.shields.io/github/go-mod/go-version/rasa/compat)](go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/rasa/compat)](https://goreportcard.com/report/github.com/rasa/compat)
[![Go Reference](https://pkg.go.dev/badge/github.com/rasa/compat.svg)](https://pkg.go.dev/github.com/rasa/compat)
[![Known Vulnerabilities](https://snyk.io/test/github/rasa/compat/badge.svg)](https://snyk.io/test/github/rasa/compat)
[![Keep a Changelog](https://img.shields.io/badge/changelog-Keep%20a%20Changelog-%23E05735)](CHANGELOG.md)
[![LICENSE](https://img.shields.io/github/license/rasa/compat)](LICENSE)

<!--ts-->
* [Overview](#overview)
* [Usage](#usage)
* [Stat](#stat)
* [Installing](#installing)
* [FileInfo Functions](#fileinfo-functions)
* [Other Functions](#other-functions)
* [Contributing](#contributing)
* [License](#license)
<!--te-->

# Overview

compat is a pure-Go library providing unified access to file and device metadata, atomic file operations, process priority, etc. on all major operating systems, include Andoid, iOS, Linux, macOS, Windows, and many others.

# Usage

The documentation is available at https://pkg.go.dev/github.com/rasa/compat

## Stat

Here's an example of calling `compat.Stat()`:

```go
package main

import (
  "fmt"
  "github.com/rasa/compat"
)

func main() {
  fi, err := compat.Stat(os.Executable())
  if err != nil {
    fmt.Println(err)
    return
  }
        // Same functions as os.Stat() and os.Lstat():
  fmt.Printf("Name()    =%v\n", fi.Name())     // base name of the file
  fmt.Printf("Size()    =%v\n", fi.Size())     // length in bytes
  fmt.Printf("Mode()    =0o%o\n", fi.Mode())   // file mode bits
  fmt.Printf("ModTime() =%v\n", fi.ModTime())  // last modified
  fmt.Printf("IsDir()   =%v\n", fi.IsDir())    // is a directory
  fmt.Printf("Sys()     =%+v\n", fi.Sys())     // underlying data source
        // New functions provided by this compat library:
  fmt.Printf("PartID()  =%v\n", fi.PartitionID()) // partition (device) ID
  fmt.Printf("FileID()  =%v\n", fi.FileID())   // file (inode) ID
  fmt.Printf("Links()   =%v\n", fi.Links())    // number of hard links
  fmt.Printf("ATime()   =%v\n", fi.ATime())    // last accessed
  fmt.Printf("BTime()   =%v\n", fi.BTime())    // created (birthed)
  fmt.Printf("CTime()   =%v\n", fi.CTime())    // status/metadata changed
  fmt.Printf("MTime()   =%v\n", fi.MTime())    // alias for ModTime
  fmt.Printf("UID()     =%v\n", fi.UID())      // user ID
  fmt.Printf("GID()     =%v\n", fi.GID())      // group ID
}
```

which, on Linux, produced:

```text
Name()    =cmd
Size()    =1597624
Mode()    =0o775
ModTime() =2025-05-08 22:11:01.353744514 -0700 PDT
IsDir()   =false
Sys()     =&{Dev:64512 Ino:56893266 Nlink:1 Mode:33277 Uid:1000 Gid:1000 X__pad0:0 Rdev:0 Size:1597624 Blksize:4096 Blocks:3128 Atim:{Sec:1746767461 Nsec:354744521} Mtim:{Sec:1746767461 Nsec:353744514} Ctim:{Sec:1746767461 Nsec:353744514} X__unused:[0 0 0]}
PartID()  =64512
FileID()  =56893266
Links()   =1
ATime()   =2025-05-08 22:11:01.354744521 -0700 PDT
BTime()   =0001-01-01 00:00:00 +0000 UTC
CTime()   =2025-05-08 22:11:01.353744514 -0700 PDT
MTime()   =2025-05-08 22:11:01.353744514 -0700 PDT
UID()     =1000
GID()     =1000
```

# Installing

To install compat, use `go get`:

  `go get github.com/rasa/compat`

# FileInfo Functions

The `Stat()` and `Lstat()` functions return a `FileInfo` object.
The table below lists the operating system support for each of the `FileInfo` functions:

| OS      | PartitionID()/ <br/>FileID()* | Links()* | ATime()* | BTime()* | CTime()* | UID()* / <br/>GID()* |
|---------|--------|--------|------|--------|------|-------|
| AIX     | ✅     | ✅    | ✅   | ❌    | ✅   | ✅   |
| Android | ✅     | ✅    | ✅   | ✅    | ✅   | ✅   |
| Dragonfly | ✅   | ✅    | ✅   | ✖️    | ✅   | ✅   |
| FreeBSD | ✅     | ✅    | ✅   | ✅    | ✅   | ✅   |
| Illumos | ✅     | ✅    | ✅   | ✖️    | ✅   | ✅   |
| iOS     | ✅     | ✅    | ✅   | ✅    | ✅   | ✅   |
| Linux   | ✅     | ✅    | ✅   | ✅    | ✅   | ✅   |
| macOS   | ✅     | ✅    | ✅   | ✅    | ✅   | ✅   |
| NetBSD  | ✅     | ✅    | ✅   | ✅    | ✅   | ✅   |
| OpenBSD | ✅     | ✅    | ✅   | ✖️    | ✅   | ✅   |
| Plan9   | ✅     | ❌    | ✅   | ❌    | ❌   | ☑️   |
| Solaris | ✅     | ✅    | ✅   | ✖️    | ✅   | ✅   |
| WebAssembly | ✅ | ✅    | ✅   | ❌    | ✅   | ✅   |
| Windows | ✅     | ✅    | ✅   | ✅    | ✅   | ✅†  |
<!--      | PartID+ | Links | ATime | BTime | CTime | UID+ | -->

\* Support will depend on the underlying file system. See [Comparison of file systems](https://wikipedia.org/wiki/Comparison_of_file_systems#Metadata) for details.
† Uses the same logic as in Cygwin/MSYS2 to map Windows SIDs to UIDs/GIDs.

✅ fully supported.<br/>
☑️ the UID() and GID() values are 64-bit hashes of the user and group names.<br/>
✖️ not implemented (but it appears the OS supports it, so we could add support).<br/>
❌ not implemented (as it appears the OS doesn't support it).<br/>
<!-- 🚧 planned to be implemented.<br/> -->

# Other Functions

All other functions provided by this library are fully supported by all the above operating systems.

# Contributing

Please feel free to submit issues, fork the repository and send pull requests!

# License

This project is [MIT licensed](LICENSE).
