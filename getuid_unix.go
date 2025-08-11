// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build js || unix || wasip1

// unix == aix || darwin || dragonfly || freebsd || illumos || linux || netbsd || openbsd || solaris

package compat

import (
	"syscall"
)

// Getuid returns the User ID as a uint64. On Windows, the user's SID is
// converted to it's POSIX equivalent, which is compatiable with Cygwin and
// Git for Windows. On Plan9, the User ID is a 64-bit hash of the user's name.
func Getuid() (uint64, error) {
	return uint64(syscall.Getuid()), nil //nolint:gosec // quiet linter
}

// Getgid returns the Group ID as a uint64. On Windows, the user's primary group's
// SID is converted to its POSIX equivalent, which is compatiable with Cygwin and
// Git for Windows. On Plan9, the Getgid returns the value returned by Getuid().
func Getgid() (uint64, error) {
	return uint64(syscall.Getgid()), nil //nolint:gosec // quiet linter
}
