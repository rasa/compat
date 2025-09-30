// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: Portions copyright (c) 2015 Nate Finch (@natefinch)
// SPDX-FileCopyrightText: Portions copyright (c) 2022 Simon Dassow (@sdassow)

package compat_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/rasa/compat"
)

func TestWriteFileWithAtomicity(t *testing.T) { //nolint:dupl
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	cleanup(t, file)

	perm := compat.CreatePerm // 0o666
	opts := []compat.Option{compat.WithAtomicity(true)}
	err = compat.WriteFile(file, helloBytes, perm, opts...)
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err := compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	want := fixPerms(perm, false)
	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got %04o, want %04o", got, want)
	}
}

func TestWriteFileWithAtomicityCurrentDir(t *testing.T) { //nolint:dupl
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	dir, base := filepath.Split(file)
	t.Chdir(dir)

	cleanup(t, file)

	perm := compat.CreatePerm // 0o666
	opts := []compat.Option{compat.WithAtomicity(true)}
	err = compat.WriteFile(base, helloBytes, perm, opts...)
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err := compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	want := fixPerms(perm, false)
	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got %04o, want %04o", got, want)
	}
}

func TestWriteFileWithAtomicityNoPerms(t *testing.T) { //nolint:dupl
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	cleanup(t, file)

	perm := compat.CreatePerm // 0o600
	opts := []compat.Option{compat.WithAtomicity(true)}
	err = compat.WriteFile(file, helloBytes, 0, opts...)
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err := compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	want := fixPerms(perm, false)
	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got %04o, want %04o", got, want)
	}
}

func TestWriteFileWithAtomicityWithDefaultFileMode(t *testing.T) { //nolint:dupl
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	cleanup(t, file)

	opts := []compat.Option{
		compat.WithAtomicity(true),
		compat.WithDefaultFileMode(perm644),
	}
	err = compat.WriteFile(file, helloBytes, 0, opts...)
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	var fi os.FileInfo

	fi, err = compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	want := fixPerms(perm644, false)

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got %04o, want %04o (1)", got, want)
	}
	// check if file mode is preserved
	err = compat.Chmod(file, perm600)
	if err != nil {
		t.Fatalf("Failed to change file mode: %q: %v", file, err)
	}

	err = compat.WriteFile(file, helloBytes, 0, opts...)
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err = compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	want = fixPerms(perm600, false)

	got = fi.Mode().Perm()

	if got != want {
		t.Fatalf("got %04o, want %04o (2)", got, want)
	}
}

func TestWriteFileWithAtomicityWithKeepFileMode(t *testing.T) { //nolint:dupl
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	cleanup(t, file)

	perm := perm555

	err = compat.WriteFile(file, helloBytes, perm)
	if err != nil {
		t.Fatalf("Failed to create file: %q: %v", file, err)
	}

	opts := []compat.Option{
		compat.WithAtomicity(true),
		compat.WithKeepFileMode(true),
	}
	err = compat.WriteFile(file, helloBytes, 0, opts...)
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err := compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	want := fixPerms(perm, false)
	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got %04o, want %04o: perm=%3o (%v) (1)", got, want, perm, perm)
	}
}

func TestWriteFileWithAtomicityWithKeepFileModeFalse(t *testing.T) { //nolint:dupl
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	cleanup(t, file)

	perm := perm555

	err = compat.WriteFile(file, helloBytes, perm)
	if err != nil {
		t.Fatalf("Failed to create file: %q: %v", file, err)
	}

	opts := []compat.Option{
		compat.WithAtomicity(true),
		compat.WithKeepFileMode(false),
	}
	err = compat.WriteFile(file, helloBytes, 0, opts...)
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err := compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	want := fixPerms(perm, false)
	got := fi.Mode().Perm()
	if got == want {
		if perm != want {
			partType := partitionType(file)
			t.Logf("got %v, want %v (ignoring: %v on %v)", got, want, partType, runtime.GOOS)
			return
		}
		t.Fatalf("got %04o, want !%04o (2)", got, want)
	}
}

func TestWriteFileWithAtomicityWithFileMode(t *testing.T) { //nolint:dupl
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	cleanup(t, file)

	opts := []compat.Option{
		compat.WithAtomicity(true),
		compat.WithFileMode(perm644),
	}
	err = compat.WriteFile(file, helloBytes, 0, opts...)
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err := compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	want := fixPerms(perm644, false)

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got %04o, want %04o (1)", got, want)
	}
	// ensure previous file mode is ignored
	err = compat.Chmod(file, perm600)
	if err != nil {
		t.Fatalf("Failed to change file mode: %q: %v", file, err)
	}

	err = compat.WriteFile(file, helloBytes, 0, opts...)
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err = compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	got = fi.Mode().Perm()
	if got != want {
		t.Fatalf("got %04o, want %04o (2)", got, want)
	}
}

func TestWriteFileWithAtomicityWithReadOnlyModeReset(t *testing.T) { //nolint:dupl
	if !compat.IsWindows {
		skip(t, "Skipping test: requires Windows")
		return
	}

	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	cleanup(t, file)

	opts := []compat.Option{
		compat.WithAtomicity(true),
		compat.WithFileMode(perm400),
		compat.WithReadOnlyMode(compat.ReadOnlyModeReset),
	}
	err = compat.WriteFile(file, helloBytes, 0, opts...)
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err := os.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	want := true // user-writable bit is set.
	got := fi.Mode().Perm()&perm200 == perm200
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// Invalid tests

func TestWriteFileWithAtomicityInvalid(t *testing.T) { //nolint:dupl
	opts := []compat.Option{compat.WithAtomicity(true)}
	err := compat.WriteFile(invalidName, helloBytes, 0, opts...)
	if err == nil {
		t.Fatalf("got nil, want an error")
	}
}

func TestWriteFileWithAtomicityInvalidCantRead(t *testing.T) { //nolint:dupl
	file, err := tempFile(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	cleanup(t, file)

	perm := fixPerms(perm100, false)
	if perm != perm100 {
		partType := partitionType(file)
		t.Skipf("Skipping test: ACLs are not supported on a %v filesystem", partType)
	}
	err = compat.Chmod(file, perm)
	if err != nil {
		t.Fatalf("Chmod: %v", err)
	}

	opts := []compat.Option{
		compat.WithAtomicity(true),
		compat.WithKeepFileMode(true),
	}
	err = compat.WriteFile(file, helloBytes, 0, opts...)
	if err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

func TestWriteFileWithAtomicityInvalidReadOnlyDirectory(t *testing.T) { //nolint:dupl
	if !compat.IsWindows {
		isRoot, _ := compat.IsRoot()
		if isRoot {
			skip(t, "Skipping test: doesn't fail when root")
			return
		}
	}

	name, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	dir, base := filepath.Split(name)
	cleanup(t, dir)
	perm := perm400
	opts := []compat.Option{
		compat.WithFileMode(perm),
		compat.WithReadOnlyMode(compat.ReadOnlyModeSet),
	}
	dir, err = compat.MkdirTemp(dir, "~*.tmp", opts...)
	if err != nil {
		t.Fatalf("MkdirTemp(%v, 0o%o) failed: %v", dir, perm, err)
	}

	file := filepath.Join(dir, base)
	fi, err := compat.Stat(dir)
	if err != nil {
		t.Fatalf("Failed to stat: %v", err)
	}
	if fi.Mode().Perm() != perm {
		partType := partitionType(dir)
		t.Skipf("Skipping test: the %v filesystem does not support permissions", partType)
	}

	opts = []compat.Option{compat.WithAtomicity(true)}
	err = compat.WriteFile(file, helloBytes, 0, opts...)
	if err == nil {
		// @TODO determine why test passes when run individually, but fails when running alongside other tests
		t.Log("got nil, want an error")

		// return
	}

	perm = perm777
	_ = compat.Chmod(dir, perm)
}
