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

		t.Fatal(err)
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
		skip(t, err)
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

		t.Fatal(err)
	}

	for n := 0; n >= compat.MinNice; n-- {
		err = compat.Renice(n)
		if err != nil {
			if !compat.SupportsNice() {
				skipf(t, "Skipping test: Renice() is not supported on %v/%v", runtime.GOOS, runtime.GOARCH)

				return // tinygo doesn't support t.Skip
			}

			// under act, "permission denied" is returned, even though we root.
			skip(t, err)

			return // tinygo doesn't support t.Skip
		}
	}

	err = compat.Renice(nice)
	if err != nil {
		skip(t, err)
	}
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
			skipf(t, "got no error calling Renice with %v (ignoring: doesn't fail on %v)", invalidNice, runtime.GOOS)

			return
		}
		fatalf(t, "got no error calling Renice with %v", invalidNice)

		return // tinygo doesn't support t.Skip
	}
}

func TestNiceErrors(t *testing.T) {
	err := errors.New("Test")
	t.Log(&compat.NiceError{err})
	t.Log(&compat.InvalidNiceError{1024})
	t.Log(&compat.ReniceError{1024, err})
}
