// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

// IsAdmin returns true if the user is root, or has Windows administrator rights.
//
// Deprecated: Use IsRoot() instead.
func IsAdmin() (bool, error) {
	return IsRoot()
}
