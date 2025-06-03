// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !windows || !debug

package compat

import (
	"os"
)

func dumpMasks(perm os.FileMode, ownerMask uint32, groupMask uint32, worldMask uint32) {
}
