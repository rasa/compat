// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"testing"

	"github.com/rasa/compat"
)

func TestGetuid(t *testing.T) {
	uid, err := compat.Getuid()
	if err != nil {
		t.Errorf("Getuid: got %v, want nil", err)

		return // tinygo doesn't support t.Fatal
	}

	if compat.IsWasip1 {
		// Wasip1 returns -1 for UID
		return
	}

	if uid == compat.UnknownID {
		t.Fatalf("Getuid: got %v (UnknownID), want a valid ID", compat.UnknownID)
	}
}

func TestGetgid(t *testing.T) {
	gid, err := compat.Getgid()
	if err != nil {
		t.Errorf("Getgid: got %v, want nil", err)

		return // tinygo doesn't support t.Fatal
	}

	if compat.IsWasip1 {
		// Wasip1 returns -1 for GID
		return
	}

	if gid == compat.UnknownID {
		t.Fatalf("Getgid: got %v (UnknownID), want a valid ID", compat.UnknownID)
	}
}
