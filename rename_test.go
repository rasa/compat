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
	cleanup(t, old, new)
	err = compat.Rename(old, new)
	if err != nil {
		t.Fatalf("renaming '%v' to '%v': %v", old, new, err)
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
		t.Fatalf("got no error renaming '%v' to '%v'", old, new)
	}
}

func TestRenameEmptyNew(t *testing.T) {
	old, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}
	cleanup(t, old)
	new := ""

	err = compat.Rename(old, new)
	if err == nil {
		t.Fatalf("got no error renaming '%v' to '%v'", old, new)
	}
}

func TestRenameInvalidOld(t *testing.T) {
	old := invalidName

	new, err := tempName(t)
	if err != nil {
		t.Fatal(err)
	}

	err = compat.Rename(old, new)
	if err == nil {
		t.Fatalf("got no error renaming '%v' to '%v'", old, new)
	}
}

func TestRenameInvalidNew(t *testing.T) {
	old, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}
	cleanup(t, old)
	new := invalidName

	err = compat.Rename(old, new)
	if err == nil {
		t.Fatalf("got no error renaming '%v' to '%v'", old, new)
	}
}

func TestRenameCantRead(t *testing.T) {
	old, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	cleanup(t, old)

	perm := fixPerms(perm100, false)
	if perm != perm100 {
		partType := partitionType(old)
		skipf(t, "Skipping test: permissions are not supported on a %v filesystem", partType)

		return
	}

	err = compat.Chmod(old, perm)
	if err != nil {
		t.Fatalf("Chmod: %v", err)
	}

	new := old + ".new"
	cleanup(t, new)
	err = compat.Rename(old, new)
	if err != nil {
		fatalf(t, "renaming '%v' to '%v': %v", old, new, err)

		return // Tinygo doesn't support T.Fatal
	}
}
