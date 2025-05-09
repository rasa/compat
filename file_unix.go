// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build aix || android || dragonfly || linux || openbsd || solaris

package compat

import (
	"time"
)

// not supported: SupportsBTime.
const supports SupportsType = SupportsLinks | SupportsATime | SupportsCTime | SupportsUID | SupportsGID

func (fs *fileStat) times() {
	fs.atime = time.Unix(int64(fs.sys.Atim.Sec), int64(fs.sys.Atim.Nsec)) //nolint:unconvert // needed conversion
	fs.ctime = time.Unix(int64(fs.sys.Ctim.Sec), int64(fs.sys.Ctim.Nsec)) //nolint:unconvert // needed conversion
}

// https://github.com/golang/go/blob/d000963d/src/os/types_unix.go#L28
