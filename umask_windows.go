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

func initUmask() error {
	umask := strings.TrimSpace(os.Getenv("UMASK"))

	umask = strings.TrimPrefix(strings.ToLower(umask), "0o")
	if umask != "" {
		n, err := strconv.ParseUint(umask, 8, 32)
		if err != nil {
			return err
		}

		startingUmask = uint32(n) & permMask
	}

	currentUmask.Store(startingUmask)

	return nil
}

// Umask sets the umask to umask, and returns the previous value.
// On Windows, the initial umask value is 022 octal, and can be changed by
// setting environmental variable UMASK, to an octal value. For example:
// `set UMASK=002`. Leading zeros and 'o's are allowed, and ignored.
// On Plan9 and Wasip1, the function does nothing, and always returns zero.
func Umask(newMask int) int {
	old := currentUmask.Swap(uint32(newMask) & permMask) //nolint:gosec

	return int(old)
}

// GetUmask returns the current umask value.
//
// On Plan9 and Wasip1, the function always returns zero.
func GetUmask() int {
	return int(currentUmask.Load())
}
