// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/rasa/compat"
)

func TestFstat(t *testing.T) {
	if !compat.SupportsFstat() {
		skipf(t, "Skipping test: Fstat() is not supported on %v/%v", runtime.GOOS, runtime.GOARCH)

		return // tinygo doesn't support t.Skip
	}

	name, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}
	cleanup(t, name)

	f, err := os.Open(name)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Fstat(f)
	if err != nil {
		t.Fatalf("Fstat: got %v, want nil", err)
	}

	got := fi.Name()
	want := filepath.Base(name)
	if got != want {
		t.Fatalf("Fstat: got %v, want %v", got, want)
	}
}

func TestFstatInvalid(t *testing.T) {
	if !compat.SupportsFstat() {
		skipf(t, "Skipping test: Fstat() is not supported on %v/%v", runtime.GOOS, runtime.GOARCH)

		return // tinygo doesn't support t.Skip
	}

	_, err := compat.Fstat(nil)
	if err == nil {
		t.Fatal("Fstat: got nil, want an error")
	}
}
