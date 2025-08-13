// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/rasa/compat"
)

const allowedTimeVariance = 1 * time.Second

func TestStatStat(t *testing.T) {
	now := time.Now()

	name, err := createTemp(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
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

	if got := fi.Mode().Perm(); got != compat.CreateTempPerm {
		t.Errorf("Mode(): got 0o%o, want 0o%o", got, compat.CreateTempPerm)
	}

	if got := fi.IsDir(); got != false {
		t.Errorf("IsDir(): got %v, want %v", got, false)
	}

	if got := fi.ModTime(); !timesClose(got, now) {
		t.Errorf("ModTime(): got %v, want %v", got, now)
	}
}

func TestStatLinks(t *testing.T) {
	if !compat.Supports(compat.Links) {
		skip(t, "Skipping test: Links() not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	name, err := createTemp(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
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

	fi, err = compat.Stat(name)
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

	fi, err = compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	want = 1
	if got := fi.Links(); got != want {
		t.Fatalf("Links(): got %v, want %v", got, want)
	}
}

func TestStatATime(t *testing.T) {
	if !compat.Supports(compat.ATime) {
		skip(t, "Skipping test: ATime() not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	now := time.Now()

	name, err := createTemp(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
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

	err = os.Chtimes(name, atime, atime)
	if err != nil {
		t.Fatal(err)
	}

	fi, err = compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	if got := fi.ATime(); !timesClose(got, atime) {
		t.Fatalf("ATime(): got %v, want %v", got, atime)
	}
}

func TestStatBTime(t *testing.T) {
	if !compat.Supports(compat.BTime) {
		skip(t, "Skipping test: BTime() not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	now := time.Now()

	name, err := createTemp(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	if got := fi.BTime(); !timesClose(got, now) {
		t.Fatalf("BTime(): got %v, want %v", got, now)
	}
}

func TestStatCTime(t *testing.T) {
	if !compat.Supports(compat.CTime) {
		skip(t, "Skipping test: CTime() not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	now := time.Now()

	name, err := createTemp(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	if got := fi.CTime(); !timesClose(got, now) {
		t.Fatalf("CTime(): got %v, want %v", got, now)
	}
}

func TestStatMTime(t *testing.T) {
	now := time.Now()

	name, err := createTemp(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
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

	err = os.Chtimes(name, mtime, mtime)
	if err != nil {
		t.Fatal(err)
	}

	fi, err = compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	if got := fi.MTime(); !timesClose(got, mtime) {
		t.Fatalf("MTime(): got %v, want %v", got, mtime)
	}
}

func TestStatUID(t *testing.T) {
	if !compat.Supports(compat.UID) {
		skip(t, "Skipping test: UID() not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	name, err := createTemp(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
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

func TestStatGID(t *testing.T) {
	if !compat.Supports(compat.GID) {
		skip(t, "Skipping test: GID() not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	name, err := createTemp(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
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

func TestStatUser(t *testing.T) {
	if !compat.Supports(compat.UID) {
		skip(t, "Skipping test: User() not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	if compat.IsTinygo {
		// tinygo: Current requires cgo or $USER, $HOME set in environment
		skip(t, "Skipping test: User() not supported on tinygo")

		return // tinygo doesn't support t.Skip
	}

	name, err := createTemp(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
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

func TestStatGroup(t *testing.T) {
	if !compat.Supports(compat.GID) {
		skip(t, "Skipping test: Group() not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	if compat.IsTinygo {
		skip(t, "Skipping test: Group() not supported on tinygo")

		return // tinygo doesn't support t.Skip
	}

	name, err := createTemp(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
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

func TestStatSamePartition(t *testing.T) {
	name, err := createTemp(t)
	if err != nil {
		t.Fatal(err)
	}

	fi1, err := compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	fi2, err := compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	if got := compat.SamePartition(fi1, fi2); !got {
		t.Fatalf("SamePartition(): got %v, want true", got)
	}
}

func TestStatSamePartitions(t *testing.T) {
	name, err := createTemp(t)
	if err != nil {
		t.Fatal(err)
	}

	if got := compat.SamePartitions(name, name); !got {
		t.Fatalf("SamePartitions(): got %v, want true", got)
	}
}

func TestStatSameFile(t *testing.T) {
	name, err := createTemp(t)
	if err != nil {
		t.Fatal(err)
	}

	fi1, err := compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	fi2, err := compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	if got := compat.SameFile(fi1, fi2); !got {
		t.Fatalf("SameFile(): got %v, want true", got)
	}
}

func TestStatSameFiles(t *testing.T) {
	name, err := createTemp(t)
	if err != nil {
		t.Fatal(err)
	}

	if got := compat.SameFiles(name, name); !got {
		t.Fatalf("SameFiles(): got %v, want true", got)
	}
}

func TestStatDiffFile(t *testing.T) {
	name1, err := createTemp(t)
	if err != nil {
		t.Fatal(err)
	}

	name2, err := createTemp(t)
	if err != nil {
		t.Fatal(err)
	}

	fi1, err := compat.Stat(name1)
	if err != nil {
		t.Fatal(err)
	}

	fi2, err := compat.Stat(name2)
	if err != nil {
		t.Fatal(err)
	}

	if got := compat.SameFile(fi1, fi2); got {
		t.Fatalf("SameFile(): got %v, want true", got)
	}
}

func TestStatDiffFiles(t *testing.T) {
	name1, err := createTemp(t)
	if err != nil {
		t.Fatal(err)
	}

	name2, err := createTemp(t)
	if err != nil {
		t.Fatal(err)
	}

	if got := compat.SameFiles(name1, name2); got {
		t.Fatalf("SameFiles(): got %v, want true", got)
	}
}

func createTemp(t *testing.T) (string, error) {
	t.Helper()

	f, err := compat.CreateTemp(t.TempDir(), "*")
	if err != nil {
		return "", err
	}

	name := f.Name()

	// oldUmask := syscall.Umask(0)
	// defer syscall.Umask(oldUmask)

	_, err = f.Write(helloBytes)
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
	return a.Sub(b).Abs() < allowedTimeVariance
}

func compareNames(got string, want string) bool {
	if compat.IsWasip1 {
		if got == "" && want == "daemon" {
			return true
		}
	}

	if !compat.IsWindows {
		return got == want
	}

	if got == "" || want == "" {
		return false
	}
	gotDomain, gotName := parseName(got)
	wantDomain, wantName := parseName(want)
	if gotName == wantName {
		if gotDomain == wantDomain || gotDomain == "" || wantDomain == "" {
			return true
		}
	}

	return false
}

func parseName(name string) (string, string) {
	parts := strings.Split(name, `\`)
	switch {
	case len(parts) == 1:
		return "", strings.ToLower(parts[0])
	default:
		return strings.ToLower(parts[0]), strings.ToLower(parts[1])
	}
}
