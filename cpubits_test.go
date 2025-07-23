// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"runtime"
	"testing"

	"github.com/rasa/compat"
)

func TestCPUBits(t *testing.T) {
	if compat.IsWasip1 {
		t.Log("Skipping test on wasip1: operation not supported")
		return
	}
	if compat.IsPlan9 {
		t.Skip("Not supported on " + runtime.GOOS + "/" + runtime.GOARCH)
	}
	_, err := compat.CPUBits()
	if err != nil {
		t.Fatalf("CPUBits: got %v, want nil", err)
	}
}
