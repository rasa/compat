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
		return nil, &os.PathError{Op: "stat", Path: "", Err: os.ErrInvalid}
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, &os.PathError{Op: "stat", Path: f.Name(), Err: err}
	}

	fd := int(f.Fd())

	link := "/proc/self/fd/" + strconv.Itoa(fd)
	path, err := os.Readlink(link)
	if err != nil {
		return nil, &os.PathError{Op: "stat", Path: f.Name(), Err: err}
	}
	path = filepath.Clean(path)

	return stat(fi, path, false)
}
