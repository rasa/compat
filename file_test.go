// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/rasa/compat"
)

var (
	hello = []byte("hello")
	mode  os.FileMode
)

func init() {
	if runtime.GOOS == "windows" {
		mode = 0o666
	} else {
		mode = 0o644
	}
}

func Test_Stat(t *testing.T) {
	dir := t.TempDir()
	base := "stat.txt"
	now := time.Now()
	name, err := write(t, dir, base)
	if err != nil {
		t.Error(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Error(err)
	}

	if got := fi.Name(); got != base {
		t.Errorf("Name(): got %v, want %v", got, base)
	}

	want := int64(len(hello))
	if got := fi.Size(); got != want {
		t.Errorf("Size(): got %v, want %v", got, want)
	}

	if got := fi.Mode(); got != mode {
		t.Errorf("Mode(): got 0o%o, want 0o%o", got, mode)
	}

	if got := fi.IsDir(); got != false {
		t.Errorf("IsDir(): got %v, want %v", got, false)
	}

	if got := fi.ModTime(); !timesClose(got, now) {
		t.Errorf("ModTime(): got %v, want %v", got, now)
	}
}

func Test_Links(t *testing.T) {
	if !compat.Supports(compat.SupportsLinks) {
		t.Skip("Links() not supported on " + runtime.GOOS)
	}

	dir := t.TempDir()
	base := "links.txt"
	name, err := write(t, dir, base)
	if err != nil {
		t.Error(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Error(err)
	}

	var want uint64 = 1
	if got := fi.Links(); got != want {
		t.Errorf("Links(): got %v, want %v", got, want)
	}

	link := path.Join(dir, "link.txt")
	err = os.Link(name, link)
	if err != nil {
		t.Error(err)
	}

	fi, err = compat.Stat(name)
	if err != nil {
		t.Error(err)
	}

	want = 2
	if got := fi.Links(); got != want {
		t.Errorf("Links(): got %v, want %v", got, want)
	}

	err = os.Remove(link)
	if err != nil {
		t.Error(err)
	}

	fi, err = compat.Stat(name)
	if err != nil {
		t.Error(err)
	}

	want = 1
	if got := fi.Links(); got != want {
		t.Errorf("Links(): got %v, want %v", got, want)
	}
}

func Test_ATime(t *testing.T) {
	if !compat.Supports(compat.SupportsATime) {
		t.Skip("ATime() not supported on " + runtime.GOOS)
	}

	dir := t.TempDir()
	base := "atime.txt"
	now := time.Now()
	name, err := write(t, dir, base)
	if err != nil {
		t.Error(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Error(err)
	}

	if got := fi.ATime(); !timesClose(got, now) {
		t.Errorf("ATime(): got %v, want %v", got, now)
	}

	atime := time.Now().Add(-24 * time.Hour)
	err = os.Chtimes(name, atime, atime)
	if err != nil {
		t.Error(err)
	}

	fi, err = compat.Stat(name)
	if err != nil {
		t.Error(err)
	}

	if got := fi.ATime(); !timesClose(got, atime) {
		t.Errorf("ATime(): got %v, want %v", got, atime)
	}
}

func Test_BTime(t *testing.T) {
	if !compat.Supports(compat.SupportsBTime) {
		t.Skip("BTime() not supported on " + runtime.GOOS)
	}

	dir := t.TempDir()
	base := "btime.txt"
	now := time.Now()
	name, err := write(t, dir, base)
	if err != nil {
		t.Error(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Error(err)
	}

	if got := fi.BTime(); !timesClose(got, now) {
		t.Errorf("BTime(): got %v, want %v", got, now)
	}
}

func Test_CTime(t *testing.T) {
	if !compat.Supports(compat.SupportsCTime) {
		t.Skip("CTime() not supported on " + runtime.GOOS)
	}

	dir := t.TempDir()
	base := "ctime.txt"
	now := time.Now()
	name, err := write(t, dir, base)
	if err != nil {
		t.Error(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Error(err)
	}

	if got := fi.CTime(); !timesClose(got, now) {
		t.Errorf("CTime(): got %v, want %v", got, now)
	}
}

func Test_MTime(t *testing.T) {
	dir := t.TempDir()
	base := "mtime.txt"
	now := time.Now()
	name, err := write(t, dir, base)
	if err != nil {
		t.Error(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Error(err)
	}

	if got := fi.MTime(); !timesClose(got, now) {
		t.Errorf("MTime(): got %v, want %v", got, now)
	}

	mtime := time.Now().Add(-24 * time.Hour)
	err = os.Chtimes(name, mtime, mtime)
	if err != nil {
		t.Error(err)
	}

	fi, err = compat.Stat(name)
	if err != nil {
		t.Error(err)
	}

	if got := fi.MTime(); !timesClose(got, mtime) {
		t.Errorf("MTime(): got %v, want %v", got, mtime)
	}
}

func Test_UID(t *testing.T) {
	if !compat.Supports(compat.SupportsUID) {
		t.Skip("UID() not supported on " + runtime.GOOS)
	}

	dir := t.TempDir()
	base := "uid.txt"
	name, err := write(t, dir, base)
	if err != nil {
		t.Error(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Error(err)
	}

	want := uint64(os.Getuid()) //nolint:gosec // G115: conversion int -> uint64
	if got := fi.UID(); got != want {
		t.Errorf("UID(): got %v, want %v", got, want)
	}
}

func Test_GID(t *testing.T) {
	if !compat.Supports(compat.SupportsGID) {
		t.Skip("GID() not supported on " + runtime.GOOS)
	}

	dir := t.TempDir()
	base := "gid.txt"
	name, err := write(t, dir, base)
	if err != nil {
		t.Error(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Error(err)
	}

	want := uint64(os.Getgid()) //nolint:gosec // G115: conversion int -> uint64
	if got := fi.GID(); got != want {
		t.Errorf("GID(): got %v, want %v", got, want)
	}
}

func Test_SameDevice(t *testing.T) {
	dir := t.TempDir()
	base := "samedevice.txt"
	name, err := write(t, dir, base)
	if err != nil {
		t.Error(err)
	}

	fi1, err := compat.Stat(name)
	if err != nil {
		t.Error(err)
	}

	fi2, err := compat.Stat(name)
	if err != nil {
		t.Error(err)
	}

	if got := compat.SameDevice(fi1, fi2); !got {
		t.Errorf("SameDevice(): got %v, want true", got)
	}
}

func Test_SameDevices(t *testing.T) {
	dir := t.TempDir()
	base := "samedevices.txt"
	name, err := write(t, dir, base)
	if err != nil {
		t.Error(err)
	}

	if got := compat.SameDevices(name, name); !got {
		t.Errorf("SameDevices(): got %v, want true", got)
	}
}

func Test_SameFile(t *testing.T) {
	dir := t.TempDir()
	base := "samefile.txt"
	name, err := write(t, dir, base)
	if err != nil {
		t.Error(err)
	}

	fi1, err := compat.Stat(name)
	if err != nil {
		t.Error(err)
	}

	fi2, err := compat.Stat(name)
	if err != nil {
		t.Error(err)
	}

	if got := compat.SameFile(fi1, fi2); !got {
		t.Errorf("SameFile(): got %v, want true", got)
	}
}

func Test_SameFiles(t *testing.T) {
	dir := t.TempDir()
	base := "samefiles.txt"
	name, err := write(t, dir, base)
	if err != nil {
		t.Error(err)
	}

	if got := compat.SameFiles(name, name); !got {
		t.Errorf("SameFiles(): got %v, want true", got)
	}
}

func Test_DiffFile(t *testing.T) {
	dir := t.TempDir()
	base1 := "difffimes1.txt"
	name1, err := write(t, dir, base1)
	if err != nil {
		t.Error(err)
	}

	base2 := "difffimes2.txt"
	name2, err := write(t, dir, base2)
	if err != nil {
		t.Error(err)
	}

	fi1, err := compat.Stat(name1)
	if err != nil {
		t.Error(err)
	}

	fi2, err := compat.Stat(name2)
	if err != nil {
		t.Error(err)
	}

	if got := compat.SameFile(fi1, fi2); got {
		t.Errorf("SameFile(): got %v, want true", got)
	}
}

func Test_DiffFiles(t *testing.T) {
	dir := t.TempDir()
	base1 := "difffimes1.txt"
	name1, err := write(t, dir, base1)
	if err != nil {
		t.Error(err)
	}

	base2 := "difffimes2.txt"
	name2, err := write(t, dir, base2)
	if err != nil {
		t.Error(err)
	}

	if got := compat.SameFiles(name1, name2); got {
		t.Errorf("SameFiles(): got %v, want true", got)
	}
}

func write(t *testing.T, dir, base string) (string, error) {
	t.Helper()

	name := path.Join(dir, base)

	// oldUmask := syscall.Umask(0)
	// defer syscall.Umask(oldUmask)

	err := os.WriteFile(name, hello, mode)
	if err != nil {
		return "", err
	}

	return name, nil
}

func timesClose(a, b time.Time) bool {
	return a.Sub(b).Abs() < 100*time.Millisecond
}
