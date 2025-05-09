// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !windows

package compat

import (
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/unix"
)

// IsWSL returns true if running under Windows Subsytem for Linux (WSL),
// otherwise false.
func IsWSL() bool {
	var uts unix.Utsname
	err := unix.Uname(&uts)
	if err == nil {
		release := byteToString(uts.Release[:])
		if strings.Contains(strings.ToLower(release), "microsoft") {
			return true
		}
	}
	data, err := os.ReadFile("/proc/sys/kernel/osrelease")
	if err == nil {
		return strings.Contains(strings.ToLower(string(data)), "microsoft")
	}
	data, err = os.ReadFile("/proc/version")
	if err == nil {
		return strings.Contains(strings.ToLower(string(data)), "microsoft")
	}
	path, err := exec.LookPath("wslpath")
	if err != nil {
		return false
	}
	return path == "/usr/bin/wslpath"
}

// Convert byte array to string.
func byteToString(b []byte) string {
	n := len(b)
	for i := 0; i < n; i++ {
		if b[i] == 0 {
			n = i
			break
		}
	}
	return string(b[:n])
}
