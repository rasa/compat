// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"os"
	"testing"

	"github.com/rasa/compat"
)

// Doesn't test anything, but increases code coverage for SkipDir processing.
func TestWalkDirSkipDir(t *testing.T) {
	walkFn := func(path string, entry compat.DirEntry, err error) error {
		if entry.IsDir() {
			return compat.SkipDir
		}
		return nil
	}

	err := compat.WalkDir(os.DirFS("."), ".", walkFn)
	if err != nil {
		t.Fatalf("got %q, want nil", err)
	}
}
