// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: Portions copyright (c) 2015 Nate Finch (@natefinch)
// SPDX-FileCopyrightText: Portions copyright (c) 2022 Simon Dassow (@sdassow)

//go:build plan9

package compat_test

import (
	"errors"
	"testing"

	"github.com/rasa/compat"
)

func TestWriteReaderWithAtomicityPlan9(t *testing.T) { //nolint:dupl
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	cleanup(t, file)

	perm := compat.CreatePerm // 0o666
	opts := []compat.Option{compat.WithAtomicity(true)}
	err = compat.WriteReader(file, helloBuf, perm, opts...)
	if err != errors.ErrUnsupported {
		t.Fatalf("got %v, want ErrUnsupported", err)
	}
}
