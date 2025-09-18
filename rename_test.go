// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"os"
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

func TestRenameInvalidOld(t *testing.T) {
	old := invalidName

	new, err := tempName(t)
	if err != nil {
		t.Fatal(err)
	}

	err = compat.Rename(old, new)
	if err == nil {
		t.Fatalf("got no error renaming %q to %q", old, new)
	}
}

func TestRenameInvalidNew(t *testing.T) {
	old, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}
	new := invalidName

	err = compat.Rename(old, new)
	if err == nil {
		t.Fatalf("got no error renaming %q to %q", old, new)
	}
}

func TestRenameCantRead(t *testing.T) {
	old, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = compat.Chmod(old, perm600)
		_ = os.Remove(old)
	})

	perm := fixPerms(perm100, false)
	if perm != perm100 {
		partType := partitionType(old)
		skipf(t, "Skipping test: ACLs are not supported on a %v filesystem", partType)

		return
	}

	err = compat.Chmod(old, perm)
	if err != nil {
		t.Fatalf("Chmod: %v", err)
	}

	new := old + ".new"
	err = compat.Rename(old, new)
	if err != nil {
		fatalf(t, "Rename: %v", err)

		return // Tinygo doesn't support T.Fatal
	}
}
