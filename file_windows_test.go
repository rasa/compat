// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat_test

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/rasa/compat"
)

var perms []os.FileMode

const perm555 = os.FileMode(0o555)

func init() {
	// @TODO(rasa): test different umask settings
	compat.Umask(0)

	testing.Init()
	flag.Parse()
	if testing.Short() {
		perms = []os.FileMode{perm555}
		return
	}

	perms = make([]os.FileMode, 0, 8*8*8)

	for u := 7; u >= 0; u-- {
		for g := 7; g >= 0; g-- {
			for o := 7; o >= 0; o-- {
				mode := os.FileMode(u<<0o6 | g<<0o3 | o) //nolint:gosec // quiet linter
				perms = append(perms, mode)
			}
		}
	}
}

func TestFileWindowsChmod(t *testing.T) {
	name, err := tmpfile(t)
	if err != nil {
		t.Fatal(err)
	}

	for _, perm := range perms {
		err = compat.Chmod(name, perm)
		checkPerm(t, name, perm)
		if err != nil {
			t.Fatalf("Chmod(%04o) failed: %v", perm, err)
		}
	}
}

func TestFileWindowsCreate(t *testing.T) {
	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)
	}
	perm := compat.CreatePerm
	fh, err := compat.Create(name)
	checkPerm(t, name, perm)
	if err != nil {
		t.Fatal(err)
	}
	_ = fh.Close()
	err = compat.Remove(name)
	checkDeleted(t, name, err)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileWindowsCreateEx(t *testing.T) {
	for _, perm := range perms {
		t.Logf("perm=%3o (%v)", perm, perm)
		name, err := tmpname(t)
		if err != nil {
			t.Fatal(err)
		}
		fh, err := compat.CreateEx(name, perm, 0)
		checkPerm(t, name, perm)
		if err != nil {
			t.Fatal(err)
		}
		_ = fh.Close()
		err = compat.Remove(name)
		checkDeleted(t, name, err)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsCreateTemp(t *testing.T) {
	dir := t.TempDir()
	fh, err := compat.CreateTemp(dir, "")
	perm := compat.CreateTempPerm
	checkPerm(t, "", perm)
	if err != nil {
		t.Fatal(err)
	}
	name := fh.Name()
	checkPerm(t, name, perm)
	_ = fh.Close()
	err = compat.Remove(name)
	checkDeleted(t, name, err)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileWindowsMkdir(t *testing.T) {
	for _, perm := range perms {
		name, err := tmpname(t)
		if err != nil {
			t.Fatal(err)
		}
		err = compat.Mkdir(name, perm)
		checkPerm(t, name, perm)
		if err != nil {
			t.Fatal(err)
		}
		err = compat.Remove(name)
		checkDeleted(t, name, err)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsMkdirAll(t *testing.T) {
	for _, perm := range perms {
		name, err := tmpname(t)
		if err != nil {
			t.Fatal(err)
		}
		err = compat.MkdirAll(name, perm)
		checkPerm(t, name, perm)
		if err != nil {
			t.Fatal(err)
		}
		err = compat.Remove(name)
		checkDeleted(t, name, err)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsMkdirTemp(t *testing.T) {
	dir := t.TempDir()
	name, err := compat.MkdirTemp(dir, "")
	perm := compat.MkdirTempPerm
	checkPerm(t, name, perm)
	if err != nil {
		t.Fatal(err)
	}
	err = compat.Remove(name)
	checkDeleted(t, name, err)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileWindowsOpenFile(t *testing.T) {
	for _, perm := range perms {
		name, err := tmpname(t)
		if err != nil {
			t.Fatal(err)
		}
		fh, err := compat.OpenFile(name, os.O_CREATE, perm)
		checkPerm(t, name, perm)
		if err != nil {
			t.Fatal(err)
		}
		_ = fh.Close()
		err = compat.Remove(name)
		checkDeleted(t, name, err)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsRemove(t *testing.T) {
	for _, perm := range perms {
		name, err := tmpfile(t)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%v (%03o): %v", perm, perm, name)
		err = compat.Chmod(name, perm)
		checkPerm(t, name, perm)
		if err != nil {
			t.Fatalf("Chmod(%04o) failed: %v", perm, err)
		}

		perm = os.FileMode(0o777) // CreatePerm
		err = compat.Chmod(name, perm)
		checkPerm(t, name, perm)
		if err != nil {
			t.Fatalf("Chmod(%04o) failed: %v", perm, err)
		}
		err = compat.Remove(name)
		checkDeleted(t, name, err)
		if err != nil {
			t.Fatalf("Remove failed: %v: %v", name, err)
		}
	}
}

func TestFileWindowsWriteFile(t *testing.T) {
	for _, perm := range perms {
		name, err := tmpname(t)
		if err != nil {
			t.Fatal(err)
		}
		err = compat.WriteFile(name, helloBytes, perm)
		checkPerm(t, name, perm)
		if err != nil {
			t.Fatal(err)
		}
		err = compat.Remove(name)
		checkDeleted(t, name, err)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileWindowsWriteFileEx(t *testing.T) {
	for _, perm := range perms {
		name, err := tmpname(t)
		if err != nil {
			t.Fatal(err)
		}
		err = compat.WriteFileEx(name, helloBytes, perm, 0)
		checkPerm(t, name, perm)
		if err != nil {
			t.Fatal(err)
		}
		err = compat.Remove(name)
		checkDeleted(t, name, err)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func checkPerm(t *testing.T, name string, perm os.FileMode) {
	t.Helper()

	if name == "" {
		return
	}

	got, err := compat.ExportStat(name) // acl.GetExplicitFileAccessMode(name)
	if err != nil {
		t.Fatalf("Stat(%v) failed: %v (%x)", name, err, errno(err))
	}

	if got != perm {
		logACLs(t, name, false)
		t.Fatalf("got 0o%03o (%v), want 0o%03o (%v)", got, got, perm, perm)
		return
	}
}

func checkDeleted(t *testing.T, name string, err error) {
	t.Helper()

	if name == "" || err == nil {
		return
	}

	_, err = compat.ExportStat(name) // acl.GetExplicitFileAccessMode(name)
	if err != nil {
		t.Logf("Stat(%v) failed: %v (%x)", name, err, errno(err))
	}

	logACLs(t, name, false)
}

func logACLs(t *testing.T, name string, doDir bool) {
	t.Helper()

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

	cmd := exec.Command(exe, args...) //nolint:noctx // quiet linter
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Error running %v: %v", exe, err)
	}
	s := "\n" + string(out)
	t.Log(s)

	return nil
}
