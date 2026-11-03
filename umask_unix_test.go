// SPDX-FileCopyrightText: Copyright © 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build unix

package compat_test

import (
	"errors"
	"runtime"
	"syscall"
	"testing"

	"github.com/rasa/compat"
)

func TestUmaskChanged(t *testing.T) {
	if !compat.SupportsUmask() {
		skipf(t, "Skipping test: Umask() not supported on %v/%v", runtime.GOOS, runtime.GOARCH)

		return
	}

	umask, err := compat.GetUmask()
	if err != nil {
		t.Errorf("GetUmask: %s", err)
	}

	umask++
	if umask > 0o777 {
		umask = 0
	}
	old := syscall.Umask(umask)

	t.Cleanup(func() {
		_ = compat.Umask(old)
	})

	newUmask, err := compat.GetUmask()

	if !errors.Is(err, compat.ErrUmaskChanged) {
		t.Fatalf("got %v, want %v", err, compat.ErrUmaskChanged)
	}

	if newUmask != umask {
		t.Fatalf("got %v, want %v", newUmask, umask)
	}
}
