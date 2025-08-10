// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"os"
	"testing"

	"github.com/rasa/compat"
)

func TestWSLIsWSL(t *testing.T) {
	// IsWSL() always returns false in Windows builds, even if the executable
	// is run inside a WSL environment.
	want := !compat.IsWindows && os.Getenv("WSL_DISTRO_NAME") != ""

	got := compat.IsWSL()

	if got != want {
		t.Fatalf("IsWSL(): got %v, want %v", got, want)
	}
}
