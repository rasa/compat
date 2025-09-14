// SPDX-FileCopyrightText: Copyright (c) 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build tinygo

package compat_test

import "errors"

func osLink(_, _ string) error {
	return errors.New("operation not supported")
}
