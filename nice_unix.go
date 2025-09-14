// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build android || (!ios && !linux && !plan9 && !wasm && !windows)

package compat

import (
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
	err := unix.Setpriority(unix.PRIO_PROCESS, 0, nice)
	if err != nil {
		return &ReniceError{nice, err}
	}

	return nil
}
