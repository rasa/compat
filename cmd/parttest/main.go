// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !(openbsd && ppc64) && !(netbsd && 386) && !(freebsd && riscv64) && !(cgo && aix && ppc64)

package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/shirou/gopsutil/v4/disk"
)

func main() {
	ctx := context.Background()

	parts, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		log.Print(err)

		return
	}

	fmt.Printf("len(parts)=%v\n", len(parts))

	for i, p := range parts {
		if strings.HasPrefix(p.Mountpoint, "/snap") {
			continue
		}

		fmt.Printf("parts[%d]=%#v\n", i, p)
	}
}
