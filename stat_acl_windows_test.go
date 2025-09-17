// SPDX-FileCopyrightText: Copyright (c) 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
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
