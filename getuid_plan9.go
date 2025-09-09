// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9

package compat

import (
	"os/user"

	"github.com/OneOfOne/xxhash"
	// was "github.com/cespare/xxhash"
)

// Getuid returns the User ID for the current user.
// On Windows, the user's SID is converted to its POSIX equivalent, which is
// compatible with Cygwin, Git for Windows, MSYS2, etc.
// On Plan9, Getuid returns a 32-bit hash of the user's name.
func Getuid() (int, error) {
	u, err := user.Current()
	if err != nil {
		return UnknownID, err
	}

	uid := int(xxhash.Checksum32([]byte(u.Username)))

	return uid, nil
}

// Getgid returns the default Group ID for the current user.
// On Windows, the user's primary group's SID is converted to its POSIX
// equivalent, which is compatible with Cygwin, Git for Windows, MSYS2, etc.
// On Plan9, Getgid returns a 32-bit hash of the user's group's name, as
// provided by golang's os/user package.
func Getgid() (int, error) {
	u, err := user.Current()
	if err != nil {
		return UnknownID, err
	}

	gid := int(xxhash.Checksum32([]byte(u.Gid)))

	return gid, nil
}

// Geteuid returns the effective User ID for the current user.
// On Windows, the user's SID is converted to its POSIX equivalent, which is
// compatible with Cygwin, Git for Windows, MSYS2, etc.
// On Plan9, Geteuid returns a 32-bit hash of the user's name.
func Geteuid() (int, error) {
	return Getuid()
}

// Getegid returns the effective default Group ID for the current user.
// On Windows, the user's primary group's SID is converted to its POSIX
// equivalent, which is compatible with Cygwin, Git for Windows, MSYS2, etc.
// On Plan9, Getegid returns a 32-bit hash of the user's group's name, as
// provided by golang's os/user package.
func Getegid() (int, error) {
	return Getgid()
}
