// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"io/fs"
	"os"
)

const (
	// CreatePerm is the FileMode used by Create().
	CreatePerm os.FileMode = 0o666
	// CreateTempPerm is the FileMode used by CreateTemp().
	CreateTempPerm os.FileMode = 0o600
	// MkdirTempPerm is the FileMode used by MkdirTemp()..
	MkdirTempPerm os.FileMode = 0o700

	// Verify we don't conflict with any of the values listed at
	//  https://github.com/golang/go/blob/77f911e3/src/syscall/types_windows.go#L37-L55

	// O_DELETE deletes the file when closed.
	O_DELETE = 0x8000000

	// https://github.com/golang/go/blob/e282cbb1/src/os/file.go#L77

	// The following constants are not used by the compat library, but are
	// provided to make code migration easier.

	O_RDONLY = os.O_RDONLY // open the file read-only. //nolint:revive // quiet linter
	O_WRONLY = os.O_WRONLY // open the file write-only. //nolint:revive // quiet linter
	O_RDWR   = os.O_RDWR   // open the file read-write. //nolint:revive // quiet linter
	// The remaining values may be or'ed in to control behavior.
	O_APPEND = os.O_APPEND // append data to the file when writing.
	O_CREATE = os.O_CREATE // create a new file if none exists.
	O_EXCL   = os.O_EXCL   // used with O_CREATE, file must not exist.
	O_SYNC   = os.O_SYNC   // open for synchronous I/O.
	O_TRUNC  = os.O_TRUNC  // truncate regular writable file when opened.

	ModeDir        = fs.ModeDir        // d: is a directory
	ModeAppend     = fs.ModeAppend     // a: append-only
	ModeExclusive  = fs.ModeExclusive  // l: exclusive use
	ModeTemporary  = fs.ModeTemporary  // T: temporary file; Plan 9 only
	ModeSymlink    = fs.ModeSymlink    // L: symbolic link
	ModeDevice     = fs.ModeDevice     // D: device file
	ModeNamedPipe  = fs.ModeNamedPipe  // p: named pipe (FIFO)
	ModeSocket     = fs.ModeSocket     // S: Unix domain socket
	ModeSetuid     = fs.ModeSetuid     // u: setuid
	ModeSetgid     = fs.ModeSetgid     // g: setgid
	ModeCharDevice = fs.ModeCharDevice // c: Unix character device, when ModeDevice is set
	ModeSticky     = fs.ModeSticky     // t: sticky
	ModeIrregular  = fs.ModeIrregular  // ?: non-regular file; nothing else is known about this file

	// ModeType is a mask for the type bits. For regular files, none will be set.
	ModeType = fs.ModeType

	// ModePerm is a mask for the Unix permission bits, 0o777.
	ModePerm = fs.ModePerm
)

type FileMode = os.FileMode
