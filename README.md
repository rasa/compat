# compat

[![Build](https://github.com/rasa/compat/actions/workflows/build-ubuntu.yml/badge.svg)](https://github.com/rasa/compat/actions/workflows/build-ubuntu.yml)
[![MegaLinter](https://github.com/rasa/compat/actions/workflows/mega-linter.yml/badge.svg)](https://github.com/rasa/compat/actions/workflows/mega-linter.yml)
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
	"os"

	"github.com/rasa/compat"
)

const mode = os.FileMode(0o654)

func main() {
	_ := compat.WriteFile("hello.txt", []byte("Hello World"), mode)

  fi, _ := compat.Stat("hello.txt")

  print(fi.String())
}

```
which, on Linux, produces:
```text
Name:   hello.txt
Size:   11
Mode:   0o654 (-rw-r-xr--)
ModTime:2025-08-14 09:25:27.190602462 -0700 PDT // last Modified
ATime:  2025-08-14 09:25:27.190602462 -0700 PDT // last Accessed
BTime:  2025-08-14 09:25:27.190602462 -0700 PDT // Birthed/created
CTime:  2025-08-14 09:25:27.190602462 -0700 PDT // metadata last Changed
IsDir:  false
Links:  1                                       // number of hard links
UID:    1000 (ross)                             // user ID
GID:    1000 (ross)                             // group ID
PartID: 64512                                   // unique partition (device) ID
FileID: 18756713                                // unique file ID on the partition
```
and on Windows, produces:
```text
Name:   hello.txt
Size:   11
Mode:   0o654 (-rw-r-xr--)
ModTime:2025-08-14 09:28:50.4214934 -0700 PDT
ATime:  2025-08-14 09:28:50.4214934 -0700 PDT
BTime:  2025-08-14 09:28:50.4209614 -0700 PDT
CTime:  2025-08-14 09:28:50.4214934 -0700 PDT
IsDir:  false
Links:  1
UID:    197609 (domain\ross)
GID:    197121 (domain\None)
PartID: 8
FileID: 844424931319952
```
icacls shows:
```
icacls hello.txt

hello.txt domain\ross:(R,W,D)
          domain\None:(RX)
          Everyone:(R)
```
powershell shows:
```
powershell -command "Get-Acl hello.txt | Format-List"

Path   : Microsoft.PowerShell.Core\FileSystem::C:\path\to\hello.txt
Owner  : domain\ross
Group  : domain\None
Access : Everyone Allow  Read, Synchronize
         domain\None Allow  ReadAndExecute, Synchronize
         domain\ross Allow  Write, Delete, Read, Synchronize
Audit  :
Sddl   : O:S-1-5-21-2970224322-3395479738-1485484954-1001G:S-1-5-21-2970224322-3395479738-1485484954-513D:P(A;;FR;;;WD)(A;;0x1200a9;;;S-1-5-21-2970224322-3395479738-1485484954-5
         19f;;;S-1-5-21-2970224322-3395479738-1485484954-1001)
```
Cygwin's stat shows:
```
$ stat hello.txt
  File: hello.txt
  Size: 11              Blocks: 1          IO Block: 65536  regular file
Device: 0,8     Inode: 844424931319952  Links: 1
Access: (0754/-rwxr-xr--)  Uid: (197609/    ross)   Gid: (197121/    None)
Access: 2025-08-14 09:28:50.421493400 -0700
Modify: 2025-08-14 09:28:50.421493400 -0700
Change: 2025-08-14 09:28:50.421493400 -0700
 Birth: 2025-08-14 09:28:50.420961400 -0700
```
Git for Windows's stat shows:
```
$ stat hello.txt
  File: hello.txt
  Size: 11              Blocks: 1          IO Block: 65536  regular file
Device: 8h/8d   Inode: 844424931319952  Links: 1
Access: (0644/-rw-r--r--)  Uid: (197609/    ross)   Gid: (197121/ UNKNOWN)
Access: 2025-08-14 09:30:06.729683700 -0700
Modify: 2025-08-14 09:28:50.421493400 -0700
Change: 2025-08-14 09:28:50.421493400 -0700
 Birth: 2025-08-14 09:28:50.420961400 -0700
```

# Installing

To install compat, use `go get`:

  `go get github.com/rasa/compat`

# FileInfo Functions

The `Stat()` and `Lstat()` functions return a `FileInfo` object.
The table below lists the operating system support for each of the `FileInfo` functions:

| OS           | PartitionID()/ <br/>FileID()* | Links()* | ATime()*<br/>(last<br/>*A*ccessed) | BTime()*<br/>(*B*irthed/<br/>created) | CTime()*<br/>(metadata<br/>last *C*hanged) | UID()/GID() |
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
| Wasip1/<br/>WASM | ✅ | ✅    | ✅†  | ❌    | ✅  | ✅   |
| Windows      | ✅     | ✅    | ✅   | ✅    | ✅   | ✅‡  |
<!--           | PartID+ | Links | ATime | BTime | CTime | UID+ | -->

\* Support will depend on the underlying file system. See [Comparison of file systems](https://wikipedia.org/wiki/Comparison_of_file_systems#Metadata) for details.<br/>
† Not supported if compiled using the Tinygo compiler.<br/>
‡ Provides the same integer values as Cygwin/MSYS2/Git for Windows in mapping Windows SIDs (Security Identifiers).

✅ fully supported.<br/>
☑️ the UID() and GID() values are 32-bit hashes of the user and group names.<br/>
✖️ not implemented (but if the OS supports it, so we could add support).<br/>
❌ not implemented (as it appears the OS doesn't support it).<br/>
<!-- 🚧 planned to be implemented.<br/> -->

# Other Functions

All other functions provided by this library are fully supported by all the above operating systems.

# Contributing

Please feel free to submit issues, fork the repository and send pull requests!

# License

This project is [MIT licensed](LICENSE).
