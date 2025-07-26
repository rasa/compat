// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !(plan9 || windows)

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
