// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build wasip1

package compat

import (
	"os"
	"syscall"
	"time"
)

// const (
// 	defaultFileMode = os.FileMode(0o600)
// 	defaultDirMode  = os.FileMode(0o700)
// )

func (fs *fileStat) BTime() time.Time { return fs.btime }
func (fs *fileStat) CTime() time.Time { return fs.ctime }

func (fs *fileStat) UID() int { return fs.uid }
func (fs *fileStat) GID() int { return fs.gid }

func (fs *fileStat) User() string  { return fs.user }
func (fs *fileStat) Group() string { return fs.group }

func stat(fi os.FileInfo, name string, _ bool) (FileInfo, error) {
	if fi == nil {
		return nil, &os.PathError{Op: "stat", Path: name, Err: os.ErrInvalid}
	}

	var fs fileStat

	fs.path = name
	fs.name = fi.Name()
	fs.size = fi.Size()
	fs.mode = fi.Mode()
	fs.mtime = fi.ModTime()
	fs.sys = *fi.Sys().(*syscall.Stat_t)

	fs.partID = uint64(fs.sys.Dev) //nolint:gosec,unconvert,nolintlint // intentional int32 → uint64 conversion
	fs.fileID = fs.sys.Ino
	fs.links = uint(fs.sys.Nlink) //nolint:gosec,unconvert,nolintlint // intentional int32 → uint conversion
	fs.uid = int(fs.sys.Uid)
	fs.gid = int(fs.sys.Gid)

	// See https://github.com/golang/go/blob/5045fdd8/src/os/stat_wasip1.go#L35
	// This code doesn't seem to be needed any more.
	// if fs.mode == 0 {
	// 	if fs.sys.Mode == syscall.S_IFDIR {
	// 		fs.mode = defaultDirMode | os.ModeDir
	// 	} else {
	// 		fs.mode = defaultFileMode
	// 	}
	// }

	// https://github.com/golang/go/blob/5045fdd8/src/syscall/syscall_wasip1.go#L356
	if fs.uid == 0 {
		fs.uid = os.Getuid()
	}
	if fs.gid == 0 {
		fs.gid = os.Getgid()
	}

	fs.times()

	return &fs, nil
}
