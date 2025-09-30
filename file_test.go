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

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o (%v), want 0%03o (%v)", got, got, want, want)

		return
	}
}

func TestFilePosixCreate(t *testing.T) {
	want := fixPosixPerms(compat.CreatePerm, false) // 0o666

	name, err := tempName(t)
	if err != nil {
		t.Fatalf("tempName failed: %v", err)

		return
	}

	fh, err := compat.Create(name)
	if err != nil {
		t.Fatalf("Create failed: %v", err)

		return
	}

	err = fh.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)

		return
	}

	fi, err := os.Stat(name)
	if err != nil {
		if compat.IsTinygo && errors.Is(err, os.ErrNotExist) {
			skip(t, "Skipping test: file is disappearing on tinygo")

			return // tinygo doesn't support t.Skip
		}

		t.Fatalf("Stat failed: %v", err)

		return
	}

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o (%v), want 0%03o (%v)", got, got, want, want)

		return
	}
}

func TestFilePosixCreateWithFileMode(t *testing.T) {
	perm := compat.CreatePerm
	want := fixPosixPerms(perm, false) // 0o666

	name, err := tempName(t)
	if err != nil {
		t.Fatalf("tempName failed: %v", err)

		return
	}

	fh, err := compat.Create(name, compat.WithFileMode(perm))
	if err != nil {
		t.Fatalf("Create failed: %v", err)

		return
	}

	err = fh.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)

		return
	}

	fi, err := os.Stat(name)
	if err != nil {
		if compat.IsTinygo && errors.Is(err, os.ErrNotExist) {
			skip(t, "Skipping test: file is disappearing on tinygo")

			return // tinygo doesn't support t.Skip
		}

		t.Fatalf("Stat failed: %v", err)

		return
	}

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o (%v), want 0%03o (%v)", got, got, want, want)

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

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o (%v), want 0%03o (%v)", got, got, want, want)

		return
	}
}

func TestFilePosixFchmod(t *testing.T) {
	perm := os.FileMode(0o644)
	want := fixPosixPerms(perm, false)

	name, err := tempFile(t)
	if err != nil {
		t.Fatal(err)

		return
	}

	f, err := os.Open(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	defer fclose(f)

	err = compat.Fchmod(f, perm)
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
		t.Fatalf("got 0%03o (%v), want 0%03o (%v)", got, got, want, want)

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

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o (%v), want 0%03o (%v)", got, got, want, want)

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

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o (%v), want 0%03o (%v)", got, got, want, want)

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

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o (%v), want 0%03o (%v)", got, got, want, want)

		return
	}
}

func TestFilePosixMkdirTempWithFileMode(t *testing.T) {
	perm := compat.MkdirTempPerm // 0o700
	want := fixPosixPerms(perm, true)
	dir := tempDir(t)
	pattern := ""

	name, err := compat.MkdirTemp(dir, pattern, compat.WithFileMode(perm))
	if err != nil {
		t.Fatal(err)

		return
	}

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o (%v), want 0%03o (%v)", got, got, want, want)

		return
	}
}

func TestFilePosixOpenFile(t *testing.T) {
	perm := os.FileMode(0o666)
	want := fixPosixPerms(perm, false)

	name, err := tempName(t)
	if err != nil {
		t.Fatalf("tempName failed: %v", err)

		return
	}

	fh, err := compat.OpenFile(name, os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		t.Fatalf("OpenFile failed: %v", err)

		return
	}

	err = fh.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)

		return
	}

	fi, err := os.Stat(name)
	if err != nil {
		if compat.IsTinygo && errors.Is(err, os.ErrNotExist) {
			skip(t, "Skipping test: file is disappearing on tinygo")

			return // tinygo doesn't support t.Skip
		}

		t.Fatalf("Stat failed: %v", err)

		return
	}

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o (%v), want 0%03o (%v)", got, got, want, want)

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

func TestFilePosixRemove(t *testing.T) {
	name, err := tempFile(t)
	if err != nil {
		t.Fatal(err)

		return
	}

	err = compat.Remove(name)
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

func TestFilePosixRemoveAll(t *testing.T) {
	name, err := tempFile(t)
	if err != nil {
		t.Fatal(err)

		return
	}

	err = compat.RemoveAll(name)
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

func TestFilePosixSymlink(t *testing.T) {
	if !supportsSymlinks(t) {
		return
	}

	old, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	new := old + ".link"
	err = compat.Symlink(old, new)
	if err != nil {
		t.Fatalf("Symlink: %q to %q: %v", old, new, err)
	}
}

func TestFilePosixChmodInvalid(t *testing.T) {
	err := compat.Chmod(invalidName, compat.CreatePerm)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestFilePosixCreateInvalid(t *testing.T) {
	_, err := compat.Create(invalidName)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestFilePosixCreateTempInvalid(t *testing.T) {
	_, err := compat.CreateTemp(invalidName, invalidName)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestFilePosixFchmodInvalid(t *testing.T) {
	err := compat.Fchmod(nil, compat.CreatePerm)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestFilePosixMkdirInvalid(t *testing.T) {
	err := compat.Mkdir(invalidName, compat.MkdirTempPerm)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestFilePosixMkdirAllInvalid(t *testing.T) {
	err := compat.MkdirAll(invalidName, compat.MkdirTempPerm)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestFilePosixMkdirTempInvalid(t *testing.T) {
	_, err := compat.MkdirTemp(invalidName, invalidName)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestFilePosixOpenFileInvalid(t *testing.T) {
	_, err := compat.OpenFile(invalidName, os.O_CREATE, compat.CreatePerm)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestFilePosixSymlinkInvalidOld(t *testing.T) {
	if !supportsSymlinks(t) {
		return
	}

	old := invalidName
	new := old + ".link"
	err := compat.Symlink(old, new)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestFilePosixSymlinkInvalidNew(t *testing.T) {
	if !supportsSymlinks(t) {
		return
	}

	old, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	new := invalidName
	err = compat.Symlink(old, new)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}
