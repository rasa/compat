// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

import (
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/capnspacehook/go-acl"
)

const supports supportsType = supportsLinks | supportsATime | supportsBTime | supportsCTime | supportsSymlinks

const userIDSource UserIDSourceType = UserIDSourceIsSID

// A fileStat is the implementation of FileInfo returned by Stat and Lstat.
// See https://github.com/golang/go/blob/8cd6d68a/src/os/types_windows.go#L18
type fileStat struct {
	name  string
	size  int64
	mode  os.FileMode
	mtime time.Time
	// See https://github.com/golang/go/blob/cad1fc52/src/os/types_windows.go#L276
	sys    syscall.Win32FileAttributeData
	partID uint64
	fileID uint64
	links  uint64
	atime  time.Time
	btime  time.Time
	ctime  time.Time
	uid    int
	gid    int
	user   string
	group  string
	path   string
	// btimed bool unused
	ctimed bool
	usered bool
	// grouped bool unused
	followSymlinks bool
	err            error
	mux            sync.Mutex // only on Windows
}

// See https://github.com/golang/go/blob/cad1fc52/src/runtime/os_windows.go#L448
var canUseLongPaths bool

func init() {
	canUseLongPaths = isWindowsAtLeast(10, 0, 15063) //nolint:mnd // quiet linter
}

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

// Portions of the following code is:
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

////////////////////////////////////////////////////////////////////////////////
// Originally copied from https://github.com/golang/go/blob/d65c209b/src/os/types_windows.go#L287
////////////////////////////////////////////////////////////////////////////////

func stat(fi os.FileInfo, name string, followSymlinks bool) (FileInfo, error) {
	var fs fileStat

	fs.followSymlinks = followSymlinks

	name = fixLongPath(name)
	pathp, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return nil, &os.PathError{Op: "stat", Path: name, Err: err}
	}

	fs.mux.Lock()
	defer fs.mux.Unlock()

	attrs := uint32(syscall.FILE_FLAG_BACKUP_SEMANTICS)

	if !followSymlinks {
		attrs |= syscall.FILE_FLAG_OPEN_REPARSE_POINT
	}

	h, err := windows.CreateFile(pathp, 0, 0, nil, windows.OPEN_EXISTING, attrs, 0)
	if err != nil {
		return nil, &os.PathError{Op: "stat", Path: name, Err: err}
	}
	defer windows.CloseHandle(h) //nolint:errcheck // quiet linter

	var i windows.ByHandleFileInformation
	err = windows.GetFileInformationByHandle(h, &i)
	if err != nil {
		return nil, &os.PathError{Op: "stat", Path: name, Err: err}
	}

	fs.path = name
	fs.name = fi.Name()
	fs.size = fi.Size()
	fs.mode = fi.Mode()
	fs.mtime = fi.ModTime()
	fs.sys = *fi.Sys().(*syscall.Win32FileAttributeData)

	perm, err := _stat(name)
	if err != nil {
		return nil, &os.PathError{Op: "stat", Path: name, Err: err}
	}
	fs.mode &= os.FileMode(^uint32(0o777)) //nolint:mnd // quiet
	fs.mode |= perm.Perm()

	fs.partID = uint64(i.VolumeSerialNumber)                             // uint32
	fs.fileID = (uint64(i.FileIndexHigh) << 32) + uint64(i.FileIndexLow) //nolint:mnd // quiet linter
	fs.links = uint64(i.NumberOfLinks)
	fs.atime = time.Unix(0, fs.sys.LastAccessTime.Nanoseconds())
	fs.btime = time.Unix(0, fs.sys.CreationTime.Nanoseconds())

	return &fs, nil
}

func (fs *fileStat) BTime() time.Time {
	return fs.btime
}

func (fs *fileStat) CTime() time.Time {
	if fs.ctimed {
		return fs.ctime
	}

	fs.ctimed = true

	pathp, err := windows.UTF16PtrFromString(fs.path)
	if err != nil {
		fs.err = &os.PathError{Op: "stat", Path: fs.path, Err: err}
		return fs.ctime
	}

	fs.mux.Lock()
	defer fs.mux.Unlock()

	attrs := uint32(syscall.FILE_FLAG_BACKUP_SEMANTICS)

	if !fs.followSymlinks {
		attrs |= syscall.FILE_FLAG_OPEN_REPARSE_POINT
	}

	h, err := windows.CreateFile(pathp, 0, 0, nil, windows.OPEN_EXISTING, attrs, 0)
	if err != nil {
		fs.err = &os.PathError{Op: "stat", Path: fs.path, Err: err}
		return fs.ctime
	}
	defer windows.CloseHandle(h) //nolint:errcheck // quiet linter

	var bi FILE_BASIC_INFO
	err = windows.GetFileInformationByHandleEx(h, windows.FileBasicInfo, (*byte)(unsafe.Pointer(&bi)), uint32(unsafe.Sizeof(bi)))
	if err != nil {
		fs.err = &os.PathError{Op: "stat", Path: fs.path, Err: err}
		return fs.ctime
	}

	// ChangedTime is 100-nanosecond intervals since January 1, 1601.
	nsec := bi.ChangedTime
	// Change starting time to the Epoch (00:00:00 UTC, January 1, 1970).
	nsec -= 116444736000000000
	// Convert into nanoseconds.
	nsec *= 100
	fs.ctime = time.Unix(0, nsec)

	return fs.ctime
}

func (fs *fileStat) UID() int {
	if !fs.usered {
		fs.usered = true
		var err error
		fs.uid, fs.gid, fs.user, fs.group, err = getUserGroup(fs.path)
		if err != nil {
			fs.err = err
		}
	}

	return fs.uid
}

func (fs *fileStat) GID() int {
	if !fs.usered {
		fs.usered = true
		var err error
		fs.uid, fs.gid, fs.user, fs.group, err = getUserGroup(fs.path)
		if err != nil {
			fs.err = err
		}
	}

	return fs.gid
}

func (fs *fileStat) User() string {
	if !fs.usered {
		fs.usered = true
		var err error
		fs.uid, fs.gid, fs.user, fs.group, err = getUserGroup(fs.path)
		if err != nil {
			fs.err = err
		}
	}

	return fs.user
}

func (fs *fileStat) Group() string {
	if !fs.usered {
		fs.usered = true
		var err error
		fs.uid, fs.gid, fs.user, fs.group, err = getUserGroup(fs.path)
		if err != nil {
			fs.err = err
		}
	}

	return fs.group
}

func _stat(name string) (os.FileMode, error) {
	return acl.GetExplicitFileAccessMode(name)
}

///////////////////////////////////////////////////////////////////////////////
// The following code is:
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
///////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// Copied from https://github.com/golang/go/blob/dbaa2d3e/src/internal/syscall/windows/syscall_windows.go#L162
////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////
// Copied from https://github.com/golang/go/blob/cad1fc52/src/os/path_windows.go#L100
////////////////////////////////////////////////////////////////////////////////

func fixLongPath(path string) string {
	if canUseLongPaths {
		return path
	}
	return addExtendedPrefix(path)
}

////////////////////////////////////////////////////////////////////////////////
// Copied from https://github.com/golang/go/blob/cad1fc52/src/os/path_windows.go#L107
////////////////////////////////////////////////////////////////////////////////

// addExtendedPrefix adds the extended path prefix (\\?\) to path.
func addExtendedPrefix(path string) string { //nolint:gocyclo // quiet linter
	if len(path) >= 4 { //nolint:mnd // quiet linter
		if path[:4] == `\??\` {
			// Already extended with \??\
			return path
		}
		if os.IsPathSeparator(path[0]) && os.IsPathSeparator(path[1]) && path[2] == '?' && os.IsPathSeparator(path[3]) {
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
	if !filepath.IsAbs(path) {
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

	if pathLength < 248 { //nolint:mnd // quiet linter
		// Don't fix. (This is how Go 1.7 and earlier worked,
		// not automatically generating the \\?\ form)
		return path
	}

	var isUNC, isDevice bool
	if len(path) >= 2 && os.IsPathSeparator(path[0]) && os.IsPathSeparator(path[1]) {
		if len(path) >= 4 && path[2] == '.' && os.IsPathSeparator(path[3]) {
			// Starts with //./
			isDevice = true
		} else {
			// Starts with //
			isUNC = true
		}
	}
	var prefix []uint16
	if isUNC { //nolint:gocritic // quiet linter
		// UNC path, prepend the \\?\UNC\ prefix.
		prefix = []uint16{'\\', '\\', '?', '\\', 'U', 'N', 'C', '\\'}
	} else if isDevice { //nolint:revive // quiet linter //revive:disable-line
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
	n := uint32(pathLength) + 1 //nolint:gosec // quiet linter
	var buf []uint16
	for {
		buf = make([]uint16, n+uint32(len(prefix))) //nolint:gosec // quiet linter
		n, err = syscall.GetFullPathName(&p[0], n, &buf[len(prefix)], nil)
		if err != nil {
			return path
		}
		if n <= uint32(len(buf)-len(prefix)) { //nolint:gosec // quiet linter
			buf = buf[:n+uint32(len(prefix))] //nolint:gosec // quiet linter
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

////////////////////////////////////////////////////////////////////////////////
// Copied from https://github.com/golang/go/blob/cad1fc52/src/os/getwd.go#L13
////////////////////////////////////////////////////////////////////////////////

var getwdCache struct {
	sync.Mutex
	dir string
}
