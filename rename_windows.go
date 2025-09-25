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

	"github.com/rasa/compat/golang"
	"github.com/rasa/compat/robustio"
)

func rename(src, dst string, opts ...Option) error {
	fopts := Options{}

	for _, opt := range opts {
		opt(&fopts)
	}

	if fopts.retrySeconds <= 0 {
		return moveFile(src, dst)
	}

	return robustio.Retry(func() (err error, mayRetry bool) {
		err = moveFile(src, dst)
		return err, robustio.IsEphemeralError(err)
	}, fopts.retrySeconds)
}

func moveFile(src, dst string) error {
	longsrc := golang.FixLongPath(src)

	src16, err := syscall.UTF16PtrFromString(longsrc)
	if err != nil {
		return &os.LinkError{Op: "rename", Old: src, New: dst, Err: err}
	}
	longdst := golang.FixLongPath(dst)
	dst16, err := syscall.UTF16PtrFromString(longdst)
	if err != nil {
		return &os.LinkError{Op: "rename", Old: src, New: dst, Err: err}
	}

	var attrs uint32 = windows.MOVEFILE_REPLACE_EXISTING | windows.MOVEFILE_WRITE_THROUGH
	// see http://msdn.microsoft.com/en-us/library/windows/desktop/aa365240(v=vs.85).aspx
	err = windows.MoveFileEx(src16, dst16, attrs)
	if err != nil {
		return &os.LinkError{Op: "rename", Old: src, New: dst, Err: err}
	}
	return nil
}
