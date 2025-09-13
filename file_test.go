// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"errors"
	"os"
	"testing"

	"github.com/rasa/compat"
)

func TestFilePosixChmod(t *testing.T) {
	perm := os.FileMode(0o644)
	want := fixPosixPerms(perm, false)

	name, err := tempFile(t)
	if err != nil {
		t.Fatal(err)

		return
	}

	err = compat.Chmod(name, perm)
	if err != nil {
		t.Fatal(err)

		return
	}

	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)

		return
	}
}

func TestFilePosixCreate(t *testing.T) {
	want := fixPosixPerms(compat.CreatePerm, false) // 0o666

	name, err := tempName(t)
	if err != nil {
		t.Fatal(err)

		return
	}

	fh, err := compat.Create(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	err = fh.Close()
	if err != nil {
		t.Fatal(err)

		return
	}

	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)

		return
	}
}

func TestFilePosixCreateEx(t *testing.T) {
	perm := compat.CreatePerm
	want := fixPosixPerms(perm, false)

	name, err := tempName(t)
	if err != nil {
		t.Fatal(err)

		return
	}

	fh, err := compat.CreateEx(name, perm, 0)
	if err != nil {
		t.Fatal(err)

		return
	}

	err = fh.Close()
	if err != nil {
		t.Fatal(err)

		return
	}

	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)

		return
	}
}

func TestFilePosixCreateExDelete(t *testing.T) {
	name, err := tempName(t)
	if err != nil {
		t.Fatal(err)

		return
	}

	fh, err := compat.CreateEx(name, compat.CreatePerm, os.O_CREATE|compat.O_FILE_FLAG_DELETE_ON_CLOSE)
	if err != nil {
		t.Fatal(err)

		return
	}

	err = fh.Close()
	if err != nil {
		t.Fatal(err)

		return
	}

	_, err = os.Stat(name)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatal("File exists, should not")

		return
	}
}

func TestFilePosixCreateTemp(t *testing.T) {
	want := fixPosixPerms(compat.CreateTempPerm, false) // 0o600

	dir := tempDir(t)

	fh, err := compat.CreateTemp(dir, "")
	if err != nil {
		t.Fatal(err)

		return
	}

	name := fh.Name()

	err = fh.Close()
	if err != nil {
		t.Fatal(err)

		return
	}

	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)

		return
	}
}

func TestFilePosixCreateTempEx(t *testing.T) {
	want := fixPosixPerms(compat.CreateTempPerm, false) // 0o600

	dir := tempDir(t)

	fh, err := compat.CreateTempEx(dir, "", 0)
	if err != nil {
		t.Fatal(err)

		return
	}

	name := fh.Name()

	err = fh.Close()
	if err != nil {
		t.Fatal(err)

		return
	}

	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)

		return
	}
}

func TestFilePosixCreateTempExDelete(t *testing.T) {
	dir := tempDir(t)

	fh, err := compat.CreateTempEx(dir, "", compat.O_FILE_FLAG_DELETE_ON_CLOSE)
	if err != nil {
		t.Fatal(err)

		return
	}

	name := fh.Name()

	err = fh.Close()
	if err != nil {
		t.Fatal(err)

		return
	}

	_, err = os.Stat(name)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatal("File exists, should not")

		return
	}
}

func TestFilePosixMkdir(t *testing.T) {
	perm := os.FileMode(0o777)
	want := fixPosixPerms(perm, true)

	name, err := tempName(t)
	if err != nil {
		t.Fatal(err)

		return
	}

	err = compat.Mkdir(name, perm)
	if err != nil {
		t.Fatal(err)

		return
	}

	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)

		return
	}
}

func TestFilePosixMkdirAll(t *testing.T) {
	perm := os.FileMode(0o777)
	want := fixPosixPerms(perm, true)

	name, err := tempName(t)
	if err != nil {
		t.Fatal(err)

		return
	}

	err = compat.MkdirAll(name, perm)
	if err != nil {
		t.Fatal(err)

		return
	}

	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)

		return
	}
}

func TestFilePosixMkdirTemp(t *testing.T) {
	want := fixPosixPerms(compat.MkdirTempPerm, true) // 0o700
	dir := tempDir(t)
	pattern := ""

	name, err := compat.MkdirTemp(dir, pattern)
	if err != nil {
		t.Fatal(err)

		return
	}

	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)

		return
	}
}

func TestFilePosixOpenFile(t *testing.T) {
	perm := os.FileMode(0o666)
	want := fixPosixPerms(perm, false)

	name, err := tempName(t)
	if err != nil {
		t.Fatal(err)

		return
	}

	fh, err := compat.OpenFile(name, os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		t.Fatal(err)

		return
	}

	err = fh.Close()
	if err != nil {
		t.Fatal(err)

		return
	}

	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)

		return
	}
}

func TestFilePosixOpenFileDelete(t *testing.T) {
	name, err := tempName(t)
	if err != nil {
		t.Fatal(err)

		return
	}

	fh, err := compat.OpenFile(name, os.O_RDWR|os.O_CREATE|compat.O_FILE_FLAG_DELETE_ON_CLOSE, os.FileMode(0o666))
	if err != nil {
		// workaround:
		// https://github.com/rasa/compat/actions/runs/16542086538/job/46784707170#step:6:48
		if compat.IsApple {
			t.Skip(err)
		} else {
			t.Fatal(err)
		}

		return
	}

	err = fh.Close()
	if err != nil {
		t.Fatal(err)

		return
	}

	_, err = os.Stat(name)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatal("File exists, should not")

		return
	}
}

func TestFilePosixWriteFile(t *testing.T) {
	perm := os.FileMode(0o666)
	want := fixPosixPerms(perm, false)

	name, err := tempName(t)
	if err != nil {
		t.Fatal(err)

		return
	}

	err = compat.WriteFile(name, helloBytes, want)
	if err != nil {
		t.Fatal(err)

		return
	}

	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)

		return
	}
}

func TestFilePosixWriteFileEx(t *testing.T) {
	perm := os.FileMode(0o666)
	want := fixPosixPerms(perm, false)

	name, err := tempName(t)
	if err != nil {
		t.Fatal(err)

		return
	}

	err = compat.WriteFileEx(name, helloBytes, perm, 0)
	if err != nil {
		t.Fatal(err)

		return
	}

	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)

		return
	}
}
