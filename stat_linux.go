// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build linux

package compat

import (
	"time"

	"golang.org/x/sys/unix"
)

// @TODO(rasa): determine why BTime is not working.
const supported SupportedType = Links | ATime | CTime | UID | GID

func (fs *fileStat) times() {
	fs.atime = time.Unix(int64(fs.sys.Atim.Sec), int64(fs.sys.Atim.Nsec)) //nolint:unconvert // needed conversion
	fs.ctime = time.Unix(int64(fs.sys.Ctim.Sec), int64(fs.sys.Ctim.Nsec)) //nolint:unconvert // needed conversion

	var stx unix.Statx_t
	err := unix.Statx(unix.AT_FDCWD, "file.txt", unix.AT_SYMLINK_NOFOLLOW, unix.STATX_BTIME, &stx)
	if err != nil {
		return
	}
	if stx.Mask&unix.STATX_BTIME == 0 {
		return
	}
	fs.btime = time.Unix(int64(stx.Btime.Sec), int64(stx.Btime.Sec)) //nolint:unconvert // needed conversion
}

// See https://github.com/golang/go/blob/d000963d/src/os/types_unix.go#L28
