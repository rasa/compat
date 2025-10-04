// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !windows

package compat

import (
	"os"
	"os/exec"
	"strings"
	"sync"
)

var isWSLOnce struct {
	sync.Once
	isWSL bool
}

func isWSL() bool {
	data, err := os.ReadFile("/proc/sys/kernel/osrelease")
	if err != nil {
		data, err = os.ReadFile("/proc/version")
	}
	if err == nil {
		return strings.Contains(strings.ToLower(string(data)), "microsoft")
	}

	_, err = exec.LookPath("wslpath")
	return err == nil
}
