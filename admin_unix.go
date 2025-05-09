// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !windows

package compat

import (
	"os"
)

// IsAdmin returns true if the user is root, or has Windows Administrator rights.
func IsAdmin() (bool, error) {
	return os.Getuid() == 0, nil
}
