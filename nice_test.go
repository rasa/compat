// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"errors"
	"runtime"
	"testing"

	"github.com/rasa/compat"
)

func TestNice(t *testing.T) {
	_, err := compat.Nice()
	if err != nil {
		if !compat.SupportsNice() {
			skipf(t, "Skipping test: Nice() is not supported on %v/%v", runtime.GOOS, runtime.GOARCH)
			return
		}

		t.Fatalf("Nice; got %v, want nil", err)
	}
}

func TestNiceRenice(t *testing.T) {
	err := compat.Renice(compat.MaxNice)
	if err != nil {
		if !compat.SupportsNice() {
			skipf(t, "Skipping test: Renice() is not supported on %v/%v", runtime.GOOS, runtime.GOARCH)
			return
		}

		// Don't fail on "permission denied" on Linux
		t.Skipf("Renice: got %v, want nil", err)
	}
}

func TestNiceReniceIfRootValid(t *testing.T) {
	if !compat.IsWindows {
		isRoot, _ := compat.IsRoot()
		if !isRoot {
			skip(t, "Skipping test: we aren't the root/admin user")
			return
		}
	}

	nice, err := compat.Nice()
	if err != nil {
		if !compat.SupportsNice() {
			skipf(t, "Skipping test: Nice() is not supported on %v/%v", runtime.GOOS, runtime.GOARCH)
			return
		}

		t.Fatalf("Nice; got %v, want nil", err)
	}

	for n := 0; n >= compat.MinNice; n-- {
		err = compat.Renice(n)
		if err != nil {
			if !compat.SupportsNice() {
				skipf(t, "Skipping test: Renice() is not supported on %v/%v", runtime.GOOS, runtime.GOARCH)
				return
			}

			// under act, "permission denied" is returned, even though we root.
			t.Skipf("Renice: got %v, want nil", err)
		}
	}

	_ = compat.Renice(nice)
}

func TestNiceReniceIfRootInvalid(t *testing.T) {
	if !compat.IsWindows {
		isRoot, _ := compat.IsRoot()
		if !isRoot {
			skip(t, "Skipping test: we aren't the root/admin user")
			return
		}
	}

	const invalidNice = compat.MinNice - 1024

	err := compat.Renice(invalidNice)
	if err == nil {
		if !compat.SupportsNice() {
			skipf(t, "Skipping test: Nice() is not supported on %v/%v", runtime.GOOS, runtime.GOARCH)
			return
		}

		if !compat.IsWindows && !compat.IsPlan9 {
			t.Skipf("Renice(%v): got nil, want error (ignoring: doesn't fail on %v)", invalidNice, runtime.GOOS)
		}

		t.Fatalf("Renice(%v): got nil, want error", invalidNice)
	}
}

func TestNiceErrors(t *testing.T) {
	err := errors.New("Test")

	e1 := &compat.NiceError{err}
	if e1.Error() == "" {
		t.Fatal("NiceError: got '', want non-empty string")
	}

	e2 := &compat.InvalidNiceError{1024}
	if e2.Error() == "" {
		t.Fatal("InvalidNiceError: got '', want non-empty string")
	}

	e3 := &compat.ReniceError{1024, err}
	if e3.Error() == "" {
		t.Fatal("ReniceError: got '', want non-empty string")
	}
}
