// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"io/fs"
	"os"
)

const (
	// CreatePerm is the FileMode used by os.Create() (and compat.Create()).
	CreatePerm os.FileMode = 0o666
	// CreateTempPerm is the FileMode used by os.CreateTemp() (and
	// compat.CreateTemp()).
	CreateTempPerm os.FileMode = 0o600
	// MkdirTempPerm is the FileMode used by os.MkdirTemp() (and
	// compat.MkdirTemp()).
	MkdirTempPerm os.FileMode = 0o700

	// DefaultAppleDirPerm is the FileMode returned for directories by
	// golang's os.Stat() function on Apple based systems
	// when the directory is on a filesystem that doesn't support
	// macOS/iOS permissions, such as exFAT, or FAT32.
	DefaultAppleDirPerm os.FileMode = 0o700
	// DefaultAppleFilePerm is the FileMode returned for files by
	// golang's os.Stat() function on Apple based systems
	// when the file is on a filesystem that doesn't support
	// macOS/iOS permissions, such as exFAT, or FAT32.
	DefaultAppleFilePerm os.FileMode = 0o700

	// DefaultUnixDirPerm is the FileMode returned for directories by
	// golang's os.Stat() function on non-Apple/non-Windows based systems
	// when the directory is on a filesystem that doesn't support
	// Unix permissions, such as exFAT, or FAT32.
	DefaultUnixDirPerm os.FileMode = 0o777
	// DefaultUnixFilePerm is the FileMode returned for files by
	// golang's os.Stat() function on non-Apple/non-Windows based systems
	// when the file is on a filesystem that doesn't support
	// Unix permissions, such as exFAT, or FAT32.
	DefaultUnixFilePerm os.FileMode = 0o777

	// DefaultWindowsDirPerm is the FileMode returned for directories by
	// golang's os.Stat() function on Windows based systems
	// when the directory is on a filesystem that doesn't support Windows'
	// Access Control Lists (ACLS), such as exFAT, or FAT32.
	DefaultWindowsDirPerm os.FileMode = 0o777
	// DefaultWindowsFilePerm is the FileMode returned for files by
	// golang's os.Stat() function on Windows based systems
	// when the file is on a filesystem that doesn't support Windows'
	// Access Control Lists (ACLS), such as exFAT, or FAT32.
	DefaultWindowsFilePerm os.FileMode = 0o666

	// Verify we don't conflict with any of the values listed at
	// https://github.com/golang/go/blob/ac803b59/src/syscall/types_windows.go#L37-L55

	// O_FILE_FLAG_DELETE_ON_CLOSE deletes the file when closed.
	O_FILE_FLAG_DELETE_ON_CLOSE = 0x04000000
	// O_FILE_FLAG_NO_RO_ATTR skips setting a file's read-only attribute on Windows.
	O_FILE_FLAG_NO_RO_ATTR = 0x00010000

	// The following constants are not used by the compat library, but are
	// provided to make code migration easier.

	// Source: https://github.com/golang/go/blob/ac803b59/src/os/file.go#L81-L89
	O_RDONLY = os.O_RDONLY // open the file read-only. //nolint:revive
	O_WRONLY = os.O_WRONLY // open the file write-only. //nolint:revive
	O_RDWR   = os.O_RDWR   // open the file read-write. //nolint:revive
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

// ReadOnlyMode defines how to handle a file's read-only attribute on Windows.
type ReadOnlyMode int

const (
	// ReadOnlyModeIgnore does not set a file's read-only attribute, and ignores
	// if it's set (Windows only).
	ReadOnlyModeIgnore ReadOnlyMode = 0 + iota
	// ReadOnlyMaskSet set a file's read-only attribute, if the specified
	// perm FileMode has the user writable bit (0o200) set. Otherwise, it will
	// resets (clears) it. (Windows only).
	ReadOnlyModeSet
	// ReadOnlyMaskReset does not set a file's read-only attribute, and if it's
	// set, it resets (clears) it. (Windows only).
	ReadOnlyModeReset
)
