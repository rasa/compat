// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !(aix || android || darwin || dragonfly || freebsd || illumos || ios || js || linux || netbsd || openbsd || plan9 || solaris || wasip1 || windows)

package compat

func Umask(newMask int) int {
	// this will intentionally not compile to alert us to a new build platform.
}

func GetUmask() (int, error) {
	// this will intentionally not compile to alert us to a new build platform.
}
