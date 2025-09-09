// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	PolicyAccountDomainInformation = 5
	POLICY_VIEW_LOCAL_INFORMATION  = 0x00000001
	OWNER_SECURITY_INFORMATION     = 0x00000001
	GROUP_SECURITY_INFORMATION     = 0x00000002
)

type LSA_UNICODE_STRING struct {
	Length        uint16
	MaximumLength uint16
	Buffer        *uint16
}

type LSA_OBJECT_ATTRIBUTES struct {
	Length                   uint32
	RootDirectory            uintptr
	ObjectName               *LSA_UNICODE_STRING
	Attributes               uint32
	SecurityDescriptor       uintptr
	SecurityQualityOfService uintptr
}

type LSA_POLICY_ACCOUNT_DOMAIN_INFO struct {
	DomainName LSA_UNICODE_STRING
	DomainSid  *windows.SID
}

var (
	modadvapi32                   = windows.NewLazySystemDLL("advapi32.dll")
	procGetNamedSecurityInfoW     = modadvapi32.NewProc("GetNamedSecurityInfoW")
	procLsaOpenPolicy             = modadvapi32.NewProc("LsaOpenPolicy")
	procLsaQueryInformationPolicy = modadvapi32.NewProc("LsaQueryInformationPolicy")
	procLsaFreeMemory             = modadvapi32.NewProc("LsaFreeMemory")
	procLsaClose                  = modadvapi32.NewProc("LsaClose")
)

func getFileOwnerAndGroupSIDs(name string) (*windows.SID, *windows.SID, error) {
	var owner, group *windows.SID
	pPath, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid file path: %w", err)
	}
	r0, _, _ := procGetNamedSecurityInfoW.Call(
		uintptr(unsafe.Pointer(pPath)),
		1, // SE_FILE_OBJECT
		OWNER_SECURITY_INFORMATION|GROUP_SECURITY_INFORMATION,
		uintptr(unsafe.Pointer(&owner)),
		uintptr(unsafe.Pointer(&group)),
		0, 0, 0,
	)
	if r0 != 0 {
		return nil, nil, fmt.Errorf("GetNamedSecurityInfo failed: %w", syscall.Errno(r0))
	}

	return owner, group, nil
}

func lsaOpenPolicy(systemName *uint16, access uint32) (handle syscall.Handle, err error) {
	var objectAttrs LSA_OBJECT_ATTRIBUTES
	r0, _, _ := procLsaOpenPolicy.Call(
		uintptr(unsafe.Pointer(systemName)),
		uintptr(unsafe.Pointer(&objectAttrs)),
		uintptr(access),
		uintptr(unsafe.Pointer(&handle)),
	)
	if r0 != 0 {
		return syscall.InvalidHandle, fmt.Errorf("LsaOpenPolicy failed: %w", syscall.Errno(r0))
	}

	return handle, nil
}

func getPrimaryDomainSID() (*windows.SID, error) {
	handle, err := lsaOpenPolicy(nil, POLICY_VIEW_LOCAL_INFORMATION)
	if err != nil {
		return nil, err
	}
	defer procLsaClose.Call(uintptr(handle)) //nolint:errcheck

	var buffer uintptr
	r0, _, _ := procLsaQueryInformationPolicy.Call(
		uintptr(handle),
		uintptr(PolicyAccountDomainInformation),
		uintptr(unsafe.Pointer(&buffer)),
	)
	if r0 != 0 {
		return nil, fmt.Errorf("LsaQueryInformationPolicy failed: %w", syscall.Errno(r0))
	}
	defer procLsaFreeMemory.Call(buffer) //nolint:errcheck

	info := (*LSA_POLICY_ACCOUNT_DOMAIN_INFO)(unsafe.Pointer(buffer))

	return info.DomainSid, nil
}

func getRID(sid *windows.SID) (int, error) {
	count := uint32(sid.SubAuthorityCount())
	if count == 0 {
		return UnknownID, fmt.Errorf("no subauthorities found for %q", sid.String())
	}

	return int(sid.SubAuthority(count - 1)), nil
}

func isSameDomainSID(sid1, sid2 *windows.SID) bool {
	if sid1 == nil || sid2 == nil {
		return false
	}
	// Compare domain portion (strip RID)
	s1 := sid1.String()
	s2 := sid2.String()
	last1 := strings.LastIndex(s1, "-")
	last2 := strings.LastIndex(s2, "-")

	return last1 > 0 && last2 > 0 && s1[:last1] == s2[:last2]
}

// See https://cygwin.com/cygwin-ug-net/ntsec.html
func sidToPOSIXID(sid *windows.SID, primaryDomainSid *windows.SID) (int, error) {
	sidStr := sid.String()

	switch {
	case strings.HasPrefix(sidStr, "S-1-5-5-"):
		return 0xFFF, nil //nolint:mnd
	case strings.HasPrefix(sidStr, "S-1-5-32-"):
		rid, err := getRID(sid)
		if err != nil {
			return UnknownID, err
		}
		return 0x20000 + rid, nil //nolint:mnd
	case strings.HasPrefix(sidStr, "S-1-5-21-"):
		rid, err := getRID(sid)
		if err != nil {
			return UnknownID, err
		}
		if isSameDomainSID(sid, primaryDomainSid) {
			return 0x40000 + rid, nil //nolint:mnd
		}

		return 0x30000 + rid, nil //nolint:mnd
	default:

		return UnknownID, fmt.Errorf("unsupported SID: %s", sidStr)
	}
}

func nameFromSID(sid *windows.SID) (string, error) {
	name16 := make([]uint16, 256)      //nolint:mnd
	domain16 := make([]uint16, 256)    //nolint:mnd
	nameLen := uint32(len(name16))     //nolint:gosec
	domainLen := uint32(len(domain16)) //nolint:gosec
	var sidUse uint32

	err := windows.LookupAccountSid(
		nil, // Local system
		sid,
		&name16[0],
		&nameLen,
		&domain16[0],
		&domainLen,
		&sidUse,
	)
	if err != nil {
		return "", fmt.Errorf("cannot get name from SID %q: %w", sid.String(), err)
	}

	name := syscall.UTF16ToString(name16[:nameLen])
	domain := syscall.UTF16ToString(domain16[:domainLen])

	if domain != "" {
		name = domain + `\` + name
	}

	return name, nil
}

func getUserGroup(path string) (int, int, string, string, error) {
	ownerSID, groupSID, err := getFileOwnerAndGroupSIDs(path)
	if err != nil {
		return UnknownID, UnknownID, "", "", err
	}

	primaryDomainSID, err := getPrimaryDomainSID()
	if err != nil {
		return UnknownID, UnknownID, "", "", err
	}

	uid, err := sidToPOSIXID(ownerSID, primaryDomainSID)
	if err != nil {
		return UnknownID, UnknownID, "", "", err
	}

	gid, err := sidToPOSIXID(groupSID, primaryDomainSID)
	if err != nil {
		return UnknownID, UnknownID, "", "", err
	}

	user, err := nameFromSID(ownerSID)
	if err != nil {
		return UnknownID, UnknownID, "", "", err
	}

	group, err := nameFromSID(groupSID)
	if err != nil {
		return UnknownID, UnknownID, "", "", err
	}

	return uid, gid, user, group, nil
}
