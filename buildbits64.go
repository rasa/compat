// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !(386 || arm || mips || mipsle)

package compat

import "math/bits"

// BuildBits returns the number of CPU bits for the build target.
// For 386, arm, mips, and mipsle, it's 32. For all other targets, it's 64.
func BuildBits() int {
	if IsWasip1Target {
		return bits.UintSize
	}
	return 64 //nolint:mnd // quiet linter
}
