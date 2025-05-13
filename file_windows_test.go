// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"fmt"
	"maps"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/rasa/compat"
)

const (
	// https://learn.microsoft.com/en-us/windows-server/identity/ad-ds/manage/understand-security-identifiers#well-known-sids
	ownerSIDString          = "S-1-3-2"
	groupSIDString          = "S-1-3-3"
	worldSIDString          = "S-1-1-0"
	localSystemSIDString    = "S-1-5-18"
	administratorsSIDString = "S-1-5-32-544"
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
	t.Skip("Skipping until we diagnose why this test is failing")
	name, err := tmpfile(t)
	if err != nil {
		t.Fatal(err)
	}

	for _, perm := range perms {
		err = compat.Chmod(name, perm)
		if err != nil {
			t.Fatalf("Chmod(%04o): %v", perm, err)
		}

		checkPerm(t, name, perm, false)
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
	checkPerm(t, name, compat.CreatePerm, false)
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
	perm := o600
	fh, err := compat.CreateEx(name, perm, 0)
	if err != nil {
		t.Fatal(err)
	}
	_ = fh.Close()
	checkPerm(t, name, perm, false)
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileWindowsCreateTemp(t *testing.T) {
	perm := compat.CreateTempPerm

	dir := t.TempDir()
	fh, err := compat.CreateTemp(dir, "", 0)
	if err != nil {
		t.Fatal(err)
	}
	name := fh.Name()
	_ = fh.Close()
	checkPerm(t, name, perm, true)
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
	perm := o700
	err = compat.Mkdir(name, perm)
	if err != nil {
		t.Fatal(err)
	}
	checkPerm(t, name, perm, true)
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
	perm := o700
	err = compat.MkdirAll(name, perm)
	if err != nil {
		t.Fatal(err)
	}
	checkPerm(t, name, perm, true)
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
	checkPerm(t, name, perm, true)
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
	perm := o600
	fh, err := compat.OpenFile(name, os.O_CREATE, perm)
	if err != nil {
		t.Fatal(err)
	}
	_ = fh.Close()
	checkPerm(t, name, perm, false)
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
	perm := o600
	err = compat.WriteFile(name, data, perm)
	if err != nil {
		t.Fatal(err)
	}
	checkPerm(t, name, perm, false)
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
	perm := o600
	err = compat.WriteFileEx(name, data, perm, 0)
	if err != nil {
		t.Fatal(err)
	}
	checkPerm(t, name, perm, false)
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
}

func checkPerm(t *testing.T, name string, perm os.FileMode, isDir bool) { //nolint:funlen,gocyclo // quiet linter
	t.Helper()

	// Get current user's SID
	usr, err := user.Current()
	if err != nil {
		t.Fatalf("cannot get current user: %v", err)
	}
	userSID, _, _, err := syscall.LookupSID("", usr.Username)
	if err != nil {
		t.Fatalf("cannot lookup SID for %s: %v", usr.Username, err)
	}
	userSIDString, err := userSID.String()
	if err != nil {
		t.Fatalf("userSID.String() failed: %v", err)
	}

	token := windows.Token(0)
	err = windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)
	if err != nil {
		t.Fatalf("failed to get process token for %s: %v", usr.Username, err)
	}
	defer token.Close()

	var size uint32
	// First call to get required buffer size
	_ = windows.GetTokenInformation(token, windows.TokenPrimaryGroup, nil, 0, &size)

	buf := make([]byte, size)
	err = windows.GetTokenInformation(token, windows.TokenPrimaryGroup, &buf[0], size, &size)
	if err != nil {
		t.Fatalf("failed to get token information for %s: %v", usr.Username, err)
	}

	// Interpret buffer as TOKEN_PRIMARY_GROUP struct
	type tokenPrimaryGroup struct {
		PrimaryGroup *windows.SID
	}
	group := (*tokenPrimaryGroup)(unsafe.Pointer(&buf[0]))
	pgroupSIDString := group.PrimaryGroup.String()

	// Get Security Descriptor via x/sys/windows high-level wrapper
	sd, err := windows.GetNamedSecurityInfo(
		name,
		windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION,
	)
	if err != nil {
		t.Fatalf("GetNamedSecurityInfo failed: %v", err)
	}

	// Get DACL from descriptor
	dacl, _, err := sd.DACL()
	if err != nil {
		t.Fatalf("DACL() failed: %v", err)
	}
	if dacl == nil {
		t.Fatal("file has no DACL")
	}

	var mode uint32

	// https://learn.microsoft.com/en-us/windows/win32/fileio/file-security-and-access-rights

	// https://github.com/golang/sys/blob/3d9a6b80792a3911da1fa665c959a5ede3abf476/windows/syscall_windows_test.go#L485
	read := uint32(windows.FILE_READ_DATA | windows.FILE_READ_ATTRIBUTES)
	write := uint32(windows.FILE_WRITE_DATA | windows.FILE_APPEND_DATA | windows.FILE_WRITE_ATTRIBUTES | windows.FILE_WRITE_EA)
	execute := uint32(windows.FILE_READ_DATA | windows.FILE_EXECUTE)

	// Walk ACEs
	raw := uintptr(unsafe.Pointer(dacl)) + unsafe.Sizeof(*dacl)
	for i := uint16(0); i < dacl.AceCount; i++ {
		ace := (*windows.ACCESS_ALLOWED_ACE)(unsafe.Pointer(raw))

		if ace.Header.AceType != windows.ACCESS_ALLOWED_ACE_TYPE {
			t.Fatalf("ACE %d is not ACCESS_ALLOWED_ACE", i)
		}

		asid := (*windows.SID)(unsafe.Pointer(&ace.SidStart))
		aceSID := asid.String()

		mask := uint32(ace.Mask)
		var bits uint32
		if mask&read == read {
			bits |= 4
		}
		if mask&write == write {
			bits |= 2
		}
		if mask&execute == execute {
			bits |= 1
		}

		switch aceSID {
		case ownerSIDString:
			fallthrough
		case userSIDString:
			mode |= bits << 6
		case pgroupSIDString:
			fallthrough
		case groupSIDString:
			mode |= bits << 3
		case worldSIDString:
			mode |= bits
		case localSystemSIDString:
			// ignoring
		case administratorsSIDString:
			// ignoring
		default:
			t.Fatalf("unknown SID: %q", aceSID)
		}

		if testing.Verbose() {
			label := "n/a"
			switch aceSID {
			case userSIDString:
				label = "user"
			case ownerSIDString:
				label = "owner"
			case pgroupSIDString:
				label = "pgroup"
			case groupSIDString:
				label = "group"
			case worldSIDString:
				label = "world"
			case localSystemSIDString:
				label = "local"
			case administratorsSIDString:
				label = "admin"
			}

			mask := aceMask(uint32(ace.Mask))
			ftype := "file"
			if isDir {
				ftype = "dir "
			}
			t.Logf("%04o: %v: %-6v: %-20v: SID: %v\n", mode, ftype, label, mask, aceSID)
		}

		raw += uintptr(ace.Header.AceSize)
	}

	got := os.FileMode(mode)
	if got != perm {
		dumpACLs(t, name, true)
		t.Fatalf("got %04o, want %04o", got, perm)
	}
}

func dumpACLs(t *testing.T, name string, doDir bool) {
	t.Helper()

	cmd := exec.Command("icacls.exe", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Error running icacls: %v\n", err)
	}
	t.Log(string(out))

	if doDir {
		dir, _ := filepath.Split(name)
		dumpACLs(t, dir, false)
	}
}

// https://learn.microsoft.com/en-us/windows/win32/fileio/file-access-rights-constants
const FILE_DELETE_CHILD = 0x40

var aceMap2 = map[uint32]string{
	// https://github.com/golang/sys/blob/3d9a6b80792a3911da1fa665c959a5ede3abf476/windows/types_windows.go#L68
	windows.FILE_READ_DATA:        "R",    // 0x001
	windows.FILE_READ_ATTRIBUTES:  "RA",   // 0x080
	windows.FILE_READ_EA:          "REA",  // 0x008
	windows.FILE_WRITE_DATA:       "W",    // 0x002
	windows.FILE_WRITE_ATTRIBUTES: "WA",   // 0x100
	windows.FILE_WRITE_EA:         "WEA",  // 0x010
	windows.FILE_APPEND_DATA:      "AD",   // 0x004
	windows.FILE_EXECUTE:          "X",    // 0x020
	windows.READ_CONTROL:          "RC",   // 0x020000
	windows.WRITE_DAC:             "WDAC", // 0x040000
	windows.WRITE_OWNER:           "WO",   // 0x080000
	windows.SYNCHRONIZE:           "S",    // 0x100000
	// https://github.com/golang/sys/blob/3d9a6b80792a3911da1fa665c959a5ede3abf476/windows/security_windows.go#L992
	windows.DELETE: "D",
	//
	FILE_DELETE_CHILD: "DC", // 0x040
	// windows.FILE_LIST_DIRECTORY: "RD", // 0x001 (dup)
	// windows.FILE_TRAVERSE: "XT", // 0x020 (dup)
}

var aceMap = map[uint32]string{
	windows.STANDARD_RIGHTS_READ | windows.FILE_READ_DATA | windows.FILE_READ_ATTRIBUTES | windows.FILE_READ_EA | windows.SYNCHRONIZE:                                "GR",
	windows.STANDARD_RIGHTS_WRITE | windows.FILE_WRITE_DATA | windows.FILE_WRITE_ATTRIBUTES | windows.FILE_WRITE_EA | windows.FILE_APPEND_DATA | windows.SYNCHRONIZE: "GW",
	windows.STANDARD_RIGHTS_EXECUTE | windows.FILE_READ_ATTRIBUTES | windows.FILE_EXECUTE | windows.SYNCHRONIZE:                                                      "GE",
}

type aceMask uint32

func (a aceMask) String() string {
	mask := uint32(a)
	rv := ""
	rights := map[string]uint32{}
	for k, v := range aceMap {
		if mask&k == k {
			rights[v] = k
		}
	}
	if len(rights) == 0 {
		return "N"
	}
	keys := slices.Collect(maps.Keys(rights))
	slices.Sort(keys)
	rv += strings.Join(keys, ",")

	if mask == 0 {
		return rv
	}

	for _, v := range rights {
		mask &^= v
	}

	rights2 := map[string]uint32{}
	for k, v := range aceMap2 {
		if mask&k == k {
			rights2[v] = k
		}
	}

	if len(rights2) == 0 {
		return rv
	}

	keys = slices.Collect(maps.Keys(rights2))
	slices.Sort(keys)
	rv += "," + strings.Join(keys, ",")

	for _, v := range rights2 {
		mask &^= v
	}

	if mask != 0 {
		rv += "," + fmt.Sprintf("0x%x", mask)
	}
	return rv
}
