// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"testing"
	"time"

	"github.com/rasa/compat"
)

func TestRenameWindowsRetry(t *testing.T) {
	src, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	dst := src + ".new"
	cleanup(t, src, dst)

	err = compat.Rename(src, dst, compat.WithRetryTimeout(2*time.Second))
	if err != nil {
		t.Fatalf("renaming %q to %q: %v", src, dst, err)
	}
}
