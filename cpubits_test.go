// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"testing"

	"github.com/rasa/compat"
)

func TestCPUBits(t *testing.T) {
	_, err := compat.CPUBits()
	if err != nil {
		t.Fatalf("CPUBits: got %v, want nil", err)
	}
}
