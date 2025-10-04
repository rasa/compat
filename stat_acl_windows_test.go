// SPDX-FileCopyrightText: Copyright (c) 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"context"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"testing"

	"golang.org/x/sys/windows"

	"github.com/rasa/compat"
)

func TestACLWindowsSupportsACLsInvalid(t *testing.T) {
	_, err := compat.SupportsACLs(invalidName)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestACLWindowsSupportsACLsCachedInvalid(t *testing.T) {
	_, err := compat.SupportsACLsCached(nil)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestACLWindowsSupportsACLsHandleInvalid(t *testing.T) {
	_, err := compat.SupportsACLsHandle(windows.InvalidHandle)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestACLWindowsOpenForQueryInvalid(t *testing.T) {
	_, err := compat.OpenForQuery(invalidName)

	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestACLWindowsGetFinalPathNameByHandleGUIDInvalid(t *testing.T) {
	_, err := compat.GetFinalPathNameByHandleGUID(windows.InvalidHandle)

	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestACLWindowsGetVolumePathNamesForVolumeNameInvalid(t *testing.T) {
	_, err := compat.GetVolumePathNamesForVolumeName("")

	if err == nil {
		t.Fatal("got nil, want an error")
	}

	_, err = compat.GetVolumePathNamesForVolumeName(invalidName)

	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestACLWindowsGetVolumeInfoByHandleInvalid(t *testing.T) {
	_, _, err := compat.GetVolumeInfoByHandle(windows.InvalidHandle)

	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestACLWindowsResolveCanonicalRootFromHandleInvalid(t *testing.T) {
	_, _, err := compat.ResolveCanonicalRootFromHandle(windows.InvalidHandle)

	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestACLWindowsResolveCanonicalRootFromHandleUNC(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}

	dir := tempDir(t)
	ctx := context.Background()
	sharename := randomBase36String(8)
	args := []string{"share", sharename + "=" + dir, "/grant:" + usr.Username + ",READ"}
	err = exec.CommandContext(ctx, "net.exe", args...).Run()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		args := []string{"share", sharename, "/del", "/yes"}
		err = exec.CommandContext(ctx, "net.exe", args...).Run()
		if err != nil {
			t.Fatal(err)
		}
	}()

	fh, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	path := `\\?\UNC\127.0.0.1\` + sharename + `\` + filepath.Base(fh.Name())
	_ = fh.Close()

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()

	h := windows.Handle(f.Fd())
	defer func() { _ = windows.CloseHandle(h) }()

	_, _, err = compat.ResolveCanonicalRootFromHandle(h)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

/*
func TestACLWindowsMultiSZToStrings(t *testing.T) {
	if err == nil {
		t.Fatal("got nil, want an error")
	}

}

func TestACLWindowsIsDriveLetterRoot(t *testing.T) {
	if err == nil {
		t.Fatal("got nil, want an error")
	}

}
*/

func TestACLWindowsNormalizeRoot(t *testing.T) {
	got := compat.NormalizeRoot("")
	want := ""
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	got = compat.NormalizeRoot("c:")
	want = `C:\`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	got = compat.NormalizeRoot(`c:\`)
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	got = compat.NormalizeRoot("c:/")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
