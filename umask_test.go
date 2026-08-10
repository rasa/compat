// SPDX-FileCopyrightText: Copyright © 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
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
	if compat.IsPlan9 || (compat.IsWasip1 && compat.IsTinygo) {
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
