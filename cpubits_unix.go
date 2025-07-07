// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !plan9 && !wasm && !windows

package compat

import (
	"strings"

	"golang.org/x/sys/unix"
)

// CPUBits returns the number of bits on the CPU. Currently, on plan9, and wasm,
// zero is returned.
func CPUBits() (int, error) {
	var uts unix.Utsname
	err := unix.Uname(&uts)
	if err != nil {
		return 0, err
	}

	machine := make([]byte, len(uts.Machine))
	for i, v := range uts.Machine {
		if v == 0 {
			machine = machine[:i]
			break
		}
		machine[i] = v
	}
	arch := strings.TrimSpace(string(machine))

	if strings.HasSuffix(arch, "64") {
		return 64, nil //nolint:mnd // quiet linter
	}

	return 32, nil //nolint:mnd // quiet linter
}
