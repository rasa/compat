// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"os"

	"github.com/rasa/compat/consts"
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
	O_FILE_FLAG_DELETE_ON_CLOSE = consts.O_FILE_FLAG_DELETE_ON_CLOSE
	// O_FILE_FLAG_NO_RO_ATTR skips setting a file's read-only attribute on Windows.
	O_FILE_FLAG_NO_RO_ATTR = consts.O_FILE_FLAG_NO_RO_ATTR
)

// A FileMode represents a file's mode and permission bits.
// The bits have the same definition on all systems, so that
// information about files can be moved from one system
// to another portably. Not all bits apply to all systems.
// The only required bit is [ModeDir] for directories.
type FileMode = os.FileMode

// ReadOnlyMode defines how to handle a file's read-only attribute on Windows.
type ReadOnlyMode int

const (
	// ReadOnlyModeIgnore does not set a file's read-only attribute, and ignores
	// if it's set (Windows only).
	ReadOnlyModeIgnore ReadOnlyMode = 0 + iota
	// ReadOnlyModeSet set a file's read-only attribute, if the specified
	// perm FileMode has the user writable bit (0o200) set. Otherwise, it will
	// resets (clears) it. (Windows only).
	ReadOnlyModeSet
	// ReadOnlyModeReset does not set a file's read-only attribute, and if it's
	// set, it resets (clears) it. (Windows only).
	ReadOnlyModeReset
)
