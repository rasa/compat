// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build js

package compat

import (
	"time"
)

// Not supported: BTime | Nice.
const supports supportsType = supportsATime | supportsCTime | supportsLinks | supportsSymlinks

const userIDSource UserIDSourceType = UserIDSourceIsNone

func (fs *fileStat) times() {
	fs.atime = time.Unix(fs.sys.Atime, int64(fs.sys.AtimeNsec))
	fs.ctime = time.Unix(fs.sys.Ctime, int64(fs.sys.CtimeNsec))
}

func (fs *fileStat) BTime() time.Time { return fs.btime }
