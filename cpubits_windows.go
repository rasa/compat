// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

import (
	"errors"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modkernel32        = windows.NewLazySystemDLL("kernel32.dll")
	procIsWow64Process = modkernel32.NewProc("IsWow64Process")
)

func cpuBits() (int, error) {
	var isWow64 uint32
	handle := windows.CurrentProcess()
	r1, _, err := procIsWow64Process.Call(uintptr(handle), uintptr(unsafe.Pointer(&isWow64)))

	if r1 == 0 {
		if !errors.Is(err, windows.Errno(0)) {
			return 0, err
		}

		return 0, errors.New("IsWow64Process call failed")
	}
	if isWow64 != 0 {
		return 64, nil //nolint:mnd
	}

	return 32, nil //nolint:mnd
}
