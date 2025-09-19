// SPDX-FileCopyrightText: Copyright (c) 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/rasa/compat"
)

func TestStatPosixWindowsGetFileOwnerAndGroupSIDs(t *testing.T) {
	name, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	ownerSID, groupSID, err := compat.GetFileOwnerAndGroupSIDs(name)
	if err != nil {
		t.Fatalf("got %q, want nil", err)
	}
	if !compat.IsValidSid(ownerSID) {
		t.Fatalf("got an invalid owner SID: %v", ownerSID.String())
	}
	if !compat.IsValidSid(groupSID) {
		t.Fatalf("got an invalid group SID: %v", groupSID.String())
	}
}

func TestStatPosixWindowsGetFileOwnerAndGroupSIDsInvalid(t *testing.T) {
	_, _, err := compat.GetFileOwnerAndGroupSIDs(invalidName)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestStatPosixWindowsLSAOpenPolicy(t *testing.T) {
	_, err := compat.LSAOpenPolicy(nil, compat.POLICY_VIEW_LOCAL_INFORMATION)
	if err != nil {
		t.Fatalf("got %q, want nil", err)
	}
}

func TestStatPosixWindowsLSAOpenPolicyInvalid(t *testing.T) {
	access := ^uint32(0) // all bits set
	_, err := compat.LSAOpenPolicy(nil, access)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestStatPosixWindowsGetPrimaryDomainSID(t *testing.T) {
	sid, err := compat.GetPrimaryDomainSID()
	if err != nil {
		t.Fatalf("got %q, want nil", err)
	}
	if !compat.IsValidSid(sid) {
		t.Fatalf("got an invalid SID: %v", sid.String())
	}
}

func TestStatPosixWindowsGetRid(t *testing.T) {
	// Build a SID manually: S-1-5-32-544 (Administrators)
	raw := []byte{
		1,                // Revision
		2,                // SubAuthorityCount (2 subauthorities)
		0, 0, 0, 0, 0, 5, // IdentifierAuthority = 5 (SECURITY_NT_AUTHORITY)
		32, 0, 0, 0, // SubAuthority[0] = 32 (BUILTIN domain)
		0x20, 0x02, 0, 0, // SubAuthority[1] = 544 (Administrators)
	}

	sid := (*windows.SID)(unsafe.Pointer(&raw[0]))

	rid, err := compat.GetRID(sid)
	if err != nil {
		t.Fatalf("got %q, want nil", err)
	}
	if rid != 544 {
		t.Fatalf("got %d, want RID 544", rid)
	}
}

func TestStatPosixWindowsGetRidInvalid(t *testing.T) {
	raw := []byte{
		1,                // Revision
		0,                // SubAuthorityCount (this is what we want to test)
		0, 0, 0, 0, 0, 0, // IdentifierAuthority (6 bytes)
		// no SubAuthorities follow because count = 0
	}

	sid := (*windows.SID)(unsafe.Pointer(&raw[0]))

	_, err := compat.GetRID(sid)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestStatPosixWindowsSIDToPOSIXID(t *testing.T) {
	tests := []struct {
		name       string
		sidStr     string
		primaryStr string
		wantIDMin  int
		wantIDMax  int
		equalDom   bool
		wantErr    bool
	}{
		{
			name:      "logon session (S-1-5-5-)",
			sidStr:    "S-1-5-5-0-1",
			wantIDMin: 0xFFF,
			wantIDMax: 0xFFF,
		},
		{
			name:      "builtin group (S-1-5-32-)",
			sidStr:    "S-1-5-32-544", // Administrators
			wantIDMin: 0x20000,
			wantIDMax: 0x2FFFF,
		},
		{
			name:       "domain user (same domain)",
			sidStr:     "S-1-5-21-111-222-333-1000",
			primaryStr: "S-1-5-21-111-222-333",
			wantIDMin:  0x40000,
			wantIDMax:  0x4FFFF,
			equalDom:   true,
		},
		{
			name:       "domain user (other domain)",
			sidStr:     "S-1-5-21-111-222-333-1001",
			primaryStr: "S-1-5-21-999-888-777",
			wantIDMin:  0x30000,
			wantIDMax:  0x3FFFF,
			equalDom:   false,
		},
		{
			name:      "Everyone",
			sidStr:    "S-1-1-0",
			wantIDMin: 0x30201,
			wantIDMax: 0x30201,
		},
		{
			name:    "S-1-5-9999999999",
			sidStr:  "S-1-5-9999999999",
			wantErr: true,
		},
		{
			name:    "S-0-0-0",
			sidStr:  "S-0-0-0",
			wantErr: true,
		},
		{
			name:    "S-2-0-0",
			sidStr:  "S-2-0-0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sid, primary *windows.SID
			var err error

			if tt.sidStr != "" {
				sid, err = windows.StringToSid(tt.sidStr)
				if err != nil {
					t.Fatalf("StringToSid(%q) failed: %v", tt.sidStr, err)
				}
			}
			if tt.primaryStr != "" {
				primary, err = windows.StringToSid(tt.primaryStr)
				if err != nil {
					t.Fatalf("StringToSid(%q) failed: %v", tt.primaryStr, err)
				}
			}
			got, err := compat.SIDToPOSIXID(sid, primary)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("got nil, want an error: POSIX id=0x%x (%d)", got, got)
				}

				return
			}
			if err != nil {
				t.Fatalf("got %q, want nil", err)
			}
			if primary == nil {
				return
			}

			if got < tt.wantIDMin {
				t.Errorf("got 0x%x (%d), want >= 0x%x (%d)", got, got, tt.wantIDMin, tt.wantIDMin)
			}
			if got > tt.wantIDMax {
				t.Errorf("got 0x%x (%d), want <= 0x%x (%d)", got, got, tt.wantIDMax, tt.wantIDMax)
			}
			equalDom, _ := compat.EqualDomainSid(sid, primary)
			if tt.equalDom != equalDom {
				t.Errorf("got 0x%x (%d), got %v, want %v", got, got, equalDom, tt.equalDom)
			}
		})
	}
}

func TestStatPosixWindowsNameFromSID(t *testing.T) {
	sid, err := windows.StringToSid("S-1-5-18") // NT AUTHORITY\\SYSTEM
	if err != nil {
		t.Fatalf("failed to create SID: %v", err)
	}

	_, err = compat.NameFromSID(sid)
	if err != nil {
		t.Fatalf("got %q, want nil", err)
	}
}

func TestStatPosixWindowsNameFromSIDInvalid(t *testing.T) {
	_, err := compat.NameFromSID(nil)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestStatPosixWindowsGetUserGroup(t *testing.T) {
	name, err := createTempFile(t)
	if err != nil {
		t.Fatal(err)
	}

	_, _, _, _, err = compat.GetUserGroup(name) //nolint:dogsled
	if err != nil {
		t.Fatalf("got %q, want nil", err)
	}
}

func TestStatPosixWindowsGetUserGroupInvalid(t *testing.T) {
	_, _, _, _, err := compat.GetUserGroup(invalidName) //nolint:dogsled
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestStatPosixWindowsSIDToPOSIXIDInvalid(t *testing.T) {
	_, err := compat.SIDToPOSIXID(nil, nil)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestStatPosixWindowsEqualDomainSidInvalid(t *testing.T) {
	_, err := compat.EqualDomainSid(nil, nil)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func TestStatPosixWindowsIsValidSidInvalid(t *testing.T) {
	b := compat.IsValidSid(nil)
	if b {
		t.Fatal("got true, want false")
	}
}
