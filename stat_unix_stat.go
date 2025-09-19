// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build js || unix

// unix == aix || darwin || dragonfly || freebsd || illumos || linux || netbsd || openbsd || solaris

package compat

import (
	"errors"
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"
)

func stat(fi os.FileInfo, name string, followSymlinks bool) (FileInfo, error) {
	if fi == nil {
		err := errors.New("fileInfo is nil")
		return nil, &os.PathError{Op: "stat", Path: name, Err: err}
	}

	var fs fileStat

	fs.path = name
	fs.followSymlinks = followSymlinks
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

	fs.times()

	return &fs, nil
}

func (fs *fileStat) CTime() time.Time { return fs.ctime }

func (fs *fileStat) UID() int { return fs.uid }
func (fs *fileStat) GID() int { return fs.gid }

func (fs *fileStat) User() string {
	if !fs.usered {
		u, err := user.LookupId(strconv.Itoa(fs.uid))
		if err != nil {
			fs.err = err
		} else {
			fs.user = u.Username
		}
		fs.usered = true
	}

	return fs.user
}

func (fs *fileStat) Group() string {
	if !fs.grouped {
		g, err := user.LookupGroupId(strconv.Itoa(fs.gid))
		if err != nil {
			fs.err = err
		} else {
			fs.group = g.Name
		}
		fs.grouped = true
	}

	return fs.group
}
