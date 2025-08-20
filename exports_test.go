// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"os"
)

// no longer used:
// func ExportChmod(name string, perm os.FileMode, mode ReadOnlyMode) error {
// 	return chmod(name, perm, mode)
// }

func ExportStat(name string) (os.FileMode, error) {
	return _stat(name)
}
