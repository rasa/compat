// SPDX-FileCopyrightText: Copyright (c) 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"errors"
	"fmt"
	"os"
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

type UnimplementedError struct {
	Op string
}

func (e *UnimplementedError) Error() string {
	return fmt.Sprintf("%s: unimplemented on %s/%s", e.Op, runtime.GOOS, runtime.GOARCH)
}

func (e *UnimplementedError) Unwrap() error {
	return errors.ErrUnsupported
}

func unsupportedError(prefix string) error { //nolint:unused
	return &UnsupportedError{prefix}
}

func unimplementedError(prefix string) error { //nolint:unused
	return &UnimplementedError{prefix}
}

func chmodError(path string, err error) error { //nolint:unused
	return &os.PathError{Op: "chmod", Path: path, Err: err}
}

func createError(path string, err error) error { //nolint:unused
	return &os.PathError{Op: "create", Path: path, Err: err}
}

func createTempError(path string, err error) error { //nolint:unused
	return &os.PathError{Op: "createtemp", Path: path, Err: err}
}

func mkdirError(path string, err error) error { //nolint:unused
	return &os.PathError{Op: "mkdir", Path: path, Err: err}
}

func mkdirTempError(path string, err error) error { //nolint:unused
	return &os.PathError{Op: "mkdirall", Path: path, Err: err}
}

func openError(path string, err error) error { //nolint:unused
	return &os.PathError{Op: "open", Path: path, Err: err}
}

func renameError(old_, new_ string, err error) error { //nolint:unused
	return &os.LinkError{Op: "rename", Old: old_, New: new_, Err: err}
}

func statError(path string, err error) error { //nolint:unused
	return &os.PathError{Op: "stat", Path: path, Err: err}
}

func symlinkError(old_, new_ string, err error) error { //nolint:unused
	return &os.LinkError{Op: "symlink", Old: old_, New: new_, Err: err}
}

func writeError(name string, err error) error { //nolint:unused
	return &os.PathError{Op: "write", Path: name, Err: err}
}
