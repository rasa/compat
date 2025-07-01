// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"errors"
	"os"
	"testing"

	"github.com/rasa/compat"
)

const (
	o666 = os.FileMode(0o666)
	o600 = os.FileMode(0o600)

	o777 = os.FileMode(0o777)
	o700 = os.FileMode(0o700)
)

var data = []byte("hello")

func init() {
	// @TODO(rasa): test different umask settings
	compat.Umask(0)
}

func TestFilePosixChmod(t *testing.T) {
	want := o666

	name, err := tmpfile(t)
	if err != nil {
		t.Fatal(err)
	}
	err = compat.Chmod(name, want)
	if err != nil {
		t.Fatal(err)
	}
	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)
	}
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFilePosixCreate(t *testing.T) {
	want := compat.CreatePerm

	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)
	}
	fh, err := compat.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	err = fh.Close()
	if err != nil {
		t.Fatal(err)
	}
	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)
	}
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFilePosixCreateEx(t *testing.T) {
	want := compat.CreatePerm

	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)
	}
	fh, err := compat.CreateEx(name, want, 0)
	if err != nil {
		t.Fatal(err)
	}
	err = fh.Close()
	if err != nil {
		t.Fatal(err)
	}
	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)
	}
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFilePosixCreateExDelete(t *testing.T) {
	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)
	}
	fh, err := compat.CreateEx(name, compat.CreatePerm, compat.O_CREATE|compat.O_DELETE)
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

func TestFilePosixCreateTemp(t *testing.T) {
	want := compat.CreateTempPerm
	if compat.IsWindows {
		want = os.FileMode(0o666)
	}

	dir := t.TempDir()
	fh, err := compat.CreateTemp(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	name := fh.Name()
	err = fh.Close()
	if err != nil {
		t.Fatal(err)
	}
	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)
	}
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFilePosixCreateTempEx(t *testing.T) {
	want := compat.CreateTempPerm
	if compat.IsWindows {
		want = os.FileMode(0o666)
	}

	dir := t.TempDir()
	fh, err := compat.CreateTempEx(dir, "", 0)
	if err != nil {
		t.Fatal(err)
	}
	name := fh.Name()
	err = fh.Close()
	if err != nil {
		t.Fatal(err)
	}
	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)
	}
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFilePosixCreateTempExDelete(t *testing.T) {
	dir := t.TempDir()
	fh, err := compat.CreateTempEx(dir, "", compat.O_DELETE)
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

func TestFilePosixMkdir(t *testing.T) {
	want := o777

	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)
	}
	err = compat.Mkdir(name, want)
	if err != nil {
		t.Fatal(err)
	}
	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)
	}
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFilePosixMkdirAll(t *testing.T) {
	want := o777

	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)
	}
	err = compat.MkdirAll(name, want)
	if err != nil {
		t.Fatal(err)
	}
	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)
	}
	err = os.RemoveAll(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFilePosixMkdirTemp(t *testing.T) {
	want := compat.MkdirTempPerm
	if compat.IsWindows {
		want = os.FileMode(0o777)
	}

	dir := t.TempDir()
	pattern := ""
	name, err := compat.MkdirTemp(dir, pattern)
	if err != nil {
		t.Fatal(err)
	}
	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)
	}
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFilePosixOpenFile(t *testing.T) {
	want := o666

	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)
	}
	fh, err := os.OpenFile(name, compat.O_CREATE, want)
	if err != nil {
		t.Fatal(err)
	}
	err = fh.Close()
	if err != nil {
		t.Fatal(err)
	}
	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)
	}
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFilePosixOpenFileDelete(t *testing.T) {
	want := o666

	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)
	}
	fh, err := compat.OpenFile(name, compat.O_CREATE|compat.O_DELETE, want)
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

func TestFilePosixWriteFile(t *testing.T) {
	want := o666

	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)
	}
	err = compat.WriteFile(name, data, want)
	if err != nil {
		t.Fatal(err)
	}
	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)
	}
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFilePosixWriteFileEx(t *testing.T) {
	want := o666

	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)
	}
	err = compat.WriteFileEx(name, data, want, 0)
	if err != nil {
		t.Fatal(err)
	}
	fs, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	got := fs.Mode().Perm()
	if got != want {
		t.Fatalf("got 0%03o, want 0%03o", got, want)
	}
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
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
