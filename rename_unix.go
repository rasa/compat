// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !windows

package compat

import (
	"os"
)

// rename atomically replaces the destination file or directory with the
// source. It is guaranteed to either replace the target file entirely, or not
// change either file.
func rename(source, destination string, _ ...Option) error {
	return os.Rename(source, destination)
}
