// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build js || unix || wasip1

// https://github.com/golang/go/blob/8ad27fb6/src/cmd/dist/build.go#L1070
// unix == aix || android || darwin || dragonfly || freebsd || illumos || ios || linux || netbsd || openbsd || solaris

package compat

import (
	"syscall"
)

// Getuid returns the User ID for the current user.
// On Windows, the user's SID is converted to its POSIX equivalent, which is
// compatible with Cygwin, Git for Windows, MSYS2, etc.
// On Plan9, Getuid returns a 32-bit hash of the user's name.
func Getuid() (int, error) {
	return syscall.Getuid(), nil
}

// Getgid returns the default Group ID for the current user.
// On Windows, the user's primary group's SID is converted to its POSIX
// equivalent, which is compatible with Cygwin, Git for Windows, MSYS2, etc.
// On Plan9, Getgid returns a 32-bit hash of the user's group's name, as
// provided by golang's os/user package.
func Getgid() (int, error) {
	return syscall.Getgid(), nil
}

// Geteuid returns the effective User ID for the current user.
// On Windows, the user's SID is converted to its POSIX equivalent, which is
// compatible with Cygwin, Git for Windows, MSYS2, etc.
// On Plan9, Geteuid returns a 32-bit hash of the user's name.
func Geteuid() (int, error) {
	return syscall.Geteuid(), nil
}

// Getegid returns the effective default Group ID for the current user.
// On Windows, the user's primary group's SID is converted to its POSIX
// equivalent, which is compatible with Cygwin, Git for Windows, MSYS2, etc.
// On Plan9, Getegid returns a 32-bit hash of the user's group's name, as
// provided by golang's os/user package.
func Getegid() (int, error) {
	return syscall.Getegid(), nil
}
