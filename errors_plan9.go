// SPDX-FileCopyrightText: Copyright (c) 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9

package compat

import (
	"errors"
	"fmt"
	"syscall"
)

func normalizeUnsupportedError(err error) error {
	if errors.Is(err, syscall.EPLAN9) {
		return fmt.Errorf("%w: %w", errors.ErrUnsupported, err)
	}

	return err
}
