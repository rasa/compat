// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build js || unix

// unix == aix || darwin || dragonfly || freebsd || illumos || linux || netbsd || openbsd || solaris

package compat

import (
	"os"
	"syscall"
)

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
