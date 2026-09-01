// SPDX-FileCopyrightText: Copyright (c) 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"errors"
	"fmt"
	"os"
	"runtime"
)

var ErrInvalidNice = errors.New("invalid nice value")

func unsupportedError(op string) error { //nolint:unused
	return fmt.Errorf("%s: %w", op, errors.ErrUnsupported)
}

func unimplementedError(op string) error { //nolint:unused
	return fmt.Errorf(
		"%s: unimplemented on %s/%s: %w",
		op,
		runtime.GOOS,
		runtime.GOARCH,
		errors.ErrUnsupported,
	)
}

func invalidNiceError(nice int) error {
	return fmt.Errorf(
		"value %d outside range %d..%d: %w",
		nice,
		MinNice,
		MaxNice,
		ErrInvalidNice,
	)
}

func unexpectedNiceError(nice int) error { //nolint:unused
	return fmt.Errorf(
		"BUG: value %d is unexpected: %w",
		nice,
		ErrInvalidNice,
	)
}

func chmodError(path string, err error) error {
	return &os.PathError{Op: "chmod", Path: path, Err: err}
}

func createError(path string, err error) error {
	return &os.PathError{Op: "create", Path: path, Err: err}
}

func createTempError(path string, err error) error {
	return &os.PathError{Op: "createtemp", Path: path, Err: err}
}

func mkdirError(path string, err error) error {
	return &os.PathError{Op: "mkdir", Path: path, Err: err}
}

func mkdirallError(path string, err error) error {
	return &os.PathError{Op: "mkdir", Path: path, Err: err}
}

func mkdirTempError(path string, err error) error {
	return &os.PathError{Op: "mkdirtemp", Path: path, Err: err}
}

func niceError(msg string, err error) error {
	return fmt.Errorf("nice: %v: %w", msg, err)
}

func openError(path string, err error) error {
	return &os.PathError{Op: "open", Path: path, Err: err}
}

func renameError(old, gnu string, err error) error {
	return &os.LinkError{Op: "rename", Old: old, New: gnu, Err: err}
}

func statError(path string, err error) error {
	return &os.PathError{Op: "stat", Path: path, Err: err}
}

func symlinkError(old, gnu string, err error) error {
	return &os.LinkError{Op: "symlink", Old: old, New: gnu, Err: err}
}

func writeError(name string, err error) error {
	return &os.PathError{Op: "write", Path: name, Err: err}
}
