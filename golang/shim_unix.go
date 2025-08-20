// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !windows

package golang

import (
	"os"
)

var (
	Mkdir    = os.Mkdir
	OpenFile = os.OpenFile
)
