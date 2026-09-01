// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"runtime"
	"testing"

	"github.com/rasa/compat"
)

func TestRename(t *testing.T) {
	src, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	dst := src + ".new"
	cleanup(t, src, dst)

	err = compat.Rename(src, dst)
	if err != nil {
		t.Fatalf("renaming %q to %q: %v", src, dst, err)
	}
}

func TestRenameEmptysrc(t *testing.T) {
	src := ""
	dst := tempName(t)

	err := compat.Rename(src, dst)
	if err == nil {
		t.Fatalf("got no error renaming %q to %q", src, dst)
	}
}

func TestRenameEmptyNew(t *testing.T) {
	src, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	cleanup(t, src)

	dst := ""

	err = compat.Rename(src, dst)
	if err == nil {
		t.Fatalf("got no error renaming %q to %q", src, dst)
	}
}

func TestRenameInvalidsrc(t *testing.T) {
	src := invalidName
	dst := tempName(t)

	err := compat.Rename(src, dst)
	if err == nil {
		t.Fatalf("got no error renaming %q to %q", src, dst)
	}
}

func TestRenameInvalidNew(t *testing.T) {
	src, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	cleanup(t, src)

	dst := invalidName

	err = compat.Rename(src, dst)
	if err == nil {
		t.Fatalf("got no error renaming %q to %q", src, dst)
	}
}

func TestRenameCantRead(t *testing.T) {
	src, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	cleanup(t, src)

	perm := fixPerms(perm100, false)
	if perm != perm100 {
		partType := partitionType(src)
		t.Skipf("Skipping test: permissions are not supported on a %v filesystem", partType)
	}

	err = compat.Chmod(src, perm)
	if err != nil {
		t.Fatalf("Chmod: %v", err)
	}

	dst := src + ".new"
	cleanup(t, dst)

	err = compat.Rename(src, dst)
	if err != nil {
		t.Fatalf("renaming %q to %q: %v", src, dst, err)
	}
}

func TestRenameWithAtomicity(t *testing.T) {
	if !compat.SupportsAtomicReplace() {
		skipf(t, "Skipping test: WithAtomicity() not supported on %v", runtime.GOOS)

		return
	}

	src, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	dst := src + ".new"
	cleanup(t, src, dst)

	err = compat.Rename(src, dst, compat.WithAtomicity(true))
	if err != nil {
		t.Fatalf("renaming %q to %q: %v", src, dst, err)
	}
}

func TestRenameWithNonAtomicReplace(t *testing.T) {
	src, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	dst := src + ".new"
	cleanup(t, src, dst)

	err = compat.Rename(src, dst, compat.WithNonAtomicReplace(true))
	if err != nil {
		t.Fatalf("renaming %q to %q: %v", src, dst, err)
	}
}
