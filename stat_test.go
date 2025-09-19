// SPDX-FileCopyrightText: Copyright (c) 2025 Ross Smith II <ross@smithii.com>
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

func TestStatStat(t *testing.T) {
	now := time.Now()

	name, err := createTempFile(t)
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

	size := int64(len(helloBytes))
	if got := fi.Size(); got != size {
		t.Errorf("Size(): got %v, want %v", got, size)
	}

	perm := compat.CreateTempPerm
	want := fixPerms(perm, false)
	if got := fi.Mode().Perm(); got != want {
		t.Errorf("Mode(): got 0o%o, want 0o%o", got, want)
	}

	if got := fi.Mode().Type(); got != 0 {
		t.Errorf("fi.Mode().Type(): got 0o%o, want 0o%o", got, 0)
	}

	if got := fi.IsDir(); got != false {
		t.Errorf("IsDir(): got %v, want %v", got, false)
	}

	if got := fi.ModTime(); !compareTimes(got, now, testEnv.mtimeGranularity) {
		fatalTimes(t, "ModTime()", got, now, testEnv.mtimeGranularity)
	}

	if got := fi.Sys(); got == nil {
		t.Error("Sys(): got nil, want not-nil")
	}
}

func TestStatLinks(t *testing.T) {
	if !supportsHardLinks(t) {
		return
	}

	name, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	var want uint = 1
	if got := fi.Links(); got != want {
		t.Fatalf("Links(): got %v, want %v", got, want)
	}

	dir, _ := filepath.Split(name)
	link := filepath.Join(dir, "link.txt")

	err = osLink(name, link)
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

func TestStatATime(t *testing.T) { //nolint:dupl
	if !compat.SupportsATime() {
		skip(t, "Skipping test: ATime() not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	now := time.Now()

	name, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	if got := fi.ATime(); !compareTimes(got, now, testEnv.atimeGranularity) {
		fatalTimes(t, "ATime()", got, now, testEnv.atimeGranularity)
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

	if got := fi.ATime(); !compareTimes(got, atime, testEnv.atimeGranularity) {
		fatalTimes(t, "ATime()", got, atime, testEnv.atimeGranularity)
	}
}

func TestStatBTime(t *testing.T) {
	if !compat.SupportsBTime() {
		skip(t, "Skipping test: BTime() not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	now := time.Now()

	name, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	if got := fi.BTime(); !compareTimes(got, now, testEnv.btimeGranularity) {
		fatalTimes(t, "BTime()", got, now, testEnv.atimeGranularity)
	}
}

func TestStatCTime(t *testing.T) {
	if !compat.SupportsCTime() {
		skip(t, "Skipping test: CTime() not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	now := time.Now()

	name, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	if got := fi.CTime(); !compareTimes(got, now, testEnv.ctimeGranularity) {
		fatalTimes(t, "CTime()", got, now, testEnv.ctimeGranularity)
	}
}

func TestStatMTime(t *testing.T) { //nolint:dupl
	now := time.Now()

	name, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	if got := fi.MTime(); !compareTimes(got, now, testEnv.mtimeGranularity) {
		fatalTimes(t, "MTime()", got, now, testEnv.mtimeGranularity)
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

	if got := fi.MTime(); !compareTimes(got, mtime, testEnv.mtimeGranularity) {
		fatalTimes(t, "MTime()", got, mtime, testEnv.mtimeGranularity)
	}
}

func TestStatUID(t *testing.T) {
	name, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	got := fi.UID()

	if compat.IsWindows {
		if !testEnv.noACLs && got == compat.UnknownID {
			t.Fatalf("UID(): got %v", got)
		}

		return
	}

	want := os.Geteuid()
	if got != want {
		partType := partitionType(name)
		if compat.IsApple && (partType == "exfat" || partType == "msdos") {
			t.Logf("UID(): got %v, want %v (ignoring: %v on %v)", got, want, partType, runtime.GOOS)

			return
		}

		t.Fatalf("UID(): got %v, want %v", got, want)
	}
}

func TestStatGID(t *testing.T) {
	name, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	got := fi.GID()

	if compat.IsWindows {
		if !testEnv.noACLs && got == compat.UnknownID {
			t.Fatalf("GID(): got %v", got)
		}

		return
	}

	isRoot, _ := compat.IsRoot()

	want := os.Getegid()
	if got != want {
		if compat.IsApple && isRoot {
			t.Logf("GID(): got %v, want %v (ignoring: we are root on %v)", got, want, runtime.GOOS)

			return
		}

		t.Fatalf("GID(): got %v, want %v", got, want)
	}
}

func TestStatUser(t *testing.T) {
	if compat.IsTinygo {
		// tinygo: Current requires cgo or $USER, $HOME set in environment
		skip(t, "Skipping test: User() not supported on tinygo")

		return // tinygo doesn't support t.Skip
	}

	name, err := createTempFile(t)
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
		partType := partitionType(name)
		if compat.IsApple && (partType == "exfat" || partType == "msdos") {
			t.Logf("User(): got %v, want %v (ignoring: %v on %v)", got, want, partType, runtime.GOOS)

			return
		}

		t.Fatalf("User(): got %v, want %v", got, want)
	}
}

func TestStatGroup(t *testing.T) {
	if compat.IsTinygo {
		skip(t, "Skipping test: Group() not supported on tinygo")

		return // tinygo doesn't support t.Skip
	}

	name, err := createTempFile(t)
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

	isRoot, _ := compat.IsRoot()

	want := g.Name
	if !compareNames(got, want) {
		if compat.IsApple && isRoot {
			t.Logf("Group(): got %v, want %v (ignoring: we are root on %v)", got, want, runtime.GOOS)

			return
		}

		t.Fatalf("Group(): got %v, want %v", got, want)
	}
}

func TestStatError(t *testing.T) { //nolint:dupl
	name, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	err = fi.Error()
	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

func TestStatFileID(t *testing.T) { //nolint:dupl
	name, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	got := fi.FileID()
	if got == 0 {
		t.Fatal("got 0, want !0")
	}
}

func TestStatPartitionID(t *testing.T) { //nolint:dupl
	name, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	got := fi.PartitionID()
	if got == 0 {
		t.Fatal("got 0, want !0")
	}
}

var stringPrefixes = []string{
	"Name:",
	"Size:",
	"Mode:",
	"ModTime:",
	"ATime:",
	"BTime:",
	"CTime:",
	"IsDir:",
	"Links:",
	"UID:",
	"GID:",
	"PartID:",
	"FileID:",
}

func TestStatString(t *testing.T) { //nolint:dupl
	name, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	got := fi.String()
	for _, prefix := range stringPrefixes {
		if !strings.Contains(got, prefix) {
			t.Fatalf("got %q, want to contained in %q", got, prefix)
		}
	}
}

func TestStatInfo(t *testing.T) { //nolint:dupl
	name, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	got, err := fi.Info()
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("got nil, want a value")
	}
}

func TestStatSamePartition(t *testing.T) {
	name, err := createTempFile(t)
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
	name, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	if got := compat.SamePartitions(name, name); !got {
		t.Fatalf("SamePartitions(): got %v, want true", got)
	}
}

func TestStatSameFile(t *testing.T) {
	name, err := createTempFile(t)
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
	name, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	if got := compat.SameFiles(name, name); !got {
		t.Fatalf("SameFiles(): got %v, want true", got)
	}
}

func TestStatDiffFile(t *testing.T) {
	name1, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	name2, err := createTempFile(t)
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
	name1, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	name2, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	if got := compat.SameFiles(name1, name2); got {
		t.Fatalf("SameFiles(): got %v, want true", got)
	}
}

const (
	userIDSourceMin = compat.UserIDSourceIsInt
	userIDSourceMax = compat.UserIDSourceIsNone
)

func TestStatUserIDSource(t *testing.T) { //nolint:dupl
	src := compat.UserIDSource()
	if src < userIDSourceMin || src > userIDSourceMax {
		t.Fatalf("got %v, want between %v and %v", src, userIDSourceMin, userIDSourceMax)
	}
}

func TestStatStatInvalid(t *testing.T) {
	_, err := compat.Stat(invalidName)
	if err == nil {
		t.Fatalf("got %q, want nil", err)
	}
}

func TestStatSamePartitionInvalid1(t *testing.T) {
	name1 := invalidName

	name2, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	fi1, _ := compat.Stat(name1)
	fi2, _ := compat.Stat(name2)

	got := compat.SamePartition(fi1, fi2)
	if got {
		t.Fatalf("got %v, want false", got)
	}
}

func TestStatSamePartitionInvalid2(t *testing.T) {
	name1, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	name2 := invalidName

	fi1, _ := compat.Stat(name1)
	fi2, _ := compat.Stat(name2)

	got := compat.SamePartition(fi1, fi2)
	if got {
		t.Fatalf("got %v, want false", got)
	}
}

func TestStatSamePartitionsInvalid1(t *testing.T) {
	name1 := invalidName

	name2, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	got := compat.SamePartitions(name1, name2)
	if got {
		t.Fatalf("got %v, want false", got)
	}
}

func TestStatSamePartitionsInvalid2(t *testing.T) {
	name1, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	name2 := invalidName

	got := compat.SamePartitions(name1, name2)
	if got {
		t.Fatalf("got %v, want false", got)
	}
}

func TestStatSameFileInvalid1(t *testing.T) {
	name1 := invalidName

	name2, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	fi1, _ := compat.Stat(name1)

	fi2, _ := compat.Stat(name2)

	got := compat.SameFile(fi1, fi2)
	if got {
		t.Fatalf("got %v, want false", got)
	}
}

func TestStatSameFileInvalid2(t *testing.T) {
	name1, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	name2 := invalidName

	fi1, _ := compat.Stat(name1)

	fi2, _ := compat.Stat(name2)

	got := compat.SameFile(fi1, fi2)
	if got {
		t.Fatalf("got %v, want false", got)
	}
}

func TestStatSameFilesInvalid1(t *testing.T) {
	name1 := invalidName

	name2, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	got := compat.SameFiles(name1, name2)
	if got {
		t.Fatalf("got %v, want false", got)
	}
}

func TestStatSameFilesInvalid2(t *testing.T) {
	name1, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	name2 := invalidName

	got := compat.SameFiles(name1, name2)
	if got {
		t.Fatalf("got %v, want false", got)
	}
}
func TestStatExportedStatInvalidFileInfo(t *testing.T) {
	name, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	_, err := compat.ExportedStat(nil, name, false)
	if err == nil {
		t.Fatalf("got %q, want nil", err)
	}
}

func TestStatExportedStatInvalidName(t *testing.T) {
	name, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	_, err := compat.ExportedStat(fi, invalidName, false)
	if err == nil {
		t.Fatalf("got %q, want nil", err)
	}
}

func createTempFile(t *testing.T) (string, error) {
	t.Helper()

	f, err := compat.CreateTemp(tempDir(t), "*")
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
