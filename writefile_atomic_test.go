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
	file := "foo.txt"
	content := []byte("foo")
	t.Cleanup(func() {
		_ = os.Remove(file)
	})
	if err := compat.WriteFileAtomic(file, content); err != nil {
		t.Errorf("Failed to write file: %q: %v", file, err)
	}
	fi, err := compat.Stat(file)
	if err != nil {
		t.Errorf("Failed to stat file: %q: %v", file, err)
	}
	want := compat.CreateTempPerm // 0o600
	got := fi.Mode().Perm()
	if got != want {
		t.Errorf("got %04o, want %04o", got, want)
	}
}

func TestWriteFileAtomicDefaultFileMode(t *testing.T) {
	file := "bar.txt"
	content := []byte("bar")
	t.Cleanup(func() {
		_ = os.Remove(file)
	})
	err := compat.WriteFileAtomic(file, content, compat.DefaultFileMode(0o644))
	if err != nil {
		t.Errorf("Failed to write file: %q: %v", file, err)
	}
	var fi os.FileInfo
	fi, err = compat.Stat(file)
	if err != nil {
		t.Errorf("Failed to stat file: %q: %v", file, err)
	}
	want := want644
	if compat.IsTinygo && compat.IsWasip1 {
		want = 0o600
	}
	got := fi.Mode().Perm()
	if got != want {
		t.Errorf("got %04o, want %04o (1)", got, want)
	}
	// check if file mode is preserved
	err = compat.Chmod(file, 0o600)
	if err != nil {
		t.Errorf("Failed to change file mode: %q: %v", file, err)
	}
	err = compat.WriteFileAtomic(file, content, compat.DefaultFileMode(0o644))
	if err != nil {
		t.Errorf("Failed to write file: %q: %v", file, err)
	}
	fi, err = compat.Stat(file)
	if err != nil {
		t.Errorf("Failed to stat file: %q: %v", file, err)
	}
	want = os.FileMode(0o600)
	got = fi.Mode().Perm()
	if got != want {
		t.Errorf("got %04o, want %04o (2)", got, want)
	}
}

func TestWriteFileAtomicMode(t *testing.T) {
	file := "baz.txt"
	content := []byte("baz")
	t.Cleanup(func() {
		_ = os.Remove(file)
	})
	err := compat.WriteFileAtomic(file, content, compat.FileMode(0o644))
	if err != nil {
		t.Errorf("Failed to write file: %q: %v", file, err)
	}
	fi, err := compat.Stat(file)
	if err != nil {
		t.Errorf("Failed to stat file: %q: %v", file, err)
	}
	want := want644
	if compat.IsTinygo && compat.IsWasip1 {
		want = 0o600
	}
	got := fi.Mode().Perm()
	if got != want {
		t.Errorf("got %04o, want %04o (1)", got, want)
	}
	// ensure previous file mode is ignored
	err = compat.Chmod(file, 0o600)
	if err != nil {
		t.Errorf("Failed to change file mode: %q: %v", file, err)
	}
	err = compat.WriteFileAtomic(file, content, compat.FileMode(0o644))
	if err != nil {
		t.Errorf("Failed to write file: %q: %v", file, err)
	}
	fi, err = compat.Stat(file)
	if err != nil {
		t.Errorf("Failed to stat file: %q: %v", file, err)
	}
	got = fi.Mode().Perm()
	if got != want {
		t.Errorf("got %04o, want %04o (2)", got, want)
	}
}
