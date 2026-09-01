// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build (!unix && !windows) || (openbsd && ppc64) || (netbsd && 386) || (freebsd && riscv64) || (cgo && aix && ppc64)

package volume

import (
	"errors"
	"fmt"
)

func Mounts() ([]Mount, error) {
	mounts := []Mount{}

	return mounts, fmt.Errorf("mounts: %w", errors.ErrUnsupported)
}
