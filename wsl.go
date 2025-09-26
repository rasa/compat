// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"os"
	"os/exec"
	"strings"
)

func iswsl() bool {
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
