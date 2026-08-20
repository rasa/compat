// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"runtime"
	"testing"

	"github.com/rasa/compat"
)

func TestRename(t *testing.T) {
	old, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	gnu := old + ".new"
	cleanup(t, old, gnu)

	err = compat.Rename(old, gnu)
	if err != nil {
		t.Fatalf("renaming '%v' to '%v': %v", old, gnu, err)
	}
}

func TestRenameEmptyOld(t *testing.T) {
	old := ""
	gnu := tempName(t)

	err := compat.Rename(old, gnu)
	if err == nil {
		t.Fatalf("got no error renaming '%v' to '%v'", old, gnu)
	}
}

func TestRenameEmptyNew(t *testing.T) {
	old, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	cleanup(t, old)

	gnu := ""

	err = compat.Rename(old, gnu)
	if err == nil {
		t.Fatalf("got no error renaming '%v' to '%v'", old, gnu)
	}
}

func TestRenameInvalidOld(t *testing.T) {
	old := invalidName
	gnu := tempName(t)

	err := compat.Rename(old, gnu)
	if err == nil {
		t.Fatalf("got no error renaming '%v' to '%v'", old, gnu)
	}
}

func TestRenameInvalidNew(t *testing.T) {
	old, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	cleanup(t, old)

	gnu := invalidName

	err = compat.Rename(old, gnu)
	if err == nil {
		t.Fatalf("got no error renaming '%v' to '%v'", old, gnu)
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
		t.Skipf("Skipping test: permissions are not supported on a %v filesystem", partType)
	}

	err = compat.Chmod(old, perm)
	if err != nil {
		t.Fatalf("Chmod: %v", err)
	}

	gnu := old + ".new"
	cleanup(t, gnu)

	err = compat.Rename(old, gnu)
	if err != nil {
		t.Fatalf("renaming '%v' to '%v': %v", old, gnu, err)
	}
}

func TestRenameWithAtomicity(t *testing.T) {
	if !compat.SupportsAtomicReplace() {
		skipf(t, "Skipping test: WithAtomicity() not supported on %v", runtime.GOOS)

		return
	}

	old, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	gnu := old + ".new"
	cleanup(t, old, gnu)

	err = compat.Rename(old, gnu, compat.WithAtomicity(true))
	if err != nil {
		t.Fatalf("renaming '%v' to '%v': %v", old, gnu, err)
	}
}

func TestRenameWithNonAtomicReplace(t *testing.T) {
	old, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	gnu := old + ".new"
	cleanup(t, old, gnu)

	err = compat.Rename(old, gnu, compat.WithNonAtomicReplace(true))
	if err != nil {
		t.Fatalf("renaming '%v' to '%v': %v", old, gnu, err)
	}
}
