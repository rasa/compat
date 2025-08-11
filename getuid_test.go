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
		t.Fatalf("Getuid: got %v, want nil", err)

		return
	}

	if uid == compat.UnknownID {
		t.Fatalf("Getuid: got %v (UnknownID), want a valid ID", compat.UnknownID)
	}
}

func TestGetgid(t *testing.T) {
	gid, err := compat.Getgid()
	if err != nil {
		t.Fatalf("Getgid: got %v, want nil", err)

		return
	}

	if gid == compat.UnknownID {
		t.Fatalf("Getgid: got %v (UnknownID), want a valid ID", compat.UnknownID)
	}
}
