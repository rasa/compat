// SPDX-FileCopyrightText: Copyright © 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"os"
	"testing"

	"github.com/rasa/compat"
)

func TestUmask(t *testing.T) {
	original := compat.GetUmask()

	t.Cleanup(func() {
		compat.Umask(original)
	})

	masks := []int{
		0o000,
		0o002,
		0o022,
		0o027,
		0o077,
		0o777,
	}

	// Umask is unsupported on Plan 9 and TinyGo/wasip1.
	if !compat.SupportsUmask() {
		for _, mask := range masks {
			if got := compat.Umask(mask); got != 0 {
				t.Errorf("Umask(0o%03o): got previous mask 0o%03o, want 0", mask, got)
			}

			if got := compat.GetUmask(); got != 0 {
				t.Errorf("GetUmask() after Umask(0o%03o): got 0o%03o, want 0", mask, got)
			}
		}

		return
	}

	previous := original

	for _, mask := range masks {
		if got := compat.Umask(mask); got != previous {
			t.Errorf(
				"Umask(0o%03o): got previous mask 0o%03o, want 0o%03o",
				mask,
				got,
				previous,
			)
		}

		if got := compat.GetUmask(); got != mask {
			t.Errorf(
				"GetUmask() after Umask(0o%03o): got 0o%03o, want 0o%03o",
				mask,
				got,
				mask,
			)
		}

		previous = mask
	}
}

func TestUmaskInitUmasK(t *testing.T) {
	if !compat.IsWindows {
		skip(t, "Skipping test: requires Windows")

		return
	}
	type test struct {
		umask string
		want  bool
	}

	tests := []test{
		{"", true},
		{"0", true},
		{"0o", true},
		{"o", false},
		{"9", false},
		{"-1", false},
		{"18446744073709551617", false},
	}
	savedUmask := os.Getenv("UMASK")
	for _, tst := range tests {
		os.Setenv("UMASK", tst.umask)
		err := compat.InitUmask()
		got := err == nil
		if got != tst.want {
			t.Errorf("%v: got %v, want %v", tst.umask, got, tst.want)
		}
	}
	os.Setenv("UMASK", savedUmask)
}
