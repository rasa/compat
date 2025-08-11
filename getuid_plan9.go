// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9

package compat

import (
	"os/user"

	"github.com/cespare/xxhash"
)

// Getuid returns the User ID as a uint64. On Windows, the user's SID is
// converted to it's POSIX equivalent, which is compatible with Cygwin and
// Git for Windows. On Plan9, the User ID is a 64-bit hash of the user's name.
func Getuid() (uint64, error) {
	u, err := user.Current()
	if err != nil {
		return UnknownID, err
	}

	uid := xxhash.Sum64([]byte(u.Username))

	return uid, nil
}

// Getgid returns the Group ID as a uint64. On Windows, the user's primary group's
// SID is converted to its POSIX equivalent, which is compatible with Cygwin and
// Git for Windows. On Plan9, the Getgid returns the value returned by Getuid().
func Getgid() (uint64, error) {
	return Getuid()
}
