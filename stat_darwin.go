// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build darwin

// The darwin build flag includes ios

package compat

import (
	"time"
)

const userIDSource UserIDSourceType = UserIDSourceIsInt

func (fs *fileStat) times() {
	fs.atime = time.Unix(int64(fs.sys.Atimespec.Sec), int64(fs.sys.Atimespec.Nsec))         //nolint:unconvert // needed conversion
	fs.btime = time.Unix(int64(fs.sys.Birthtimespec.Sec), int64(fs.sys.Birthtimespec.Nsec)) //nolint:unconvert // needed conversion
	fs.ctime = time.Unix(int64(fs.sys.Ctimespec.Sec), int64(fs.sys.Ctimespec.Nsec))         //nolint:unconvert // needed conversion
}

func (fs *fileStat) BTime() time.Time { return fs.btime }
