// SPDX-FileCopyrightText: Copyright (c) 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !plan9

package compat

import (
	"errors"
	"fmt"
)

func normalizeUnsupportedError(err error) error {
	return err
}
