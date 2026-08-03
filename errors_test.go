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

func TestErrorsIsUsupportedError(t *testing.T) {
	err := &compat.UnsupportedError{Op: "test"} //nolint:goconst
	got := compat.IsUnsupportedError(err)
	want := true
	if got != want {
		t.Fatalf("IsUsupportedError: got %v, want %v", got, want)
	}
}

func TestErrorsUnsupportedError(t *testing.T) {
	err := &compat.UnsupportedError{Op: "test"} //nolint:goconst
	got := err.Error()
	want := "test: unsupported"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("UnsupportedError: got %q; want %q", got, want)
	}
	if !errors.Is(err, errors.ErrUnsupported) {
		t.Fatalf("UnsupportedError: got %v, want ErrUnsupported", err)
	}
}

func TestErrorsUnimplementedError(t *testing.T) {
	err := &compat.UnimplementedError{Op: "test"} //nolint:goconst
	got := err.Error()
	want := "test: unimplemented"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("UnimplementedError: got %q; want %q", got, want)
	}
	if !errors.Is(err, errors.ErrUnsupported) {
		t.Fatalf("UnimplementedError: got %v; want ErrUnsupported", err)
	}
}

func TestErrorsExportedUnsupportedError(t *testing.T) {
	got := compat.ExportedUnsupportedError("test").Error() //nolint:goconst
	want := "test: unsupported"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("ExportedUnsupportedError: got %q; want %q", got, want)
	}
}

func TestErrorsExportedUnimplementedError(t *testing.T) {
	got := compat.ExportedUnimplementedError("test").Error() //nolint:goconst
	want := "test: unimplemented"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("ExportedUnimplementedError: got %q; want %q", got, want)
	}
}

func TestErrorsChmodError(t *testing.T) {
	got := compat.ChmodError("path", os.ErrInvalid).Error() //nolint:goconst
	want := "chmod path:"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("ChmodError: got %q; want %q", got, want)
	}
}

func TestErrorsCreateError(t *testing.T) {
	got := compat.CreateError("path", os.ErrInvalid).Error() //nolint:goconst
	want := "create path:"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("CreateError: got %q; want %q", got, want)
	}
}

func TestErrorsCreateTempError(t *testing.T) {
	got := compat.CreateTempError("path", os.ErrInvalid).Error() //nolint:goconst
	want := "createtemp path:"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("CreateTempError: got %q; want %q", got, want)
	}
}

func TestErrorsMkdirError(t *testing.T) {
	got := compat.MkdirError("path", os.ErrInvalid).Error() //nolint:goconst
	want := "mkdir path:"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("MkdirError: got %q; want %q", got, want)
	}
}

func TestErrorsMkdirallError(t *testing.T) {
	got := compat.MkdirallError("path", os.ErrInvalid).Error() //nolint:goconst
	want := "mkdir path:"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("MkdirallError: got %q; want %q", got, want)
	}
}

func TestErrorsMkdirTempError(t *testing.T) {
	got := compat.MkdirTempError("path", os.ErrInvalid).Error() //nolint:goconst
	want := "mkdirtemp path:"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("MkdirTempError: got %q; want %q", got, want)
	}
}

func TestErrorsOpenError(t *testing.T) {
	got := compat.OpenError("path", os.ErrInvalid).Error() //nolint:goconst
	want := "open path:"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("OpenError: got %q; want %q", got, want)
	}
}

func TestErrorsRenameError(t *testing.T) {
	got := compat.RenameError("old", "new", os.ErrInvalid).Error()
	want := "rename old new:"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("RenameError: got %q; want %q", got, want)
	}
}

func TestErrorsStatError(t *testing.T) {
	got := compat.StatError("path", os.ErrInvalid).Error() //nolint:goconst
	want := "stat path:"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("StatError: got %q; want %q", got, want)
	}
}

func TestErrorsSymlinkError(t *testing.T) {
	got := compat.SymlinkError("old", "new", os.ErrInvalid).Error()
	want := "symlink old new:"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("SymlinkError: got %q; want %q", got, want)
	}
}

func TestErrorsWriteError(t *testing.T) {
	got := compat.WriteError("path", os.ErrInvalid).Error() //nolint:goconst
	want := "write path:"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("WriteError: got %q; want %q", got, want)
	}
}
