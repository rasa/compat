// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !(js || plan9 || unix || wasip1 || windows)

// https://github.com/golang/go/blob/8ad27fb6/src/cmd/dist/build.go#L1070
// unix == aix || android || darwin || dragonfly || freebsd || illumos || ios || linux || netbsd || openbsd || solaris

package compat

// Umask sets the umask to umask, and returns the previous value.
// On Windows, the initial umask value is 022 octal, and can be changed by
// setting environmental variable UMASK, to an octal value. For example:
// `set UMASK=002`. Leading zeros and 'o's are allowed, and ignored.
// On Plan9 and Wasip1, the function does nothing, and always returns zero.
func Umask(newMask int) int {
	// this will intentionally not compile to alert us to a new build platform.
}

// GetUmask returns the current umask value.
// On Plan9 and Wasip1, the function always returns zero.
func GetUmask() int {
	// this will intentionally not compile to alert us to a new build platform.
}
