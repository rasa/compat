// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"os"
	"testing"

	"github.com/rasa/compat"
)

func TestFilePosixWriteFile(t *testing.T) {
	perm := os.FileMode(0o666)
	want := fixPosixPerms(perm, false)

	name, err := tempName(t)
	if err != nil {
		t.Fatal(err)

		return
	}

	err = compat.WriteFile(name, helloBytes, want)
	if err != nil {
		t.Fatal(err)

		return
	}

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o (%v), want 0%03o (%v)", got, got, want, want)

		return
	}
}

func TestFilePosixWriteFileInvalid(t *testing.T) {
	err := compat.WriteFile(invalidName, helloBytes, compat.CreatePerm)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}
