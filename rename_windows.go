// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: Portions copyright (c) 2015 Nate Finch (@natefinch)
// SPDX-FileCopyrightText: Portions copyright (c) 2022 Simon Dassow (@sdassow)

//go:build windows

package compat

import (
	"os"
	"syscall"

	"golang.org/x/sys/windows"
)

func rename(source, destination string) error {
	sourcep := fixLongPath(source)

	src, err := syscall.UTF16PtrFromString(sourcep)
	if err != nil {
		return &os.LinkError{Op: "rename", Old: source, New: destination, Err: err}
	}
	destinationp := fixLongPath(destination)
	dest, err := syscall.UTF16PtrFromString(destinationp)
	if err != nil {
		return &os.LinkError{Op: "rename", Old: source, New: destination, Err: err}
	}

	var attrs uint32 = windows.MOVEFILE_REPLACE_EXISTING | windows.MOVEFILE_WRITE_THROUGH
	// see http://msdn.microsoft.com/en-us/library/windows/desktop/aa365240(v=vs.85).aspx
	if err := windows.MoveFileEx(src, dest, attrs); err != nil {
		return &os.LinkError{Op: "rename", Old: source, New: destination, Err: err}
	}

	return nil
}
