// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

import (
	"os"

	"github.com/rasa/compat/golang"
)

func fstat(f *os.File) (FileInfo, error) {
	if f == nil {
		return nil, statError("", os.ErrInvalid)
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, statError(f.Name(), err)
	}

	path, err := golang.Filepath(f)
	if err != nil {
		return nil, statError(f.Name(), err)
	}

	return stat(fi, path, false)
}
