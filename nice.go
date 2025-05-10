// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"fmt"
)

const (
	// MaxNice is the maximum value returned by Nice().
	MaxNice = 19
	// MinNice is the minimum value returned by Nice().
	MinNice = -20

	myPID = 0
)

// NiceError is returned when system call failed.
type NiceError struct {
	err error
}

func (e *NiceError) Error() string {
	return fmt.Errorf("nice: %w", e.err).Error()
}

// InvalidNiceError is returned when the niceness value passed by the user is
// invalid.
type InvalidNiceError struct {
	nice int
}

func (e *InvalidNiceError) Error() string {
	return fmt.Sprintf("nice: invalid nice value %d", e.nice)
}

// ReniceError is returned when the system failed to set the OS's niceness level.
type ReniceError struct {
	nice int
	err  error
}

func (e *ReniceError) Error() string {
	return fmt.Errorf("nice: failed to set nice to %d: %w", e.nice, e.err).Error()
}
