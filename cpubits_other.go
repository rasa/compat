// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build js || plan9 || wasip1

package compat

import (
	"errors"
	"runtime"
)

// CPUBits returns the number of bits in an integer on the CPU. Currently, on
// plan9, and wasm, zero is returned.
func CPUBits() (int, error) {
	return 0, errors.New("Not supported on " + runtime.GOOS + "/" + runtime.GOARCH)
}
