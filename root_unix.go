// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !windows

package compat

import (
	"os"
)

// IsRoot returns true if the user is (effectively) root on a non-Windows
// system, or is running with elevated privileges (administrator rights) on a
// Windows system.
func IsRoot() (bool, error) {
	return os.Geteuid() == 0, nil
}
