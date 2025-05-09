// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !windows

package compat_test

import (
	"os"
	"testing"

	"github.com/rasa/compat"
)

func TestIsAdmin(t *testing.T) {
	isAdmin := os.Getuid() == 0
	got, err := compat.IsAdmin()
	if err != nil {
		t.Errorf("IsAdmin() returned: %v", err)
	}
	if got != isAdmin {
		// Report result, but don't fail, as the user may not be an admin.
		t.Skipf("IsAdmin(): got %v, want %v", got, isAdmin)
	}
}
