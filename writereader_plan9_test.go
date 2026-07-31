// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: Portions copyright (c) 2015 Nate Finch (@natefinch)
// SPDX-FileCopyrightText: Portions copyright (c) 2022 Simon Dassow (@sdassow)

//go:build plan9

package compat_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/rasa/compat"
)

func TestWriteReaderAtomicReplacePlan9(t *testing.T) {
	file, err := tempName(t)
	if err != nil {
		t.Fatal(err)
	}

	cleanup(t, file)

	oldData := []byte("old")

	// The destination must already exist to test unsupported
	// atomic replacement.
	if err := os.WriteFile(file, oldData, 0o666); err != nil {
		t.Fatal(err)
	}

	err = compat.WriteReader(
		file,
		bytes.NewReader(helloBytes),
		0o666,
		compat.WithAtomicity(true),
	)

	if !isUnsupportedError(err) {
		t.Fatalf(
			"got %v, want unsupported error",
			err,
		)
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

func TestWriteReaderAtomicCreatePlan9(t *testing.T) {
	file, err := tempName(t)
	if err != nil {
		t.Fatal(err)
	}

	cleanup(t, file)

	// The destination does not exist, so atomic creation should succeed.
	err = compat.WriteReader(
		file,
		bytes.NewReader(helloBytes),
		0o666,
		compat.WithAtomicity(true),
	)
	if err != nil {
		t.Fatal(err)
	}
}
