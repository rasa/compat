// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

import (
	"fmt"
	"sync"

	"golang.org/x/sys/windows"
)

var getuidOnce struct {
	sync.Once
	uid int
	err error
}

// Getuid returns the User ID for the current user.
// On Windows, the user's SID is converted to its POSIX equivalent, which is
// compatible with Cygwin, Git for Windows, MSYS2, etc.
// On Plan9, Getuid returns a 32-bit hash of the user's name.
func Getuid() (int, error) {
	getuidOnce.Do(func() {
		getuidOnce.uid, getuidOnce.err = getuid()
	})
	return getuidOnce.uid, getuidOnce.err
}

func getuid() (int, error) {
	var token windows.Token
	err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)
	if err != nil {
		return UnknownID, fmt.Errorf("failed to open process token: %w", err)
	}
	defer token.Close()

	// Get the token's user
	tokenUser, err := token.GetTokenUser()
	if err != nil {
		return UnknownID, fmt.Errorf("failed to get token user: %w", err)
	}

	primaryDomainSID, err := getPrimaryDomainSID()
	if err != nil {
		return UnknownID, fmt.Errorf("failed to get primary domain SID: %w", err)
	}

	uid, err := sidToPOSIXID(tokenUser.User.Sid, primaryDomainSID)
	if err != nil {
		return UnknownID, fmt.Errorf("failed to convert SID to POSIX ID: %w", err)
	}

	return uid, nil
}

var getgidOnce struct {
	sync.Once
	gid int
	err error
}

// Getgid returns the default Group ID for the current user.
// On Windows, the user's primary group's SID is converted to its POSIX
// equivalent, which is compatible with Cygwin, Git for Windows, MSYS2, etc.
// On Plan9, Getgid returns a 32-bit hash of the user's group's name, as
// provided by golang's os/user package.
func Getgid() (int, error) {
	getgidOnce.Do(func() {
		getgidOnce.gid, getgidOnce.err = getgid()
	})
	return getgidOnce.gid, getgidOnce.err
}

func getgid() (int, error) {
	var token windows.Token
	err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)
	if err != nil {
		return UnknownID, fmt.Errorf("failed to open process token: %w", err)
	}
	defer token.Close()

	primaryDomainSID, err := getPrimaryDomainSID()
	if err != nil {
		return UnknownID, fmt.Errorf("failed to get primary domain SID: %w", err)
	}

	groupSID, err := getPrimaryGroupSID(token)
	if err != nil {
		return UnknownID, fmt.Errorf("failed to get primary group SID: %w", err)
	}

	gid, err := sidToPOSIXID(groupSID, primaryDomainSID)
	if err != nil {
		return UnknownID, fmt.Errorf("failed to convert SID to POSIX ID: %w", err)
	}

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
