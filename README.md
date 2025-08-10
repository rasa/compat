# compat

[![Build](https://github.com/rasa/compat/actions/workflows/build.yml/badge.svg)](https://github.com/rasa/compat/actions/workflows/build.yml)
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

	"github.com/rasa/compat"
)

func main() {
	name := "hello.txt"
	err := compat.WriteFile(name, []byte("Hello World"), 0o654)
	if err != nil {
		log.Fatal(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Name()   =%v\n", fi.Name())
	fmt.Printf("Size()   =%v\n", fi.Size())
	fmt.Printf("Mode()   =0o%o (%v)\n", fi.Mode(), fi.Mode())
	fmt.Printf("ModTime()=%v\n", fi.ModTime())
	fmt.Printf("IsDir()  =%v\n", fi.IsDir())
	fmt.Printf("Sys()    =%+v\n", fi.Sys())
	fmt.Printf("PartID() =%v\n", fi.PartitionID())
	fmt.Printf("FileID() =%v (0x%x)\n", fi.FileID(), fi.FileID())
	fmt.Printf("Links()  =%v\n", fi.Links())
	fmt.Printf("ATime()  =%v\n", fi.ATime())
	fmt.Printf("BTime()  =%v\n", fi.BTime())
	fmt.Printf("CTime()  =%v\n", fi.CTime())
	fmt.Printf("MTime()  =%v\n", fi.MTime())
	fmt.Printf("UID()    =%v (0x%x)\n", fi.UID(), fi.UID())
	fmt.Printf("GID()    =%v (0x%x)\n", fi.GID(), fi.GID())
}

```
which, on Linux, produces:
```text
Name()   =hello.txt
Size()   =11
Mode()   =0o654 (-rw-r-xr--)
ModTime()=2025-08-09 18:46:55.360909223 -0700 PDT
IsDir()  =false
Sys()    =&{Dev:64512 Ino:18756660 Nlink:1 Mode:33204 Uid:1000 Gid:1000 X__pad0:0 Rdev:0 Size:11 Blksize:4096 Blocks:8 Atim:{Sec:1754790298 Nsec:87132300} Mtim:{Sec:1754790415 Nsec:360909223} Ctim:{Sec:1754790415 Nsec:360909223} X__unused:[0 0 0]}
PartID() =64512
FileID() =18756660 (0x11e3434)
Links()  =1
ATime()  =2025-08-09 18:44:58.0871323 -0700 PDT
BTime()  =2025-08-09 10:57:37.073054326 -0700 PDT
CTime()  =2025-08-09 18:46:55.360909223 -0700 PDT
MTime()  =2025-08-09 18:46:55.360909223 -0700 PDT
UID()    =1000 (0x3e8)
GID()    =1000 (0x3e8)
```
and on Windows, produces:
```text
Name()   =hello.txt
Size()   =11
Mode()   =0o654 (-rw-r-xr--)
ModTime()=2025-08-09 18:44:58.0871323 -0700 PDT
IsDir()  =false
Sys()    =&{FileAttributes:32 CreationTime:{LowDateTime:4072684994 HighDateTime:31197526} LastAccessTime:{LowDateTime:1626920091 HighDateTime:31197592} LastWriteTime:{LowDateT
HighDateTime:31197592} FileSizeHigh:0 FileSizeLow:11}
PartID() =8
FileID() =19421773393914848 (0x450000000d6be0)
Links()  =1
ATime()  =2025-08-09 18:44:58.0871323 -0700 PDT
BTime()  =2025-08-09 10:56:35.879469 -0700 PDT
CTime()  =2025-08-09 18:44:58.0871323 -0700 PDT
MTime()  =2025-08-09 18:44:58.0871323 -0700 PDT
UID()    =197609 (0x303e9)
GID()    =197121 (0x30201)
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
† Uses the same logic as in Cygwin/MSYS2 to map Windows string-based SIDs to 64-bit integer UIDs/GIDs.

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
