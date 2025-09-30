// SPDX-FileCopyrightText: Copyright (c) 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build unix

package compat_test

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/sys/unix"
)

var uname struct {
	sync.Once
	macOSMajor int
}

func getMacOSMajor() int {
	uname.Once.Do(func() {
		var u unix.Utsname
		err := unix.Uname(&u)
		if err != nil {
			return
		}
		rel := unix.ByteSliceToString(u.Release[:])
		ver := strings.Split(rel, ".")
		maj, _ := strconv.Atoi(ver[0])
		uname.macOSMajor = maj
	})
	return uname.macOSMajor
}

// See https://en.wikipedia.org/wiki/MacOS_version_history#Overview
var macOSMap = map[int]ver{
	20: {11, 0, 0},
	21: {12, 0, 0},
	22: {13, 0, 0},
	23: {14, 0, 0},
	24: {15, 0, 0},
	25: {26, 0, 0},
}

func osVersion() (v ver, err error) { //nolint:gocyclo
	major := getMacOSMajor()

	val, ok := macOSMap[major]
	if ok {
		return val, nil
	}

	var u unix.Utsname
	err = unix.Uname(&u)
	if err != nil {
		return v, err
	}
	rel := unix.ByteSliceToString(u.Release[:])

	return v, fmt.Errorf("cannot parse '%v'", rel)
}

/*
func osVersion() (v ver, err error) { //nolint:gocyclo
	info, err := host.Info()
	if err != nil {
		return v, err
	}

	parts := strings.Split(info.PlatformVersion, ".")
	if len(parts) == 0 {
		return v, fmt.Errorf("unable to parse %q", info.PlatformVersion)
	}

	var major, minor, patch int

	major, _ = strconv.Atoi(parts[0])
	if len(parts) > 1 {
		minor, _ = strconv.Atoi(parts[1])
	}
	if len(parts) > 2 {
		patch, _ = strconv.Atoi(parts[2])
	}

	if compat.IsApple {
		val, ok := macOSMap[major]
		if ok {
			return val, nil
		}

		return v, fmt.Errorf("unknown MacOS version %q", info.PlatformVersion)
	}

	if compat.IsWindows {
		if len(parts) < 3 {
			return v, fmt.Errorf("unable to parse %q", info.PlatformVersion)
		}

		switch {
		case major == 6 && minor == 1:
			return ver{7, 0, 0}, nil
		case major == 6 && (minor == 2 || minor == 3):
			return ver{8, 0, 0}, nil
		case major == 10 && patch < 22000:
			return ver{10, 0, 0}, nil
		case major == 10 && patch >= 22000:
			return ver{11, 0, 0}, nil
		case major >= 11:
			return ver{major, minor, patch}, nil
		default:
			return v, fmt.Errorf("unknown Windows OS version %q", info.PlatformVersion)
		}
	}

	return ver{major, minor, patch}, nil
}
*/
