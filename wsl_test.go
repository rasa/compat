// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/rasa/compat"
)

func TestIsWSL(t *testing.T) {
	// IsWSL() always returns false in Windows builds, even if the executable
	// is run inside a WSL environment.
	isWSL := false
	if runtime.GOOS != "windows" {
		isWSL = os.Getenv("WSL_DISTRO_NAME") != ""
	}
	got := compat.IsWSL()
	if got != isWSL {
		t.Errorf("IsWSL(): got %v, expected %v", got, isWSL)
	}
}
