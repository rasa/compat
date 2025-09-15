// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat_test

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/rasa/compat"
)

var perms []os.FileMode

func init() {
	testing.Init()
	flag.Parse()
	loadPerms()
}

func loadPerms() {
	if testing.Short() {
		perms = []os.FileMode{perm555}
		return
	}

	perms = make([]os.FileMode, 0, 8*8*8)

	for u := 7; u >= 0; u-- {
		for g := 7; g >= 0; g-- {
			for o := 7; o >= 0; o-- {
				mode := os.FileMode(u<<0o6 | g<<0o3 | o) //nolint:gosec
				// @TODO(rasa) support 0o0 perms on Windows
				if mode == perm000 {
					break
				}
				perms = append(perms, mode)
			}
		}
	}
}

func TestFileWindowsChmod(t *testing.T) {
	name, err := tempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	for _, perm := range perms {
		err = compat.Chmod(name, perm)
		checkPerm(t, name, perm, false)
		if err != nil {
			t.Fatalf("Chmod(%04o) failed: %v", perm, err)
		}
	}
}

func TestFileWindowsChmodIgnoreNotSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}
		// ReadOnlyModeIgnore: do not set a file's RO attribute, and ignore if it's set.
		err = compat.Chmod(name, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeIgnore))
		if err != nil {
			t.Fatal(err)
		}

		fi, err := os.Stat(name)
		if err != nil {
			t.Fatal(err)
		}

		want := true // user-writable bit should be set.
		got := fi.Mode().Perm()&perm200 == perm200
		if want != got {
			t.Fatalf("got %v, want %v: perm=%03o (%v): %v", got, want, perm, perm, name)
		}

		err = compat.Remove(name)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsChmodIgnoreSet(t *testing.T) { //nolint:dupl
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		// Set the RO attribute.
		err = os.Chmod(name, perm400)
		if err != nil {
			t.Fatal(err)
		}

		// ReadOnlyModeIgnore: do not set a file's RO attribute, and ignore if it's set.
		err = compat.Chmod(name, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeIgnore))
		if err != nil {
			t.Fatalf("perm=%03o (%v): %v", perm, perm, err)
		}

		fi, err := os.Stat(name)
		if err != nil {
			t.Fatal(err)
		}

		want := false // user-writable bit should not be set.
		got := fi.Mode().Perm()&perm200 == perm200
		if want != got {
			t.Fatalf("got %v, want %v: perm=%03o (%v): %v", got, want, perm, perm, name)
		}

		err = compat.Remove(name)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsChmodSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		// ReadOnlyMaskSet: set a file's RO attribute if the file's FileMode has the
		// user writable bit set.
		err = compat.Chmod(name, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeSet))
		if err != nil {
			if perm&perm200 != perm200 {
				debugf(t, "perm=%03o (%v): %v (ignoring)", perm, perm, err)
				continue
			}

			t.Fatalf("perm=%03o (%v): %v", perm, perm, err)
		}

		fi, err := os.Stat(name)
		if err != nil {
			t.Fatal(err)
		}

		want := perm&perm200 == perm200
		got := fi.Mode().Perm()&perm200 == perm200
		if want != got {
			if perm&perm200 != perm200 {
				debugf(t, "got %v, want %v: perm=%03o (%v): %v (ignoring)", got, want, perm, perm, name)
				continue
			}

			t.Fatalf("got %v, want %v: perm=%03o (%v): %v, %03o (%v)", got, want, perm, perm, name, fi.Mode().Perm(), fi.Mode().Perm())
		}

		err = compat.Remove(name)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsChmodResetNotSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		// Reset the RO attribute.
		err = os.Chmod(name, perm600)
		if err != nil {
			t.Fatal(err)
		}

		// ReadOnlyModeReset: do not set a file's RO attribute, and if it's set, reset it.
		err = compat.Chmod(name, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeReset))
		if err != nil {
			if perm&perm200 != perm200 {
				debugf(t, "perm=%03o (%v): %v (ignoring)", perm, perm, err)
				continue
			}

			t.Fatalf("perm=%03o (%v): %v", perm, perm, err)
		}

		fi, err := os.Stat(name)
		if err != nil {
			t.Fatal(err)
		}

		want := false // user-writable bit should not be set.
		got := fi.Mode().Perm()&perm200 == perm200
		if want != got {
			if perm&perm200 != perm200 {
				debugf(t, "got %v, want %v: perm=%03o (%v): %v (ignoring)", got, want, perm, perm, name)
				continue
			}

			t.Fatalf("got %v, want %v: perm=%03o (%v): %v", got, want, perm, perm, name)
		}

		err = compat.Remove(name)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsChmodResetSet(t *testing.T) { //nolint:dupl
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		// Set the RO attribute.
		err = os.Chmod(name, perm400)
		if err != nil {
			t.Fatal(err)
		}

		// ReadOnlyModeReset: do not set a file's RO attribute, and if it's set, reset it.
		err = compat.Chmod(name, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeReset))
		if err != nil {
			t.Fatalf("perm=%03o (%v): %v", perm, perm, err)
		}

		fi, err := os.Stat(name)
		if err != nil {
			t.Fatal(err)
		}

		want := false // user-writable bit should not be set.
		got := fi.Mode().Perm()&perm200 == perm200
		if want != got {
			t.Fatalf("got %v, want %v: perm=%03o (%v): %v", got, want, perm, perm, name)
		}

		err = compat.Remove(name)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsCreate(t *testing.T) {
	name, err := tempName(t)
	if err != nil {
		t.Fatal(err)
	}
	perm := compat.CreatePerm
	fh, err := compat.Create(name)
	checkPerm(t, name, perm, false)
	if err != nil {
		t.Fatal(err)
	}
	_ = fh.Close()
	err = compat.Remove(name)
	checkDeleted(t, name, perm, err)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileWindowsCreateEx(t *testing.T) {
	for _, perm := range perms {
		name, err := tempName(t)
		if err != nil {
			t.Fatalf("perm=%3o (%v): %v", perm, perm, err)
		}
		fh, err := compat.CreateEx(name, perm, 0)
		checkPerm(t, name, perm, false)
		if err != nil {
			t.Fatalf("perm=%3o (%v): %v", perm, perm, err)
		}
		_ = fh.Close()
		err = compat.Remove(name)
		checkDeleted(t, name, perm, err)
		if err != nil {
			t.Fatalf("perm=%3o (%v): %v", perm, perm, err)
		}
	}
}

func TestFileWindowsCreateTemp(t *testing.T) {
	dir := tempDir(t)
	fh, err := compat.CreateTemp(dir, "")
	perm := compat.CreateTempPerm
	checkPerm(t, "", perm, false)
	if err != nil {
		t.Fatal(err)
	}
	name := fh.Name()
	checkPerm(t, name, perm, false)
	_ = fh.Close()
	err = compat.Remove(name)
	checkDeleted(t, name, perm, err)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileWindowsFchmod(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		f, err := os.Open(name)
		if err != nil {
			t.Fatal(err)
		}
		defer fclose(f)

		err = compat.Fchmod(f, perm)
		checkPerm(t, name, perm, false)
		if err != nil {
			t.Fatalf("Chmod(%04o) failed: %v", perm, err)
		}
	}
}

func TestFileWindowsFchmodIgnoreNotSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		f, err := os.Open(name)
		if err != nil {
			t.Fatal(err)
		}
		defer fclose(f)

		// ReadOnlyModeIgnore: do not set a file's RO attribute, and ignore if it's set.
		err = compat.Fchmod(f, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeIgnore))
		if err != nil {
			t.Fatal(err)
		}

		_ = f.Close()

		fi, err := os.Stat(name)
		if err != nil {
			t.Fatal(err)
		}

		want := true // user-writable bit should be set.
		got := fi.Mode().Perm()&perm200 == perm200
		if want != got {
			t.Fatalf("got %v, want %v: perm=%03o (%v): %v", got, want, perm, perm, name)
		}

		err = compat.Remove(name)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsFchmodIgnoreSet(t *testing.T) { //nolint:dupl
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		// Set the RO attribute.
		err = os.Chmod(name, perm400)
		if err != nil {
			t.Fatal(err)
		}

		f, err := os.Open(name)
		if err != nil {
			t.Fatal(err)
		}
		defer fclose(f)

		// ReadOnlyModeIgnore: do not set a file's RO attribute, and ignore if it's set.
		err = compat.Fchmod(f, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeIgnore))
		if err != nil {
			t.Fatalf("perm=%03o (%v): %v", perm, perm, err)
		}

		_ = f.Close()

		fi, err := os.Stat(name)
		if err != nil {
			t.Fatal(err)
		}

		want := false // user-writable bit should not be set.
		got := fi.Mode().Perm()&perm200 == perm200
		if want != got {
			t.Fatalf("got %v, want %v: perm=%03o (%v): %v", got, want, perm, perm, name)
		}

		err = compat.Remove(name)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsFchmodSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		f, err := os.Open(name)
		if err != nil {
			t.Fatal(err)
		}
		defer fclose(f)

		// ReadOnlyMaskSet: set a file's RO attribute if the file's FileMode has the
		// user writable bit set.
		err = compat.Fchmod(f, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeSet))
		if err != nil {
			if perm&perm200 != perm200 {
				debugf(t, "perm=%03o (%v): %v (ignoring)", perm, perm, err)
				continue
			}

			t.Fatalf("perm=%03o (%v): %v", perm, perm, err)
		}

		_ = f.Close()

		fi, err := os.Stat(name)
		if err != nil {
			t.Fatal(err)
		}

		want := perm&perm200 == perm200
		got := fi.Mode().Perm()&perm200 == perm200
		if want != got {
			if perm&perm200 != perm200 {
				debugf(t, "got %v, want %v: perm=%03o (%v): %v (ignoring)", got, want, perm, perm, name)
				continue
			}

			t.Fatalf("got %v, want %v: perm=%03o (%v): %v, %03o (%v)", got, want, perm, perm, name, fi.Mode().Perm(), fi.Mode().Perm())
		}

		err = compat.Remove(name)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsFchmodResetNotSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		f, err := os.Open(name)
		if err != nil {
			t.Fatal(err)
		}
		defer fclose(f)

		// Reset the RO attribute.
		err = os.Chmod(name, perm600)
		if err != nil {
			t.Fatal(err)
		}

		// ReadOnlyModeReset: do not set a file's RO attribute, and if it's set, reset it.
		err = compat.Fchmod(f, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeReset))
		if err != nil {
			if perm&perm200 != perm200 {
				debugf(t, "perm=%03o (%v): %v (ignoring)", perm, perm, err)
				continue
			}

			t.Fatalf("perm=%03o (%v): %v", perm, perm, err)
		}

		_ = f.Close()

		fi, err := os.Stat(name)
		if err != nil {
			t.Fatal(err)
		}

		want := false // user-writable bit should not be set.
		got := fi.Mode().Perm()&perm200 == perm200
		if want != got {
			if perm&perm200 != perm200 {
				debugf(t, "got %v, want %v: perm=%03o (%v): %v (ignoring)", got, want, perm, perm, name)
				continue
			}

			t.Fatalf("got %v, want %v: perm=%03o (%v): %v", got, want, perm, perm, name)
		}

		err = compat.Remove(name)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsFchmodResetSet(t *testing.T) { //nolint:dupl
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		f, err := os.Open(name)
		if err != nil {
			t.Fatal(err)
		}
		defer fclose(f)

		// Set the RO attribute.
		err = os.Chmod(name, perm400)
		if err != nil {
			t.Fatal(err)
		}

		// ReadOnlyModeReset: do not set a file's RO attribute, and if it's set, reset it.
		err = compat.Fchmod(f, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeReset))
		if err != nil {
			t.Fatalf("perm=%03o (%v): %v", perm, perm, err)
		}

		_ = f.Close()

		fi, err := os.Stat(name)
		if err != nil {
			t.Fatal(err)
		}

		want := false // user-writable bit should not be set.
		got := fi.Mode().Perm()&perm200 == perm200
		if want != got {
			t.Fatalf("got %v, want %v: perm=%03o (%v): %v", got, want, perm, perm, name)
		}

		err = compat.Remove(name)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsMkdir(t *testing.T) {
	for _, perm := range perms {
		name, err := tempName(t)
		if err != nil {
			t.Fatal(err)
		}
		err = compat.Mkdir(name, perm)
		checkPerm(t, name, perm, true)
		if err != nil {
			t.Fatal(err)
		}
		err = compat.Remove(name)
		checkDeleted(t, name, perm, err)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsMkdirAll(t *testing.T) {
	for _, perm := range perms {
		name, err := tempName(t)
		if err != nil {
			t.Fatal(err)
		}
		err = compat.MkdirAll(name, perm)
		checkPerm(t, name, perm, true)
		if err != nil {
			t.Fatal(err)
		}
		err = compat.Remove(name)
		checkDeleted(t, name, perm, err)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsMkdirTemp(t *testing.T) {
	dir := tempDir(t)
	name, err := compat.MkdirTemp(dir, "")
	perm := compat.MkdirTempPerm
	checkPerm(t, name, perm, true)
	if err != nil {
		t.Fatal(err)
	}
	err = compat.Remove(name)
	checkDeleted(t, name, perm, err)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileWindowsOpenFile(t *testing.T) {
	for _, perm := range perms {
		name, err := tempName(t)
		if err != nil {
			t.Fatal(err)
		}
		fh, err := compat.OpenFile(name, os.O_CREATE, perm)
		checkPerm(t, name, perm, false)
		if err != nil {
			t.Fatal(err)
		}
		_ = fh.Close()
		err = compat.Remove(name)
		checkDeleted(t, name, perm, err)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsRemove(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		err = compat.Chmod(name, perm)
		checkPerm(t, name, perm, false)
		if err != nil {
			t.Fatalf("Chmod(%04o) failed: %v", perm, err)
		}

		perm = perm777
		err = compat.Chmod(name, perm)
		checkPerm(t, name, perm, false)
		if err != nil {
			t.Fatalf("Chmod(%04o) failed: %v", perm, err)
		}
		err = compat.Remove(name)
		checkDeleted(t, name, perm, err)
		if err != nil {
			t.Fatalf("Remove failed: %v: %v", name, err)
		}
	}
}

func TestFileWindowsRemoveAll(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		err = compat.Chmod(name, perm)
		checkPerm(t, name, perm, false)
		if err != nil {
			t.Fatalf("Chmod(%04o) failed: %v", perm, err)
		}

		perm = perm777
		err = compat.Chmod(name, perm)
		checkPerm(t, name, perm, false)
		if err != nil {
			t.Fatalf("Chmod(%04o) failed: %v", perm, err)
		}
		err = compat.RemoveAll(name)
		checkDeleted(t, name, perm, err)
		if err != nil {
			t.Fatalf("RemoveAll failed: %v: %v", name, err)
		}
	}
}

func TestFileWindowsWithReadOnlyModeIgnore(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}
		// ReadOnlyModeIgnore do not set a file's RO attribute, and ignore if it's set.
		err = compat.Chmod(name, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeIgnore))
		if err != nil {
			t.Fatal(err)
		}

		fi, err := os.Stat(name)
		if err != nil {
			t.Fatal(err)
		}

		want := true // user-writable bit should be set.
		got := fi.Mode().Perm()&perm200 == perm200
		if want != got {
			t.Fatalf("WithReadOnlyMode(ReadOnlyModeIgnore): got %v, want %v: perm=%03o (%v): %v", got, want, perm, perm, name)
		}

		err = compat.Remove(name)
		if err != nil {
			t.Fatal(err)
		}
	}

	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		// Set the RO attribute.
		err = os.Chmod(name, perm400)
		if err != nil {
			t.Fatal(err)
		}

		// ReadOnlyModeIgnore do not set a file's RO attribute, and ignore if it's set.
		err = compat.Chmod(name, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeIgnore))
		if err != nil {
			t.Fatal(err)
		}

		fi, err := os.Stat(name)
		if err != nil {
			t.Fatal(err)
		}

		want := false // user-writable bit should not be set.
		got := fi.Mode().Perm()&perm200 == perm200
		if want != got {
			t.Fatalf("WithReadOnlyMode(ReadOnlyModeIgnore): got %v, want %v: perm=%03o (%v): %v", got, want, perm, perm, name)
		}

		err = compat.Remove(name)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsWithReadOnlyModeSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		// ReadOnlyMaskSet set a file's RO attribute if the file's FileMode has the
		// user writable bit set.
		err = compat.Chmod(name, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeSet))
		if err != nil {
			if perm&perm200 != perm200 {
				debugf(t, "perm=%03o (%v): %v (ignoring)", perm, perm, err)
				continue
			}

			t.Fatalf("perm=%3o (%v): %v", perm, perm, err)
		}

		fi, err := os.Stat(name)
		if err != nil {
			t.Fatalf("perm=%3o (%v): %v", perm, perm, err)
		}

		want := perm&perm200 == perm200
		got := fi.Mode().Perm()&perm200 == perm200
		if want != got {
			t.Fatalf("WithReadOnlyMode(ReadOnlyModeSet): got %v, want %v: perm=%03o (%v): %v", got, want, perm, perm, name)
		}

		err = compat.Remove(name)
		if err != nil {
			t.Fatalf("perm=%3o (%v): %v", perm, perm, err)
		}
	}
}

func TestFileWindowsWithReadOnlyModeResetSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		// Set the RO attribute.
		err = os.Chmod(name, perm400)
		if err != nil {
			t.Fatal(err)
		}

		// ReadOnlyModeReset do not set a file's RO attribute, and if it's set, reset it.
		err = compat.Chmod(name, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeReset))
		if err != nil {
			if perm&perm200 != perm200 {
				debugf(t, "perm=%03o (%v): %v (ignoring)", perm, perm, err)
				continue
			}

			t.Fatal(err)
		}

		fi, err := os.Stat(name)
		if err != nil {
			t.Fatal(err)
		}

		want := false // user-writable bit should not be set.
		got := fi.Mode().Perm()&perm200 == perm200
		if want != got {
			t.Fatalf("WithReadOnlyMode(ReadOnlyModeReset): got %v, want %v: perm=%03o (%v): %v", got, want, perm, perm, name)
		}

		err = compat.Remove(name)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsWithReadOnlyModeResetReset(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatalf("perm=%3o (%v): %v", perm, perm, err)
		}

		// Reset the RO attribute.
		err = os.Chmod(name, perm600)
		if err != nil {
			t.Fatalf("perm=%3o (%v): %v", perm, perm, err)
		}

		// ReadOnlyModeReset do not set a file's RO attribute, and if it's set, reset it.
		err = compat.Chmod(name, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeReset))
		if err != nil {
			if perm&perm200 != perm200 {
				debugf(t, "perm=%03o (%v): %v (ignoring)", perm, perm, err)
				continue
			}

			t.Fatalf("perm=%3o (%v): %v", perm, perm, err)
		}

		fi, err := os.Stat(name)
		if err != nil {
			t.Fatalf("perm=%3o (%v): %v", perm, perm, err)
		}

		want := false // user-writable bit should not be set.
		got := fi.Mode().Perm()&perm200 == perm200
		if want != got {
			t.Fatalf("WithReadOnlyMode(ReadOnlyModeReset): got %v, want %v: perm=%03o (%v) fi.Mode().Perm()=%03o (%v): %v", got, want, perm, perm, fi.Mode().Perm(), fi.Mode().Perm(), name)
		}

		err = compat.Remove(name)
		if err != nil {
			t.Fatalf("perm=%3o (%v): %v", perm, perm, err)
		}
	}
}

func TestFileWindowsWriteFile(t *testing.T) {
	for _, perm := range perms {
		name, err := tempName(t)
		if err != nil {
			t.Fatal(err)
		}
		err = compat.WriteFile(name, helloBytes, perm)
		checkPerm(t, name, perm, false)
		if err != nil {
			t.Fatal(err)
		}
		err = compat.Remove(name)
		checkDeleted(t, name, perm, err)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsWriteFileEx(t *testing.T) {
	for _, perm := range perms {
		name, err := tempName(t)
		if err != nil {
			t.Fatal(err)
		}
		err = compat.WriteFileEx(name, helloBytes, perm, 0)
		checkPerm(t, name, perm, false)
		if err != nil {
			t.Fatal(err)
		}
		err = compat.Remove(name)
		checkDeleted(t, name, perm, err)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsCreateReadOnlyModeSet(t *testing.T) {
	perm := perm400

	name, err := tempName(t)
	if err != nil {
		t.Fatal(err)

		return
	}

	defer removeIt(name)

	fh, err := compat.Create(name, compat.WithFileMode(perm), compat.WithReadOnlyMode(compat.ReadOnlyModeSet))
	if err != nil {
		t.Fatal(err)

		return
	}

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

	want := false // the user-writable bit is not set.
	got := fi.Mode().Perm()&perm200 == perm200
	if got != want {
		t.Fatalf("got %v, want %v", got, want)

		return
	}
}

func checkPerm(t *testing.T, name string, perm os.FileMode, isDir bool) {
	t.Helper()

	if name == "" {
		return
	}

	fi, err := compat.Stat(name)
	if err != nil {
		t.Fatalf("Stat() failed: perm=%03o (%v): error %x: %v", perm, perm, errno(err), err)
	}

	got := fi.Mode().Perm()
	want := fixPerms(perm, isDir)
	if got != want {
		logACLs(t, name, false)
		t.Fatalf("got 0o%03o (%v), want 0o%03o (%v): %v", got, got, want, want, name)
		return
	}
}

func checkDeleted(t *testing.T, name string, perm os.FileMode, err error) {
	t.Helper()

	if name == "" || err == nil {
		return
	}

	_, err = compat.Stat(name)
	if err != nil {
		t.Fatalf("Stat() failed: perm=%03o (%v): error %x: %v", perm, perm, errno(err), err)
	}

	logACLs(t, name, false)
}

func logACLs(t *testing.T, name string, doDir bool) {
	t.Helper()

	if !strings.Contains(compatDebug, "ACLS") {
		return
	}

	args := []string{name}
	_ = logOutput(t, "attrib.exe", args)

	args = []string{name, "/q"}
	_ = logOutput(t, "icacls.exe", args)

	command := fmt.Sprintf("Get-Acl '%s' | Format-List", name)
	args = []string{"-Command", command}
	exe, err := exec.LookPath("pwsh.exe")
	if err != nil {
		exe, err = exec.LookPath("powershell.exe")
	}
	if err == nil {
		_ = logOutput(t, exe, args)
	}

	if doDir {
		dir, _ := filepath.Split(name)
		logACLs(t, dir, false)
	}
}

func logOutput(t *testing.T, exe string, args []string) error {
	t.Helper()

	exe, err := exec.LookPath(exe)
	if err != nil {
		if testing.Verbose() {
			t.Logf("Command not found: %v", err)
		}

		return err
	}

	cmd := exec.Command(exe, args...) //nolint:noctx
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Error running %v: %v", exe, err)
	}
	s := "\n" + string(out)
	t.Log(s)

	return nil
}

func errno(err error) uint32 { //nolint:unused
	if err == nil {
		return 0
	}
	var errno syscall.Errno
	if errors.As(err, &errno) {
		return uint32(errno)
	}

	return ^uint32(0)
}

func fclose(f *os.File) {
	if f != nil {
		_ = f.Close()
	}
}
