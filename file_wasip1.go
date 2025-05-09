// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build wasip1

package compat

import (
	"time"
)

// Not supported: BTime.
const supported SupportedType = Links | ATime | CTime | UID | GID

func (fs *fileStat) times() {
	fs.atime = time.Unix(0, int64(fs.sys.Atime))
	fs.ctime = time.Unix(0, int64(fs.sys.Ctime))
}
