// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build linux

package compat

import (
	"time"

	"golang.org/x/sys/unix"
)

const supports supportsType = supportsLinks | supportsBTime | supportsCTime | supportsSymlinks

const userIDSource UserIDSourceType = UserIDSourceIsNumeric

func (fs *fileStat) times() {
	fs.atime = time.Unix(int64(fs.sys.Atim.Sec), int64(fs.sys.Atim.Nsec)) //nolint:unconvert // needed conversion
	fs.ctime = time.Unix(int64(fs.sys.Ctim.Sec), int64(fs.sys.Ctim.Nsec)) //nolint:unconvert // needed conversion
}

func (fs *fileStat) BTime() time.Time {
	if !fs.btimed {
		fs.btimed = true

		var stx unix.Statx_t

		var flags int
		if !fs.followSymlinks {
			flags = unix.AT_SYMLINK_NOFOLLOW
		}
		err := unix.Statx(unix.AT_FDCWD, fs.path, flags, unix.STATX_BTIME, &stx)
		if err != nil {
			fs.err = err
			return fs.btime
		}

		if stx.Mask&unix.STATX_BTIME == 0 {
			return fs.btime
		}

		fs.btime = time.Unix(int64(stx.Btime.Sec), int64(stx.Btime.Nsec)) //nolint:unconvert // needed conversion
	}

	return fs.btime
}

// See https://github.com/golang/go/blob/d000963d/src/os/types_unix.go#L28
