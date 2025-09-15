// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"runtime"
	"testing"

	"github.com/rasa/compat"
)

func TestNice(t *testing.T) {
	if !compat.SupportsNice() {
		skipf(t, "Skipping test: Nice() is not supported on %v/%v", runtime.GOOS, runtime.GOARCH)

		return // tinygo doesn't support t.Skip
	}

	_, err := compat.Nice()
	if err != nil {
		t.Fatal(err)
	}
}

func TestNiceRenice(t *testing.T) {
	if !compat.SupportsNice() {
		skipf(t, "Skipping test: Nice() is not supported on %v/%v", runtime.GOOS, runtime.GOARCH)

		return // tinygo doesn't support t.Skip
	}

	err := compat.Renice(compat.MaxNice)
	if err != nil {
		// Don't fail on "permission denied" on Linux
		skip(t, err)
	}
}

func TestNiceReniceIfRootValid(t *testing.T) {
	if !compat.SupportsNice() {
		skipf(t, "Skipping test: Nice() is not supported on %v/%v", runtime.GOOS, runtime.GOARCH)

		return // tinygo doesn't support t.Skip
	}

	isRoot, _ := compat.IsRoot()

	if !compat.IsWindows && !isRoot {
		skip(t, "Skipping test: we aren't the root/admin user")

		return // tinygo doesn't support t.Skip
	}

	nice, err := compat.Nice()
	if err != nil {
		t.Fatal(err)
	}

	for n := 0; n >= compat.MinNice; n-- {
		err = compat.Renice(n)
		if err != nil {
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
	if !compat.SupportsNice() {
		skipf(t, "Skipping test: Nice() is not supported on %v/%v", runtime.GOOS, runtime.GOARCH)

		return // tinygo doesn't support t.Skip
	}

	isRoot, _ := compat.IsRoot()

	if !compat.IsWindows && !isRoot {
		skip(t, "Skipping test: we aren't the root/admin user")

		return // tinygo doesn't support t.Skip
	}

	const invalidNice = compat.MinNice - 1024

	err := compat.Renice(invalidNice)
	if err == nil {
		if compat.IsBSD {
			skipf(t, "got no error calling Renice with %v on %v (ignoring)", invalidNice, runtime.GOOS)
		}
		fatalf(t, "got no error calling Renice with %v", invalidNice)

		return // tinygo doesn't support t.Skip
	}
}
