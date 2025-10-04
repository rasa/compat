// SPDX-FileCopyrightText: Copyright (c) 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !unix

package compat_test

func getOSVersion() (ver semanticVersion, err error) { //nolint:gocyclo
	return ver, err
}

/*

import "github.com/shirou/gopsutil/v4/host"

func getOSVersion() (v ver, err error) { //nolint:gocyclo
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
