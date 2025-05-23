// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build freebsd || netbsd

package compat

import (
	"time"
)

// Not supported: BTime.
const supported SupportedType = Links | ATime | CTime | UID | GID

func (fs *fileStat) times() {
	fs.atime = time.Unix(int64(fs.sys.Atimespec.Sec), int64(fs.sys.Atimespec.Nsec))
	// fs.btime not supported
	fs.ctime = time.Unix(int64(fs.sys.Ctimespec.Sec), int64(fs.sys.Ctimespec.Nsec))
}
