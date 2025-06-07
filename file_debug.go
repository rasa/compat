// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !debug || !windows

package compat

import (
	"os"
)

func dumpMasks(perm os.FileMode, ownerMask uint32, groupMask uint32, worldMask uint32) { //nolint:nolintlint,unused // quiet linter
}
