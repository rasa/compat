// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

import (
	"os"
	"strconv"
	"strings"
	"sync/atomic"
)

var (
	// Default umask on *nix: remove write for group and others.
	startingUmask uint32 = 0o022
	currentUmask  atomic.Uint32
	// These are all the bits we care about on Windows (for now?).
	permMask uint32 = 0o777
)

func init() {
	umask := os.Getenv("UMASK")
	if umask != "" {
		umask = strings.Trim(umask, " \t")
		umask = strings.TrimLeft(umask, "0")
		umask = strings.Trim(umask, "o")
		umask = strings.TrimLeft(umask, "0")
		ui64, err := strconv.ParseInt(umask, 8, 32)
		if err == nil {
			// ignore errors
			startingUmask = uint32(ui64) & permMask //nolint:gosec // quiet linter
		}
	}
	currentUmask.Store(startingUmask)
}

// Umask sets the umask to umask, and returns the previous value.
// On Windows, the initial umask value is 022 octal, and can be changed by
// setting environmental variable UMASK, to an octal value. For example:
// `set UMASK=002`. Leading zeros and 'o's are allowed, and ignored.
// On Plan9 and Wasip1, the function does nothing, and always returns zero.
func Umask(newMask int) int {
	old := currentUmask.Swap(uint32(newMask) & permMask) //nolint:gosec // quiet linter
	return int(old)
}

// GetUmask returns the current umask value.
// On Plan9 and Wasip1, the function always returns zero.
func GetUmask() int {
	return int(currentUmask.Load())
}
