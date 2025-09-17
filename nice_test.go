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

			return // tinygo doesn't support t.Skip
		}

		fatalf(t, "Nice; got %v, want nil", err)
	}
}

func TestNiceRenice(t *testing.T) {
	err := compat.Renice(compat.MaxNice)
	if err != nil {
		if !compat.SupportsNice() {
			skipf(t, "Skipping test: Renice() is not supported on %v/%v", runtime.GOOS, runtime.GOARCH)

			return // tinygo doesn't support t.Skip
		}

		// Don't fail on "permission denied" on Linux
		skipf(t, "Renice: got %v, want nil", err)
	}
}

func TestNiceReniceIfRootValid(t *testing.T) {
	isRoot, _ := compat.IsRoot()

	if !compat.IsWindows && !isRoot {
		skip(t, "Skipping test: we aren't the root/admin user")

		return // tinygo doesn't support t.Skip
	}

	nice, err := compat.Nice()
	if err != nil {
		if !compat.SupportsNice() {
			skipf(t, "Skipping test: Nice() is not supported on %v/%v", runtime.GOOS, runtime.GOARCH)

			return // tinygo doesn't support t.Skip
		}

		fatalf(t, "Nice; got %v, want nil", err)
	}

	for n := 0; n >= compat.MinNice; n-- {
		err = compat.Renice(n)
		if err != nil {
			if !compat.SupportsNice() {
				skipf(t, "Skipping test: Renice() is not supported on %v/%v", runtime.GOOS, runtime.GOARCH)

				return // tinygo doesn't support t.Skip
			}

			// under act, "permission denied" is returned, even though we root.
			skipf(t, "Renice: got %v, want nil", err)

			return // tinygo doesn't support t.Skip
		}
	}

	_ = compat.Renice(nice)
}

func TestNiceReniceIfRootInvalid(t *testing.T) {
	isRoot, _ := compat.IsRoot()

	if !compat.IsWindows && !isRoot {
		skip(t, "Skipping test: we aren't the root/admin user")

		return // tinygo doesn't support t.Skip
	}

	const invalidNice = compat.MinNice - 1024

	err := compat.Renice(invalidNice)
	if err == nil {
		if !compat.SupportsNice() {
			skipf(t, "Skipping test: Nice() is not supported on %v/%v", runtime.GOOS, runtime.GOARCH)

			return // tinygo doesn't support t.Skip
		}

		if !compat.IsWindows && !compat.IsPlan9 {
			skipf(t, "Renice(%v): got nil, want error (ignoring: doesn't fail on %v)", invalidNice, runtime.GOOS)

			return
		}

		fatalf(t, "Renice(%v): got nil, want error", invalidNice)

		return // tinygo doesn't support t.Skip
	}
}

func TestNiceErrors(t *testing.T) {
	err := errors.New("Test")

	e1 := &compat.NiceError{err}
	if e1.Error() == "" {
    		fatal(t, "NiceError: got '', want non-empty string")
	}

	e2 := &compat.InvalidNiceError{1024}
	if e2.Error() == "" {
    		fatal(t, "InvalidNiceError: got '', want non-empty string")
	}

	e3 := &compat.ReniceError{1024, err}
	if e3.Error() == "" {
    		fatal(t, "ReniceError: got '', want non-empty string")
	}
}
