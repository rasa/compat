// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"sync"
)

var cpuBitsOnce struct {
	sync.Once
	cpuBits int
	err     error
}

// CPUBits returns the number of bits on the CPU. Currently, on plan9, and wasm,
// BuildBits() is returned.
func CPUBits() (int, error) {
	cpuBitsOnce.Do(func() {
		cpuBitsOnce.cpuBits, cpuBitsOnce.err = cpuBits()
	})
	return cpuBitsOnce.cpuBits, cpuBitsOnce.err
}
