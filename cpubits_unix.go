// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !plan9 && !wasm && !windows

package compat

import (
	"strings"

	"golang.org/x/sys/unix"
)

const (
	bits32 = 32
	bits64 = 64
)

func cpuBits() (int, error) {
	var uts unix.Utsname

	err := unix.Uname(&uts)
	if err != nil {
		return 0, err
	}

	machine := make([]byte, len(uts.Machine))

	for idx, val := range uts.Machine {
		if val == 0 {
			machine = machine[:idx]

			break
		}

		machine[idx] = val
	}

	arch := strings.TrimSpace(string(machine))

	if strings.HasSuffix(arch, "64") {
		return bits64, nil
	}

	return bits32, nil
}
