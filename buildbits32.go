// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build 386 || arm || mips || mipsle || tinygo

package compat

// BuildBits returns the number of CPU bits for the build target.
// For 386, arm, mips, and mipsle, it's 32. For all other targets, it's 64.
func BuildBits() int {
	return 32 //nolint:mnd // quiet linter
}
