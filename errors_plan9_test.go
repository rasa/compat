// SPDX-FileCopyrightText: Copyright (c) 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9

package compat_test

import (
	"errors"
	"syscall"
	"testing"

	"github.com/rasa/compat"
)

func TestErrorsNormalizeUnsupportedErrorPlan9(t *testing.T) {
	err := compat.NormalizeUnsupportedError(syscall.EPLAN9)
	if !errors.Is(err, errors.ErrUnsupported) {
		t.Fatalf("got %v, want %v", err, errors.ErrUnsupported)
	}
}
