// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build linux

package compat

import (
	"os"
	"path/filepath"
	"strconv"
)

func fstat(file *os.File) (FileInfo, error) {
	if file == nil {
		return nil, statError("", os.ErrInvalid)
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, statError(file.Name(), err)
	}

	fd := int(file.Fd())

	link := "/proc/self/fd/" + strconv.Itoa(fd)

	path, err := os.Readlink(link)
	if err != nil {
		return nil, statError(file.Name(), err)
	}

	path = filepath.Clean(path)

	return stat(fi, path, false)
}
