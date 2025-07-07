// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build aix || darwin || dragonfly || freebsd || illumos || js || linux || netbsd || openbsd || solaris || (wasip1 && !tinygo)

package compat

import (
	"sync"
	"syscall"
)

var (
	currentUmask = 0o022
	umaskMutex   sync.Mutex
)

func init() {
	_ = GetUmask()
}

// Umask sets the umask to umask, and returns the previous value.
// On Windows, the initial umask value is 022 octal, and can be changed by
// setting environmental variable UMASK, to an octal value. For example:
// `set UMASK=002`. Leading zeros and 'o's are allowed, and ignored.
// On Plan9 and Wasip1, the function does nothing, and always returns zero.
func Umask(newMask int) int {
	umaskMutex.Lock()
	defer umaskMutex.Unlock()

	currentUmask = newMask
	return syscall.Umask(currentUmask)
}

// GetUmask returns the current umask value.
// On Plan9 and Wasip1, the function always returns zero.
func GetUmask() int {
	umaskMutex.Lock()
	defer umaskMutex.Unlock()

	now := syscall.Umask(currentUmask)
	if now != currentUmask {
		_ = syscall.Umask(now)
		currentUmask = now
	}

	return now
}
