// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: Portions copyright (c) 2015 Nate Finch (@natefinch)
// SPDX-FileCopyrightText: Portions copyright (c) 2022 Simon Dassow (@sdassow)

package compat_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/rasa/compat"
)

var helloBuf = bytes.NewBuffer(helloBytes)

func TestWriteReaderAtomic(t *testing.T) {
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	err = compat.WriteReaderAtomic(file, helloBuf)
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
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

func TestWriteReaderAtomicDefaultFileMode(t *testing.T) {
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	err = compat.WriteReaderAtomic(file, helloBuf, compat.WithDefaultFileMode(perm644))
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
	err = compat.Chmod(file, 0o600)
	if err != nil {
		t.Fatalf("Failed to change file mode: %q: %v", file, err)
	}

	err = compat.WriteReaderAtomic(file, helloBuf, compat.WithDefaultFileMode(perm644))
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

func TestWriteReaderAtomicKeepFileMode(t *testing.T) { //nolint:dupl
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

	err = compat.WriteReaderAtomic(file, helloBuf, compat.KeepFileMode(true))
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
		t.Fatalf("got %04o, want %04o: perm=%03o (%v) (1)", got, want, perm, perm)
	}

	err = compat.WriteReaderAtomic(file, helloBuf, compat.KeepFileMode(false))
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err = compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	got = fi.Mode().Perm()
	if got == want {
		if perm != want {
			t.Logf("got %v, want %v (ignoring: %v on %v)", got, want, partType, runtime.GOOS)
			returm
		}
		t.Fatalf("got %04o, want !%04o (2)", got, want)
	}
}

func TestWriteReaderAtomicWithFileMode(t *testing.T) {
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	err = compat.WriteReaderAtomic(file, helloBuf, compat.WithFileMode(perm644))
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

	err = compat.WriteReaderAtomic(file, helloBuf, compat.WithFileMode(perm644))
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
