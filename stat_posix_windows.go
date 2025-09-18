// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

import (
	"fmt"
	"os"
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
	procCopySid                   = modadvapi32.NewProc("CopySid")
	procEqualDomainSid            = modadvapi32.NewProc("EqualDomainSid")
	procGetNamedSecurityInfoW     = modadvapi32.NewProc("GetNamedSecurityInfoW")
	procIsValidSid                = modadvapi32.NewProc("IsValidSid") //nolint:unused
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
		return nil, nil, fmt.Errorf("failed to get named security info: %w", syscall.Errno(r0))
	}

	return owner, group, nil
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
		return nil, fmt.Errorf("failed to query information policy: %w", syscall.Errno(r0))
	}
	defer procLsaFreeMemory.Call(buffer) //nolint:errcheck

	info := (*LSA_POLICY_ACCOUNT_DOMAIN_INFO)(unsafe.Pointer(buffer))

	sidCopy, err := copySid(info.DomainSid)
	if err != nil {
		return nil, err
	}

	return sidCopy, nil
}

func getRID(sid *windows.SID) (int, error) {
	count := uint32(sid.SubAuthorityCount())
	if count == 0 {
		return UnknownID, fmt.Errorf("no subauthorities found for %q", sid.String())
	}

	return int(sid.SubAuthority(count - 1)), nil
}

var wellKnownSIDToGID = map[string]int{
	"S-1-1-0":      0x30201, //nolint:mnd // 197121 // Everyone
	"S-1-5-5-":     0xfff,   //nolint:mnd
	"S-1-5-32-544": 544,     //nolint:mnd // Administrators
	"S-1-5-32-545": 545,     //nolint:mnd // Users
	"S-1-5-32-546": 546,     //nolint:mnd // Guests
	"S-1-5-32-547": 547,     //nolint:mnd // Power Users
	"S-1-5-32-548": 548,     //nolint:mnd // Account Operators
	"S-1-5-32-549": 549,     //nolint:mnd // Server Operators
	"S-1-5-32-550": 550,     //nolint:mnd // Print Operators
	"S-1-5-32-551": 551,     //nolint:mnd // Backup Operators
	"S-1-5-32-552": 552,     //nolint:mnd // Replicators
}

// See https://cygwin.com/cygwin-ug-net/ntsec.html
func sidToPOSIXID(sid *windows.SID, primaryDomainSid *windows.SID) (int, error) {
	if sid == nil {
		return 0, os.ErrInvalid
	}

	sidStr := sid.String()

	switch {
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
		b, err := equalDomainSid(sid, primaryDomainSid)
		if err != nil {
			return UnknownID, err
		}
		if b {
			return 0x40000 + rid, nil //nolint:mnd
		}

		return 0x30000 + rid, nil //nolint:mnd
	default:
		// fallthrough
	}

	for prefix, gid := range wellKnownSIDToGID {
		if strings.HasPrefix(sidStr, prefix) {
			return gid, nil
		}
	}

	return UnknownID, fmt.Errorf("unsupported SID: %s", sidStr)
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

func equalDomainSid(sid1, sid2 *windows.SID) (bool, error) {
	if sid1 == nil || sid2 == nil {
		return false, nil
	}

	var equal int32
	r1, _, e1 := syscall.SyscallN(
		procEqualDomainSid.Addr(),
		uintptr(unsafe.Pointer(sid1)),
		uintptr(unsafe.Pointer(sid2)),
		uintptr(unsafe.Pointer(&equal)),
	)
	if r1 == 0 {
		if e1 != 0 {
			return false, error(e1)
		}

		return false, syscall.EINVAL
	}

	return equal != 0, nil
}

// Used only in the tests, for now.
func isValidSid(sid *windows.SID) bool { //nolint:unused
	if sid == nil {
		return false
	}
	r1, _, _ := syscall.SyscallN(
		procIsValidSid.Addr(),
		uintptr(unsafe.Pointer(sid)),
	)

	return r1 != 0
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
		return syscall.InvalidHandle, fmt.Errorf("failed to open LSA policy: %w", syscall.Errno(r0))
	}

	return handle, nil
}

func copySid(src *windows.SID) (*windows.SID, error) {
	if src == nil {
		return nil, syscall.EINVAL
	}

	length := windows.GetLengthSid(src)
	dst := make([]byte, length)

	ret, _, err := procCopySid.Call(
		uintptr(length),
		uintptr(unsafe.Pointer(&dst[0])),
		uintptr(unsafe.Pointer(src)),
	)
	if ret == 0 {
		return nil, err
	}
	return (*windows.SID)(unsafe.Pointer(&dst[0])), nil
}
