// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build fstat && (dragonfly || netbsd)

// was go:build dragonfly || netbsd

package compat

import (
	"os"
	"path/filepath"
	"strconv"
)

func fstat(f *os.File) (FileInfo, error) {
	if f == nil {
		return nil, statError("", os.ErrInvalid)
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, statError(f.Name(), err)
	}

	pid := os.Getpid()
	fd := int(f.Fd())
	link := "/proc/" + strconv.Itoa(pid) + "/fd/" + strconv.Itoa(fd)
	path, err := os.Readlink(link)
	if err != nil {
		return nil, statError(f.Name(), err)
	}
	path = filepath.Clean(path)

	return stat(fi, path, false)
}
