// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: Portions copyright (c) 2015 Nate Finch (@natefinch)
// SPDX-FileCopyrightText: Portions copyright (c) 2022 Simon Dassow (@sdassow)

package compat_test

import (
	"os"
	"testing"

	"github.com/rasa/compat"
)

func TestWriteFileAtomic(t *testing.T) {
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	err = compat.WriteFileAtomic(file, helloBytes)
	if err != nil {
		fatalf(t, "Failed to write file: %q: %v", file, err)

		return // Tinygo doesn't support T.Fatal
	}

	fi, err := compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	perm := compat.CreateTempPerm // 0o600
	want := fixPerms(perm, false)

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got %04o, want %04o", got, want)
	}
}

func TestWriteFileAtomicDefaultFileMode(t *testing.T) {
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	err = compat.WriteFileAtomic(file, helloBytes, compat.WithDefaultFileMode(perm644))
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

	err = compat.WriteFileAtomic(file, helloBytes, compat.WithDefaultFileMode(perm644))
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

func TestWriteFileAtomicKeepFileMode(t *testing.T) {
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	perm := perm555

	err = compat.WriteFile(file, helloBytes, perm)
	if err != nil {
		t.Fatalf("Failed to create file: %q: %v", file, err)
	}

	err = compat.WriteFileAtomic(file, helloBytes, compat.KeepFileMode(true))
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
		t.Fatalf("got %04o, want %04o (1)", got, want)
	}

	err = compat.WriteFileAtomic(file, helloBytes, compat.KeepFileMode(false))
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err = compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	got = fi.Mode().Perm()
	if got == want {
		t.Fatalf("got %04o, want %04o (2)", got, want)
	}
}

func TestWriteFileAtomicWithFileMode(t *testing.T) {
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	err = compat.WriteFileAtomic(file, helloBytes, compat.WithFileMode(perm644))
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

	err = compat.WriteFileAtomic(file, helloBytes, compat.WithFileMode(perm644))
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
