// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"math/bits"
	"testing"

	"github.com/rasa/compat"
)

func TestCPUBits(t *testing.T) {
	want := bits.UintSize
	got := compat.CPUBits
	if got != want {
		t.Fatalf("CPUBits: got %v, want %v", got, want)
	}
}
