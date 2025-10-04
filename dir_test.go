// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/rasa/compat"
)

// Source: https://github.com/golang/go/blob/ac803b59/src/os/read_test.go#L104-L144

func TestReadDir(t *testing.T) {
	// t.Parallel()
	if compat.IsTinygo {
		skip(t, "Skipping test: fdopendir /tmp/TestReadDir256423683/000/foo: errno 8")

		return // tinygo doesn't support t.Skip
	}

	dirname := "rumpelstilzchen"
	if _, err := compat.ReadDir(dirname); err == nil { // compat: s|ReadDir|compat.ReadDir|
		t.Fatalf("ReadDir %s: error expected, none found", dirname)
	}

	filename := filepath.Join(t.TempDir(), "foo")
	f, err := os.Create(filename) //nolint:govet // compat: s|Create|os.Create|
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	if list, err := compat.ReadDir(filename); list != nil || !errors.Is(err, syscall.ENOTDIR) { //nolint:govet // compat: s|ReadDir|compat.ReadDir|
		t.Fatalf("ReadDir %s: (nil, ENOTDIR) expected, got (%v, %v)", filename, list, err)
	}

	dirname = "."
	list, err := compat.ReadDir(dirname) //nolint:govet // compat: s|ReadDir|compat.ReadDir|
	if err != nil {
		t.Fatalf("ReadDir %s: %v", dirname, err)
	}

	foundFile := false
	foundSubDir := false
	for _, dir := range list {
		switch {
		case !dir.IsDir() && dir.Name() == "dir_test.go":
			foundFile = true
		case dir.IsDir() && dir.Name() == "golang":
			foundSubDir = true
		}
	}
	if !foundFile {
		t.Fatalf("ReadDir %s: dir_test.go file not found", dirname)
	}
	if !foundSubDir {
		t.Fatalf("ReadDir %s: golang directory not found", dirname)
	}
}

func TestDirEntry(t *testing.T) {
	dirname := "."
	list, err := compat.ReadDir(dirname)
	if err != nil {
		t.Fatalf("ReadDir %s: %v", dirname, err)
	}

	for _, dir := range list {
		info, err := dir.Info()
		if err != nil {
			t.Fatalf("ReadDir %s: %v", dirname, err)
		}
		if info == nil {
			t.Fatalf("ReadDir %s: %v", dirname, "info is nil")
		}
		if info.Name() != dir.Name() {
			t.Fatalf("ReadDir %s: Name(): got %v; want %v", dirname, info.Name(), dir.Name())
		}
		if info.IsDir() != dir.IsDir() {
			t.Fatalf("ReadDir %s: IsDir(): got %v; want %v", dirname, info.IsDir(), dir.IsDir())
		}
		if info.Mode().Type() != dir.Type() {
			t.Fatalf("ReadDir %s: Type(): got %v; want %v", dirname, info.Mode().Type(), dir.Type())
		}
	}
}
