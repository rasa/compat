// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9 || wasm

package compat

import (
	"errors"
	"runtime"
)

func CPUBits() (int, error) {
	return 0, errors.New("Not supported on " + runtime.GOOS + "/" + runtime.GOARCH)
}
