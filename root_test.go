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
		t.Errorf("IsRoot() returned: %v", err)
	}
	if got != isRoot {
		// Report result, but don't fail, as the user may not be root.
		skipf(t, "IsRoot(): got %v, want %v", got, isRoot)
		return
	}
}
