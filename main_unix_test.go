// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build (!darwin && !linux && !windows && !tinygo) || android

package compat_test

import (
	"fmt"
	"runtime"
	"testing"
)

func testMain(m *testing.M, _, nativeFSType, _ string) int { // fsToTest
	fmt.Printf("Testing on a %v filesystem on %v/%v\n", nativeFSType, runtime.GOOS, runtime.GOARCH)

	return m.Run()
}
