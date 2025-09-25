// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"testing"

	"github.com/rasa/compat"
)

func TestRenameWindowsRetry(t *testing.T) {
	old, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}
	new := old + ".new"
	cleanup(t, old, new)
	err = compat.Rename(old, new, compat.WithRetrySeconds(2))
	if err != nil {
		t.Fatalf("renaming '%v' to '%v': %v", old, new, err)
	}
}
