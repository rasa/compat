// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build (!fstat && (dragonfly || netbsd || openbsd)) || aix || illumos || js || plan9 || solaris || wasip1

package compat

import (
	"os"
	"path/filepath"
)

func fstat(f *os.File) (FileInfo, error) {
	if f == nil {
		return nil, statError("", os.ErrInvalid)
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, statError(f.Name(), err)
	}

	return stat(fi, filepath.Clean(f.Name()), false)
}
