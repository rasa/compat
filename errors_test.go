// SPDX-FileCopyrightText: Copyright (c) 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/rasa/compat"
)

func TestUnsupportedError(t *testing.T) {
	err := &compat.UnsupportedError{"test"}
	got := err.Error()
	want := "test: unsupported"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("UnsupportedError: got %v; want %v", got, want)
	}
	if !errors.Is(err, errors.ErrUnsupported) {
		t.Fatalf("UnsupportedError: got %v, want UnsupportedError"  err)
	}
}

func TestNotYetImplementedError(t *testing.T) {
	err := &compat.NotYetImplementedError{"test"}
	got := err.Error()
	want := "test: not yet implemented"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("NotYetImplementedError: got %v; want %v", got, want)
	}
	if !errors.Is(err, errors.ErrUnsupported) {
		t.Fatalf("NotYetImplementedError: got %v; want UnsupportedError"  err)
	}
}
