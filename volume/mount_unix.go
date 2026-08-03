// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build unix && !(openbsd && ppc64) && !(netbsd && 386) && !(freebsd && riscv64) && !(cgo && aix && ppc64)

package volume

import (
	"context"
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v4/disk"
)

func Mounts() ([]Mount, error) {
	mounts := []Mount{}

	ctx := context.Background()
	parts, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		return mounts, fmt.Errorf("PartitionsWithContext: %w", err)
	}

	for _, p := range parts {
		mount := Mount{
			Device:     p.Device,
			Mountpoint: p.Mountpoint,
			Fstype:     strings.ToLower(p.Fstype),
			Opts:       p.Opts,
		}

		mounts = append(mounts, mount)
	}

	return mounts, nil
}
