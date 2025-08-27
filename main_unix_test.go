// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !android && !darwin && !linux && !windows && !tinygo

package compat_test

import (
	"fmt"
	"testing"
)

func testMain(m *testing.M, _, nativeFSType, _ string) int { // fsToTest
	fmt.Printf("Testing on a %v filesystem\n", nativeFSType)

	return m.Run()
}
