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
)

// NiceError is returned when the system call failed.
type NiceError struct {
	Err error
}

func (e *NiceError) Error() string {
	return fmt.Errorf("nice: %w", e.Err).Error()
}

// InvalidNiceError is returned when the niceness value passed by the user is
// invalid.
type InvalidNiceError struct {
	Nice int
}

func (e *InvalidNiceError) Error() string {
	return fmt.Sprintf("nice: invalid nice value %d", e.Nice)
}

// ReniceError is returned when the system failed to set the OS's niceness level.
type ReniceError struct {
	Nice int
	Err  error
}

func (e *ReniceError) Error() string {
	return fmt.Errorf("nice: failed to set nice to %d: %w", e.Nice, e.Err).Error()
}
