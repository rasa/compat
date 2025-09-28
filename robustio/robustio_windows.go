// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// See also https://github.com/golang/go/blob/34e67623/src/testing/testing_windows.go#L17

// SPDX-FileCopyrightText: Copyright 2019 The Go Authors.
// SPDX-License-Identifier: BSD-3

// Source: https://github.com/golang/go/blob/f15cd63e/src/cmd/internal/robustio/robustio_windows.go

// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package robustio

import (
	"errors"
	// "internal/syscall/windows" // compat: s|"internal|// "internal|.
	"syscall"

	"golang.org/x/sys/windows" // compat: added
)

const errFileNotFound = syscall.ERROR_FILE_NOT_FOUND

// isEphemeralError returns true if err may be resolved by waiting.
func isEphemeralError(err error) bool {
	var errno syscall.Errno
	if errors.As(err, &errno) {
		switch errno {
		case syscall.ERROR_ACCESS_DENIED,
			syscall.ERROR_FILE_NOT_FOUND,
			windows.ERROR_SHARING_VIOLATION:
			return true
		}
	}
	return false
}
