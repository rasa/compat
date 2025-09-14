// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build linux && !android

package compat

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// Nice gets the CPU process priority. The return value is in a range from
// -20 (least nice), to 19 (most nice), even on non-Unix systems such as
// Windows, plan9, etc. If not supported by the operating system, an error is
// returned.
func Nice() (int, error) {
	nice, err := unix.Getpriority(unix.PRIO_PROCESS, 0)
	if err != nil {
		return 0, &NiceError{err}
	}

	return nice, nil
}

// Renice sets the CPU process priority. The nice parameter can range from
// -20 (least nice), to 19 (most nice), even on non-Unix systems such as
// Windows, plan9, etc.
func Renice(nice int) error {
	// Move ourselves to a new process group so that we can use the process
	// group variants of Setpriority etc to affect all of our threads in one
	// go. If this fails, bail, so that we don't affect things we shouldn't.
	// If we are already the leader of our own process group, do nothing.
	//
	// Oh and this is because Linux doesn't follow the POSIX threading model
	// where setting the niceness of the process would actually set the
	// niceness of the process, instead it just affects the current thread
	// so we need this workaround...
	pgid, err := unix.Getpgid(0)
	if err != nil {
		// This error really shouldn't happen
		return fmt.Errorf("nice: get process group: %w", err)
	}

	if pgid != os.Getpid() {
		// We are not process group leader. Elevate!
		err = unix.Setpgid(0, 0)
		if err != nil {
			return fmt.Errorf("nice: set process group: %w", err)
		}
	}

	err = unix.Setpriority(unix.PRIO_PROCESS, 0, nice)
	if err != nil {
		return &ReniceError{nice, err}
	}

	return nil
}
