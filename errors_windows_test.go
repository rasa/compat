// SPDX-FileCopyrightText: Copyright (c) 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat_test

import (
	"errors"
	"syscall"
	"testing"

	"github.com/rasa/compat"
)

func TestErrorsNormalizeUnsupportedErrorWindows(t *testing.T) {
	err := compat.NormalizeUnsupportedError(syscall.EWINDOWS)
	if !errors.Is(err, errors.ErrUnsupported) {
		t.Fatalf("got %v, want %v", err, errors.ErrUnsupported)
	}
}
