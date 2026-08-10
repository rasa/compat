// SPDX-FileCopyrightText: Copyright © 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build tinygo

package compat

// IsTinygo is true if the go compiler is tinygo.
const IsTinygo = true // was runtime.Compiler == "tinygo"
