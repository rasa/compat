// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build 386 || arm || mips || mipsle

package compat

// Is32bit reports whether the build target is 32bit.
func Is32bit() bool {
  return true
}

// Is64bit reports whether the build target is 64bit.
func Is64bit() bool {
  return false
}
