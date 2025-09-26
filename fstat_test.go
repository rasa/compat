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
	name, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}
	cleanup(t, name)

	f, err := os.Open(name)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close() //nolint:errcheck

	fi, err := compat.Fstat(f)
	if err != nil {
		if !compat.SupportsFstat() {
			return
		}
		t.Fatalf("Fstat: got %v, want nil", err)
	}

	if !compat.SupportsFstat() {
		t.Fatalf("Fstat: got nil, want an error")
	}

	got := fi.Name()
	want := filepath.Base(name)
	if got != want {
		t.Fatalf("Fstat: got %v, want %v", got, want)
	}
}

func TestFstatInvalid(t *testing.T) {
	_, err := compat.Fstat(nil)
	if err == nil {
		t.Fatal("Fstat: got nil, want an error")
	}
}
