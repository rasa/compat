// SPDX-FileCopyrightText: Copyright (c) 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

import (
	"errors"
	"fmt"
	"syscall"
)

func normalizeUnsupportedError(err error) error {
	if errors.Is(err, syscall.EWINDOWS) {
		return fmt.Errorf("%w: %w", errors.ErrUnsupported, err)
	}

	return err
}
