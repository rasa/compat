// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"os"
)

func ExportChmod(name string, perm os.FileMode) error {
	return chmod(name, perm)
}

func ExportStat(name string) (os.FileMode, error) {
	return _stat(name)
}
