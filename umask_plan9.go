// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9

package compat

// Umask sets the umask to umask, and returns the previous value.
// On Windows, the intial umask value is 022 octal, and can be changed by
// setting environmental variable UMASK, to an octal value. For example:
// `set UMASK=022` . Leading zeros and 'o's are allowed, and ignored.
// On Plan9, the function does nothing, and always returns zero.
func Umask(_ int) int {
	return 0
}

// GetUmask returns the current umask value.
// On Plan9, the function always returns zero.
func GetUmask() int {
	return 0
}
