// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !(openbsd && ppc64) && !(netbsd && 386) && !(freebsd && riscv64) && !(cgo && aix && ppc64)

package compat

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/shirou/gopsutil/v4/disk"
)

// PartitionType returns the filesystem type (e.g., "apfs", "btrfs", "exfat",
// "ext4", "f2fs", "fat32", "hfs+", "ntfs", "refs", "udf", "xfs", etc.)
// for the disk partition that contains path. The value returned is always
// lower case.
func PartitionType(ctx context.Context, path string) (string, error) {
	var err error

	absPath := path
	if !filepath.IsAbs(path) {
		absPath, err = filepath.Abs(path)
		if err != nil {
			return "", fmt.Errorf("cannot convert '%v' to an absolute path: %w", path, err)
		}
	}

	normalizedPath := normalizePath(absPath)

	parts, err := disk.PartitionsWithContext(ctx, true)
	if err != nil {
		return "", err
	}

	for _, p := range parts {
		normalizedMountpoint := normalizePath(p.Mountpoint)
		if SamePartitions(normalizedPath, normalizedMountpoint) {
			return strings.ToLower(p.Fstype), nil
		}
	}

	return "", fmt.Errorf("no mountpoint found for '%v'", path)
}

func normalizePath(path string) string {
	path = filepath.Clean(path)

	if !IsWindows {
		return path
	}

	// strip \\?\ prefix
	path = strings.TrimPrefix(path, `\\?\`)

	// normalize \\?\UNC\server\share -> \\server\share
	if strings.HasPrefix(path, `UNC\`) {
		path = `\` + strings.TrimPrefix(path, `UNC`)
	}

	// c:. => c: (as Clean() changes c: to c:.)
	if len(path) == 3 && strings.HasSuffix(path, ".") {
		path = path[:2]
	}

	// c: => c:\
	if len(path) == 2 && path[1] == ':' {
		path += `\`
	}

	return path
}
