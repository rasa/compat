// SPDX-FileCopyrightText: Copyright © 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"errors"
)

const initialUmask = 0o022

var ErrUmaskChanged = errors.New("process umask changed outside compat")

func init() {
	_ = initUmask()
}
