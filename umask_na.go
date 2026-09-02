// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9 || wasip1

package compat

import (
	"errors"
)

func initUmask() error {
	return nil
}

// Umask sets the umask to umask, and returns the previous value.
// On Windows, the initial umask value is 022 octal, and can be changed by
// setting environmental variable UMASK, to an octal value. For example:
// `set UMASK=002`. A leading '0 and an 'o' are allowed, and ignored.
// On Plan9 and Wasip1, the function does nothing, and always returns zero.
func Umask(_ int) int {
	return 0
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
	return 0, errors.ErrUnsupported
}
