// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"errors"
	"os"
	"testing"

	"github.com/rasa/compat"
)

func TestFileOptionsCreateDelete(t *testing.T) {
	name, err := tempName(t)
	if err != nil {
		t.Fatal(err)
	}

	fh, err := compat.Create(name, compat.WithFlags(compat.O_FILE_FLAG_DELETE_ON_CLOSE))
	if err != nil {
		t.Fatal(err)
	}

	err = fh.Close()
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(name)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatal("File exists, should not")
	}
}

func TestFileOptionsCreateExcl(t *testing.T) {
	name, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	fh, err := compat.Create(name, compat.WithFlags(os.O_EXCL))
	if err == nil {
		_ = fh.Close()
		t.Fatal("got no error")
	}
}

func TestFileOptionsCreateTempDelete(t *testing.T) {
	dir := tempDir(t)

	fh, err := compat.CreateTemp(dir, "", compat.WithFlags(compat.O_FILE_FLAG_DELETE_ON_CLOSE))
	if err != nil {
		t.Fatal(err)
	}

	name := fh.Name()

	err = fh.Close()
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(name)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatal("File exists, should not")
	}
}

func TestFileOptionsCreateTempFileMode(t *testing.T) {
	want := fixPosixPerms(0o777, false)

	dir := tempDir(t)

	fh, err := compat.CreateTemp(dir, "", compat.WithFileMode(want))
	if err != nil {
		t.Fatal(err)
	}

	name := fh.Name()

	err = fh.Close()
	if err != nil {
		t.Fatal(err)
	}

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)
	}
}

func TestFileOptionsMkdirTempFileMode(t *testing.T) {
	want := fixPosixPerms(0o777, true)
	dir := tempDir(t)
	pattern := ""

	name, err := compat.MkdirTemp(dir, pattern, compat.WithFileMode(want))
	if err != nil {
		t.Fatal(err)
	}

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)
	}
}

func TestFileOptionsOpenFileDelete(t *testing.T) {
	name, err := tempName(t)
	if err != nil {
		t.Fatal(err)
	}

	fh, err := compat.OpenFile(name, os.O_RDWR|os.O_CREATE, os.FileMode(0o666), compat.WithFlags(compat.O_FILE_FLAG_DELETE_ON_CLOSE))
	if err != nil {
		// workaround:
		// https://github.com/rasa/compat/actions/runs/16542086538/job/46784707170#step:6:48
		if compat.IsApple {
			t.Skip(err)
		}
		t.Fatal(err)
	}

	err = fh.Close()
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(name)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatal("File exists, should not")
	}
}

func TestFileOptionsOpenFileFileMode(t *testing.T) {
	perm := os.FileMode(0o666)
	want := fixPosixPerms(perm, false)

	name, err := tempName(t)
	if err != nil {
		t.Fatal(err)
	}

	fh, err := compat.OpenFile(name, os.O_RDWR|os.O_CREATE, 0, compat.WithFileMode(want))
	if err != nil {
		t.Fatal(err)
	}

	err = fh.Close()
	if err != nil {
		t.Fatal(err)
	}

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)
	}
}

func TestFileOptionsWriteFileFileMode(t *testing.T) {
	perm := os.FileMode(0o764)
	want := fixPosixPerms(perm, false)

	name, err := tempName(t)
	if err != nil {
		t.Fatal(err)
	}

	err = compat.WriteFile(name, helloBytes, 0, compat.WithFileMode(want))
	if err != nil {
		t.Fatal(err)
	}

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)
	}
}
