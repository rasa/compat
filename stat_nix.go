// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || wasm

package compat

import (
	"os"
	"syscall"
	"time"
)

// A fileStat is the implementation of FileInfo returned by Stat and Lstat.
// See https://github.com/golang/go/blob/8cd6d68a/src/os/types_unix.go#L15
type fileStat struct {
	name   string
	size   int64
	mode   os.FileMode
	mtime  time.Time
	sys    syscall.Stat_t
	partID uint64
	fileID uint64
	links  uint64
	atime  time.Time
	btime  time.Time
	ctime  time.Time
	uid    uint64
	gid    uint64
	path   string
}

func stat(fi os.FileInfo, name string) (FileInfo, error) {
	var fs fileStat

	fs.path = name
	fs.name = fi.Name()
	fs.size = fi.Size()
	fs.mode = fi.Mode()
	fs.mtime = fi.ModTime()
	fs.sys = *fi.Sys().(*syscall.Stat_t)

	fs.partID = uint64(fs.sys.Dev) //nolint:gosec,unconvert,nolintlint // intentional int32 → uint64 conversion
	fs.fileID = fs.sys.Ino
	fs.links = uint64(fs.sys.Nlink) //nolint:gosec,unconvert,nolintlint // intentional int32 → uint64 conversion
	fs.uid = uint64(fs.sys.Uid)
	fs.gid = uint64(fs.sys.Gid)
	fs.times()

	return &fs, nil
}

func _stat(name string) (os.FileMode, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return 0, err
	}

	return fi.Mode(), nil
}
