// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

import (
	"errors"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Getuid returns the User ID for the current user.
// On Windows, the user's SID is converted to its POSIX equivalent, which is
// compatible with Cygwin, Git for Windows, MSYS2, etc.
// On Plan9, Getuid returns a 32-bit hash of the user's name.
func Getuid() (int, error) {
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

// Getgid returns the default Group ID for the current user.
// On Windows, the user's primary group's SID is converted to its POSIX
// equivalent, which is compatible with Cygwin, Git for Windows, MSYS2, etc.
// On Plan9, Getgid returns a 32-bit hash of the user's group's name, as
// provided by golang's os/user package.
func Getgid() (int, error) {
	primaryDomainSID, err := getPrimaryDomainSID()
	if err != nil {
		return UnknownID, fmt.Errorf("failed to get primary domain SID: %w", err)
	}

	groupSID, err := getPrimaryGroupSID()
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

// @TODO(rasa) improve this logic per
// https://github.com/golang/go/blob/cc8a6780/src/os/user/lookup_windows.go#L351
func getPrimaryGroupSID() (*windows.SID, error) {
	var token windows.Token
	err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)
	if err != nil {
		return nil, fmt.Errorf("failed to open process token: %w", err)
	}
	defer token.Close()

	bufSize := initialBufSize
	for {
		var newBufSize uint32
		buf := make([]byte, bufSize)
		if bufSize <= 1 {
			err = windows.GetTokenInformation(
				token,
				windows.TokenPrimaryGroup,
				nil,
				0,
				&newBufSize)
		} else {
			err = windows.GetTokenInformation(
				token,
				windows.TokenPrimaryGroup,
				&buf[0],
				bufSize,
				&newBufSize)
		}
		if err == nil {
			pg := (*windows.Tokenprimarygroup)(unsafe.Pointer(&buf[0]))
			return pg.PrimaryGroup, nil
		}
		if !errors.Is(err, windows.ERROR_INSUFFICIENT_BUFFER) {
			return nil, fmt.Errorf("failed to get token information: %w", err)
		}
		if newBufSize > bufSize {
			bufSize = newBufSize
		} else {
			bufSize *= 2
		}
	}
}
