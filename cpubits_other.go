// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build js || plan9 || wasip1

package compat

// CPUBits returns the number of bits on the CPU. Currently, on plan9, and wasm,
// BuildBits() is returned.
func CPUBits() (int, error) {
	return BuildBits(), nil
}
