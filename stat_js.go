// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build js

package compat

import (
	"time"
)

// Not supported: BTime.
const supported SupportedType = Links | ATime | CTime | UID | GID

func (fs *fileStat) times() {
	fs.atime = time.Unix(fs.sys.Atime, int64(fs.sys.AtimeNsec))
	fs.ctime = time.Unix(fs.sys.Ctime, int64(fs.sys.CtimeNsec))
}
