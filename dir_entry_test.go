// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"testing"

	"github.com/rasa/compat"
)

func TestDirEntryFileInfoToDirEntryNil(t *testing.T) {
	de := compat.FileInfoToDirEntry(nil, "")
	if de != nil {
		t.Fatalf("expected nil, got a %T", de)
	}
}

func TestDirEntryOSDirEntryToDirEntryNil(t *testing.T) {
	de := compat.OSDirEntryToDirEntry(nil, "")
	if de != nil {
		t.Fatalf("expected nil, got a %T", de)
	}
}

func TestDirEntryFSDirEntryToDirEntryNil(t *testing.T) {
	de := compat.FSDirEntryToDirEntry(nil, "")
	if de != nil {
		t.Fatalf("expected nil, got a %T", de)
	}
}

func TestDirEntryFSFileInfoToDirEntryNil(t *testing.T) {
	de := compat.FSFileInfoToDirEntry(nil, "")
	if de != nil {
		t.Fatalf("expected nil, got a %T", de)
	}
}
