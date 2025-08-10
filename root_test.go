// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !windows

package compat_test

import (
	"os"
	"testing"

	"github.com/rasa/compat"
)

func TestRootIsRoot(t *testing.T) {
	isRoot := os.Getuid() == 0

	got, err := compat.IsRoot()
	if err != nil {
		t.Fatalf("IsRoot() returned: %v", err)
	}

	if !isRoot {
		skipf(t, "Skipping test: we aren't the root/admin user")

		return
	}

	if got != isRoot {
		t.Fatalf("IsRoot(): got %v, want %v", got, isRoot)

		return
	}
}
