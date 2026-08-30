// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

const (
	// MaxNice is the maximum value returned by Nice().
	MaxNice = 19
	// MinNice is the minimum value returned by Nice().
	MinNice = -20
)

func validateNice(nice int) error {
	if nice < MinNice || nice > MaxNice {
		return invalidNiceError(nice)
	}

	return nil
}
