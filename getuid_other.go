// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !(plan9 || unix || wasm || windows)

package compat

func Getuid() (int, error) {
	// this will intentionally not compile to alert us to a new build platform.
}

func Getgid() (int, error) {
	// this will intentionally not compile to alert us to a new build platform.
}

func Geteuid() (int, error) {
	// this will intentionally not compile to alert us to a new build platform.
}

func Getegid() (int, error) {
	// this will intentionally not compile to alert us to a new build platform.
}
