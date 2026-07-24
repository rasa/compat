// SPDX-FileCopyrightText: Copyright © 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !windows && !darwin

package robustio

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestRetryReturnsFirstResult(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("boom")
	calls := 0
	err := Retry(func() (error, bool) {
		calls++
		return wantErr, true
	}, 2)

	if !errors.Is(err, wantErr) {
		t.Fatalf("Retry(): got %v, want %v", err, wantErr)
	}
	if calls != 1 {
		t.Fatalf("Retry(): got %v calls, want 1", calls)
	}
}

func TestIsEphemeralErrorAlwaysFalse(t *testing.T) {
	t.Parallel()

	if IsEphemeralError(errors.New("x")) {
		t.Fatal("IsEphemeralError(): got true, want false")
	}
	if IsEphemeralError(nil) {
		t.Fatal("IsEphemeralError(nil): got true, want false")
	}
}

func TestRenameReadFileRemoveAll(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	oldPath := filepath.Join(dir, "old.txt")
	newPath := filepath.Join(dir, "new.txt")
	content := []byte("hello robustio")

	if err := os.WriteFile(oldPath, content, 0o600); err != nil {
		t.Fatalf("setup WriteFile(): %v", err)
	}

	if err := Rename(oldPath, newPath); err != nil {
		t.Fatalf("Rename(): %v", err)
	}

	data, err := ReadFile(newPath)
	if err != nil {
		t.Fatalf("ReadFile(): %v", err)
	}
	if string(data) != string(content) {
		t.Fatalf("ReadFile(): got %q, want %q", data, content)
	}

	if err := RemoveAll(dir); err != nil {
		t.Fatalf("RemoveAll(): %v", err)
	}
	if _, err := os.Stat(dir); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("RemoveAll(): directory still exists: %v", err)
	}
}
