// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9

package compat

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func fstat(f *os.File) (FileInfo, error) {
	if f == nil {
		return nil, &os.PathError{Op: "stat", Path: "", Err: os.ErrInvalid}
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, &os.PathError{Op: "stat", Path: f.Name(), Err: err}
	}

	pid := os.Getpid()
	fd := int(f.Fd())
	fdinfo := "/proc/" + strconv.Itoa(pid) + "/fdinfo/" + strconv.Itoa(fd)

	data, err := os.ReadFile(fdinfo)
	if err != nil {
		return nil, &os.PathError{Op: "stat", Path: f.Name(), Err: err}
	}

	lines := strings.SplitN(string(data), "\n", 2)
	if len(lines) == 0 {
		return nil, &os.PathError{Op: "stat", Path: f.Name(), Err: os.ErrInvalid}
	}
	// First line is usually the path
	path := filepath.Clean(lines[0])

	return stat(fi, path, false)
}
