// SPDX-FileCopyrightText: Copyright (c) 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build tinygo

package compat

import "errors"

// Link creates newname as a hard link to the oldname file.
// If there is an error, it will be of type *LinkError.
func Link(_, _ string) error {
	// See https://github.com/tinygo-org/tinygo/blob/3869f768/src/os/errors.go#L29
	return errors.New("operation not implemented")
}
