// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

import (
	"os"
	"strconv"
	"strings"
)

func initUmask() error {
	var startingUmask uint64 = initialUmask
	umask := strings.TrimSpace(os.Getenv("UMASK"))
	umask = strings.TrimPrefix(strings.ToLower(umask), "0o")

	if umask != "" {
		intVal, err := strconv.ParseUint(umask, 8, 32)
		if err != nil {
			umaskErr = err
		} else {
			startingUmask = intVal & umaskMask
			validateUmask(intVal)
		}
	}

	savedUmask.Store(startingUmask)
	return umaskErr
}

// Umask sets the umask to umask, and returns the previous value.
// On Windows, the initial umask value is 022 octal, and can be changed by
// setting environmental variable UMASK, to an octal value. For example:
// `set UMASK=002`. A leading '0 and an 'o' are allowed, and ignored.
// On Plan9 and Wasip1, the function does nothing, and always returns zero.
func Umask(mask int) int {
	umaskMutex.Lock()
	defer umaskMutex.Unlock()

	umaskErr = nil
	oldUmask := savedUmask.Swap(uint64(mask) & umaskMask) //nolint:gosec
	return int(oldUmask)
}

// GetUmask returns the current process umask.
//
// compat keeps track of the umask set through [Umask]. On Unix-like systems,
// there is no system call for reading the umask without setting it, so
// GetUmask verifies the value by setting it to the value last known to compat.
//
// If the process umask was changed outside compat, GetUmask detects the change,
// immediately restores the externally set value, updates its cached value, and
// returns it.
//
// The umask is process-wide. Calls to Umask and GetUmask through compat are
// serialized, but compat cannot synchronize direct syscall.Umask calls or file
// creation performed concurrently by other code in the process.
//
// On Plan9 and Wasip1, the function always returns zero, and ErrUnsupported.
func GetUmask() (int, error) {
	umaskMutex.Lock()
	defer umaskMutex.Unlock()

	return int(savedUmask.Load()), umaskErr
}
