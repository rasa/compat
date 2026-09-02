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
	"strconv"
	"strings"
	"testing"

	"golang.org/x/sys/windows"

	"github.com/rasa/compat"
)

const (
	permRW           = perm600
	permRO           = perm400
	userWritableMask = perm200 // aka windows.S_IWRITE
)

var perms []os.FileMode

func init() {
	testing.Init()
	flag.Parse()
	loadPerms()
}

func loadPerms() {
	p := os.Getenv("COMPAT_DEBUG_PERM")
	if p != "" {
		o, err := parseOctal(p)
		if err == nil {
			perms = []os.FileMode{os.FileMode(o)}
		}

		return
	}

	if testing.Short() {
		perms = []os.FileMode{perm555}

		return
	}

	perms = make([]os.FileMode, 0, 8*8*8)

	for u := 7; u >= 0; u-- {
		for g := 7; g >= 0; g-- {
			for o := 7; o >= 0; o-- {
				mode := os.FileMode(u<<0o6 | g<<0o3 | o)
				// @TODO(rasa) support 0o0 perms on Windows
				if mode == perm000 {
					break
				}

				perms = append(perms, mode)
			}
		}
	}
}

func parseOctal(s string) (uint64, error) {
	// Normalize: trim 0o or leading 0 if present
	s = strings.TrimPrefix(strings.ToLower(s), "0o")
	s = strings.TrimPrefix(s, "0")

	// If the whole string was "0" or empty, treat it as "0".
	if s == "" {
		return 0, nil
	}

	v, err := strconv.ParseUint(s, 8, 32)
	if err != nil {
		return 0, err
	}

	return v, nil
}

func TestFileWindowsChmod(t *testing.T) {
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
	}
}

func TestFileWindowsChmodReadOnlyModeIgnoreNotSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		cleanup(t, name)

		err = compat.Chmod(name, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeIgnore))
		if err != nil {
			t.Fatal(err)
		}

		// Check that the RO attribute is not set
		checkChmod(t, name, perm, permRW, false, nil)
	}
}

func TestFileWindowsChmodReadOnlyModeIgnoreSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		cleanup(t, name)

		// Set the file's RO attribute (attrib +R, u-w)
		err = os.Chmod(name, perm400)
		if err != nil {
			t.Fatal(err)
		}

		err = compat.Chmod(name, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeIgnore))
		if err != nil {
			t.Fatalf("perm=%03o (%v): %v", perm, perm, err)
		}

		// Check that the RO attribute is set
		checkChmod(t, name, perm, permRO, false, nil)
	}
}

func TestFileWindowsChmodReadOnlyModeFromPermissionsNotSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		cleanup(t, name)

		err = compat.Chmod(name, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeFromPermissions))

		// Check that the RO attribute is not set
		checkChmod(t, name, perm, perm&permRW, true, err)
	}
}

func TestFileWindowsChmodReadOnlyModeFromPermissionsSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		cleanup(t, name)

		// Set the file's RO attribute (attrib +R, u-w)
		err = os.Chmod(name, perm400)
		if err != nil {
			t.Fatal(err)
		}

		err = compat.Chmod(name, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeFromPermissions))

		// Check that the RO attribute is set to the value in perm
		checkChmod(t, name, perm, perm&permRW, true, err)
	}
}

func TestFileWindowsChmodReadOnlyModeClearNotSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		cleanup(t, name)

		err = compat.Chmod(name, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeClear))

		checkChmod(t, name, perm, perm&permRW, true, err)
	}
}

func TestFileWindowsChmodReadOnlyModeClearSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		cleanup(t, name)

		// Set the file's RO attribute (attrib +R, u-w)
		err = os.Chmod(name, perm400)
		if err != nil {
			t.Fatal(err)
		}

		err = compat.Chmod(name, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeClear))

		// Check that the RO attribute is not set
		checkChmod(t, name, perm, userWritableMask, true, err)
	}
}

func TestFileWindowsCreate(t *testing.T) {
	name := tempName(t)

	cleanup(t, name)

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

func TestFileWindowsCreateReadOnlyModeFromPermissions(t *testing.T) {
	perm := perm400

	name := tempName(t)

	cleanup(t, name)

	fh, err := compat.Create(name, compat.WithFileMode(perm), compat.WithReadOnlyMode(compat.ReadOnlyModeFromPermissions))
	if err != nil {
		t.Fatal(err)
	}

	err = fh.Close()
	if err != nil {
		t.Fatal(err)
	}

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	want := false // the user-writable bit is not set.

	got := fi.Mode().Perm()&userWritableMask == userWritableMask
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestFileWindowsCreateTemp(t *testing.T) {
	dir := tempDir(t)

	fh, err := compat.CreateTemp(dir, "")
	if err != nil {
		t.Fatal(err)
	}

	name := fh.Name()
	cleanup(t, name)

	perm := compat.CreateTempPerm
	checkPerm(t, name, perm, false)

	err = fh.Close()
	if err != nil {
		t.Fatal(err)
	}

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

		cleanup(t, name)

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

func TestFileWindowsFchmodReadOnlyModeIgnoreNotSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		cleanup(t, name)

		f, err := os.Open(name)
		if err != nil {
			t.Fatal(err)
		}
		defer fclose(f)

		err = compat.Fchmod(f, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeIgnore))
		if err != nil {
			t.Fatal(err)
		}

		_ = f.Close()

		// Check that the RO attribute is not set
		checkChmod(t, name, perm, permRW, false, nil)
	}
}

func TestFileWindowsFchmodReadOnlyModeIgnoreSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		cleanup(t, name)

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

		err = compat.Fchmod(f, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeIgnore))
		if err != nil {
			t.Fatalf("perm=%03o (%v): %v", perm, perm, err)
		}

		_ = f.Close()

		// Check that the RO attribute is set
		checkChmod(t, name, perm, permRO, false, nil)
	}
}

func TestFileWindowsFchmodReadOnlyModeFromPermissionsNotSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		cleanup(t, name)

		f, err := os.Open(name)
		if err != nil {
			t.Fatal(err)
		}
		defer fclose(f)

		err = compat.Fchmod(f, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeFromPermissions))
		_ = f.Close()

		checkChmod(t, name, perm, perm&permRW, true, err)
	}
}

func TestFileWindowsFchmodReadOnlyModeFromPermissionsSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		cleanup(t, name)

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

		err = compat.Fchmod(f, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeFromPermissions))

		_ = f.Close()

		checkChmod(t, name, perm, perm&permRW, true, err)
	}
}

func TestFileWindowsFchmodReadOnlyModeClearNotSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		cleanup(t, name)

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

		err = compat.Fchmod(f, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeClear))

		_ = f.Close()

		checkChmod(t, name, perm, perm&permRW, true, err)
	}
}

func TestFileWindowsFchmodReadOnlyModeClearSet(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		cleanup(t, name)

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

		err = compat.Fchmod(f, perm, compat.WithReadOnlyMode(compat.ReadOnlyModeClear))

		_ = f.Close()

		checkChmod(t, name, perm, userWritableMask, true, err)
	}
}

func TestFileWindowsMkdir(t *testing.T) {
	for _, perm := range perms {
		name := tempName(t)
		cleanup(t, name)
		err := compat.Mkdir(name, perm)
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
		name := tempName(t)
		cleanup(t, name)
		err := compat.MkdirAll(name, perm)
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
	cleanup(t, name)

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
		name := tempName(t)
		cleanup(t, name)
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

		cleanup(t, name)

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

		cleanup(t, name)

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

func TestFileWindowsRemoveAllRetry(t *testing.T) {
	for _, perm := range perms {
		name, err := tempFile(t)
		if err != nil {
			t.Fatal(err)
		}

		cleanup(t, name)

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

		err = compat.RemoveAll(name, compat.WithRetrySeconds(2))
		checkDeleted(t, name, perm, err)

		if err != nil {
			t.Fatalf("RemoveAll failed: %v: %v", name, err)
		}
	}
}

func TestFileWindowsWriteFile(t *testing.T) {
	for _, perm := range perms {
		name := tempName(t)
		cleanup(t, name)
		err := compat.WriteFile(name, helloBytes, perm)
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

func TestFileWindowsCurrentUsername(t *testing.T) {
	username := compat.CurrentUsername()
	if username == "" {
		t.Fatal("currentUsername: got '', want a value")
	}
}

var seTakeOwnershipPrivilegeW, _ = windows.UTF16PtrFromString("SeTakeOwnershipPrivilege")

func TestFileWindowsEnablePrivilegeInvalidName(t *testing.T) {
	var tok windows.Token

	err := windows.OpenProcessToken(
		windows.CurrentProcess(),
		windows.TOKEN_ADJUST_PRIVILEGES|windows.TOKEN_QUERY,
		&tok,
	)
	if err != nil {
		t.Fatalf("failed to open process token: %v", err)
	}
	defer tok.Close()

	// Enable SeTakeOwnershipPrivilege (required to take ownership when you don't own it)
	err = compat.EnablePrivilege(tok, nil)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestFileWindowsEnablePrivilegeInvalidToken(t *testing.T) {
	var tok windows.Token
	defer tok.Close()

	// Enable SeTakeOwnershipPrivilege (required to take ownership when you don't own it)
	err := compat.EnablePrivilege(tok, seTakeOwnershipPrivilegeW)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestFileWindowsSaFromPermFalse(t *testing.T) {
	_, err := compat.SaFromPerm(perm000, false)
	if err != nil {
		t.Fatalf("got %q, want nil", err)
	}
}

func TestFileWindowsSetOwnerToCurrentUserInvalid(t *testing.T) {
	err := compat.SetOwnerToCurrentUser(invalidName)
	if err == nil {
		t.Fatal("got nil, want an error")
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

	cmd := exec.Command(exe, args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Error running %v: %v", exe, err)
	}

	s := "\n" + string(out)
	t.Log(s)

	return nil
}

func errno(err error) uint32 {
	if err == nil {
		return 0
	}

	var errno windows.Errno
	if errors.As(err, &errno) {
		return uint32(errno)
	}

	return ^uint32(0)
}

func checkChmod(t *testing.T, name string, perm, want os.FileMode, ignore bool, err error) {
	t.Helper()

	if err != nil {
		if perm&userWritableMask != userWritableMask {
			debugf(t, "perm=%03o (%v): %v (ignoring: we can't set RO bit if u-w)", perm, perm, err)

			return
		}

		t.Fatalf("perm=%03o (%v): %v", perm, perm, err)
	}

	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	got := fi.Mode().Perm()

	// Mask all bits but the user's write bits
	if (want & userWritableMask) != (got & userWritableMask) {
		if ignore && perm&userWritableMask != userWritableMask {
			debugf(t, "perm=%03o (%v): got %03o (%v), want %03o (%v) (ignoring: we can't set RO bit if u-w)", perm, perm, got, got, want, want)

			return
		}
		t.Fatalf("perm=%03o (%v): got %03o (%v), want %03o (%v)", perm, perm, got, got, want, want)
	}
}
