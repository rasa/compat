// SPDX-FileCopyrightText: Copyright (c) 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"errors"
	"fmt"
	"runtime"
)

type UnsupportedError struct {
	Op string
}

func (e *UnsupportedError) Error() string {
	return e.Op + ": unsupported"
}

func (e *UnsupportedError) Unwrap() error {
	return errors.ErrUnsupported
}

type NotYetImplementedError struct {
	Op string
}

func (e *NotYetImplementedError) Error() string {
	return fmt.Sprintf("%s: not yet implemented on %s/%s", e.Op, runtime.GOOS, runtime.GOARCH)
}

func (e *NotYetImplementedError) Unwrap() error {
	return errors.ErrUnsupported
}
