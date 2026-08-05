// SPDX-FileCopyrightText: Copyright © 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !tinygo

package compat

// SupportsATimeSetting returns true setting a files's atime is supported by the OS.
func SupportsATimeSetting() bool {
	return supportsATimeSetting
}
