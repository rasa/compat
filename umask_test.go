// SPDX-FileCopyrightText: Copyright © 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"runtime"
	"testing"

	"github.com/rasa/compat"
)

func TestUmask(t *testing.T) {
	if !compat.SupportsUmask() {
		skipf(t, "Skipping test: not supported on %v/%v", runtime.GOOS, runtime.GOARCH)

		return
	}

	original, err := compat.GetUmask()
	if err != nil {
		t.Errorf("GetUmask: %v", err)
	}

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

	previous := original

	for _, mask := range masks {
		got := compat.Umask(mask)
		if got != previous {
			t.Errorf(
				"Umask(0o%03o): got previous mask 0o%03o, want 0o%03o",
				mask,
				got,
				previous,
			)
		}

		got, err := compat.GetUmask()
		if err != nil {
			t.Errorf("GetUmask: %v", err)
		}

		if got != mask {
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

	for _, tst := range tests {
		t.Setenv("UMASK", tst.umask)

		err := compat.InitUmask()

		got := err == nil
		if got != tst.want {
			t.Errorf("%v: got %v, want %v", tst.umask, got, tst.want)
		}
	}
}
