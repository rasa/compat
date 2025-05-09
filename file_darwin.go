// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build darwin || ios

package compat

import (
	"time"
)

const supports SupportsType = SupportsLinks | SupportsUID | SupportsGID | SupportsATime | SupportsBTime | SupportsCTime

func (fs *fileStat) times() {
	fs.atime = time.Unix(int64(fs.sys.Atimespec.Sec), int64(fs.sys.Atimespec.Nsec))         //nolint:unconvert // needed conversion
	fs.btime = time.Unix(int64(fs.sys.Birthtimespec.Sec), int64(fs.sys.Birthtimespec.Nsec)) //nolint:unconvert // needed conversion
	fs.ctime = time.Unix(int64(fs.sys.Ctimespec.Sec), int64(fs.sys.Ctimespec.Nsec))         //nolint:unconvert // needed conversion
}
