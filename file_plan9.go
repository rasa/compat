// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9

package compat

import (
	"errors"
	"os"
	"syscall"
	"time"

	"github.com/cespare/xxhash"
)

// Not supported: SupportsLinks | SupportsBTime | SupportsCTime
const supports SupportsType = SupportsUID | SupportsGID | SupportsATime

// A fileStat is the implementation of FileInfo returned by Stat and Lstat.
// See https://github.com/golang/go/blob/8cd6d68a/src/os/types_plan9.go#L13
type fileStat struct {
	name     string
	size     int64
	mode     os.FileMode
	mtime    time.Time
	sys      syscall.Dir
	deviceID uint64
	fileID   uint64
	links    uint64
	atime    time.Time
	btime    time.Time
	ctime    time.Time
	uid      uint64
	gid      uint64
}

func stat(name string) (FileInfo, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return &fileStat{}, err
	}

	return loadInfo(fi, name)
}

func lstat(name string) (FileInfo, error) {
	fi, err := os.Lstat(name)
	if err != nil {
		return &fileStat{}, err
	}

	return loadInfo(fi, name)
}

func loadInfo(fi os.FileInfo, _ string) (FileInfo, error) {
	var fs fileStat

	sys, ok := fi.Sys().(*syscall.Dir)
	if !ok {
		return &fs, errors.New("failed to cast fi.Sys()")
	}

	fs.name = fi.Name()
	fs.size = fi.Size()
	fs.mode = fi.Mode()
	fs.mtime = fi.ModTime()
	fs.sys = *sys

	fs.deviceID = uint64(fs.sys.Type)<<32 + uint64(fs.sys.Dev)
	fs.fileID = uint64(fs.sys.Qid.Path)
	// fs.links not supported
	fs.atime = time.Unix(int64(fs.sys.Atime), 0)
	// fs.btime not supported
	// fs.ctime not supported
	fs.uid = xxhash.Sum64([]byte(fs.sys.Uid))
	fs.gid = xxhash.Sum64([]byte(fs.sys.Gid))

	return &fs, nil
}

// https://github.com/golang/go/blob/d13da63929df73ab506314f35524ebb9b0f8a216/src/os/types_plan9.go#L26
