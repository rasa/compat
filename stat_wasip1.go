// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build wasip1

package compat

import (
	"os"
	"syscall"
)

const (
	defaultFileMode = os.FileMode(0o600)
	defaultDirMode  = os.FileMode(0o700)
)

func stat(fi os.FileInfo, name string, _ bool) (FileInfo, error) {
	var fs fileStat

	fs.path = name
	fs.name = fi.Name()
	fs.size = fi.Size()
	fs.mode = fi.Mode()
	fs.mtime = fi.ModTime()
	fs.sys = *fi.Sys().(*syscall.Stat_t)
	// See https://github.com/golang/go/blob/5045fdd8/src/os/stat_wasip1.go#L35
	if fs.mode == 0 {
		if fs.sys.Mode == syscall.S_IFDIR {
			fs.mode = defaultDirMode
		} else {
			fs.mode = defaultFileMode
		}
	}
	fs.partID = uint64(fs.sys.Dev) //nolint:gosec,unconvert,nolintlint // intentional int32 → uint64 conversion
	fs.fileID = fs.sys.Ino
	fs.links = uint64(fs.sys.Nlink) //nolint:gosec,unconvert,nolintlint // intentional int32 → uint64 conversion
	fs.uid = int(fs.sys.Uid)
	fs.gid = int(fs.sys.Gid)
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
