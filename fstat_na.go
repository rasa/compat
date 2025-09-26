// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !(darwin || linux || plan9 || windows)

package compat

import (
	"os"
)

func fstat(f *os.File) (FileInfo, error) {
	if f == nil {
		return nil, &os.PathError{Op: "stat", Path: "", Err: os.ErrInvalid}
	}

	return nil, &os.PathError{Op: "stat", Path: f.Name(), Err: os.ErrInvalid}
}
