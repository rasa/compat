// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/capnspacehook/go-acl"
	"golang.org/x/sys/windows"

	"github.com/rasa/compat/golang"
)

const perm000 = os.FileMode(0o0)

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
	links  uint
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
	mux            sync.Mutex // Windows only
}

// Portions of the following code is:
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

////////////////////////////////////////////////////////////////////////////////
// Originally copied from
// https://github.com/golang/go/blob/77f911e3/src/os/types_windows.go#L287-L336
////////////////////////////////////////////////////////////////////////////////

func stat(fi os.FileInfo, name string, followSymlinks bool) (FileInfo, error) {
	if fi == nil {
		err := errors.New("fileInfo is nil")
		return nil, &os.PathError{Op: "stat", Path: name, Err: err}
	}
	var fs fileStat

	fs.mux.Lock()
	defer fs.mux.Unlock()

	fs.followSymlinks = followSymlinks

	name = golang.FixLongPath(name)

	pathp, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return nil, &os.PathError{Op: "stat", Path: name, Err: err}
	}

	attrs := uint32(syscall.FILE_FLAG_BACKUP_SEMANTICS)

	if !followSymlinks {
		attrs |= syscall.FILE_FLAG_OPEN_REPARSE_POINT
	}

	h, err := windows.CreateFile(pathp, 0, 0, nil, windows.OPEN_EXISTING, attrs, 0)
	if err != nil {
		return nil, &os.PathError{Op: "stat", Path: name, Err: err}
	}
	defer windows.CloseHandle(h) //nolint:errcheck
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
	sys, ok := fi.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		err = fmt.Errorf("sys is not a Win32FileAttributeData, it's a %T", fi.Sys())
		return nil, &os.PathError{Op: "stat", Path: name, Err: err}
	}
	fs.sys = *sys

	fs.partID = uint64(i.VolumeSerialNumber)                             // uint32
	fs.fileID = (uint64(i.FileIndexHigh) << 32) + uint64(i.FileIndexLow) //nolint:mnd
	fs.links = uint(i.NumberOfLinks)
	fs.atime = time.Unix(0, fs.sys.LastAccessTime.Nanoseconds())
	fs.btime = time.Unix(0, fs.sys.CreationTime.Nanoseconds())

	perm, err := fs.stat()
	if err != nil {
		return nil, &os.PathError{Op: "stat", Path: name, Err: err}
	}
	fs.mode &^= ModePerm // os.FileMode(^uint32(0o777)) //nolint:mnd // quiet
	fs.mode |= perm.Perm()

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
	defer windows.CloseHandle(h) //nolint:errcheck

	var bi golang.FILE_BASIC_INFO
	err = windows.GetFileInformationByHandleEx(h, windows.FileBasicInfo, (*byte)(unsafe.Pointer(&bi)), uint32(unsafe.Sizeof(bi)))
	if err != nil {
		fs.err = &os.PathError{Op: "stat", Path: fs.path, Err: err}

		return fs.ctime
	}

	if bi.ChangedTime == 0 {
		// exFAT returns 0
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

func (fs *fileStat) stat() (os.FileMode, error) {
	b, err := supportsACLsCached(fs)
	if err == nil && !b {
		if fs.mode.IsDir() {
			return DefaultWindowsDirPerm, nil
		} else {
			return DefaultWindowsFilePerm, nil
		}
	}

	perm, err := acl.GetExplicitFileAccessMode(fs.path)
	if err != nil {
		fs.err = err
		return perm, err
	}
	if perm == perm000 {
		b, err = supportsACLs(fs.path)
		if err != nil {
			fs.err = err
			return perm, err
		}
		if !b {
			if fs.mode.IsDir() {
				return DefaultWindowsDirPerm, nil
			} else {
				return DefaultWindowsFilePerm, nil
			}
		}
	}

	return perm, nil
}
