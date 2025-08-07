// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/rasa/compat"
)

var perms []os.FileMode

func init() {
	// @TODO(rasa): test different umask settings
	compat.Umask(0)

	for u := 6; u <= 7; u++ {
		for g := 0; g <= 7; g++ {
			for o := 0; o <= 7; o++ {
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
		if err != nil {
			t.Fatalf("Chmod(%04o): %v", perm, err)
		}

		checkPerm(t, name, perm)
	}
}

func TestFileWindowsCreate(t *testing.T) {
	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)
	}
	fh, err := compat.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	_ = fh.Close()
	checkPerm(t, name, compat.CreatePerm)
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileWindowsCreateEx(t *testing.T) {
	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)
	}
	perm := want600
	fh, err := compat.CreateEx(name, perm, 0)
	if err != nil {
		t.Fatal(err)
	}
	_ = fh.Close()
	checkPerm(t, name, perm)
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileWindowsCreateTemp(t *testing.T) {
	perm := compat.CreateTempPerm

	dir := t.TempDir()
	fh, err := compat.CreateTemp(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	name := fh.Name()
	_ = fh.Close()
	checkPerm(t, name, perm)
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileWindowsMkdir(t *testing.T) {
	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)
	}
	perm := want700
	err = compat.Mkdir(name, perm)
	if err != nil {
		t.Fatal(err)
	}
	checkPerm(t, name, perm)
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileWindowsMkdirAll(t *testing.T) {
	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)
	}
	perm := want700
	err = compat.MkdirAll(name, perm)
	if err != nil {
		t.Fatal(err)
	}
	checkPerm(t, name, perm)
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileWindowsMkdirTemp(t *testing.T) {
	dir := t.TempDir()
	perm := compat.MkdirTempPerm
	name, err := compat.MkdirTemp(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	checkPerm(t, name, perm)
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileWindowsOpenFile(t *testing.T) {
	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)
	}
	perm := want600
	fh, err := compat.OpenFile(name, os.O_CREATE, perm)
	if err != nil {
		t.Fatal(err)
	}
	_ = fh.Close()
	checkPerm(t, name, perm)
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileWindowsWriteFile(t *testing.T) {
	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)
	}
	perm := want600
	err = compat.WriteFile(name, helloBytes, perm)
	if err != nil {
		t.Fatal(err)
	}
	checkPerm(t, name, perm)
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileWindowsWriteFileEx(t *testing.T) {
	name, err := tmpname(t)
	if err != nil {
		t.Fatal(err)
	}
	perm := want600
	err = compat.WriteFileEx(name, helloBytes, perm, 0)
	if err != nil {
		t.Fatal(err)
	}
	checkPerm(t, name, perm)
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func checkPerm(t *testing.T, name string, perm os.FileMode) {
	t.Helper()

	got, err := compat.ExportStat(name) // acl.GetExplicitFileAccessMode(name)
	if err != nil {
		t.Fatalf("GetExplicitFileAccessMode(%v) returned %v", name, err)
	}

	if got != perm {
		dumpACLs(t, name, true)
		t.Fatalf("got %04o, want %04o", got, perm)
	}
}

func dumpACLs(t *testing.T, name string, doDir bool) {
	t.Helper()

	exe, err := exec.LookPath("icacls.exe")
	if err != nil {
		t.Logf("Command not found: %v", err)
		return
	}

	cmd := exec.Command(exe, name, "/q")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Error running icacls: %v", err)
	}
	s := "\n" + string(out)
	t.Log(s)

	exe, err = exec.LookPath("pwsh.exe")
	if err != nil {
		exe, err = exec.LookPath("powershell.exe")
	}
	if err != nil {
		t.Logf("Command not found: %v", err)
	}

	params := fmt.Sprintf("Get-Acl '%s' | Format-List", name)
	cmd = exec.Command(exe, "-Command", params)
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Logf("Error running pwsh: %v", err)
	}
	s = "\n" + string(out)
	t.Log(s)

	if doDir {
		dir, _ := filepath.Split(name)
		dumpACLs(t, dir, false)
	}
}
