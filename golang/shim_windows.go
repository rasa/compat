// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package golang

import (
	"errors"
	"os"
	"unsafe"

	"github.com/capnspacehook/go-acl"
	"golang.org/x/sys/windows"

	"github.com/rasa/compat/consts"
)

const (
	// Redefining here to avoid a circular dependency.
	// O_FILE_FLAG_DELETE_ON_CLOSE deletes the file when closed.
	O_FILE_FLAG_DELETE_ON_CLOSE = consts.O_FILE_FLAG_DELETE_ON_CLOSE
	// O_FILE_FLAG_NO_RO_ATTR skips setting a file's read-only attribute on Windows.
	O_FILE_FLAG_NO_RO_ATTR = consts.O_FILE_FLAG_NO_RO_ATTR
)

const perm600 = os.FileMode(0o600)

const (
	CREATE_NEW                   = windows.CREATE_NEW
	EISDIR                       = windows.EISDIR
	ENOTDIR                      = windows.ENOTDIR
	ERROR_ACCESS_DENIED          = windows.ERROR_ACCESS_DENIED
	ERROR_ALREADY_EXISTS         = windows.ERROR_ALREADY_EXISTS
	ERROR_FILE_NOT_FOUND         = windows.ERROR_FILE_NOT_FOUND
	FILE_APPEND_DATA             = windows.FILE_APPEND_DATA
	FILE_ATTRIBUTE_DIRECTORY     = windows.FILE_ATTRIBUTE_DIRECTORY
	FILE_ATTRIBUTE_NORMAL        = windows.FILE_ATTRIBUTE_NORMAL
	FILE_ATTRIBUTE_READONLY      = windows.FILE_ATTRIBUTE_READONLY
	FILE_FLAG_BACKUP_SEMANTICS   = windows.FILE_FLAG_BACKUP_SEMANTICS
	FILE_FLAG_OPEN_REPARSE_POINT = windows.FILE_FLAG_OPEN_REPARSE_POINT
	FILE_SHARE_DELETE            = windows.FILE_SHARE_DELETE
	FILE_SHARE_READ              = windows.FILE_SHARE_READ
	FILE_SHARE_WRITE             = windows.FILE_SHARE_WRITE
	FILE_WRITE_ATTRIBUTES        = windows.FILE_WRITE_ATTRIBUTES
	GENERIC_READ                 = windows.GENERIC_READ
	GENERIC_WRITE                = windows.GENERIC_WRITE
	InvalidHandle                = windows.InvalidHandle
	O_CREAT                      = windows.O_CREAT
	O_CLOEXEC                    = windows.O_CLOEXEC
	// https://github.com/golang/go/blob/ac803b59/src/syscall/types_windows.go#L50
	o_DIRECTORY           = 0x04000
	OPEN_ALWAYS           = windows.OPEN_ALWAYS
	OPEN_EXISTING         = windows.OPEN_EXISTING
	S_IWRITE              = windows.S_IWRITE
	STANDARD_RIGHTS_WRITE = windows.STANDARD_RIGHTS_WRITE
	SYNCHRONIZE           = windows.SYNCHRONIZE
	// See https://github.com/golang/go/blob/ac803b59/src/syscall/types_windows.go#L114
	// and https://github.com/golang/go/blob/ac803b59/src/internal/syscall/windows/types_windows.go#L28
	_FILE_WRITE_EA = windows.FILE_WRITE_EA
	// See https://github.com/golang/go/blob/ac803b59/src/internal/syscall/windows/types_windows.go#L180
	O_FILE_FLAG_OVERLAPPED = windows.FILE_FLAG_OVERLAPPED
)

var (
	FixLongPath   = fixLongPath
	OpenFile      = openFileNolog
	OpenFileNolog = openFileNolog
	RemoveAll     = removeAll
)

var (
	CloseHandle                = windows.CloseHandle
	CreateFile                 = windows.CreateFile
	Ftruncate                  = windows.Ftruncate
	GetFileAttributes          = windows.GetFileAttributes
	GetFileInformationByHandle = windows.GetFileInformationByHandle
	// Syscall6                   = windows.Syscall6 //nolint:staticcheck.
	UTF16PtrFromString = windows.UTF16PtrFromString
)

type (
	ByHandleFileInformation = windows.ByHandleFileInformation
	Handle                  = windows.Handle
	SecurityAttributes      = windows.SecurityAttributes
)

// See https://github.com/golang/go/blob/ac803b59/src/runtime/os_windows.go#L435
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
	if flag&O_FILE_FLAG_DELETE_ON_CLOSE == O_FILE_FLAG_DELETE_ON_CLOSE {
		attrs &^= uint32(windows.FILE_ATTRIBUTE_READONLY)
		attrs |= (windows.FILE_FLAG_DELETE_ON_CLOSE | windows.FILE_ATTRIBUTE_TEMPORARY)
		sharemode |= FILE_SHARE_DELETE
	}

	if flag&O_FILE_FLAG_NO_RO_ATTR == O_FILE_FLAG_NO_RO_ATTR {
		attrs &^= uint32(windows.FILE_ATTRIBUTE_READONLY)
	}

	return attrs, sharemode
}

func mkdir(longName string, _ uint32, sa *windows.SecurityAttributes) error {
	name, err := windows.UTF16PtrFromString(longName)
	if err != nil {
		return err
	}
	err = windows.CreateDirectory(name, sa)
	if err != nil {
		_ = Remove(longName)

		return err
	}

	return nil
}

// See https://github.com/golang/go/blob/ac803b59/src/os/sticky_notbsd.go#L9

const supportsCreateWithStickyBit = true

// See https://github.com/golang/go/blob/ac803b59/src/os/file.go#L351-L357

func setStickyBit(name string) error {
	return nil
}

// See https://github.com/golang/go/blob/ac803b59/src/syscall/zsyscall_windows.go#L43
var modkernel32 = windows.NewLazyDLL("kernel32.dll")

// procCreateFileW = modkernel32.NewProc("CreateFileW").

// Emulate newFile() as f.cleanup and f.pfd are private.
// See https://github.com/golang/go/blob/ac803b59/src/os/file_windows.go#L50
func newFile(h windows.Handle, name string /*kind*/, _ string /*nonBlocking*/, _ bool) *File {
	if h == windows.InvalidHandle {
		return nil
	}

	return NewFile(uintptr(h), name)
}

func Remove(path string) error {
	var err error
	// See https://github.com/golang/go/blob/ac803b59/src/os/removeall_noat.go#L126
	err1 := os.Remove(path)
	if err1 == nil || os.IsNotExist(err1) {
		return nil
	}
	if /* runtime.GOOS == "windows" && */ os.IsPermission(err1) {
		if fs, err2 := os.Stat(path); err2 == nil {
			err = acl.Chmod(path, perm600)
			if err != nil {
				return &PathError{Op: "remove", Path: path, Err: err}
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

// Source: https://github.com/golang/go/blob/ac803b59/src/syscall/syscall_windows.go#L357-L362

func makeInheritSa(sa *SecurityAttributes) *SecurityAttributes {
	if sa == nil {
		sa = &SecurityAttributes{}
		sa.Length = uint32(unsafe.Sizeof(sa))
	}
	sa.InheritHandle = 1

	return sa
}

var procGetFinalPathNameByHandleW = modkernel32.NewProc("GetFinalPathNameByHandleW")

// Source: https://github.com/golang/go/blob/ac803b59/src/syscall/zsyscall_windows.go#L783-L790

func GetFinalPathNameByHandle(file Handle, filePath *uint16, filePathSize uint32, flags uint32) (n uint32, err error) {
	r0, _, e1 := windows.Syscall6(procGetFinalPathNameByHandleW.Addr(), 4, uintptr(file), uintptr(unsafe.Pointer(filePath)), uintptr(filePathSize), uintptr(flags), 0, 0) //nolint:gosec,mnd,staticcheck
	n = uint32(r0)
	if n == 0 || n >= filePathSize {
		err = errnoErr(e1)
	}
	return n, err
}

// Source: https://github.com/golang/go/blob/ac803b59/src/syscall/syscall_windows.go#L1295-L1312

func fdpath(fd windows.Handle, buf []uint16) ([]uint16, error) {
	const (
		FILE_NAME_NORMALIZED = 0
		VOLUME_NAME_DOS      = 0
	)
	for {
		n, err := GetFinalPathNameByHandle(fd, &buf[0], uint32(len(buf)), FILE_NAME_NORMALIZED|VOLUME_NAME_DOS) //nolint:gosec
		if err == nil {
			buf = buf[:n]
			break
		}
		if !errors.Is(err, windows.ERROR_NOT_ENOUGH_MEMORY) { // compat: was if err != _ERROR_NOT_ENOUGH_MEMORY {
			return nil, err
		}
		buf = append(buf, make([]uint16, n-uint32(len(buf)))...) //nolint:gosec
	}
	return buf, nil
}

func Filepath(f *os.File) (string, error) {
	if f == nil {
		return "", errors.New("nil file pointer")
	}
	fd := windows.Handle(f.Fd())

	// Source: https://github.com/golang/go/blob/ac803b59/src/syscall/syscall_windows.go#L1315-L1331

	var buf [windows.MAX_PATH + 1]uint16
	path, err := fdpath(fd, buf[:])
	if err != nil {
		return "", err
	}
	// When using VOLUME_NAME_DOS, the path is always prefixed by "\\?\".
	// That prefix tells the Windows APIs to disable all string parsing and to send
	// the string that follows it straight to the file system.
	// Although SetCurrentDirectory and GetCurrentDirectory do support the "\\?\" prefix,
	// some other Windows APIs don't. If the prefix is not removed here, it will leak
	// to Getwd, and we don't want such a general-purpose function to always return a
	// path with the "\\?\" prefix after Fchdir is called.
	// The downside is that APIs that do support it will parse the path and try to normalize it,
	// when it's already normalized.
	if len(path) >= 4 && path[0] == '\\' && path[1] == '\\' && path[2] == '?' && path[3] == '\\' {
		path = path[4:]
	}
	pathString := windows.UTF16ToString(path)

	return pathString, nil
}

// Source: https://github.com/golang/go/blob/ac803b59/src/syscall/types_windows.go#L92

const fileFlagsMask = 0xFFF00000

// Source: https://github.com/golang/go/blob/ac803b59/src/syscall/types_windows.go#L136-L148

const (
	// The following flags are supported by [Open]
	// and exported in [golang.org/x/sys/windows].
	_FILE_FLAG_OPEN_NO_RECALL = 0x00100000
	// FILE_FLAG_OPEN_REPARSE_POINT = 0x00200000.
	_FILE_FLAG_SESSION_AWARE   = 0x00800000
	_FILE_FLAG_POSIX_SEMANTICS = 0x01000000
	// FILE_FLAG_BACKUP_SEMANTICS   = 0x02000000.
	_FILE_FLAG_DELETE_ON_CLOSE = 0x04000000
	_FILE_FLAG_SEQUENTIAL_SCAN = 0x08000000
	_FILE_FLAG_RANDOM_ACCESS   = 0x10000000
	_FILE_FLAG_NO_BUFFERING    = 0x20000000
	// FILE_FLAG_OVERLAPPED         = 0x40000000.
	_FILE_FLAG_WRITE_THROUGH = 0x80000000
)

// Source: https://github.com/golang/go/blob/ac803b59/src/syscall/types_windows.go#L94-L104

const validFileFlagsMask = FILE_FLAG_OPEN_REPARSE_POINT |
	FILE_FLAG_BACKUP_SEMANTICS |
	windows.FILE_FLAG_OVERLAPPED |
	_FILE_FLAG_OPEN_NO_RECALL |
	_FILE_FLAG_SESSION_AWARE |
	_FILE_FLAG_POSIX_SEMANTICS |
	_FILE_FLAG_DELETE_ON_CLOSE |
	_FILE_FLAG_SEQUENTIAL_SCAN |
	_FILE_FLAG_NO_BUFFERING |
	_FILE_FLAG_RANDOM_ACCESS |
	_FILE_FLAG_WRITE_THROUGH

// Source: https://github.com/golang/go/blob/ac803b59/src/internal/oserror/errors.go#L13

var ErrInvalid = errors.New("invalid argument")
