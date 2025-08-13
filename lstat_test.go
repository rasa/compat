// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/rasa/compat"
)

func TestLstatStat(t *testing.T) {
	now := time.Now()

	name, _, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Lstat(name)
	if err != nil {
		t.Fatal(err)
	}

	_, base := filepath.Split(name)

	if got := fi.Name(); got != base {
		t.Errorf("Name(): got %v, want %v", got, base)
	}

	want := int64(len(helloBytes))
	if got := fi.Size(); got != want {
		t.Errorf("Size(): got %v, want %v", got, want)
	}

	if got := fi.Mode(); got != compat.CreateTempPerm {
		t.Errorf("Mode(): got 0o%o, want 0o%o", got, compat.CreateTempPerm)
	}

	if got := fi.IsDir(); got != false {
		t.Errorf("IsDir(): got %v, want %v", got, false)
	}

	if got := fi.ModTime(); !timesClose(got, now) {
		t.Errorf("ModTime(): got %v, want %v", got, now)
	}
}

func TestLstatLinks(t *testing.T) {
	if !compat.Supports(compat.Links) {
		skip(t, "Skipping test: Links() not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	_, name, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Lstat(name)
	if err != nil {
		t.Fatal(err)
	}

	var want uint64 = 1
	if got := fi.Links(); got != want {
		t.Fatalf("Links(): got %v, want %v", got, want)
	}

	dir, _ := filepath.Split(name)
	link := filepath.Join(dir, "link.txt")

	err = os.Link(name, link)
	if err != nil {
		t.Fatal(err)
	}

	fi, err = compat.Lstat(name)
	if err != nil {
		t.Fatal(err)
	}

	want = 2
	if got := fi.Links(); got != want {
		t.Fatalf("Links(): got %v, want %v", got, want)
	}

	err = os.Remove(link)
	if err != nil {
		t.Fatal(err)
	}

	fi, err = compat.Lstat(name)
	if err != nil {
		t.Fatal(err)
	}

	want = 1
	if got := fi.Links(); got != want {
		t.Fatalf("Links(): got %v, want %v", got, want)
	}
}

func TestLstatATime(t *testing.T) {
	if !compat.Supports(compat.ATime) {
		skip(t, "Skipping test: ATime() not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	now := time.Now()

	target, link, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Lstat(link)
	if err != nil {
		t.Fatal(err)
	}

	if got := fi.ATime(); !timesClose(got, now) {
		t.Fatalf("ATime(): got %v, want %v", got, now)
	}

	if compat.IsTinygo {
		// os.Chtimes fails with "operation not implemented" on tinygo
		return
	}

	atime := time.Now().Add(-24 * time.Hour)

	err = os.Chtimes(target, atime, atime)
	if err != nil {
		t.Fatal(err)
	}

	fi, err = compat.Lstat(link)
	if err != nil {
		t.Fatal(err)
	}

	if got := fi.ATime(); !timesClose(got, now) {
		t.Fatalf("ATime(): got %v, want %v", got, now)
	}

	fi, err = compat.Lstat(target)
	if err != nil {
		t.Fatal(err)
	}

	if got := fi.ATime(); !timesClose(got, atime) {
		t.Fatalf("ATime(): got %v, want %v", got, atime)
	}
}

func TestLstatBTime(t *testing.T) {
	if !compat.Supports(compat.BTime) {
		skip(t, "Skipping test: BTime() not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	now := time.Now()

	_, name, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Lstat(name)
	if err != nil {
		t.Fatal(err)
	}

	if got := fi.BTime(); !timesClose(got, now) {
		t.Fatalf("BTime(): got %v, want %v", got, now)
	}
}

func TestLstatCTime(t *testing.T) {
	if !compat.Supports(compat.CTime) {
		skip(t, "Skipping test: CTime() not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	now := time.Now()

	_, name, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Lstat(name)
	if err != nil {
		t.Fatal(err)
	}

	if got := fi.CTime(); !timesClose(got, now) {
		t.Fatalf("CTime(): got %v, want %v", got, now)
	}
}

func TestLstatMTime(t *testing.T) {
	now := time.Now()

	target, link, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Lstat(link)
	if err != nil {
		t.Fatal(err)
	}

	if got := fi.MTime(); !timesClose(got, now) {
		t.Fatalf("MTime(): got %v, want %v", got, now)
	}

	if compat.IsTinygo {
		// os.Chtimes fails with "operation not implemented" on tinygo
		return
	}

	mtime := time.Now().Add(-24 * time.Hour)

	err = os.Chtimes(target, mtime, mtime)
	if err != nil {
		t.Fatal(err)
	}

	fi, err = compat.Lstat(link)
	if err != nil {
		t.Fatal(err)
	}

	if got := fi.MTime(); !timesClose(got, now) {
		t.Fatalf("MTime(): got %v, want %v", got, now)
	}

	fi, err = compat.Lstat(target)
	if err != nil {
		t.Fatal(err)
	}

	if got := fi.MTime(); !timesClose(got, mtime) {
		t.Fatalf("MTime(): got %v, want %v", got, mtime)
	}
}

func TestLstatUID(t *testing.T) {
	if !compat.Supports(compat.UID) {
		skip(t, "Skipping test: UID() not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	_, name, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Lstat(name)
	if err != nil {
		t.Fatal(err)
	}

	got := fi.UID()

	if compat.IsWindows {
		if got == compat.UnknownID {
			t.Fatalf("UID(): got %v", got)
		}

		return
	}

	want := os.Getuid()
	if got != want {
		t.Fatalf("UID(): got %v, want %v", got, want)
	}
}

func TestLstatGID(t *testing.T) {
	if !compat.Supports(compat.GID) {
		skip(t, "Skipping test: GID() not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	_, name, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Lstat(name)
	if err != nil {
		t.Fatal(err)
	}

	got := fi.GID()

	if compat.IsWindows {
		if got == compat.UnknownID {
			t.Fatalf("GID(): got %v", got)
		}

		return
	}

	want := os.Getgid()
	if got != want {
		t.Fatalf("GID(): got %v, want %v", got, want)
	}
}

func TestLstatUser(t *testing.T) {
	if !compat.Supports(compat.UID) {
		skip(t, "Skipping test: User() not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	if compat.IsTinygo {
		// tinygo: Current requires cgo or $USER, $HOME set in environment
		skip(t, "Skipping test: User() not supported on tinygo")

		return // tinygo doesn't support t.Skip
	}

	if compat.IsWindows {
		skip(t, "Skipping test: symlinks are not yet supported in Windows")

		return // tinygo doesn't support t.Skip
	}

	_, name, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Lstat(name)
	if err != nil {
		t.Fatal(err)
	}

	got := fi.User()

	u, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	want := u.Username

	if !compareNames(got, want) {
		t.Fatalf("User(): got %v, want %v", got, want)
	}
}

func TestLstatGroup(t *testing.T) {
	if !compat.Supports(compat.GID) {
		skip(t, "Skipping test: Group() not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	if compat.IsTinygo {
		skip(t, "Skipping test: Group() not supported on tinygo")

		return // tinygo doesn't support t.Skip
	}

	if compat.IsWindows {
		// though this test will pass, as the group is computername\None.
		skip(t, "Skipping test: symlinks are not yet supported in Windows")

		return // tinygo doesn't support t.Skip
	}

	_, name, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Lstat(name)
	if err != nil {
		t.Fatal(err)
	}

	got := fi.Group()

	u, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}

	g, err := user.LookupGroupId(u.Gid)
	if err != nil {
		t.Fatal(err)
	}

	want := g.Name
	if !compareNames(got, want) {
		t.Fatalf("Group(): got %v, want %v", got, want)
	}
}

func TestLstatSamePartition(t *testing.T) {
	_, name, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	fi1, err := compat.Lstat(name)
	if err != nil {
		t.Fatal(err)
	}

	fi2, err := compat.Lstat(name)
	if err != nil {
		t.Fatal(err)
	}

	if got := compat.SamePartition(fi1, fi2); !got {
		t.Fatalf("SamePartition(): got %v, want true", got)
	}
}

func TestLstatSamePartitions(t *testing.T) {
	_, name, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	if got := compat.SamePartitions(name, name); !got {
		t.Fatalf("SamePartitions(): got %v, want true", got)
	}
}

func TestLstatSameFile(t *testing.T) {
	_, name, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	fi1, err := compat.Lstat(name)
	if err != nil {
		t.Fatal(err)
	}

	fi2, err := compat.Lstat(name)
	if err != nil {
		t.Fatal(err)
	}

	if got := compat.SameFile(fi1, fi2); !got {
		t.Fatalf("SameFile(): got %v, want true", got)
	}
}

func TestLstatSameFiles(t *testing.T) {
	_, name, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	if got := compat.SameFiles(name, name); !got {
		t.Fatalf("SameFiles(): got %v, want true", got)
	}
}

func TestLstatDiffFile(t *testing.T) {
	_, name1, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	_, name2, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	fi1, err := compat.Lstat(name1)
	if err != nil {
		t.Fatal(err)
	}

	fi2, err := compat.Lstat(name2)
	if err != nil {
		t.Fatal(err)
	}

	if got := compat.SameFile(fi1, fi2); got {
		t.Fatalf("SameFile(): got %v, want true", got)
	}
}

func TestLstatDiffFiles(t *testing.T) {
	_, name1, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	_, name2, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	if got := compat.SameFiles(name1, name2); got {
		t.Fatalf("SameFiles(): got %v, want true", got)
	}
}

func createTempSymlink(t *testing.T) (string, string, error) {
	t.Helper()

	f, err := compat.CreateTemp(t.TempDir(), "*")
	if err != nil {
		return "", "", err
	}

	target := f.Name()
	link := target + ".lnk"

	// oldUmask := syscall.Umask(0)
	// defer syscall.Umask(oldUmask)

	_, err = f.Write(helloBytes)
	if err != nil {
		return "", "", err
	}

	err = f.Close()
	if err != nil {
		return "", "", err
	}

	err = os.Symlink(target, link)
	if err != nil {
		return "", "", err
	}

	return target, link, nil
}
