// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

import (
	"errors"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func CPUBits() (int, error) {
	mod := syscall.NewLazyDLL("kernel32.dll")
	proc := mod.NewProc("IsWow64Process")
	var isWow64 uint32
	handle := windows.CurrentProcess()
	r1, _, err := proc.Call(uintptr(handle), uintptr(unsafe.Pointer(&isWow64)))

	if r1 == 0 {
		if err != syscall.Errno(0) {
			return 0, err
		}

		return 0, errors.New("IsWow64Process call failed")
	}
	if isWow64 != 0 {
		return 64, nil
	}

	return 32, nil
}
