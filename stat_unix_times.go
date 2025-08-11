// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build aix || dragonfly || illumos || openbsd || solaris

package compat

import (
	"time"
)

// Not supported: BTime.
const supported SupportedType = Links | ATime | CTime | UID | GID

const userIDSource UserIDSourceType = UserIDSourceIsNumeric

func (fs *fileStat) times() {
	fs.atime = time.Unix(int64(fs.sys.Atim.Sec), int64(fs.sys.Atim.Nsec)) //nolint:unconvert // needed conversion
	fs.ctime = time.Unix(int64(fs.sys.Ctim.Sec), int64(fs.sys.Ctim.Nsec)) //nolint:unconvert // needed conversion
}

func (fs *fileStat) BTime() time.Time { return fs.btime }

// See https://github.com/golang/go/blob/d000963d/src/os/types_unix.go#L28
