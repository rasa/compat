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

func TestWriteFileAtomicReplacePlan9(t *testing.T) {
	file, err := tempName(t)
	if err != nil {
		t.Fatal(err)
	}

	cleanup(t, file)

	oldData := []byte("old")
	if err := os.WriteFile(file, oldData, 0o666); err != nil {
		t.Fatal(err)
	}

	err = compat.WriteFile(
		file,
		helloBytes,
		0o666,
		compat.WithAtomicity(true),
	)

	if !isUnsupportedError(err) {
		t.Fatalf("got %v, want unsupported error", err)
	}

	got, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, oldData) {
		t.Fatalf("destination changed: got %q, want %q", got, oldData)
	}
}

func TestWriteFileAtomicCreatePlan9(t *testing.T) {
	file, err := tempName(t)
	if err != nil {
		t.Fatal(err)
	}

	cleanup(t, file)

	err = compat.WriteFile(
		file,
		helloBytes,
		0o666,
		compat.WithAtomicity(true),
	)
	if err != nil {
		t.Fatal(err)
	}
}
