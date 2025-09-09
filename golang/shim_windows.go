// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package golang

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/capnspacehook/go-acl"
	"golang.org/x/sys/windows"
)

const (
	// Redefining here to avoid a circular dependency.
	// O_DELETE deletes the file when closed.
	O_DELETE = 0x8000000
	// O_NOROATTR doesn't set a file's read-only attribute if mode.
	O_NOROATTR = 0x4000000
)

const perm600 = os.FileMode(0o600)

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
	O_CLOEXEC                    = syscall.O_CLOEXEC
	OPEN_ALWAYS                  = syscall.OPEN_ALWAYS
	OPEN_EXISTING                = syscall.OPEN_EXISTING
	S_IWRITE                     = syscall.S_IWRITE
	STANDARD_RIGHTS_WRITE        = syscall.STANDARD_RIGHTS_WRITE
	SYNCHRONIZE                  = syscall.SYNCHRONIZE
	// See https://github.com/golang/go/blob/77f911e3/src/syscall/types_windows.go#L100
	// and https://github.com/golang/go/blob/77f911e3/src/internal/syscall/windows/types_windows.go#L27
	_FILE_WRITE_EA = windows.FILE_WRITE_EA
)

var (
	FixLongPath   = fixLongPath
	OpenFile      = openFileNolog
	OpenFileNolog = openFileNolog
	RemoveAll     = removeAll
)

var (
	CloseHandle        = syscall.CloseHandle
	CreateFile         = syscall.CreateFile
	Ftruncate          = syscall.Ftruncate
	GetFileAttributes  = syscall.GetFileAttributes
	Syscall9           = syscall.Syscall9 //nolint:staticcheck
	UTF16PtrFromString = syscall.UTF16PtrFromString
)

type (
	Handle             = syscall.Handle
	SecurityAttributes = syscall.SecurityAttributes
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
	canUseLongPaths = isWindowsAtLeast(10, 0, 15063) //nolint:mnd
}

func fixAttributesAndShareMode(flag int, attrs uint32, sharemode uint32) (uint32, uint32) {
	if flag&O_DELETE == O_DELETE {
		attrs &^= uint32(windows.FILE_ATTRIBUTE_READONLY)
		attrs |= (windows.FILE_FLAG_DELETE_ON_CLOSE | windows.FILE_ATTRIBUTE_TEMPORARY)
		sharemode |= FILE_SHARE_DELETE
		flag &^= O_DELETE
	}

	if flag&O_NOROATTR == O_NOROATTR {
		attrs &^= uint32(windows.FILE_ATTRIBUTE_READONLY)
		flag &^= O_NOROATTR //nolint:ineffassign
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

func Remove(path string) error {
	var err error
	// See https://github.com/golang/go/blob/77f911e3/src/os/removeall_noat.go#L126
	err1 := os.Remove(path)
	if err1 == nil || os.IsNotExist(err1) {
		return nil
	}
	if /* runtime.GOOS == "windows" && */ os.IsPermission(err1) {
		if fs, err2 := os.Stat(path); err2 == nil {
			err = acl.Chmod(path, perm600)
			if err != nil {
				return fmt.Errorf("remove: %w (1)", err)
			}
			if err = os.Chmod(path, os.FileMode(0o200|uint32(fs.Mode().Perm()))); err == nil { //nolint:mnd
				err1 = os.Remove(path)
			}
		}
	}
	if err == nil {
		err = err1
	}

	return err
}

// Source: https://github.com/golang/go/blob/77f911e3/src/syscall/syscall_windows.go#L357-L362

func makeInheritSa(sa *SecurityAttributes) *SecurityAttributes {
	if sa == nil {
		sa = &SecurityAttributes{}
		sa.Length = uint32(unsafe.Sizeof(sa))
	}
	sa.InheritHandle = 1

	return sa
}

var procGetFinalPathNameByHandleW = modkernel32.NewProc("GetFinalPathNameByHandleW")

// Source: https://github.com/golang/go/blob/77f911e3/src/syscall/zsyscall_windows.go#L783-790

func GetFinalPathNameByHandle(file Handle, filePath *uint16, filePathSize uint32, flags uint32) (n uint32, err error) {
	r0, _, e1 := syscall.Syscall6(procGetFinalPathNameByHandleW.Addr(), 4, uintptr(file), uintptr(unsafe.Pointer(filePath)), uintptr(filePathSize), uintptr(flags), 0, 0) //nolint:gosec,mnd,staticcheck
	n = uint32(r0)
	if n == 0 || n >= filePathSize {
		err = errnoErr(e1)
	}
	return n, err
}
