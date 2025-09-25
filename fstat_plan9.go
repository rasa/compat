// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9

package compat

import (
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

func fstat(f *os.File) (FileInfo, error) {
	if f == nil {
		return nil, &os.PathError{Op: "stat", Path: "", Err: os.ErrInvalid}
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, &os.PathError{Op: "stat", Path: f.Name(), Err: err}
	}

	pid := syscall.Getpid()
	fd := f.Fd()
	link := "/proc/" + strconv.Itoa(pid) + "/fd/" + strconv.Itoa(int(fd))
	path, err := os.Readlink(link)
	if err != nil {
		return nil, &os.PathError{Op: "stat", Path: f.Name(), Err: err}
	}

	path = filepath.Clean(path)

	return stat(fi, path, false)
}
