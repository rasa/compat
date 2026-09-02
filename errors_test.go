// SPDX-FileCopyrightText: Copyright (c) 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/rasa/compat"
)

const (
	errStr  = "error message"
	newStr  = "new"
	oldStr  = "old"
	opStr   = "test"
	pathStr = "path"
)

func TestErrorsChmodError(t *testing.T) {
	got := compat.ChmodError(pathStr, os.ErrInvalid).Error()

	want := "chmod " + pathStr + ":"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("ChmodError: got %q; want %q", got, want)
	}
}

func TestErrorsCreateError(t *testing.T) {
	got := compat.CreateError(pathStr, os.ErrInvalid).Error()

	want := "create " + pathStr + ":"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("CreateError: got %q; want %q", got, want)
	}
}

func TestErrorsCreateTempError(t *testing.T) {
	got := compat.CreateTempError(pathStr, os.ErrInvalid).Error()

	want := "createtemp " + pathStr + ":"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("CreateTempError: got %q; want %q", got, want)
	}
}

func TestErrorsMkdirError(t *testing.T) {
	got := compat.MkdirError(pathStr, os.ErrInvalid).Error()

	want := "mkdir " + pathStr + ":"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("MkdirError: got %q; want %q", got, want)
	}
}

func TestErrorsMkdirTempError(t *testing.T) {
	got := compat.MkdirTempError(pathStr, os.ErrInvalid).Error()

	want := "mkdirtemp " + pathStr + ":"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("MkdirTempError: got %q; want %q", got, want)
	}
}

func TestErrorsNiceError(t *testing.T) {
	got := compat.NiceError(errStr, errors.ErrUnsupported).Error()

	want := "nice: " + errStr
	if !strings.HasPrefix(got, want) {
		t.Fatalf("NiceError: got %q; want %q", got, want)
	}
}

func TestErrorsOpenError(t *testing.T) {
	got := compat.OpenError(pathStr, os.ErrInvalid).Error()

	want := "open " + pathStr + ":"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("OpenError: got %q; want %q", got, want)
	}
}

func TestErrorsRenameError(t *testing.T) {
	got := compat.RenameError(oldStr, newStr, os.ErrInvalid).Error()

	want := "rename " + oldStr + " " + newStr + ":"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("RenameError: got %q; want %q", got, want)
	}
}

func TestErrorsStatError(t *testing.T) {
	got := compat.StatError(pathStr, os.ErrInvalid).Error()

	want := "stat " + pathStr + ":"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("StatError: got %q; want %q", got, want)
	}
}

func TestErrorsSymlinkError(t *testing.T) {
	got := compat.SymlinkError(oldStr, newStr, os.ErrInvalid).Error()

	want := "symlink " + oldStr + " " + newStr + ":"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("SymlinkError: got %q; want %q", got, want)
	}
}

func TestErrorsWriteError(t *testing.T) {
	got := compat.WriteError(pathStr, os.ErrInvalid).Error()

	want := "write " + pathStr + ":"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("WriteError: got %q; want %q", got, want)
	}
}

func TestErrorsNormalizeUnsupportedError(t *testing.T) {
	err := compat.NormalizeUnsupportedError(errors.ErrUnsupported)
	if !errors.Is(err, errors.ErrUnsupported) {
		t.Fatalf("normalizeUnsupportedError: got %v, want %v", err, errors.ErrUnsupported)
	}

	err = compat.NormalizeUnsupportedError(os.ErrInvalid)
	if errors.Is(err, errors.ErrUnsupported) {
		t.Fatalf("normalizeUnsupportedError: got %v, want %v", err, os.ErrInvalid)
	}
}
