// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !windows

package compat

import (
	"os"
)

func _stat(name string) (os.FileMode, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return 0, err
	}

	return fi.Mode(), nil
}
