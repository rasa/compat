// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build ios || wasm

package compat

import (
	"errors"
	"runtime"
)

// Nice gets the CPU process priority. The return value is in a range from
// -20 (least nice), to 19 (most nice), even on non-Unix systems such as
// Windows, plan9, etc. If not supported by the operating system, an error
// is returned.
func Nice() (int, error) {
	return 0, errors.New("nice: function not supported on " + runtime.GOOS)
}

// Renice sets the CPU process priority. The nice parameter can range from
// -20 (least nice), to 19 (most nice), even on non-Unix systems such as
// Windows, plan9, etc.
func Renice(_ int) error {
	return nil
}
