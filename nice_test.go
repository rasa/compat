// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"testing"

	"github.com/rasa/compat"
)

func TestNice(t *testing.T) {
	_, err := compat.Nice()
	if err != nil {
		t.Fatal(err)
	}
}

func TestNiceRenice(t *testing.T) {
	if compat.IsWasip1 {
		skip(t, "Skipping test: operation not supported")

		return // tinygo doesn't support t.Skip
	}

	err := compat.Renice(compat.MaxNice)
	if err != nil {
		// Don't fail on "permission denied" on Linux
		skip(t, err)
	}
}

func TestNiceReniceIfRootValid(t *testing.T) {
	if compat.IsWasip1 {
		skip(t, "Skipping test: operation not supported")

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
	if compat.IsWasip1 {
		skip(t, "Skipping test: operation not supported")

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
		fatalf(t, "got no error calling Renice with %v", invalidNice)

		return // tinygo doesn't support t.Skip
	}
}
