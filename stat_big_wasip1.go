// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build wasip1 && !tinygo

package compat

import (
	"time"
)

func (fs *fileStat) times() {
	fs.atime = time.Unix(0, int64(fs.sys.Atime))
	fs.ctime = time.Unix(0, int64(fs.sys.Ctime))
}
