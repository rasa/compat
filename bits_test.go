// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"math/bits"
	"testing"

	"github.com/rasa/compat"
)

func TestIs32bit(t *testing.T) {
	want := bits.UintSize == 32
	got := compat.Is32bit()
	if got != want {
		t.Fatalf("Is32Bit(): got %v, want %v", got, want)
	}
}

func TestIs64bit(t *testing.T) {
	want := bits.UintSize == 64
	got := compat.Is64bit()
	if got != want {
		t.Fatalf("Is64Bit(): got %v, want %v", got, want)
	}
}
