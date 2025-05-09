// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"os"
	"path/filepath"
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
		mode = 0o600
	}
}

func TestStat(t *testing.T) {
	now := time.Now()
	name, err := createTemp(t)
	if err != nil {
		t.Error(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Error(err)
	}

	_, base := filepath.Split(name)

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

func TestLinks(t *testing.T) {
	if !compat.Supports(compat.Links) {
		t.Skip("Links() not supported on " + runtime.GOOS)
	}

	name, err := createTemp(t)
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

	dir, _ := filepath.Split(name)
	link := filepath.Join(dir, "link.txt")
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

func TestATime(t *testing.T) {
	if !compat.Supports(compat.ATime) {
		t.Skip("ATime() not supported on " + runtime.GOOS)
	}

	now := time.Now()
	name, err := createTemp(t)
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

func TestBTime(t *testing.T) {
	if !compat.Supports(compat.BTime) {
		t.Skip("BTime() not supported on " + runtime.GOOS)
	}

	now := time.Now()
	name, err := createTemp(t)
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

func TestCTime(t *testing.T) {
	if !compat.Supports(compat.CTime) {
		t.Skip("CTime() not supported on " + runtime.GOOS)
	}

	now := time.Now()
	name, err := createTemp(t)
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

func TestMTime(t *testing.T) {
	now := time.Now()
	name, err := createTemp(t)
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

func TestUID(t *testing.T) {
	if !compat.Supports(compat.UID) {
		t.Skip("UID() not supported on " + runtime.GOOS)
	}

	name, err := createTemp(t)
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

func TestGID(t *testing.T) {
	if !compat.Supports(compat.GID) {
		t.Skip("GID() not supported on " + runtime.GOOS)
	}

	name, err := createTemp(t)
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

func TestSameDevice(t *testing.T) {
	name, err := createTemp(t)
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

func TestSameDevices(t *testing.T) {
	name, err := createTemp(t)
	if err != nil {
		t.Error(err)
	}

	if got := compat.SameDevices(name, name); !got {
		t.Errorf("SameDevices(): got %v, want true", got)
	}
}

func TestSameFile(t *testing.T) {
	name, err := createTemp(t)
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

func TestSameFiles(t *testing.T) {
	name, err := createTemp(t)
	if err != nil {
		t.Error(err)
	}

	if got := compat.SameFiles(name, name); !got {
		t.Errorf("SameFiles(): got %v, want true", got)
	}
}

func TestDiffFile(t *testing.T) {
	name1, err := createTemp(t)
	if err != nil {
		t.Error(err)
	}

	name2, err := createTemp(t)
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

func TestDiffFiles(t *testing.T) {
	name1, err := createTemp(t)
	if err != nil {
		t.Error(err)
	}

	name2, err := createTemp(t)
	if err != nil {
		t.Error(err)
	}

	if got := compat.SameFiles(name1, name2); got {
		t.Errorf("SameFiles(): got %v, want true", got)
	}
}

func createTemp(t *testing.T) (string, error) {
	t.Helper()

	f, err := os.CreateTemp(t.TempDir(), "*")
	if err != nil {
		return "", err
	}
	name := f.Name()

	// oldUmask := syscall.Umask(0)
	// defer syscall.Umask(oldUmask)

	_, err = f.Write(hello)
	if err != nil {
		return "", err
	}

	err = f.Close()
	if err != nil {
		return "", err
	}

	return name, nil
}

func timesClose(a, b time.Time) bool {
	return a.Sub(b).Abs() < 100*time.Millisecond
}
