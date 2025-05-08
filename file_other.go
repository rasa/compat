// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !aix && !darwin && !dragonfly && !freebsd && !js && !linux && !netbsd && !openbsd && !plan9 && !solaris && !wasip1 && !windows

package compat

import (
	"os"
	"time"
)

// not supported: SupportsLinks | SupportsUID | SupportsGID | SupportsATime | SupportsBTime | SupportsCTime
const supports SupportsType = 0

// A fileStat is the implementation of FileInfo returned by Stat and Lstat.
type fileStat struct {
	name     string
	size     int64
	mode     os.FileMode
	mtime    time.Time
	sys      any
	deviceID uint64
	fileID   uint64
	links    uint64
	atime    time.Time
	btime    time.Time
	ctime    time.Time
	uid      uint32
	gid      uint32
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
