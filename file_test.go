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
		t.Fatalf("got 0%03o, want 0%03o", got, want)

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
		t.Fatalf("got 0%03o, want 0%03o", got, want)

		return
	}
}

func TestFilePosixCreateWithFileMode(t *testing.T) {
	want := fixPosixPerms(compat.CreatePerm, false) // 0o666

	name, err := tempName(t)
	if err != nil {
		t.Fatalf("tempName failed: %v", err)

		return
	}

	fh, err := compat.Create(name, compat.WithFileMode(perm000))
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
		t.Fatalf("got 0%03o, want 0%03o", got, want)

		return
	}
}

func TestFilePosixCreateEx(t *testing.T) {
	perm := compat.CreatePerm
	want := fixPosixPerms(perm, false)

	name, err := tempName(t)
	if err != nil {
		t.Fatalf("tempName failed: %v", err)

		return
	}

	fh, err := compat.CreateEx(name, perm, 0)
	if err != nil {
		t.Fatalf("CreateEx failed: %v", err)

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

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fi.Mode().Perm()
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

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fi.Mode().Perm()
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
		t.Fatalf("got 0%03o, want 0%03o", got, want)

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

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fi.Mode().Perm()
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

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)

		return
	}
}

func TestFilePosixMkdirTempWithFileMode(t *testing.T) {
	want := fixPosixPerms(compat.MkdirTempPerm, true) // 0o700
	dir := tempDir(t)
	pattern := ""

	name, err := compat.MkdirTemp(dir, pattern, compat.WithFileMode(perm000))
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
		t.Fatalf("got 0%03o, want 0%03o", got, want)

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

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fi.Mode().Perm()
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

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)

		return
	}

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)

		return
	}
}

func TestFilePosixChmodInvalid(t *testing.T) {
	name := "/an/invalid/file/chmod"
	perm := compat.CreatePerm
	err := compat.Chmod(name, perm)
	if err == nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestFilePosixCreateInvalid(t *testing.T) {
	name := "/an/invalid/file/create"
	_, err := compat.Create(name)
	if err == nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestFilePosixCreateExInvalid(t *testing.T) {
	name := "/an/invalid/file/createex"
	perm := compat.CreatePerm
	_, err := compat.CreateEx(name, perm, 0)
	if err == nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestFilePosixCreateTempInvalid(t *testing.T) {
	dir := "/an/invalid/dir/createtemp"
	_, err := compat.CreateTemp(dir, "")
	if err == nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestFilePosixFchmodInvalid(t *testing.T) {
	perm := compat.CreatePerm
	err := compat.Fchmod(nil, perm)
	if err == nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestFilePosixMkdirInvalid(t *testing.T) {
	dir := "/an/invalid/dir/mkdir"
	perm := compat.MkdirTempPerm
	err := compat.Mkdir(dir, perm)
	if err == nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestFilePosixMkdirAllInvalid(t *testing.T) {
	dir := ""
	perm := compat.MkdirTempPerm
	err := compat.MkdirAll(dir, perm)
	if err == nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestFilePosixMkdirTempInvalid(t *testing.T) {
	dir := "/an/invalid/dir/mkdirtemp"
	_, err := compat.MkdirTemp(dir, "")
	if err == nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestFilePosixOpenFileInvalid(t *testing.T) {
	name := "/an/invalid/file/openfile"
	perm := compat.CreatePerm
	_, err := compat.OpenFile(name, os.O_CREATE, perm)
	if err == nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestFilePosixWriteFileInvalid(t *testing.T) {
	name := "/an/invalid/file/writefile"
	perm := compat.CreatePerm
	err := compat.WriteFile(name, helloBytes, perm)
	if err == nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestFilePosixWriteFileExInvalid(t *testing.T) {
	name := "/an/invalid/file/writefileex"
	perm := compat.CreatePerm
	err := compat.WriteFileEx(name, helloBytes, perm, 0)
	if err == nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
