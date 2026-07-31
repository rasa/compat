// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build (openbsd && ppc64) || (netbsd && 386) || (freebsd && riscv64) || (cgo && aix && ppc64)

package compat

import (
	"context"
)

// PartitionType returns the filesystem type (e.g., "ext4", "NTFS", "FAT32", etc.)
// for the disk partition that contains path.
func PartitionType(_ context.Context, _ string) (string, error) {
	return "", &NotYetImplementedError{Op: "partitionType"}
}
