// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build wasip1

package compat

import (
	"time"
)

// not supported: SupportsBTime
const supports SupportsType = SupportsLinks | SupportsUID | SupportsGID | SupportsATime | SupportsCTime

func (fs *fileStat) times() {
	fs.atime = time.Unix(0, int64(fs.sys.Atime))
	fs.ctime = time.Unix(0, int64(fs.sys.Ctime))
}
