// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !linux && !plan9 && !windows

package volume

import (
	"syscall"
)

func typeOf(_ string, _ string) (Type, error) {
	return TypeUnknown, syscall.ENOSYS
}
