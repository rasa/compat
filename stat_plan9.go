// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9

package compat

import (
	"os"
	"syscall"
	"time"

	"github.com/cespare/xxhash"
)

// Not supported: Links | BTime | CTime.
const supported SupportedType = ATime | UID | GID

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
	links  uint64
	atime  time.Time
	btime  time.Time
	ctime  time.Time
	uid    uint64
	gid    uint64
}

func stat(fi os.FileInfo, _ string) (FileInfo, error) {
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
	fs.uid = xxhash.Sum64([]byte(fs.sys.Uid))
	fs.gid = xxhash.Sum64([]byte(fs.sys.Gid))

	return &fs, nil
}

// See https://github.com/golang/go/blob/d13da639/src/os/types_plan9.go#L26
