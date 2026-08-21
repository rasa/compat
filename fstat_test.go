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
		skipf(t, "Fstat() not supported on %v/%v", runtime.GOOS, runtime.GOARCH)

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

func TestFstatRel(t *testing.T) {
	if !compat.SupportsFstat() {
		skipf(t, "Fstat() not supported on %v/%v", runtime.GOOS, runtime.GOARCH)

		return
	}

	name, err := tempFile(t)
	if err != nil {
		fatal(t, err)

		return
	}

	cleanup(t, name)

	dir := tempDir(t)
	t.Chdir(dir)
	name = strings.TrimPrefix(name, dir)

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

func TestFstatRelativeOK(t *testing.T) {
	if !compat.SupportsRelativeFstat() {
		skipf(t, "Relative Fstat() not supported on %v/%v", runtime.GOOS, runtime.GOARCH)

		return
	}

	name, err := tempFile(t)
	if err != nil {
		fatal(t, err)

		return
	}

	cleanup(t, name)

	dir := tempDir(t)
	newDir := filepath.Clean(filepath.Join(dir, ".."))
	t.Chdir(newDir)

	name = strings.TrimPrefix(name, dir)

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

func TestFstatRelativeFail(t *testing.T) {
	if compat.SupportsRelativeFstat() {
		skipf(t, "Relative Fstat() is supported on %v/%v", runtime.GOOS, runtime.GOARCH)

		return
	}

	name, err := tempFile(t)
	if err != nil {
		fatal(t, err)

		return
	}

	cleanup(t, name)

	dir := tempDir(t)
	newDir := filepath.Clean(filepath.Join(dir, ".."))
	t.Chdir(newDir)

	name = strings.TrimPrefix(name, dir)

	f, err := os.Open(name)
	if err != nil {
		fatal(t, err)

		return
	}
	defer f.Close()

	fi, err := compat.Fstat(f)
	if err == nil {
		fatal(t, "Relative Fstat: got nil, want error")

		return
	}

	got := fi.Name()

	want := filepath.Base(name)
	if got != want {
		fatalf(t, "Fstat: got %v, want %v", got, want)

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
