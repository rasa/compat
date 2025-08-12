// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build js || unix || wasip1

// unix == aix || darwin || dragonfly || freebsd || illumos || linux || netbsd || openbsd || solaris

package compat

import (
	"syscall"
)

// Getuid returns the User ID for the current user. On Windows, the user's SID is
// converted to its POSIX equivalent, which is compatible with Cygwin and
// Git for Windows. On Plan9, Getuid returns a 32-bit hash of the user's name.
func Getuid() (int, error) {
	return syscall.Getuid(), nil
}

// Getgid returns the default Group ID for the current user. On Windows, the
// user's primary group's SID is converted to its POSIX equivalent, which is
// compatible with Cygwin and Git for Windows. On Plan9, Getuid returns a
// 32-bit hash of the user's group's name, as provided by golang's os/user package.
func Getgid() (int, error) {
	return syscall.Getgid(), nil
}
