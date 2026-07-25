// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9

package volume

import (
	"errors"
)

func Mounts() ([]Mount, error) {
	mounts := []Mount{}

	return mounts, errors.ErrUnsupported
}
