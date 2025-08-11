// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// Getuid returns the User ID as a uint64. On Windows, the user's SID is
// converted to it's POSIX equivalent, which is compatiable with Cygwin and
// Git for Windows. On Plan9, the User ID is a 64-bit hash of the user's name.
func Getuid() (uint64, error) {
	var token windows.Token
	err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)
	if err != nil {
		return UnknownID, err
	}
	defer token.Close()

	// Get the token's user
	tokenUser, err := token.GetTokenUser()
	if err != nil {
		return UnknownID, err
	}

	primaryDomainSID, err := getPrimaryDomainSID()
	if err != nil {
		return UnknownID, err
	}

	uid, err := sidToPOSIXID(tokenUser.User.Sid, primaryDomainSID)
	if err != nil {
		return UnknownID, err
	}

	return uint64(uid), nil
}

// Getgid returns the Group ID as a uint64. On Windows, the user's primary group's
// SID is converted to its POSIX equivalent, which is compatiable with Cygwin and
// Git for Windows. On Plan9, the Getgid returns the value returned by Getuid().
func Getgid() (uint64, error) {
	primaryDomainSID, err := getPrimaryDomainSID()
	if err != nil {
		return UnknownID, err
	}

	groupSID, err := getPrimaryGroupSID()
	if err != nil {
		return UnknownID, err
	}

	gid, err := sidToPOSIXID(groupSID, primaryDomainSID)
	if err != nil {
		return UnknownID, err
	}

	return uint64(gid), nil
}

func getPrimaryGroupSID() (*windows.SID, error) {
	var token windows.Token
	err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)
	if err != nil {
		return nil, err
	}
	defer token.Close()

	// Get size for TOKEN_PRIMARY_GROUP
	var size uint32
	err = windows.GetTokenInformation(token, windows.TokenPrimaryGroup, nil, 0, &size)
	if err != windows.ERROR_INSUFFICIENT_BUFFER {
		return nil, err
	}

	buf := make([]byte, size)
	err = windows.GetTokenInformation(token, windows.TokenPrimaryGroup, &buf[0], size, &size)
	if err != nil {
		return nil, err
	}

	tpg := (*tokenPrimaryGroup)(unsafe.Pointer(&buf[0]))

	return tpg.PrimaryGroup, nil
}
