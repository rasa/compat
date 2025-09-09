// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: Portions copyright (c) 2015 Nate Finch (@natefinch)
// SPDX-FileCopyrightText: Portions copyright (c) 2022 Simon Dassow (@sdassow)

package compat_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/rasa/compat"
)

func TestWriteReaderAtomic(t *testing.T) {
	dir := tempDir(t)
	file := filepath.Join(dir, "foo.txt")
	content := bytes.NewBufferString("foo")

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	if err := compat.WriteReaderAtomic(file, content); err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err := compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	want := compat.CreateTempPerm // 0o600
	want = fixPerms(want, false)

	got := fi.Mode().Perm()

	if got != want {
		t.Fatalf("got %04o, want %04o", got, want)
	}
}

func TestWriteReaderAtomicDefaultFileMode(t *testing.T) {
	dir := tempDir(t)
	file := filepath.Join(dir, "bar.txt")
	content := bytes.NewBufferString("bar")

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	err := compat.WriteReaderAtomic(file, content, compat.WithDefaultFileMode(perm644))
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	var fi os.FileInfo

	fi, err = compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	want := fixPerms(perm644, false)
	// if compat.IsTinygo && compat.IsWasip1 {
	// 	want = perm600
	// }

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got %04o, want %04o (1)", got, want)
	}
	// check if file mode is preserved
	err = compat.Chmod(file, 0o600)
	if err != nil {
		t.Fatalf("Failed to change file mode: %q: %v", file, err)
	}

	err = compat.WriteReaderAtomic(file, content, compat.WithDefaultFileMode(perm644))
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

func TestWriteReaderAtomicMode(t *testing.T) {
	dir := tempDir(t)
	file := filepath.Join(dir, "baz.txt")
	content := bytes.NewBufferString("baz")

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	err := compat.WriteReaderAtomic(file, content, compat.WithFileMode(perm644))
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err := compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	want := fixPerms(perm644, false)
	// if compat.IsTinygo && compat.IsWasip1 {
	// 	want = perm600
	// }

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got %04o, want %04o (1)", got, want)
	}
	// ensure previous file mode is ignored
	err = compat.Chmod(file, perm600)
	if err != nil {
		t.Fatalf("Failed to change file mode: %q: %v", file, err)
	}

	err = compat.WriteReaderAtomic(file, content, compat.WithFileMode(perm644))
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
