// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"testing"

	"github.com/rasa/compat"
)

func TestRename(t *testing.T) {
	old, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	new := old + ".new"
	err = compat.Rename(old, new)
	if err != nil {
		t.Fatalf("renaming %q to %q: %v", old, new, err)
	}
}

func TestRenameEmptyOld(t *testing.T) {
	old := ""
	new, err := tempName(t)
	if err != nil {
		t.Fatal(err)
	}

	err = compat.Rename(old, new)
	if err == nil {
		t.Fatalf("got no error renaming %q to %q", old, new)
	}
}

func TestRenameEmptyNew(t *testing.T) {
	old, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}
	new := ""

	err = compat.Rename(old, new)
	if err == nil {
		t.Fatalf("got no error renaming %q to %q", old, new)
	}
}

func TestRenameBadOld(t *testing.T) {
	old := "/an/invalid/filename"

	new, err := tempName(t)
	if err != nil {
		t.Fatal(err)
	}

	err = compat.Rename(old, new)
	if err == nil {
		t.Fatalf("got no error renaming %q to %q", old, new)
	}
}

func TestRenameBadNew(t *testing.T) {
	old, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}
	new := "/an/invalid/filename"

	err = compat.Rename(old, new)
	if err == nil {
		t.Fatalf("got no error renaming %q to %q", old, new)
	}
}
