// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build js || unix

// https://github.com/golang/go/blob/8ad27fb6/src/cmd/dist/build.go#L1070
// unix == aix || android || darwin || dragonfly || freebsd || illumos || ios || linux || netbsd || openbsd || solaris

package compat

import (
	"fmt"
	"syscall"
)

func initUmask() error {
	umaskMutex.Lock()
	umask := syscall.Umask(int(initialUmask))
	_ = syscall.Umask(umask)
	umaskMutex.Unlock()
	savedUmask.Store(uint64(umask))

	return nil
}

// Umask sets the umask to mask, and returns the previous value.
//
// On Windows, the initial umask value is 0o022 octal, and can be changed by
// setting the environmental variable UMASK, to an octal value. For example:
// `set UMASK=002` or `set UMASK=0o002`. A leading '0 and an 'o' are allowed,
// and ignored. If the value is greater than 0o777, only the lower 9 bits will
// be used.
//
// On Plan9 and Wasip1, the function does nothing, and always returns zero.
func Umask(mask int) int {
	savedUmask.Store(uint64(mask) & umaskMask)
	umaskMutex.Lock()
	defer umaskMutex.Unlock()

	return syscall.Umask(mask)
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
// On Unix-like systems, if GetMask finds that the system's umask value was
// changed outside of compat, GetMask will return the current umask value, and
// [ErrUmaskChanged].
//
// On Windows, if the UMASK environment variable is not a valid octal number,
// or greater than 0o777, GetUmask will return the current umask value, and
// [ErrInvalidUmask].
//
// On Plan9 and Wasip1, the function always returns zero, and [os.ErrUnsupported].
func GetUmask() (int, error) {
	saved := int(savedUmask.Load())

	umaskMutex.Lock()
	current := syscall.Umask(saved)
	last := syscall.Umask(current)
	umaskMutex.Unlock()
	savedUmask.Store(uint64(current))

	var err error

	if current != saved || last != saved {
		err = fmt.Errorf(
			"umask changed from 0o%03o to 0o%03o outside compat: %w",
			saved,
			current,
			ErrUmaskChanged,
		)
	}
	return current, err
}
