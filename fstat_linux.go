// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build linux

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

	fd := int(f.Fd())

	link := "/proc/self/fd/" + strconv.Itoa(fd)

	path, err := os.Readlink(link)
	if err != nil {
		return nil, statError(f.Name(), err)
	}

	path = filepath.Clean(path)

	return stat(fi, path, false)
}
