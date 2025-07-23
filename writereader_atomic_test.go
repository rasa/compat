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

func TestWriteReaderAtomic(t *testing.T) {
	if compat.IsWasip1Target {
		t.Log("Skipping test on wasip1 target: operation not supported")
		return
	}
	file := "foo.txt"
	content := bytes.NewBufferString("foo")
	defer func() { _ = os.Remove(file) }()
	if err := compat.WriteReaderAtomic(file, content); err != nil {
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

func TestWriteReaderAtomicDefaultFileMode(t *testing.T) {
	if compat.IsWasip1Target {
		t.Log("Skipping test on wasip1 target: operation not supported")
		return
	}
	file := "bar.txt"
	content := bytes.NewBufferString("bar")
	defer func() { _ = os.Remove(file) }()
	err := compat.WriteReaderAtomic(file, content, compat.DefaultFileMode(0o644))
	if err != nil {
		t.Errorf("Failed to write file: %q: %v", file, err)
	}
	var fi os.FileInfo
	fi, err = compat.Stat(file)
	if err != nil {
		t.Errorf("Failed to stat file: %q: %v", file, err)
	}
	want := want644
	got := fi.Mode().Perm()
	if got != want {
		t.Errorf("got %04o, want %04o", got, want)
	}
	// check if file mode is preserved
	err = compat.Chmod(file, 0o600)
	if err != nil {
		t.Errorf("Failed to change file mode: %q: %v", file, err)
	}
	err = compat.WriteReaderAtomic(file, content, compat.DefaultFileMode(0o644))
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
		t.Errorf("got %04o, want %04o", got, want)
	}
}

func TestWriteReaderAtomicMode(t *testing.T) {
	if compat.IsWasip1Target {
		t.Log("Skipping test on wasip1 target: operation not supported")
		return
	}
	file := "baz.txt"
	content := bytes.NewBufferString("baz")
	defer func() { _ = os.Remove(file) }()
	err := compat.WriteReaderAtomic(file, content, compat.FileMode(0o644))
	if err != nil {
		t.Errorf("Failed to write file: %q: %v", file, err)
	}
	fi, err := compat.Stat(file)
	if err != nil {
		t.Errorf("Failed to stat file: %q: %v", file, err)
	}
	want := want644
	got := fi.Mode().Perm()
	if got != want {
		t.Errorf("got %04o, want %04o", got, want)
	}
	// ensure previous file mode is ignored
	err = compat.Chmod(file, 0o600)
	if err != nil {
		t.Errorf("Failed to change file mode: %q: %v", file, err)
	}
	err = compat.WriteReaderAtomic(file, content, compat.FileMode(0o644))
	if err != nil {
		t.Errorf("Failed to write file: %q: %v", file, err)
	}
	fi, err = compat.Stat(file)
	if err != nil {
		t.Errorf("Failed to stat file: %q: %v", file, err)
	}
	got = fi.Mode().Perm()
	if got != want {
		t.Errorf("got %04o, want %04o", got, want)
	}
}
