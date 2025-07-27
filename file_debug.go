// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !debug || !windows

package compat

import (
	"os"
)

func dumpMasks(_ os.FileMode, _ uint32, _ uint32, _ uint32) { //nolint:nolintlint,unused // quiet linter
}
