// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build linux

package compat

import (
	"time"

	"golang.org/x/sys/unix"
)

func (fs *fileStat) times() {
	fs.atime = time.Unix(int64(fs.sys.Atim.Sec), int64(fs.sys.Atim.Nsec)) //nolint:unconvert // needed conversion
	fs.ctime = time.Unix(int64(fs.sys.Ctim.Sec), int64(fs.sys.Ctim.Nsec)) //nolint:unconvert // needed conversion
}

func (fs *fileStat) BTime() time.Time {
	fs.btimeOnce.Do(func() {
		var stx unix.Statx_t

		var flags int
		if !fs.followSymlinks {
			flags = unix.AT_SYMLINK_NOFOLLOW
		}

		err := unix.Statx(unix.AT_FDCWD, fs.path, flags, unix.STATX_BTIME, &stx)
		if err != nil {
			fs.err = err

			return
		}

		if stx.Mask&unix.STATX_BTIME == 0 {
			return
		}

		fs.btime = time.Unix(int64(stx.Btime.Sec), int64(stx.Btime.Nsec)) //nolint:unconvert // needed conversion
	})

	return fs.btime
}
