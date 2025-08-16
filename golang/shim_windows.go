// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package golang

import (
	"os"
	"syscall"
	_ "unsafe"

	"golang.org/x/sys/windows"
)

type (
	File      = os.File
	FileMode  = os.FileMode
	PathError = os.PathError
)

var (
	ErrExist        = os.ErrExist
	FixLongPath     = fixLongPath
	IsExist         = os.IsExist
	IsNotExist      = os.IsNotExist
	IsPathSeparator = os.IsPathSeparator
	Lstat           = os.Lstat
	ModeSetgid      = os.ModeSetgid
	ModeSetuid      = os.ModeSetuid
	ModeSticky      = os.ModeSticky
	NewFile         = os.NewFile
	OpenFile        = openFileNolog
	OpenFileNolog   = openFileNolog
	PathSeparator   = os.PathSeparator
	Remove          = os.Remove
	Stat            = os.Stat
	TempDir         = os.TempDir
)

const (
	CREATE_NEW                   = syscall.CREATE_NEW
	EISDIR                       = syscall.EISDIR
	ERROR_ACCESS_DENIED          = syscall.ERROR_ACCESS_DENIED
	ERROR_ALREADY_EXISTS         = syscall.ERROR_ALREADY_EXISTS
	ERROR_FILE_NOT_FOUND         = syscall.ERROR_FILE_NOT_FOUND
	FILE_APPEND_DATA             = syscall.FILE_APPEND_DATA
	FILE_ATTRIBUTE_DIRECTORY     = syscall.FILE_ATTRIBUTE_DIRECTORY
	FILE_ATTRIBUTE_NORMAL        = syscall.FILE_ATTRIBUTE_NORMAL
	FILE_ATTRIBUTE_READONLY      = syscall.FILE_ATTRIBUTE_READONLY
	FILE_FLAG_BACKUP_SEMANTICS   = syscall.FILE_FLAG_BACKUP_SEMANTICS
	FILE_FLAG_OPEN_REPARSE_POINT = syscall.FILE_FLAG_OPEN_REPARSE_POINT
	FILE_SHARE_DELETE            = syscall.FILE_SHARE_DELETE
	FILE_SHARE_READ              = syscall.FILE_SHARE_READ
	FILE_SHARE_WRITE             = syscall.FILE_SHARE_WRITE
	FILE_WRITE_ATTRIBUTES        = syscall.FILE_WRITE_ATTRIBUTES
	GENERIC_READ                 = syscall.GENERIC_READ
	GENERIC_WRITE                = syscall.GENERIC_WRITE
	InvalidHandle                = syscall.InvalidHandle
	O_CREAT                      = syscall.O_CREAT
	OPEN_ALWAYS                  = syscall.OPEN_ALWAYS
	OPEN_EXISTING                = syscall.OPEN_EXISTING
	S_IWRITE                     = syscall.S_IWRITE
	STANDARD_RIGHTS_WRITE        = syscall.STANDARD_RIGHTS_WRITE
	SYNCHRONIZE                  = syscall.SYNCHRONIZE
)

type (
	Handle             = syscall.Handle
	SecurityAttributes = syscall.SecurityAttributes
)

var (
	CloseHandle        = syscall.CloseHandle
	CreateFile         = syscall.CreateFile
	Ftruncate          = syscall.Ftruncate
	GetFileAttributes  = syscall.GetFileAttributes
	Syscall9           = syscall.Syscall9
	UTF16PtrFromString = syscall.UTF16PtrFromString
)

// See https://github.com/golang/go/blob/77f911e3/src/syscall/types_windows.go#L100
// and https://github.com/golang/go/blob/77f911e3/src/internal/syscall/windows/types_windows.go#L27
const _FILE_WRITE_EA = windows.FILE_WRITE_EA

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
	O_DELETE = 0x40000000
)

// See https://github.com/golang/go/blob/77f911e3/src/runtime/os_windows.go#L446
var canUseLongPaths bool

func isWindowsAtLeast(major, minor, build uint32) bool {
	mg, mn, bl := windows.RtlGetNtVersionNumbers()
	if mg < major {
		return false
	}
	if mn < minor {
		return false
	}
	return bl >= build
}

func init() {
	canUseLongPaths = isWindowsAtLeast(10, 0, 15063)
}

func setDeleteAttributes(flag int, attrs uint32, sharemode uint32) (uint32, uint32) {
	if flag&O_DELETE == O_DELETE {
		attrs &^= uint32(windows.FILE_ATTRIBUTE_READONLY)
		attrs |= (windows.FILE_FLAG_DELETE_ON_CLOSE | windows.FILE_ATTRIBUTE_TEMPORARY)
		sharemode |= FILE_SHARE_DELETE
	}

	return attrs, sharemode
}

func mkdir(longName string, _ uint32, sa *syscall.SecurityAttributes) error {
	name, err := syscall.UTF16PtrFromString(longName)
	if err != nil {
		return err
	}
	err = syscall.CreateDirectory(name, sa)
	if err != nil {
		_ = Remove(longName)

		return err
	}

	return nil
}

// Source: https://github.com/golang/go/blob/77f911e3/src/os/tempfile.go#L19-L24

//go:linkname runtime_rand runtime.rand
func runtime_rand() uint64

// See https://github.com/golang/go/blob/77f911e3/src/os/sticky_notbsd.go#L9

const supportsCreateWithStickyBit = true

// See https://github.com/golang/go/blob/77f911e3//src/os/file.go#L351-L357

func setStickyBit(name string) error {
	return nil
}

var (
	// See https://github.com/golang/go/blob/77f911e3/src/syscall/zsyscall_windows.go#L43
	modkernel32     = syscall.NewLazyDLL("kernel32.dll")
	procCreateFileW = modkernel32.NewProc("CreateFileW")
)

// emulate newFile() as f.cleanup and f.pfd are private.
func newFile(h syscall.Handle, name string /*kind*/, _ string /*nonBlocking*/, _ bool) *File {
	if h == syscall.InvalidHandle {
		return nil
	}

	return NewFile(uintptr(h), name)
}
