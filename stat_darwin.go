// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build darwin

// The darwin build flag includes ios

package compat

import (
	"time"
)

func (fs *fileStat) times() {
	fs.atime = time.Unix(int64(fs.sys.Atimespec.Sec), int64(fs.sys.Atimespec.Nsec))         //nolint:unconvert
	fs.btime = time.Unix(int64(fs.sys.Birthtimespec.Sec), int64(fs.sys.Birthtimespec.Nsec)) //nolint:unconvert
	fs.ctime = time.Unix(int64(fs.sys.Ctimespec.Sec), int64(fs.sys.Ctimespec.Nsec))         //nolint:unconvert
}

func (fs *fileStat) BTime() time.Time { return fs.btime }
