// SPDX-FileCopyrightText: Copyright © 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/rasa/compat"
)

func TestWriteReaderNonAtomicReplace(t *testing.T) {
	if compat.SupportsAtomicReplace() {
		skip(t, "Skipping test: requires non-atomic rename")
		return
	}

	file, err := tempName(t)
	if err != nil {
		t.Fatal(err)
	}

	cleanup(t, file)

	oldData := []byte("old")

	// The destination must already exist to test unsupported
	// atomic replacement.
	err = os.WriteFile(file, oldData, 0o600)
	if err != nil {
		t.Fatal(err)
	}

	err = compat.WriteReader(
		file,
		bytes.NewReader(helloBytes),
		0o600,
		compat.WithAtomicity(true),
	)

	if !compat.IsUnsupportedError(err) {
		t.Fatalf("got %v, want unsupported", err)
	}

	got, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, oldData) {
		t.Fatalf(
			"destination changed: got %q, want %q",
			got,
			oldData,
		)
	}
}

func TestWriteReaderNonAtomicCreate(t *testing.T) {
	if compat.SupportsAtomicReplace() {
		skip(t, "Skipping test: requires non-atomic rename")
		return
	}

	file, err := tempName(t)
	if err != nil {
		t.Fatal(err)
	}

	cleanup(t, file)

	// The destination does not exist, so atomic creation should succeed.
	err = compat.WriteReader(
		file,
		bytes.NewReader(helloBytes),
		0o6,
		compat.WithAtomicity(true),
	)
	if err != nil {
		t.Fatal(err)
	}
}
