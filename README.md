# compat

[![Keep a Changelog](https://img.shields.io/badge/changelog-Keep%20a%20Changelog-%23E05735)](CHANGELOG.md)
[![Go Reference](https://pkg.go.dev/badge/github.com/rasa/compat.svg)](https://pkg.go.dev/github.com/rasa/compat)
[![go.mod](https://img.shields.io/github/go-mod/go-version/rasa/compat)](go.mod)
[![LICENSE](https://img.shields.io/github/license/rasa/compat)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/rasa/compat)](https://goreportcard.com/report/github.com/rasa/compat)
<!-- [![Codecov](https://codecov.io/gh/rasa/compat/branch/main/graph/badge.svg)](https://codecov.io/gh/rasa/compat) -->

# Overview

compat is a pure-Go library providing unified access to file and device metadata, atomic file operations, process priority, etc. on all major operating systems, include Windows, Linux, macOS, Android, iOS, and many others.

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
	fmt.Printf("ModTime() =%v\n", fi.ModTime())  // last modified time
	fmt.Printf("IsDir()   =%v\n", fi.IsDir())    // is a directory
	fmt.Printf("Sys()     =%+v\n", fi.Sys())     // underlying data source
        // New functions provided by this compat library:
	fmt.Printf("DeviceID()=%v\n", fi.DeviceID()) // device ID
	fmt.Printf("FileID()  =%v\n", fi.FileID())   // inode/file ID
	fmt.Printf("Links()   =%v\n", fi.Links())    // number of hard links
	fmt.Printf("ATime()   =%v\n", fi.ATime())    // last access time
	fmt.Printf("BTime()   =%v\n", fi.BTime())    // birth/created time
	fmt.Printf("CTime()   =%v\n", fi.CTime())    // metadata changed time
	fmt.Printf("MTime()   =%v\n", fi.MTime())    // alias for ModTime
	fmt.Printf("UID()     =%v\n", fi.UID())      // user ID
	fmt.Printf("GID()     =%v\n", fi.GID())      // group ID
}
```

which outputed on Linux:

```text
Name()    =cmd
Size()    =1597624
Mode()    =0o775
ModTime() =2025-05-08 22:11:01.353744514 -0700 PDT
IsDir()   =false
Sys()     =&{Dev:64512 Ino:56893266 Nlink:1 Mode:33277 Uid:1000 Gid:1000 X__pad0:0 Rdev:0 Size:1597624 Blksize:4096 Blocks:3128 Atim:{Sec:1746767461 Nsec:354744521} Mtim:{Sec:1746767461 Nsec:353744514} Ctim:{Sec:1746767461 Nsec:353744514} X__unused:[0 0 0]}
DeviceID()=64512
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

| OS      | DeviceID()    | FileID()* | Links()* | ATime()* | BTime()* | CTime()* | UID()* | GID()* |
|---------|---------------|----------|----------|----------|----------|----------|--------|--------|
| AIX     | ✅	          | ✅	     | ✅	     | ✅	     | ❌      | ✅      | ✅    |  ✅  |
| Android | ✅	          | ✅	     | ✅	     | ✅	     | ❌      | ✅      | ✅    |  ✅  |
| Darwin<br/>(macOS) | ✅ | ✅	     | ✅	     | ✅	     | ✅      | ✅      | ✅    |  ✅  |
| Dragonfly | ✅	       | ✅	     | ✅	     | ✅	     | ❌      | ✅      | ✅    |  ✅  |
| FreeBSD | ✅	          | ✅	     | ✅	     | ✅	     | ❌      | ✅      | ✅    |  ✅  |
| Illumos | ✅	          | ✅	     | ✅	     | ✅	     | ❌      | ✅      | ✅    |  ✅  |
| iOS     | ✅	          | ✅	     | ✅	     | ✅	     | ✅	     | ✅      | ✅    |  ✅  |
| Linux   | ✅	          | ✅	     | ✅	     | ✅	     | ❌      | ✅      | ✅    |  ✅  |
| NetBSD  | ✅	          | ✅	     | ✅	     | ✅	     | ❌      | ✅      | ✅    |  ✅  |
| OpenBSD | ✅	          | ✅	     | ✅	     | ✅	     | ❌      | ✅      | ✅    |  ✅  |
| Plan9   | ✅	          | ✅	     | ❌	     | ✅	     | ❌      | ❌      | 🟠    |  🟠  |
| Solaris | ✅	          | ✅	     | ✅	     | ✅	     | ❌      | ✅      | ✅    |  ✅  |
| WebAssembly<br/>(Js) | ✅	    | ✅	     | ✅	     | ✅	     | ❌      | ✅      | ✅    |  ✅  |
| WebAssembly<br/>(WAPI) | ✅	 | ✅	     | ✅	     | ✅	     | ❌      | ✅      | ✅    |  ✅  |
| Windows | ✅	          | ✅	     | ✅      | ✅ 	  | ✅      | ❌      | 🚧    |  🚧  |

* May not be supported on older filesystems, such as FAT32.

✅ fully supported.<br/>
❌ not implemented (though support could be added if the OS provides the information).<br/>
🟠 the UID() and GID() values are 64-bit hashes of the user and group names.<br/>
🚧 planned to be implemented.

# Other Functions

All other functions provided by this library are fully supported by all the above operating systems.

# Contributing

Please feel free to submit issues, fork the repository and send pull requests!

# License

This project is licensed under the terms of the [CC0](https://creativecommons.org/public-domain/cc0/) license.
