// SPDX-FileCopyrightText: Copyright © 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !windows

package golang

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateTemp(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f, err := CreateTemp(dir, "prefix*suffix", 0o600)
	if err != nil {
		t.Fatalf("CreateTemp(): %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(f.Name())
		_ = f.Close()
	})

	base := filepath.Base(f.Name())
	if !strings.HasPrefix(base, "prefix") {
		t.Fatalf("CreateTemp(): filename %q does not have expected prefix", base)
	}
	if !strings.HasSuffix(base, "suffix") {
		t.Fatalf("CreateTemp(): filename %q does not have expected suffix", base)
	}
}

func TestCreateTempPatternSeparatorError(t *testing.T) {
	t.Parallel()

	_, err := CreateTemp(t.TempDir(), "bad"+string(PathSeparator)+"name", 0o600)
	if err == nil {
		t.Fatal("CreateTemp(): got nil error, want error")
	}
	var pe *PathError
	if !errors.As(err, &pe) {
		t.Fatalf("CreateTemp(): got %T, want *PathError", err)
	}
	if pe.Err != errPatternHasSeparator {
		t.Fatalf("CreateTemp(): PathError.Err got %v, want %v", pe.Err, errPatternHasSeparator)
	}
}

func TestMkdirTemp(t *testing.T) {
	t.Parallel()

	dir, err := MkdirTemp(t.TempDir(), "prefix*suffix", 0o700)
	if err != nil {
		t.Fatalf("MkdirTemp(): %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	base := filepath.Base(dir)
	if !strings.HasPrefix(base, "prefix") {
		t.Fatalf("MkdirTemp(): dirname %q does not have expected prefix", base)
	}
	if !strings.HasSuffix(base, "suffix") {
		t.Fatalf("MkdirTemp(): dirname %q does not have expected suffix", base)
	}
}

func TestMkdirTempPatternSeparatorError(t *testing.T) {
	t.Parallel()

	_, err := MkdirTemp(t.TempDir(), "bad"+string(PathSeparator)+"name", 0o700)
	if err == nil {
		t.Fatal("MkdirTemp(): got nil error, want error")
	}
	var pe *PathError
	if !errors.As(err, &pe) {
		t.Fatalf("MkdirTemp(): got %T, want *PathError", err)
	}
	if pe.Err != errPatternHasSeparator {
		t.Fatalf("MkdirTemp(): PathError.Err got %v, want %v", pe.Err, errPatternHasSeparator)
	}
}

func TestMkdirTempMissingParent(t *testing.T) {
	t.Parallel()

	missingParent := filepath.Join(t.TempDir(), "missing", "parent")
	_, err := MkdirTemp(missingParent, "prefix*", 0o700)
	if err == nil {
		t.Fatal("MkdirTemp(): got nil error, want error")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("MkdirTemp(): got %v, want not-exist error", err)
	}
}
