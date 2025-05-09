# compat

[![Keep a Changelog](https://img.shields.io/badge/changelog-Keep%20a%20Changelog-%23E05735)](CHANGELOG.md)
[![Go Reference](https://pkg.go.dev/badge/github.com/rasa/compat.svg)](https://pkg.go.dev/github.com/rasa/compat)
[![go.mod](https://img.shields.io/github/go-mod/go-version/rasa/compat)](go.mod)
[![LICENSE](https://img.shields.io/github/license/rasa/compat)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/rasa/compat)](https://goreportcard.com/report/github.com/rasa/compat)
<!-- [![Codecov](https://codecov.io/gh/rasa/compat/branch/main/graph/badge.svg)](https://codecov.io/gh/rasa/compat) -->

# Overview

compat is a cross-platform Go library for filesystem and system compatibility, providing unified access to file metadata, volume attributes, timestamps, device types, and process priority across all major operating systems, include Windows, Linux, macOS, Android, IOS, and many others.

# Usage

The documentation is available at https://pkg.go.dev/github.com/rasa/compat

## Stat

Here's an example of calling `compat.Stat`:

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
	fmt.Printf("Name()    =%v\n", fi.Name())
	fmt.Printf("Size()    =%v\n", fi.Size())
	fmt.Printf("Mode()    =0o%o\n", fi.Mode())
	fmt.Printf("ModTime() =%v\n", fi.ModTime())
	fmt.Printf("IsDir()   =%v\n", fi.IsDir())
	fmt.Printf("Sys()     =%+v\n", fi.Sys())
	fmt.Printf("DeviceID()=%v\n", fi.DeviceID())
	fmt.Printf("FileID()  =%v\n", fi.FileID())
	fmt.Printf("Links()   =%v\n", fi.Links())
	fmt.Printf("ATime()   =%v\n", fi.ATime())
	fmt.Printf("BTime()   =%v\n", fi.BTime())
	fmt.Printf("CTime()   =%v\n", fi.CTime())
	fmt.Printf("MTime()   =%v\n", fi.MTime())
	fmt.Printf("UID()     =%v\n", fi.UID())
	fmt.Printf("GID()     =%v\n", fi.GID())
}
```

which outputs:

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

   go get github.com/rasa/compat

# Contributing

Please feel free to submit issues, fork the repository and send pull requests!

# License

This project is licensed under the terms of the [CC0](https://creativecommons.org/public-domain/cc0/) license.
