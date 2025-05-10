// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9 || wasm

package compat

import (
	"os"
	"os/exec"
	"strings"
)

// IsWSL returns true if running under Windows Subsytem for Linux (WSL),
// otherwise false.
func IsWSL() bool {
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
