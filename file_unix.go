// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !windows

package compat

import (
	"errors"
	"os"

	"github.com/rasa/compat/golang"
)

func chmod(name string, mode os.FileMode, _ ReadOnlyMode) error {
	return os.Chmod(name, mode)
}

func create(name string, perm os.FileMode, flag int) (*os.File, error) {
	if perm == 0 {
		perm = CreatePerm
	}

	flag |= O_CREATE // & ^O_EXCL
	return openFile(name, flag, perm)
}

func createTemp(dir, pattern string, perm os.FileMode, flag int) (*os.File, error) {
	if perm == 0 {
		perm = CreateTempPerm
	}

	f, err := golang.CreateTemp(dir, pattern, perm)
	if err != nil {
		return nil, err
	}

	return wrap(f.Name(), flag, f)
}

func fchmod(f *os.File, mode os.FileMode, _ ReadOnlyMode) error {
	if f == nil {
		return errors.New("nil file pointer")
	}

	return f.Chmod(mode)
}

var mkdir = os.Mkdir

var mkdirAll = os.MkdirAll

func mkdirTemp(dir, pattern string, perm os.FileMode) (string, error) {
	if perm == 0 {
		perm = MkdirTempPerm
	}

	return golang.MkdirTemp(dir, pattern, perm)
}

func openFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	// don't pass compat-only flags to os function.
	oflag := flag & ^(O_FILE_FLAG_DELETE_ON_CLOSE | O_FILE_FLAG_NO_RO_ATTR)
	f, err := os.OpenFile(name, oflag, perm)
	if err != nil {
		return nil, err
	}

	return wrap(name, flag, f)
}

var remove = os.Remove

func removeAll(path string, _ ...Option) error {
	return os.RemoveAll(path)
}

func symlink(oldname, newname string, _ bool) error {
	return os.Symlink(oldname, newname)
}

func writeFile(name string, data []byte, perm os.FileMode, _ int) error {
	return os.WriteFile(name, data, perm)
}

func wrap(name string, flag int, f *os.File) (*os.File, error) {
	if flag&O_FILE_FLAG_DELETE_ON_CLOSE == 0 {
		return f, nil
	}

	err := os.Remove(name)
	if err == nil || os.IsNotExist(err) {
		return f, nil
	}
	if f != nil {
		_ = f.Close()
	}
	_ = os.Remove(name)

	return nil, err
}
