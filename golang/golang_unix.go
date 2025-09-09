//nolint:all
// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !windows

package golang

// SPDX-FileCopyrightText: Copyright 2010 The Go Authors. All rights reserved.
// SPDX-License-Identifier: BSD-3

// The following code is:
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Snippet: https://github.com/golang/go/blob/77f911e3/src/os/tempfile.go#L35-L58

func CreateTemp(dir, pattern string, perm FileMode) (*File, error) { // compat: added: perm FileMode
	if dir == "" {
		dir = TempDir()
	}

	prefix, suffix, err := prefixAndSuffix(pattern)
	if err != nil {
		return nil, &PathError{Op: "createtemp", Path: pattern, Err: err}
	}
	prefix = joinPath(dir, prefix)

	try := 0
	for {
		name := prefix + nextRandom() + suffix
		f, err := OpenFile(name, O_RDWR|O_CREATE|O_EXCL, perm) // compat: changed 0600 to perm
		if IsExist(err) {
			if try++; try < 10000 {
				continue
			}
			return nil, &PathError{Op: "createtemp", Path: prefix + "*" + suffix, Err: ErrExist}
		}
		return f, err
	}
}

// Snippet: https://github.com/golang/go/blob/77f911e3/src/os/tempfile.go#L86-L117

func MkdirTemp(dir, pattern string, perm FileMode) (string, error) { // compat: added: perm FileMode
	if dir == "" {
		dir = TempDir()
	}

	prefix, suffix, err := prefixAndSuffix(pattern)
	if err != nil {
		return "", &PathError{Op: "mkdirtemp", Path: pattern, Err: err}
	}
	prefix = joinPath(dir, prefix)

	try := 0
	for {
		name := prefix + nextRandom() + suffix
		err := Mkdir(name, perm) // compat: changed 0700 to perm
		if err == nil {
			return name, nil
		}
		if IsExist(err) {
			if try++; try < 10000 {
				continue
			}
			return "", &PathError{Op: "mkdirtemp", Path: dir + string(PathSeparator) + prefix + "*" + suffix, Err: ErrExist}
		}
		if IsNotExist(err) {
			if _, err := Stat(dir); IsNotExist(err) {
				return "", err
			}
		}
		return "", err
	}
}
