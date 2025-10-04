// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build js || plan9 || wasip1

package compat

func cpuBits() (int, error) {
	return BuildBits(), nil
}
