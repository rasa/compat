// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !(386 || arm || mips || mipsle)

package compat

// CPUBits returns the number of bits in an integer on the build target.
// For 386, arm, mips, and mipsle, it's 32. For all other targets, it's 64.
const CPUBits = 64
