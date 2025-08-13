# compat

[![Build](https://github.com/rasa/compat/actions/workflows/build-ubuntu.yml/badge.svg)](https://github.com/rasa/compat/actions/workflows/build-ubuntu.yml)
[![Last Commit](https://img.shields.io/github/last-commit/rasa/compat.svg)](https://github.com/rasa/compat/commits)
[![Codecov](https://codecov.io/gh/rasa/compat/branch/main/graph/badge.svg)](https://codecov.io/gh/rasa/compat)
[![Release](https://img.shields.io/github/v/release/rasa/compat.svg?style=flat)](https://github.com/rasa/compat/releases)
[![go.mod](https://img.shields.io/github/go-mod/go-version/rasa/compat)](go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/rasa/compat)](https://goreportcard.com/report/github.com/rasa/compat)
[![Go Reference](https://pkg.go.dev/badge/github.com/rasa/compat.svg)](https://pkg.go.dev/github.com/rasa/compat)
<!-- @synk: The badge feature is no longer actively being maintained or developed.
[![Known Vulnerabilities](https://snyk.io/test/github/rasa/compat/badge.svg)](https://snyk.io/test/github/rasa/compat)
-->
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

compat is a pure-Go library providing unified access to file and device metadata, atomic file operations, process priority, etc. on all major operating systems, include Android, iOS, Linux, macOS, Windows, and many others.

# Usage

The documentation is available at https://pkg.go.dev/github.com/rasa/compat

## Stat

Here's an example of calling `compat.Stat()`:

```go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rasa/compat"
)

const mode = os.FileMode(0o654) // Something other than the default.

func main() {
	name := "hello.txt"
	err := compat.WriteFile(name, []byte("Hello World"), mode)
	if err != nil {
		log.Fatal(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		log.Fatal(err)
	}

	print(fi.String())
}

```
which, on Linux, produces:
```text
Name()   =hello.txt
Size()   =11
Mode()   =0o654 (-rw-r-xr--)
ModTime()=2025-08-12 22:46:17.404309557 -0700 PDT
IsDir()  =false
PartID() =64512
FileID() =18756660
Links()  =1
ATime()  =2025-08-12 09:24:50.0851287 -0700 PDT
BTime()  =2025-08-09 19:29:29.108563883 -0700 PDT
CTime()  =2025-08-12 22:46:17.404309557 -0700 PDT
UID()    =1000
GID()    =1000
```
and on Windows, produces:
```text
Name()   =hello.txt
Size()   =11
Mode()   =0o654 (-rw-r-xr--)
ModTime()=2025-08-12 22:44:50.4214598 -0700 PDT
IsDir()  =false
PartID() =8
FileID() =10414574139156045
Links()  =1
ATime()  =2025-08-12 22:44:50.4214598 -0700 PDT
BTime()  =2025-08-09 10:56:35.879469 -0700 PDT
CTime()  =2025-08-12 22:44:50.4214598 -0700 PDT
UID()    =197609
GID()    =197121
```
with icacls showing:
```
icacls hello.txt

hello.txt computername\ross:(R,W,D)
          computername\None:(RX)
          Everyone:(R)
```
and powershell showing:
```
powershell -command "Get-Acl hello.txt | Format-List"

Path   : Microsoft.PowerShell.Core\FileSystem::C:\path\to\hello.txt
Owner  : computername\ross
Group  : computername\None
Access : Everyone Allow  Read, Synchronize
         computername\None Allow  ReadAndExecute, Synchronize
         computername\ross Allow  Write, Delete, Read, Synchronize
Audit  :
Sddl   : O:S-1-5-21-2970224322-3395479738-1485484954-1001G:S-1-5-21-2970224322-3395479738-1485484954-513D:P(A;;FR;;;WD)(A;;0x1200a9;;;S-1-5-21-2970224322-3395479738-1485484954-5
         19f;;;S-1-5-21-2970224322-3395479738-1485484954-1001)
```

# Installing

To install compat, use `go get`:

  `go get github.com/rasa/compat`

# FileInfo Functions

The `Stat()` and `Lstat()` functions return a `FileInfo` object.
The table below lists the operating system support for each of the `FileInfo` functions:

| OS           | PartitionID()/ <br/>FileID()* | Links()* | ATime()* | BTime()* | CTime()* | UID()* / <br/>GID()* |
|--------------|--------|--------|------|--------|------|-------|
| AIX          | ✅     | ✅    | ✅   | ❌    | ✅   | ✅   |
| Android      | ✅     | ✅    | ✅   | ✅    | ✅   | ✅   |
| Dragonfly    | ✅     | ✅    | ✅   | ✖️    | ✅   | ✅   |
| FreeBSD      | ✅     | ✅    | ✅   | ✅    | ✅   | ✅   |
| Illumos      | ✅     | ✅    | ✅   | ✖️    | ✅   | ✅   |
| iOS          | ✅     | ✅    | ✅   | ✅    | ✅   | ✅   |
| Js/<br/>WASM | ✅     | ✅    | ✅   | ❌    | ✅   | ✅   |
| Linux        | ✅     | ✅    | ✅   | ✅    | ✅   | ✅   |
| macOS        | ✅     | ✅    | ✅   | ✅    | ✅   | ✅   |
| NetBSD       | ✅     | ✅    | ✅   | ✅    | ✅   | ✅   |
| OpenBSD      | ✅     | ✅    | ✅   | ✖️    | ✅   | ✅   |
| Plan9        | ✅     | ❌    | ✅   | ❌    | ❌   | ☑️   |
| Solaris      | ✅     | ✅    | ✅   | ✖️    | ✅   | ✅   |
| Wasip1/<br/>WASM | ✅ | ✅    | ✅   | ❌    | ✅  | ✅   |
| Windows      | ✅     | ✅    | ✅   | ✅    | ✅   | ✅†  |
<!--           | PartID+ | Links | ATime | BTime | CTime | UID+ | -->

\* Support will depend on the underlying file system. See [Comparison of file systems](https://wikipedia.org/wiki/Comparison_of_file_systems#Metadata) for details.<br/>
† Uses the same logic as in Cygwin/MSYS2 to map Windows string-based SIDs to 32-bit integer UIDs/GIDs.

✅ fully supported.<br/>
☑️ the UID() and GID() values are 32-bit hashes of the user and group names.<br/>
✖️ not implemented (but it appears the OS supports it, so we could add support).<br/>
❌ not implemented (as it appears the OS doesn't support it).<br/>
<!-- 🚧 planned to be implemented.<br/> -->

# Other Functions

All other functions provided by this library are fully supported by all the above operating systems.

# Contributing

Please feel free to submit issues, fork the repository and send pull requests!

# License

This project is [MIT licensed](LICENSE).
