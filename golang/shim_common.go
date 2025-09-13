// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package golang

import (
	"os"
	_ "unsafe" // for go:linkname
)

const (
	O_RDONLY = os.O_RDONLY // open the file read-only.
	O_WRONLY = os.O_WRONLY // open the file write-only.
	O_RDWR   = os.O_RDWR   // open the file read-write.
	// The remaining values may be or'ed in to control behavior.
	O_APPEND = os.O_APPEND // append data to the file when writing.
	O_CREATE = os.O_CREATE // create a new file if none exists.
	O_EXCL   = os.O_EXCL   // used with O_CREATE, file must not exist.
	O_SYNC   = os.O_SYNC   // open for synchronous I/O.
	O_TRUNC  = os.O_TRUNC  // truncate regular writable file when opened.
)

var (
	Chmod           = os.Chmod
	ErrExist        = os.ErrExist
	IsExist         = os.IsExist
	IsNotExist      = os.IsNotExist
	IsPathSeparator = os.IsPathSeparator
	IsPermission    = os.IsPermission
	Lstat           = os.Lstat
	ModeSetgid      = os.ModeSetgid
	ModeSetuid      = os.ModeSetuid
	ModeSticky      = os.ModeSticky
	NewFile         = os.NewFile
	PathSeparator   = os.PathSeparator
	Stat            = os.Stat
	TempDir         = os.TempDir
)

type (
	File      = os.File
	FileMode  = os.FileMode
	PathError = os.PathError
)

// Source: https://github.com/golang/go/blob/cc8a6780/src/os/tempfile.go#L19-L24

//go:linkname runtime_rand runtime.rand
func runtime_rand() uint64

var PrefixAndSuffix = prefixAndSuffix
