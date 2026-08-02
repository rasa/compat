// SPDX-FileCopyrightText: Copyright (c) 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/rasa/compat"
)

func TestIsUsupportedError(t *testing.T) {
	err := &compat.UnsupportedError{Op: "test1"}
	got := compat.IsUnsupportedError(err)
	want := true
	if got != want {
		t.Fatalf("IsUsupportedError: got %v, want %v", got, want)
	}
}

func TestUnsupportedError(t *testing.T) {
	err := &compat.UnsupportedError{Op: "test2"}
	got := err.Error()
	want := "test: unsupported"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("UnsupportedError: got %v; want %v", got, want)
	}
	if !errors.Is(err, errors.ErrUnsupported) {
		t.Fatalf("UnsupportedError: got %v, want ErrUnsupported", err)
	}
}

func TestUnimplementedError(t *testing.T) {
	err := &compat.UnimplementedError{Op: "test3"}
	got := err.Error()
	want := "test: unimplemented"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("UnimplementedError: got %v; want %v", got, want)
	}
	if !errors.Is(err, errors.ErrUnsupported) {
		t.Fatalf("UnimplementedError: got %v; want ErrUnsupported", err)
	}
}
