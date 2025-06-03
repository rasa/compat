// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"runtime"
	"testing"

	"github.com/rasa/compat"
)

func TestCPUBits(t *testing.T) {
	if runtime.GOOS == "plan9" || runtime.GOARCH == "wasm" {
		t.Skip("Not supported on " + runtime.GOOS + "/" + runtime.GOARCH)
	}
	_, err := compat.CPUBits()
	if err != nil {
		t.Fatalf("CPUBits: got %v, want nil", err)
	}
}
