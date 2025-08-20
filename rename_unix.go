// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: Portions copyright (c) 2015 Nate Finch (@natefinch)
// SPDX-FileCopyrightText: Portions copyright (c) 2022 Simon Dassow (@sdassow)

//go:build !windows

package compat

import (
	"os"
)

var rename = os.Rename
