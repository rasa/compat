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

func TestLstatStat(t *testing.T) { //nolint:dupl
	if !supportsSymlinks(t) {
		return
	}

	now := time.Now()

	_, link, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Stat(link)
	if err != nil {
		t.Fatal(err)
	}

	_, base := filepath.Split(link)

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
		// if compat.IsWindows {
		//	t.Logf("Mode(): got 0o%o, want 0o%o (ignoring for now)", got, want)
		// } else {
		t.Errorf("Mode(): got 0o%o, want 0o%o", got, want)
		// }
	}

	if got := fi.Mode()&os.ModeSymlink == 0; got != true {
		t.Errorf("Mode()&os.ModeSymlink==0: got %v, want %v", got, true)
	}

	if got := fi.IsDir(); got != false {
		t.Errorf("IsDir(): got %v, want %v", got, false)
	}

	if got := fi.ModTime(); !compareTimes(got, now, testEnv.mtimeGranularity) {
		t.Errorf("ModTime(): got %v, want %v", got, now)
	}
}

func TestLstatLstat(t *testing.T) { //nolint:dupl
	if !supportsSymlinks(t) {
		return
	}

	now := time.Now()

	_, link, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Lstat(link)
	if err != nil {
		t.Fatal(err)
	}

	_, base := filepath.Split(link)

	if got := fi.Name(); got != base {
		t.Errorf("Name(): got %v, want %v", got, base)
	}

	size := int64(len(helloBytes))
	if got := fi.Size(); got == size {
		t.Errorf("Size(): got %v, want !%v", got, size)
	}

	perm := compat.CreateTempPerm
	want := fixPerms(perm, false)
	if got := fi.Mode().Perm(); got == want {
		if testEnv.noACLs {
			t.Logf("Mode(): got 0o%o, want !0o%o", got, want)
		} else {
			t.Errorf("Mode(): got 0o%o, want !0o%o", got, want)
		}
	}

	if got := fi.Mode()&os.ModeSymlink != 0; got != true {
		t.Errorf("Mode()&os.ModeSymlink!=0: got %v, want %v", got, true)
	}

	if got := fi.IsDir(); got != false {
		t.Errorf("IsDir(): got %v, want %v", got, false)
	}

	if got := fi.ModTime(); !compareTimes(got, now, testEnv.mtimeGranularity) {
		t.Errorf("ModTime(): got %v, want %v", got, now)
	}
}

func TestLstatLinks(t *testing.T) {
	if !supportsHardLinks(t) {
		return
	}

	if !supportsSymlinks(t) {
		return
	}

	_, name, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := compat.Lstat(name)
	if err != nil {
		t.Fatal(err)
	}

	var want uint = 1
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
	if compat.IsBSDLike {
		// not sure why a hard link to a symlink doesn't count on BSD
		want = 1
	}

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

func TestLstatATime(t *testing.T) { //nolint:dupl
	if !compat.SupportsATime() {
		skip(t, "Skipping test: ATime() not supported on "+runtime.GOOS)

		return
	}

	if !supportsSymlinks(t) {
		return
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

	if got := fi.ATime(); !compareTimes(got, now, testEnv.atimeGranularity) {
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

	if got := fi.ATime(); !compareTimes(got, now, testEnv.atimeGranularity) {
		t.Fatalf("ATime(): got %v, want %v", got, now)
	}

	fi, err = compat.Lstat(target)
	if err != nil {
		t.Fatal(err)
	}

	if got := fi.ATime(); !compareTimes(got, atime, testEnv.atimeGranularity) {
		t.Fatalf("ATime(): got %v, want %v", got, atime)
	}
}

func TestLstatBTime(t *testing.T) {
	if !compat.SupportsBTime() {
		skip(t, "Skipping test: BTime() not supported on "+runtime.GOOS)

		return
	}

	if !supportsSymlinks(t) {
		return
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

	if got := fi.BTime(); !compareTimes(got, now, testEnv.btimeSymlinkGranularity) {
		t.Fatalf("BTime(): got %v, want %v", got, now)
	}
}

func TestLstatCTime(t *testing.T) {
	if !compat.SupportsCTime() {
		skip(t, "Skipping test: CTime() not supported on "+runtime.GOOS)

		return
	}

	if !supportsSymlinks(t) {
		return
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

	if got := fi.CTime(); !compareTimes(got, now, testEnv.ctimeGranularity) {
		t.Fatalf("CTime(): got %v, want %v", got, now)
	}
}

func TestLstatMTime(t *testing.T) { //nolint:dupl
	if !supportsSymlinks(t) {
		return
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

	if got := fi.MTime(); !compareTimes(got, now, testEnv.mtimeGranularity) {
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

	if got := fi.MTime(); !compareTimes(got, now, testEnv.mtimeGranularity) {
		t.Fatalf("MTime(): got %v, want %v", got, now)
	}

	fi, err = compat.Lstat(target)
	if err != nil {
		t.Fatal(err)
	}

	if got := fi.MTime(); !compareTimes(got, mtime, testEnv.mtimeGranularity) {
		t.Fatalf("MTime(): got %v, want %v", got, mtime)
	}
}

func TestLstatUID(t *testing.T) {
	if !supportsSymlinks(t) {
		return
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

	want := os.Geteuid()
	if got != want {
		t.Fatalf("UID(): got %v, want %v", got, want)
	}
}

func TestLstatGID(t *testing.T) {
	if !supportsSymlinks(t) {
		return
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

	isRoot, _ := compat.IsRoot()

	want := os.Getegid()
	if got != want {
		if compat.IsApple && isRoot {
			t.Logf("GID(): got %v, want %v (ignoring as we are root on %v)", got, want, runtime.GOOS)
		} else {
			t.Fatalf("GID(): got %v, want %v", got, want)
		}
	}
}

func TestLstatUser(t *testing.T) {
	if !supportsSymlinks(t) {
		return
	}

	if compat.IsTinygo {
		// tinygo: Current requires cgo or $USER, $HOME set in environment
		skip(t, "Skipping test: User() not supported on tinygo")

		return
	}

	if compat.IsWindows {
		// tinygo: Current requires cgo or $USER, $HOME set in environment
		skip(t, "Skipping test: User() will be indeterminate on Windows")

		return
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

	if compareNames(got, want) == compat.IsWindows {
		t.Fatalf("User(): got %v, want %v", got, want)
	}
}

func TestLstatUserSetOwner(t *testing.T) {
	if !supportsSymlinks(t) {
		return
	}

	if !compat.IsWindows {
		// tinygo: Current requires cgo or $USER, $HOME set in environment
		skip(t, "Skipping test: Windows only test")

		return
	}

	_, name, err := createTempSymlink(t, compat.WithSetSymlinkOwner(true))
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
	if !supportsSymlinks(t) {
		return
	}

	if compat.IsTinygo {
		skip(t, "Skipping test: Group() not supported on tinygo")

		return
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

	isRoot, _ := compat.IsRoot()

	want := g.Name
	if !compareNames(got, want) {
		if compat.IsApple && isRoot {
			t.Logf("Group(): got %v, want %v (ignoring as we are root on %v)", got, want, runtime.GOOS)
		} else {
			t.Fatalf("Group(): got %v, want %v", got, want)
		}
	}
}

func TestLstatGroupSetOwner(t *testing.T) {
	if !supportsSymlinks(t) {
		return
	}

	if !compat.IsWindows {
		// tinygo: Current requires cgo or $USER, $HOME set in environment
		skip(t, "Skipping test: Windows only test")

		return
	}

	if compat.IsTinygo {
		skip(t, "Skipping test: Group() not supported on tinygo")

		return
	}

	_, name, err := createTempSymlink(t, compat.WithSetSymlinkOwner(true))
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
	if !supportsSymlinks(t) {
		return
	}

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
	if !supportsSymlinks(t) {
		return
	}

	_, name, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	if got := compat.SamePartitions(name, name); !got {
		t.Fatalf("SamePartitions(): got %v, want true", got)
	}
}

func TestLstatSameFile(t *testing.T) {
	if !supportsSymlinks(t) {
		return
	}

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
	if !supportsSymlinks(t) {
		return
	}

	_, name, err := createTempSymlink(t)
	if err != nil {
		t.Fatal(err)
	}

	if got := compat.SameFiles(name, name); !got {
		t.Fatalf("SameFiles(): got %v, want true", got)
	}
}

func TestLstatDiffFile(t *testing.T) {
	if !supportsSymlinks(t) {
		return
	}

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
	if !supportsSymlinks(t) {
		return
	}

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

func createTempSymlink(t *testing.T, opts ...compat.Option) (string, string, error) {
	t.Helper()

	f, err := compat.CreateTemp(tempDir(t), "*")
	if err != nil {
		return "", "", err
	}

	target := f.Name()
	dir, base := filepath.Split(target)
	link := filepath.Join(dir, "link-"+base+".lnk")

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

	err = compat.Symlink(target, link, opts...)
	if err != nil {
		return "", "", err
	}

	return target, link, nil
}
