// SPDX-FileCopyrightText: Copyright (c) 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !plan9

package compat

import (
	"errors"
	"syscall"
)

// IsUnsupportedError returns true if err indicates the function is unsupported.
func IsUnsupportedError(err error) bool {
	return errors.Is(err, errors.ErrUnsupported) ||
		errors.Is(err, syscall.ENOTSUP) ||
		errors.Is(err, syscall.EOPNOTSUPP)
}
