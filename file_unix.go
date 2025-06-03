// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !windows

package compat

import (
	"os"
)

func chmod(name string, perm os.FileMode) error {
	return os.Chmod(name, perm)
}

func create(name string, perm os.FileMode, flag int) (*os.File, error) {
	f, err := os.Create(name)
	if perm != 0 {
		err = os.Chmod(name, perm)
		if err != nil {
			_ = f.Close()
			_ = os.Remove(name)
			return nil, err
		}
	}
	return wrap(name, flag, f, err)
}

func createTemp(dir, pattern string, flag int) (*os.File, error) {
	f, err := os.CreateTemp(dir, pattern)
	return wrap(f.Name(), flag, f, err)
}

func mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func mkdirAll(name string, perm os.FileMode) error {
	return os.MkdirAll(name, perm)
}

func mkdirTemp(dir, pattern string) (string, error) {
	return os.MkdirTemp(dir, pattern)
}

func openFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	return wrap(name, flag, f, err)
}

func writeFile(name string, data []byte, perm os.FileMode, _ int) error {
	return os.WriteFile(name, data, perm)
}

func wrap(name string, flag int, f *os.File, err error) (*os.File, error) {
	if err != nil {
		return nil, err
	}

	if flag&O_DELETE == O_DELETE {
		err = os.Remove(name)
		if err != nil {
			_ = f.Close()
			_ = os.Remove(name)
			return nil, err
		}
	}

	return f, nil
}
