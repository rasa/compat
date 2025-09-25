// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: Copyright 2019 The Go Authors.
// SPDX-License-Identifier: BSD-3

// Source: https://github.com/golang/go/blob/f15cd63e/src/cmd/internal/robustio/robustio_windows.go

// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !windows && !darwin

package robustio

import (
	"os"
)

func rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func readFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func removeAll(path string) error {
	return os.RemoveAll(path)
}

func isEphemeralError(err error) bool {
	return false
}
