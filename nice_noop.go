// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build ios || wasm

package compat

// Nice gets the CPU process priority. The return value is in a range from
// -20 (least nice), to 19 (most nice), even on non-Unix systems such as
// Windows, plan9, etc. If not supported by the operating system, 0 is returned.
func Nice() (int, error) {
	return 0, nil
}

// Renice sets the CPU process priority. The nice parameter can range from
// -20 (least nice), to 19 (most nice), even on non-Unix systems such as
// Windows, plan9, etc. If not supported by the operating system, nil is returned.
func Renice(_ int) error {
	return nil
}
