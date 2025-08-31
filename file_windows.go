// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"syscall"
	"unsafe"

	"github.com/capnspacehook/go-acl"
	"golang.org/x/sys/windows"

	"github.com/rasa/compat/golang"
)

type tokenPrimaryGroup struct {
	PrimaryGroup *windows.SID
}

type securityInfo struct {
	ownerSid *windows.SID
	groupSid *windows.SID
	acl      *windows.ACL
	perm     os.FileMode
}

func chmod(name string, perm os.FileMode, mask ReadOnlyMode) error {
	perm = perm.Perm()

	// set Windows' ACLs
	err := acl.Chmod(name, perm)
	if err != nil {
		return fmt.Errorf("%w (acl)", err)
	}

	if mask == ReadOnlyModeIgnore {
		return nil
	}

	fi, err := os.Stat(name)
	if err != nil {
		return fmt.Errorf("%w (stat)", err)
	}

	// Set or clear Windows' read-only attribute
	want := perm&syscall.S_IWRITE != 0 // 0x80 (0o200)
	got := fi.Mode().Perm()&syscall.S_IWRITE != 0

	if want == got {
		return nil
	}

	if mask == ReadOnlyModeReset {
		if !got {
			return nil
		}
		want = false
	}

	if want {
		perm |= syscall.S_IWRITE
	} else {
		perm &= ^os.FileMode(syscall.S_IWRITE)
	}
	err = os.Chmod(name, perm)
	if err != nil {
		return fmt.Errorf("%w (chmod)", err)
	}

	return nil
}

func create(name string, perm os.FileMode, flag int) (*os.File, error) {
	flag |= os.O_CREATE
	sa, err := saFromPerm(perm, true)
	if err != nil {
		return nil, err
	}

	return golang.OpenFileNolog(name, flag, perm, sa)
}

func createTemp(dir, pattern string, perm os.FileMode, flag int) (*os.File, error) {
	if perm == 0 {
		perm = CreateTempPerm
	}
	sa, err := saFromPerm(perm, true) // 0o600
	if err != nil {
		return nil, err
	}

	return golang.CreateTemp(dir, pattern, flag, sa)
}

func fchmod(f *os.File, mode os.FileMode, mask ReadOnlyMode) error {
	if f == nil {
		return errors.New("nil file pointer")
	}
	fd := syscall.Handle(f.Fd())

	// Source: https://github.com/golang/go/blob/77f911e3/src/syscall/syscall_windows.go#L1294-L1310

	var buf [syscall.MAX_PATH + 1]uint16
	path, err := fdpath(fd, buf[:])
	if err != nil {
		return err
	}
	// When using VOLUME_NAME_DOS, the path is always prefixed by "\\?\".
	// That prefix tells the Windows APIs to disable all string parsing and to send
	// the string that follows it straight to the file system.
	// Although SetCurrentDirectory and GetCurrentDirectory do support the "\\?\" prefix,
	// some other Windows APIs don't. If the prefix is not removed here, it will leak
	// to Getwd, and we don't want such a general-purpose function to always return a
	// path with the "\\?\" prefix after Fchdir is called.
	// The downside is that APIs that do support it will parse the path and try to normalize it,
	// when it's already normalized.
	if len(path) >= 4 && path[0] == '\\' && path[1] == '\\' && path[2] == '?' && path[3] == '\\' {
		path = path[4:]
	}
	pathString := syscall.UTF16ToString(path)

	return chmod(pathString, mode, mask)

}

// Source: https://github.com/golang/go/blob/77f911e3/src/syscall/syscall_windows.go#L183

const _ERROR_NOT_ENOUGH_MEMORY = syscall.Errno(8)

// Source: https://github.com/golang/go/blob/77f911e3/src/syscall/syscall_windows.go#L1274-L1291

func fdpath(fd syscall.Handle, buf []uint16) ([]uint16, error) {
	const (
		FILE_NAME_NORMALIZED = 0
		VOLUME_NAME_DOS      = 0
	)
	for {
		n, err := golang.GetFinalPathNameByHandle(fd, &buf[0], uint32(len(buf)), FILE_NAME_NORMALIZED|VOLUME_NAME_DOS)
		if err == nil {
			buf = buf[:n]
			break
		}
		if err != _ERROR_NOT_ENOUGH_MEMORY {
			return nil, err
		}
		buf = append(buf, make([]uint16, n-uint32(len(buf)))...)
	}
	return buf, nil
}

func mkdir(name string, perm os.FileMode) error {
	sa, err := saFromPerm(perm, true)
	if err != nil {
		return err
	}

	return golang.Mkdir(name, 0o700, sa) //nolint:mnd
}

func mkdirAll(name string, perm os.FileMode) error {
	sa, err := saFromPerm(perm, true)
	if err != nil {
		return err
	}

	return golang.MkdirAll(name, 0o700, sa) //nolint:mnd
}

func mkdirTemp(dir, pattern string, perm os.FileMode) (string, error) {
	sa, err := saFromPerm(perm, true) // 0o700
	if err != nil {
		return "", err
	}

	return golang.MkdirTemp(dir, pattern, sa)
}

func openFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	sa, err := saFromPerm(perm, (flag&os.O_CREATE) == os.O_CREATE)
	if err != nil {
		return nil, err
	}

	return golang.OpenFileNolog(name, flag, perm, sa)
}

func remove(name string) error {
	return golang.Remove(name)
}

func removeAll(path string) error {
	return golang.RemoveAll(path)
}

func symlink(oldname, newname string, setSymlinkOwner bool) error {
	err := os.Symlink(oldname, newname)
	if err != nil {
		return err
	}

	if !setSymlinkOwner {
		return nil
	}

	return setOwnerToCurrentUser(newname)
}

func writeFile(name string, data []byte, perm os.FileMode, flag int) error {
	flag |= os.O_CREATE

	f, err := openFile(name, flag, perm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		_ = f.Close()
		return err
	}

	return f.Close()
}

// saFromPerm converts a perm (FileMode) to an *sa (*syscall.SecurityAttributes).
// @TODO return a *windows.SecurityAttributes.
func saFromPerm(perm os.FileMode, create bool) (*syscall.SecurityAttributes, error) {
	var sa syscall.SecurityAttributes
	sa.Length = uint32(unsafe.Sizeof(sa))

	if !create {
		return &sa, nil
	}

	perm &^= os.FileMode(GetUmask()) //nolint:gosec
	sd, err := sdFromPerm(perm)
	if err != nil {
		return nil, err
	}

	sa.SecurityDescriptor = uintptr(unsafe.Pointer(sd))
	sa.InheritHandle = 0

	return &sa, nil
}

// siFromPerm converts a perm (FileMode) to an *si (*securityInfo).
func siFromPerm(perm os.FileMode) (*securityInfo, error) {
	perm &^= os.FileMode(GetUmask()) //nolint:gosec

	// Get current user's SID
	token := windows.Token(0)
	err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)
	if err != nil {
		return nil, fmt.Errorf("failed to get process token for %s: %w", currentUsername(), err)
	}
	defer token.Close()

	var size uint32
	// First call to get required buffer size
	_ = windows.GetTokenInformation(token, windows.TokenUser, nil, 0, &size)

	buf := make([]byte, size)
	err = windows.GetTokenInformation(token, windows.TokenUser, &buf[0], size, &size)
	if err != nil {
		return nil, fmt.Errorf("failed to get user token information for %s: %w", currentUsername(), err)
	}

	tu := (*windows.Tokenuser)(unsafe.Pointer(&buf[0]))
	ownerSid := tu.User.Sid

	_ = windows.GetTokenInformation(token, windows.TokenPrimaryGroup, nil, 0, &size)

	buf = make([]byte, size)
	err = windows.GetTokenInformation(token, windows.TokenPrimaryGroup, &buf[0], size, &size)
	if err != nil {
		return nil, fmt.Errorf("failed to get group token information for %s: %w", currentUsername(), err)
	}

	tg := (*tokenPrimaryGroup)(unsafe.Pointer(&buf[0]))
	groupSid := tg.PrimaryGroup

	worldSid, err := windows.CreateWellKnownSid(windows.WinWorldSid)
	if err != nil {
		return nil, fmt.Errorf("failed to create world SID: %w", err)
	}

	var ea [3]windows.EXPLICIT_ACCESS

	ownerMask := accessMask(perm, 6) //nolint:mnd
	setExplicitAccess(&ea[0], ownerSid, ownerMask, windows.TRUSTEE_IS_USER)

	groupMask := accessMask(perm, 3) //nolint:mnd
	setExplicitAccess(&ea[1], groupSid, groupMask, windows.TRUSTEE_IS_GROUP)

	worldMask := accessMask(perm, 0)
	setExplicitAccess(&ea[2], worldSid, worldMask, windows.TRUSTEE_IS_WELL_KNOWN_GROUP)

	dumpMasks(perm, ownerMask, groupMask, worldMask)

	acl, err := windows.ACLFromEntries(ea[:], nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ACLs: %w", err)
	}

	si := securityInfo{ownerSid, groupSid, acl, perm}

	return &si, nil
}

// sdFromPerm converts a perm (FileMode) to an *sd (*windows.SECURITY_DESCRIPTOR).
func sdFromPerm(perm os.FileMode) (*windows.SECURITY_DESCRIPTOR, error) {
	si, err := siFromPerm(perm)
	if err != nil {
		return nil, err
	}

	sd, err := sdFromSi(*si)
	if err != nil {
		return nil, err
	}

	return sd, err
}

// sdFromSi converts a si (securityInfo) to an *sd (*windows.SECURITY_DESCRIPTOR).
func sdFromSi(si securityInfo) (*windows.SECURITY_DESCRIPTOR, error) {
	sd, err := windows.NewSecurityDescriptor()
	if err != nil {
		return nil, fmt.Errorf("failed to create security descriptor: %w", err)
	}
	err = sd.SetOwner(si.ownerSid, false)
	if err != nil {
		return nil, fmt.Errorf("failed to set ACL owner in security descriptor: %w", err)
	}
	err = sd.SetGroup(si.groupSid, false)
	if err != nil {
		return nil, fmt.Errorf("failed to set ACL group in security descriptor: %w", err)
	}
	err = sd.SetDACL(si.acl, true, false)
	if err != nil {
		return nil, fmt.Errorf("failed to set ACL in security descriptor: %w", err)
	}

	return sd, nil
}

func currentUsername() string {
	usr, err := user.Current()
	if err != nil {
		return "n/a"
	}

	return usr.Username
}

func accessMask(mode os.FileMode, shift int) uint32 {
	perm := uint32(mode.Perm())

	var mask uint32

	if perm&(0o4<<shift) == (0o4 << shift) { //nolint:mnd
		mask |= windows.GENERIC_READ
	}
	if perm&(0o2<<shift) == (0o2 << shift) { //nolint:mnd
		mask |= windows.GENERIC_WRITE | windows.DELETE
	}
	if perm&(0o1<<shift) == (0o1 << shift) { //nolint:mnd
		mask |= windows.GENERIC_EXECUTE
	}

	return mask
}

func setExplicitAccess(ea *windows.EXPLICIT_ACCESS, sid *windows.SID, mask uint32, tt windows.TRUSTEE_TYPE) {
	ea.AccessPermissions = windows.ACCESS_MASK(mask)
	ea.AccessMode = windows.SET_ACCESS
	ea.Inheritance = windows.NO_INHERITANCE // was windows.SUB_CONTAINERS_AND_OBJECTS_INHERIT
	ea.Trustee.TrusteeForm = windows.TRUSTEE_IS_SID
	ea.Trustee.TrusteeType = tt
	ea.Trustee.TrusteeValue = windows.TrusteeValueFromSID(sid)
}

func setOwnerToCurrentUser(path string) error {
	var tok windows.Token
	err := windows.OpenProcessToken(
		windows.CurrentProcess(),
		windows.TOKEN_ADJUST_PRIVILEGES|windows.TOKEN_QUERY,
		&tok,
	)
	if err != nil {
		return fmt.Errorf("OpenProcessToken: %w", err)
	}
	defer tok.Close()

	// Current user SID (needs TOKEN_QUERY)
	tu, err := tok.GetTokenUser()
	if err != nil {
		return fmt.Errorf("GetTokenUser: %w", err)
	}
	userSID := tu.User.Sid

	// Enable SeTakeOwnershipPrivilege (required to take ownership when you don't own it)
	err = enablePrivilege(tok, seTakeOwnershipPrivilegeW)
	if err != nil {
		return fmt.Errorf("enable SeTakeOwnershipPrivilege: %w", err)
	}
	// Optional, sometimes helpful
	err = enablePrivilege(tok, seRestorePrivilegeW)
	if err != nil {
		return fmt.Errorf("seRestorePrivilegeW: %w", err)
	}

	// Set owner by name (affects target if path is a symlink)
	err = windows.SetNamedSecurityInfo(
		path,
		windows.SE_FILE_OBJECT,
		windows.OWNER_SECURITY_INFORMATION,
		userSID, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("SetNamedSecurityInfo: %w", err)
	}

	return nil
}

func enablePrivilege(tok windows.Token, name *uint16) error {
	var luid windows.LUID
	err := windows.LookupPrivilegeValue(nil, name, &luid)
	if err != nil {
		return fmt.Errorf("LookupPrivilegeValue: %w", err)
	}

	tp := windows.Tokenprivileges{
		PrivilegeCount: 1,
		Privileges: [1]windows.LUIDAndAttributes{{
			Luid:       luid,
			Attributes: windows.SE_PRIVILEGE_ENABLED,
		}},
	}

	// Must be called on a real token handle opened with TOKEN_ADJUST_PRIVILEGES.
	err = windows.AdjustTokenPrivileges(tok, false, &tp, 0, nil, nil)
	if err != nil {
		return fmt.Errorf("AdjustTokenPrivileges: %w", err)
	}
	// AdjustTokenPrivileges can "succeed" but not assign; check last error.
	if le := windows.GetLastError(); errors.Is(le, windows.ERROR_NOT_ALL_ASSIGNED) {
		return fmt.Errorf("privilege not held: %w", le)
	}

	return nil
}

var (
	seTakeOwnershipPrivilegeW, _ = windows.UTF16PtrFromString("SeTakeOwnershipPrivilege")
	seRestorePrivilegeW, _       = windows.UTF16PtrFromString("SeRestorePrivilege")
)
