// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/rasa/compat"
)

func TestFstatAbs(t *testing.T) {
	if !compat.SupportsFstat() {
		skipf(t, "skipping test: Fstat() not supported on %v/%v", runtime.GOOS, runtime.GOARCH)

		return
	}

	name, err := tempFile(t)
	if err != nil {
		fatal(t, err)

		return
	}

	cleanup(t, name)

	f, err := os.Open(name)
	if err != nil {
		fatal(t, err)

		return
	}
	defer f.Close()

	fi, err := compat.Fstat(f)
	if err != nil {
		fatalf(t, "Fstat: got %v, want nil", err)

		return
	}

	got := fi.Name()

	want := filepath.Base(name)
	if got != want {
		fatalf(t, "Fstat: got %v, want %v", got, want)

		return
	}
}

func TestFstatRelative(t *testing.T) {
	name, err := tempFile(t)
	if err != nil {
		fatal(t, err)

		return
	}

	cleanup(t, name)

	dir, _ := filepath.Split(name)
	name = strings.TrimPrefix(name, dir)
	t.Chdir(dir)

	f, err := os.Open(name)
	if err != nil {
		fatal(t, err)

		return
	}
	defer f.Close()

	fi, err := compat.Fstat(f)
	if err != nil {
		fatalf(t, "Fstat: got %v, want nil", err)

		return
	}

	got := fi.Name()

	want := filepath.Base(name)
	if got != want {
		fatalf(t, "Fstat: got %v, want %v", got, want)

		return
	}
}

func TestFstatRelativeChdir(t *testing.T) {
	name1, err := tempFile(t)
	if err != nil {
		fatal(t, err)

		return
	}

	cleanup(t, name1)

	f1, err := os.Open(name1)
	if err != nil {
		fatal(t, err)

		return
	}

	fi1, err := compat.Fstat(f1)
	if err != nil {
		fatalf(t, "Fstat: got %v, want nil", err)

		return
	}

	_ = f1.Close()

	dir, name2 := filepath.Split(name1)
	t.Chdir(dir)

	f2, err := os.Open(name2)
	if err != nil {
		fatal(t, err)

		return
	}
	defer f2.Close()

	fi2, err := compat.Fstat(f2)
	if err != nil {
		fatalf(t, "Fstat: got %v, want nil", err)

		return
	}

	want := fi1.FileID()
	got := fi2.FileID()

	if got != want {
		fatalf(t, "Fstat (Relative): got %v, want %v (for %v and %v)", got, want, name1, name2)

		return
	}
}

func TestFstatInvalidNil(t *testing.T) {
	_, err := compat.Fstat(nil)
	if err == nil {
		fatal(t, "Fstat: got nil, want an error")

		return
	}
}
