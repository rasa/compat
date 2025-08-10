// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"errors"
	"os"
	"testing"

	"github.com/rasa/compat"
)

var (
	wantCreatePerm     = compat.CreatePerm     // 0o666
	wantCreateTempPerm = compat.CreateTempPerm // 0o600
	wantMkdirTempPerm  = compat.MkdirTempPerm  // 0o700

	helloBytes = []byte("hello")
)

func init() {
	// @TODO(rasa): test different umask settings
	compat.Umask(0)

	wantCreatePerm = fixPosixPerms(wantCreatePerm, false)
	wantCreateTempPerm = fixPosixPerms(wantCreateTempPerm, false)
	wantMkdirTempPerm = fixPosixPerms(wantMkdirTempPerm, true)
}

func fixPerms(perm os.FileMode) os.FileMode {
	if compat.IsWasip1 {
		if compat.IsTinygo {
			// Tinygo's os.Stat() returns mode 0o000
			return os.FileMode(0o000)
		} else {
			return perm & 0o700
		}
	}

	return perm
}

func fixPosixPerms(perm os.FileMode, isDir bool) os.FileMode {
	if compat.IsWindows {
		if isDir {
			return os.FileMode(0o777)
		} else {
			return os.FileMode(0o666)
		}
	}
	return fixPerms(perm)
}

func TestFilePosixChmod(t *testing.T) {
	perm := os.FileMode(0o644)
	want := fixPosixPerms(perm, false)

	name, err := tmpfile(t)
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
	want := wantCreatePerm

	name, err := tmpname(t)
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

	name, err := tmpname(t)
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
	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)

		return
	}

	fh, err := compat.CreateEx(name, compat.CreatePerm, compat.O_CREATE|compat.O_DELETE)
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
	want := wantCreateTempPerm

	dir := t.TempDir()

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
	want := wantCreateTempPerm

	dir := t.TempDir()

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
	dir := t.TempDir()

	fh, err := compat.CreateTempEx(dir, "", compat.O_DELETE)
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

	name, err := tmpname(t)
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

	name, err := tmpname(t)
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
	want := wantMkdirTempPerm
	dir := t.TempDir()
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

	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)

		return
	}

	fh, err := compat.OpenFile(name, compat.O_RDWR|compat.O_CREATE, perm)
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
	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)

		return
	}

	fh, err := compat.OpenFile(name, compat.O_RDWR|compat.O_CREATE|compat.O_DELETE, os.FileMode(0o666))
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

	name, err := tmpname(t)
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

	name, err := tmpname(t)
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

func tmpfile(t *testing.T) (string, error) {
	f, err := compat.CreateTemp(t.TempDir(), "")
	if err != nil {
		return "", err
	}

	name := f.Name()

	err = f.Close()
	if err != nil {
		return "", err
	}

	return name, nil
}

func tmpname(t *testing.T) (string, error) {
	name, err := tmpfile(t)
	if err != nil {
		return "", err
	}

	err = os.Remove(name)
	if err != nil {
		return "", err
	}

	return name, nil
}
