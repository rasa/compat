// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9

package compat

import (
	"errors"
	"os"
	"syscall"
	"time"

	"github.com/OneOfOne/xxhash"
)

// Not supported: BTime | CTime | Links | Symlinks
const supports supportsType = supportsATime | supportsNice

const userIDSource UserIDSourceType = UserIDSourceIsString

// A fileStat is the implementation of FileInfo returned by Stat and Lstat.
// See https://github.com/golang/go/blob/8cd6d68a/src/os/types_plan9.go#L13
type fileStat struct {
	name   string
	size   int64
	mode   os.FileMode
	mtime  time.Time
	sys    syscall.Dir
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
	// path string // unused
	// btimed bool // unused
	// ctimed bool // unused
	usered  bool
	grouped bool
	// followSymlinks bool // unused
	err error
}

func stat(fi os.FileInfo, _ string, _ bool) (FileInfo, error) {
	if fi == nil {
		err := errors.New("fileInfo is nil")
		return nil, &os.PathError{Op: "stat", Path: name, Err: err}
	}

	var fs fileStat

	fs.name = fi.Name()
	fs.size = fi.Size()
	fs.mode = fi.Mode()
	fs.mtime = fi.ModTime()
	fs.sys = *fi.Sys().(*syscall.Dir)

	
	fs.partID = uint64(fs.sys.Type)<<32 + uint64(fs.sys.Dev)
	fs.fileID = uint64(fs.sys.Qid.Path)
	// fs.links not supported
	fs.atime = time.Unix(int64(fs.sys.Atime), 0)
	// fs.btime not supported
	// fs.ctime not supported
	fs.user = fs.sys.Uid
	fs.group = fs.sys.Gid

	return &fs, nil
}

func (fs *fileStat) BTime() time.Time { return fs.btime }
func (fs *fileStat) CTime() time.Time { return fs.ctime }

func (fs *fileStat) UID() int {
	if !fs.usered {
		fs.usered = true
		if fs.user == "" {
			fs.uid = UnknownID
		} else {
			fs.uid = int(xxhash.Checksum32([]byte(fs.user)))
		}
	}

	return fs.uid
}

func (fs *fileStat) GID() int {
	if !fs.grouped {
		fs.grouped = true
		if fs.group == "" {
			fs.gid = UnknownID
		} else {
			fs.gid = int(xxhash.Checksum32([]byte(fs.group)))
		}
	}

	return fs.gid
}

func (fs *fileStat) User() string  { return fs.user }
func (fs *fileStat) Group() string { return fs.group }

// See https://github.com/golang/go/blob/d13da639/src/os/types_plan9.go#L26
