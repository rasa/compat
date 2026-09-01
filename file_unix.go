// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !windows

package compat

import (
	"os"

	"github.com/rasa/compat/golang"
)

func chmod(name string, mode os.FileMode, _ options) error {
	return os.Chmod(name, mode)
}

func create(name string, opts options) (*os.File, error) {
	return openFile(name, opts)
}

func createTemp(dir, pattern string, opts options) (*os.File, error) {
	fp, err := golang.CreateTemp(dir, pattern, opts.fileMode)
	if err != nil {
		return nil, err
	}

	if opts.deleteOnClose {
		return deleteOnClose(fp.Name(), fp)
	}

	return fp, nil
}

func fchmod(file *os.File, mode os.FileMode, _ options) error {
	if file == nil {
		return chmodError("", os.ErrInvalid)
	}

	return file.Chmod(mode)
}

var mkdir = os.Mkdir

var mkdirAll = os.MkdirAll

func mkdirTemp(dir, pattern string, opts options) (string, error) {
	return golang.MkdirTemp(dir, pattern, opts.fileMode)
}

func openFile(name string, opts options) (*os.File, error) {
	fp, err := os.OpenFile(name, opts.flags, opts.fileMode)
	if err != nil {
		return nil, err
	}

	if opts.deleteOnClose {
		return deleteOnClose(fp.Name(), fp)
	}

	return fp, nil
}

var remove = os.Remove

func removeAll(path string, _ options) error {
	return os.RemoveAll(path)
}

func symlink(oldname, newname string, _ options) error {
	return os.Symlink(oldname, newname)
}

func deleteOnClose(name string, file *os.File) (*os.File, error) {
	err := os.Remove(name)
	if err == nil || os.IsNotExist(err) {
		return file, nil
	}

	if file != nil {
		_ = file.Close()
	}

	_ = os.Remove(name)

	return nil, err
}
