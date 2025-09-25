// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build freebsd || netbsd

package compat

import (
	"time"
)

// Not supported: none.
const supports supportsType = supportsATime | supportsBTime | supportsCTime | supportsFstat | supportsLinks | supportsNice | supportsSymlinks

const userIDSource UserIDSourceType = UserIDSourceIsInt

func (fs *fileStat) times() {
	fs.atime = time.Unix(int64(fs.sys.Atimespec.Sec), int64(fs.sys.Atimespec.Nsec))
	fs.btime = time.Unix(int64(fs.sys.Birthtimespec.Sec), int64(fs.sys.Birthtimespec.Nsec))
	fs.ctime = time.Unix(int64(fs.sys.Ctimespec.Sec), int64(fs.sys.Ctimespec.Nsec))
}

func (fs *fileStat) BTime() time.Time { return fs.btime }
