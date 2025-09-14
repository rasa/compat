//nolint:all
// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package golang

import (
	"fmt"
	"io"
	"os"
	filepathlite "path/filepath"
	"runtime"
	"syscall"

	"github.com/capnspacehook/go-acl"
)

// SPDX-FileCopyrightText: Copyright 2012 The Go Authors. All rights reserved.
// SPDX-License-Identifier: BSD-3

// The following code is:
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Snippet: https://github.com/golang/go/blob/ac803b59/src/internal/syscall/windows/syscall_windows.go#L162-L179

type FILE_BASIC_INFO struct {
	CreationTime   int64
	LastAccessTime int64
	LastWriteTime  int64
	ChangedTime    int64
	FileAttributes uint32

	// Pad out to 8-byte alignment.
	//
	// Without this padding, TestChmod fails due to an argument validation error
	// in SetFileInformationByHandle on windows/386.
	//
	// https://learn.microsoft.com/en-us/cpp/build/reference/zp-struct-member-alignment?view=msvc-170
	// says that “The C/C++ headers in the Windows SDK assume the platform's
	// default alignment is used.” What we see here is padding rather than
	// alignment, but maybe it is related.
	_ uint32
}

// Snippet: https://github.com/golang/go/blob/ac803b59/src/internal/poll/errno_windows.go#L14-L16

var errERROR_IO_PENDING error = syscall.Errno(syscall.ERROR_IO_PENDING)

// Snippet: https://github.com/golang/go/blob/ac803b59/src/internal/poll/errno_windows.go#L20-L31

func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case syscall.ERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

// Snippet: https://github.com/golang/go/blob/ac803b59/src/os/file.go#L327-L348

// compat: s|FileMode|FileMode, sa *syscall.SecurityAttributes|
func Mkdir(name string, perm FileMode, sa *syscall.SecurityAttributes) error {
	longName := fixLongPath(name)
	e := ignoringEINTR(func() error {
		return mkdir(longName, syscallMode(perm), sa) // compat: s|syscall.Mkdir|mkdir|; s|\)$|), sa)|;
	})

	if e != nil {
		return &PathError{Op: "mkdir", Path: name, Err: e}
	}

	// mkdir(2) itself won't handle the sticky bit on *BSD and Solaris
	if !supportsCreateWithStickyBit && perm&ModeSticky != 0 {
		e = setStickyBit(name)

		if e != nil {
			Remove(name)
			return e
		}
	}

	return nil
}

// Snippet: https://github.com/golang/go/blob/ac803b59/src/os/file_posix.go#L60-L73

func syscallMode(i FileMode) (o uint32) {
	o |= uint32(i.Perm())
	if i&ModeSetuid != 0 {
		o |= syscall.S_ISUID
	}
	if i&ModeSetgid != 0 {
		o |= syscall.S_ISGID
	}
	if i&ModeSticky != 0 {
		o |= syscall.S_ISVTX
	}
	// No mapping for Go's ModeTemporary (plan9 only).
	return
}

// Snippet: https://github.com/golang/go/blob/ac803b59/src/os/file_posix.go#L254-L261

func ignoringEINTR(fn func() error) error {
	for {
		err := fn()
		if err != syscall.EINTR {
			return err
		}
	}
}

// Snippet: https://github.com/golang/go/blob/ac803b59/src/os/file_windows.go#L114-L125

// compat: s|FileMode|FileMode, sa *syscall.SecurityAttributes|
func openFileNolog(name string, flag int, perm FileMode, sa *syscall.SecurityAttributes) (*File, error) {
	if name == "" {
		return nil, &PathError{Op: "open", Path: name, Err: syscall.ENOENT}
	}
	path := fixLongPath(name)
	r, err := Open(path, flag|syscall.O_CLOEXEC, syscallMode(perm), sa) // compat: s|\)$|), sa)|
	if err != nil {
		return nil, &PathError{Op: "open", Path: name, Err: err}
	}
	nonblocking := flag&O_FILE_FLAG_OVERLAPPED != 0 // compat: s|windows.||
	return newFile(r, name, "file", nonblocking), nil
}

// Snippet: https://github.com/golang/go/blob/ac803b59/src/os/path.go#L19-L66

// compat: s|FileMode|FileMode, sa *syscall.SecurityAttributes|
func MkdirAll(path string, perm FileMode, sa *syscall.SecurityAttributes) error {
	// Fast path: if we can tell whether path is a directory or file, stop with success or error.
	dir, err := Stat(path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return &PathError{Op: "mkdir", Path: path, Err: syscall.ENOTDIR}
	}

	// Slow path: make sure parent exists and then call Mkdir for path.

	// Extract the parent folder from path by first removing any trailing
	// path separator and then scanning backward until finding a path
	// separator or reaching the beginning of the string.
	i := len(path) - 1
	for i >= 0 && IsPathSeparator(path[i]) {
		i--
	}
	for i >= 0 && !IsPathSeparator(path[i]) {
		i--
	}
	if i < 0 {
		i = 0
	}

	// If there is a parent directory, and it is not the volume name,
	// recurse to ensure parent directory exists.
	if parent := path[:i]; len(parent) > len(filepathlite.VolumeName(path)) {
		err = MkdirAll(parent, perm, sa) // compat: s|\)$|), sa)|
		if err != nil {
			return err
		}
	}

	// Parent now exists; invoke Mkdir and use its result.
	err = Mkdir(path, perm, sa) // compat: s|\)$|), sa)|
	if err != nil {
		// Handle arguments like "foo/." by
		// double-checking that directory doesn't exist.
		dir, err1 := Lstat(path)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}
	return nil
}

// Snippet: https://github.com/golang/go/blob/ac803b59/src/os/path_windows.go#L100-L105

func fixLongPath(path string) string {
	if canUseLongPaths {
		return path
	}
	return addExtendedPrefix(path)
}

// Snippet: https://github.com/golang/go/blob/ac803b59/src/os/path_windows.go#L108-L202

func addExtendedPrefix(path string) string {
	if len(path) >= 4 {
		if path[:4] == `\??\` {
			// Already extended with \??\
			return path
		}
		if IsPathSeparator(path[0]) && IsPathSeparator(path[1]) && path[2] == '?' && IsPathSeparator(path[3]) {
			// Already extended with \\?\ or any combination of directory separators.
			return path
		}
	}

	// Do nothing (and don't allocate) if the path is "short".
	// Empirically (at least on the Windows Server 2013 builder),
	// the kernel is arbitrarily okay with < 248 bytes. That
	// matches what the docs above say:
	// "When using an API to create a directory, the specified
	// path cannot be so long that you cannot append an 8.3 file
	// name (that is, the directory name cannot exceed MAX_PATH
	// minus 12)." Since MAX_PATH is 260, 260 - 12 = 248.
	//
	// The MSDN docs appear to say that a normal path that is 248 bytes long
	// will work; empirically the path must be less then 248 bytes long.
	pathLength := len(path)
	if !filepathlite.IsAbs(path) {
		// If the path is relative, we need to prepend the working directory
		// plus a separator to the path before we can determine if it's too long.
		// We don't want to call syscall.Getwd here, as that call is expensive to do
		// every time fixLongPath is called with a relative path, so we use a cache.
		// Note that getwdCache might be outdated if the working directory has been
		// changed without using os.Chdir, i.e. using syscall.Chdir directly or cgo.
		// This is fine, as the worst that can happen is that we fail to fix the path.
		getwdCache.Lock()
		if getwdCache.dir == "" {
			// Init the working directory cache.
			getwdCache.dir, _ = syscall.Getwd()
		}
		pathLength += len(getwdCache.dir) + 1
		getwdCache.Unlock()
	}

	if pathLength < 248 {
		// Don't fix. (This is how Go 1.7 and earlier worked,
		// not automatically generating the \\?\ form)
		return path
	}

	var isUNC, isDevice bool
	if len(path) >= 2 && IsPathSeparator(path[0]) && IsPathSeparator(path[1]) {
		if len(path) >= 4 && path[2] == '.' && IsPathSeparator(path[3]) {
			// Starts with //./
			isDevice = true
		} else {
			// Starts with //
			isUNC = true
		}
	}
	var prefix []uint16
	if isUNC {
		// UNC path, prepend the \\?\UNC\ prefix.
		prefix = []uint16{'\\', '\\', '?', '\\', 'U', 'N', 'C', '\\'}
	} else if isDevice {
		// Don't add the extended prefix to device paths, as it would
		// change its meaning.
	} else {
		prefix = []uint16{'\\', '\\', '?', '\\'}
	}

	p, err := syscall.UTF16FromString(path)
	if err != nil {
		return path
	}
	// Estimate the required buffer size using the path length plus the null terminator.
	// pathLength includes the working directory. This should be accurate unless
	// the working directory has changed without using os.Chdir.
	n := uint32(pathLength) + 1
	var buf []uint16
	for {
		buf = make([]uint16, n+uint32(len(prefix)))
		n, err = syscall.GetFullPathName(&p[0], n, &buf[len(prefix)], nil)
		if err != nil {
			return path
		}
		if n <= uint32(len(buf)-len(prefix)) {
			buf = buf[:n+uint32(len(prefix))]
			break
		}
	}
	if isUNC {
		// Remove leading \\.
		buf = buf[2:]
	}
	copy(buf, prefix)
	return syscall.UTF16ToString(buf)
}

// Snippet: https://github.com/golang/go/blob/ac803b59/src/os/tempfile.go#L35-L58

// compat: s|int|int, sa *syscall.SecurityAttributes|
func CreateTemp(dir, pattern string, flag int, sa *syscall.SecurityAttributes) (*File, error) {
	if dir == "" {
		dir = TempDir()
	}

	prefix, suffix, err := prefixAndSuffix(pattern)
	if err != nil {
		return nil, &PathError{Op: "createtemp", Path: pattern, Err: err}
	}
	prefix = joinPath(dir, prefix)

	try := 0
	for {
		name := prefix + nextRandom() + suffix
		f, err := OpenFile(name, O_RDWR|O_CREATE|O_EXCL|flag, 0o600, sa) // compat: s|\)$|), sa)|
		if IsExist(err) {
			if try++; try < 10000 {
				continue
			}
			return nil, &PathError{Op: "createtemp", Path: prefix + "*" + suffix, Err: ErrExist}
		}
		return f, err
	}
}

// Snippet: https://github.com/golang/go/blob/ac803b59/src/os/tempfile.go#L86-L117

// compat: s|string|string, sa *syscall.SecurityAttributes|
func MkdirTemp(dir, pattern string, sa *syscall.SecurityAttributes) (string, error) {
	if dir == "" {
		dir = TempDir()
	}

	prefix, suffix, err := prefixAndSuffix(pattern)
	if err != nil {
		return "", &PathError{Op: "mkdirtemp", Path: pattern, Err: err}
	}
	prefix = joinPath(dir, prefix)

	try := 0
	for {
		name := prefix + nextRandom() + suffix
		err := Mkdir(name, 0o700, sa) // compat: s|\)$|), sa)|
		if err == nil {
			return name, nil
		}
		if IsExist(err) {
			if try++; try < 10000 {
				continue
			}
			return "", &PathError{Op: "mkdirtemp", Path: prefix + "*" + suffix, Err: ErrExist}
		}
		if IsNotExist(err) {
			if _, err := Stat(dir); IsNotExist(err) {
				return "", err
			}
		}
		return "", err
	}
}

// Snippet: https://github.com/golang/go/blob/ac803b59/src/syscall/syscall_windows.go#L364-L456

// compat: s|uint32|uint32, sa *syscall.SecurityAttributes|
func Open(name string, flag int, perm uint32, sa *syscall.SecurityAttributes) (fd Handle, err error) {
	if len(name) == 0 {
		return InvalidHandle, ERROR_FILE_NOT_FOUND
	}
	namep, err := UTF16PtrFromString(name)
	if err != nil {
		return InvalidHandle, err
	}
	accessFlags := flag & (O_RDONLY | O_WRONLY | O_RDWR)
	var access uint32
	switch accessFlags {
	case O_RDONLY:
		access = GENERIC_READ
	case O_WRONLY:
		access = GENERIC_WRITE
	case O_RDWR:
		access = GENERIC_READ | GENERIC_WRITE
	}
	if flag&O_CREAT != 0 {
		access |= GENERIC_WRITE
	}
	if flag&O_APPEND != 0 {
		// Remove GENERIC_WRITE unless O_TRUNC is set, in which case we need it to truncate the file.
		// We can't just remove FILE_WRITE_DATA because GENERIC_WRITE without FILE_WRITE_DATA
		// starts appending at the beginning of the file rather than at the end.
		if flag&O_TRUNC == 0 {
			access &^= GENERIC_WRITE
		}
		// Set all access rights granted by GENERIC_WRITE except for FILE_WRITE_DATA.
		access |= FILE_APPEND_DATA | FILE_WRITE_ATTRIBUTES | _FILE_WRITE_EA | STANDARD_RIGHTS_WRITE | SYNCHRONIZE
	}
	sharemode := uint32(FILE_SHARE_READ | FILE_SHARE_WRITE)
	// var sa *SecurityAttributes // compat: s|var|// var|
	if flag&O_CLOEXEC == 0 {
		sa = makeInheritSa(sa) // compat: s|\)$|sa)|
	}
	var attrs uint32 = FILE_ATTRIBUTE_NORMAL
	if perm&S_IWRITE == 0 {
		attrs = FILE_ATTRIBUTE_READONLY
	}
	// fileFlags contains the high 12 bits of flag.
	// This bit range can be used by the caller to specify the file flags
	// passed to CreateFile. It is an error to use if the bits can't be
	// mapped to the supported FILE_FLAG_* constants.
	if fileFlags := uint32(flag) & fileFlagsMask; fileFlags&^validFileFlagsMask == 0 {
		attrs |= fileFlags
	} else {
		return InvalidHandle, ErrInvalid // compat: s|oserror.||
	}
	switch accessFlags {
	case O_WRONLY, O_RDWR:
		// Unix doesn't allow opening a directory with O_WRONLY
		// or O_RDWR, so we don't set the flag in that case,
		// which will make CreateFile fail with ERROR_ACCESS_DENIED.
		// We will map that to EISDIR if the file is a directory.
	default:
		// We might be opening a directory for reading,
		// and CreateFile requires FILE_FLAG_BACKUP_SEMANTICS
		// to work with directories.
		attrs |= FILE_FLAG_BACKUP_SEMANTICS
	}
	if flag&O_SYNC != 0 {
		attrs |= _FILE_FLAG_WRITE_THROUGH
	}
	// We don't use CREATE_ALWAYS, because when opening a file with
	// FILE_ATTRIBUTE_READONLY these will replace an existing file
	// with a new, read-only one. See https://go.dev/issue/38225.
	//
	// Instead, we ftruncate the file after opening when O_TRUNC is set.
	var createmode uint32
	switch {
	case flag&(O_CREAT|O_EXCL) == (O_CREAT | O_EXCL):
		createmode = CREATE_NEW
		attrs |= FILE_FLAG_OPEN_REPARSE_POINT // don't follow symlinks
	case flag&O_CREAT == O_CREAT:
		createmode = OPEN_ALWAYS
	default:
		createmode = OPEN_EXISTING
	}
	attrs, sharemode = fixAttributesAndShareMode(flag, attrs, sharemode)             // compat: added
	h, err := syscall.CreateFile(namep, access, sharemode, sa, createmode, attrs, 0) // compat: s|createFile|syscall.CreateFile|
	if h == InvalidHandle {
		if err == ERROR_ACCESS_DENIED && (attrs&FILE_FLAG_BACKUP_SEMANTICS == 0) {
			// We should return EISDIR when we are trying to open a directory with write access.
			fa, e1 := GetFileAttributes(namep)
			if e1 == nil && fa&FILE_ATTRIBUTE_DIRECTORY != 0 {
				err = EISDIR
			}
		}
		return h, err
	}
	if flag&o_DIRECTORY != 0 {
		// Check if the file is a directory, else return ENOTDIR.
		var fi ByHandleFileInformation
		if err := GetFileInformationByHandle(h, &fi); err != nil {
			CloseHandle(h)
			return InvalidHandle, err
		}
		if fi.FileAttributes&FILE_ATTRIBUTE_DIRECTORY == 0 {
			CloseHandle(h)
			return InvalidHandle, ENOTDIR
		}
	}
	// Ignore O_TRUNC if the file has just been created.
	if flag&O_TRUNC == O_TRUNC &&
		(createmode == OPEN_EXISTING || (createmode == OPEN_ALWAYS && err == ERROR_ALREADY_EXISTS)) {
		err = Ftruncate(h, 0)
		if err != nil {
			CloseHandle(h)
			return InvalidHandle, err
		}
	}
	return h, nil
}

// Snippet: https://github.com/golang/go/blob/ac803b59/src/syscall/zsyscall_windows.go#L506-L513

// func createFile(name *uint16, access uint32, mode uint32, sa *SecurityAttributes, createmode uint32, attrs uint32, templatefile int32) (handle Handle, err error) {
// 	r0, _, e1 := SyscallN(procCreateFileW.Addr(), uintptr(unsafe.Pointer(name)), uintptr(access), uintptr(mode), uintptr(unsafe.Pointer(sa)), uintptr(createmode), uintptr(attrs), uintptr(templatefile))
// 	handle = Handle(r0)
// 	if handle == InvalidHandle || e1 == ERROR_ALREADY_EXISTS {
// 		err = errnoErr(e1)
// 	}
// 	return
// }

// Snippet: https://github.com/golang/go/blob/ac803b59/src/os/removeall_noat.go#L15-L142

func removeAll(path string) error {
	if path == "" {
		// fail silently to retain compatibility with previous behavior
		// of RemoveAll. See issue 28830.
		return nil
	}

	// The rmdir system call permits removing "." on Plan 9,
	// so we don't permit it to remain consistent with the
	// "at" implementation of RemoveAll.
	if endsWithDot(path) {
		return &PathError{Op: "RemoveAll", Path: path, Err: syscall.EINVAL}
	}

	// Simple case: if Remove works, we're done.
	err := Remove(path)
	if err == nil || IsNotExist(err) {
		return nil
	}

	// Otherwise, is this a directory we need to recurse into?
	dir, serr := Lstat(path)
	if serr != nil {
		if serr, ok := serr.(*PathError); ok && (IsNotExist(serr.Err) || serr.Err == syscall.ENOTDIR) {
			return nil
		}
		return serr
	}
	if !dir.IsDir() {
		// Not a directory; return the error from Remove.
		return err
	}

	// Remove contents & return first error.
	err = nil
	for {
		fd, err := os.Open(path) // compat: s|Open|os.Open|
		if err != nil {
			if IsNotExist(err) {
				// Already deleted by someone else.
				return nil
			}
			return err
		}

		const reqSize = 1024
		var names []string
		var readErr error

		for {
			numErr := 0
			names, readErr = fd.Readdirnames(reqSize)

			for _, name := range names {
				err1 := removeAll(path + string(PathSeparator) + name) // compat: s|RemoveAll|removeAll|
				if err == nil {
					err = err1
				}
				if err1 != nil {
					numErr++
				}
			}

			// If we can delete any entry, break to start new iteration.
			// Otherwise, we discard current names, get next entries and try deleting them.
			if numErr != reqSize {
				break
			}
		}

		// Removing files from the directory may have caused
		// the OS to reshuffle it. Simply calling Readdirnames
		// again may skip some entries. The only reliable way
		// to avoid this is to close and re-open the
		// directory. See issue 20841.
		fd.Close()

		if readErr == io.EOF {
			break
		}
		// If Readdirnames returned an error, use it.
		if err == nil {
			err = readErr
		}
		if len(names) == 0 {
			break
		}

		// We don't want to re-open unnecessarily, so if we
		// got fewer than request names from Readdirnames, try
		// simply removing the directory now. If that
		// succeeds, we are done.
		if len(names) < reqSize {
			err1 := Remove(path)
			if err1 == nil || IsNotExist(err1) {
				return nil
			}

			if err != nil {
				// We got some error removing the
				// directory contents, and since we
				// read fewer names than we requested
				// there probably aren't more files to
				// remove. Don't loop around to read
				// the directory again. We'll probably
				// just get the same error.
				return err
			}
		}
	}

	// Remove directory.
	err1 := Remove(path)
	if err1 == nil || IsNotExist(err1) {
		return nil
	}
	if runtime.GOOS == "windows" && IsPermission(err1) {
		if fs, err := Stat(path); err == nil {
			err = acl.Chmod(path, 0o600) // compat: added
			if err != nil {              // compat: added
				return fmt.Errorf("removeall: %w (1)", err) // compat: added
			} // compat: added
			if err = Chmod(path, FileMode(0o200|int(fs.Mode()))); err == nil {
				err1 = Remove(path)
			}
		}
	}
	if err == nil {
		err = err1
	}
	return err
}
