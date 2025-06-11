// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

// IsAdmin returns true if the user is root, or has Windows administrator rights.
//
// Deprecated: Use IsRoot() instead.
func IsAdmin() (bool, error) {
	return IsRoot()
}

// ReplaceFile atomically replaces the destination file or directory with the
// source.  It is guaranteed to either replace the target file entirely, or not
// change either file.
//
// Deprecated: Use Rename() instead.
func ReplaceFile(source, destination string) error {
	return Rename(source, destination)
}
