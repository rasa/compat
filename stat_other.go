// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// this excludes all known build targets, so it will only build on newly introduced systems:
//go:build !(aix || darwin || dragonfly || freebsd || illumos || linux || netbsd || openbsd || plan9 || solaris || wasm || windows)

package compat

import (
	"os"
	"time"
)

// Not supported: Links | ATime | BTime | CTime | UID | GID.
const supported SupportedType = 0

// A fileStat is the implementation of FileInfo returned by Stat and Lstat.
type fileStat struct {
	name   string
	size   int64
	mode   os.FileMode
	mtime  time.Time
	sys    any
	partID uint64
	fileID uint64
	links  uint64
	atime  time.Time
	btime  time.Time
	ctime  time.Time
	uid    uint32
	gid    uint32
}

func loadInfo(fi os.FileInfo) (FileInfo, error) {
	return fileStat{}, errors.New("not implemented")
}

func sameFile(_, _ *fileStat) bool {
	return false
}

func sameDevice(_, _ *fileStat) bool {
	return false
}
