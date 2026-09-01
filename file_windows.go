// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"sync"
	"syscall"
	"unsafe"

	"github.com/capnspacehook/go-acl"
	"golang.org/x/sys/windows"

	"github.com/rasa/compat/golang"
	"github.com/rasa/compat/robustio"
)

// UnknownUsername is returned when the current username is not available.
const UnknownUsername = "n/a"

// type tokenPrimaryGroup struct {
//	PrimaryGroup *windows.SID
// }

type securityInfo struct {
	ownerSid *windows.SID
	groupSid *windows.SID
	acl      *windows.ACL
	perm     os.FileMode
}

func chmod(name string, perm os.FileMode, opts options) error {
	if !opts.skipACLs {
		// acl.Chmod will panic otherwise
		_, err := windows.UTF16PtrFromString(name)
		if err != nil {
			return chmodError(name, os.ErrInvalid)
		}

		perm = opts.fileMode.Perm()

		// set Windows' ACLs
		err = acl.Chmod(name, perm)
		if err != nil {
			return chmodError(name, fmt.Errorf("%w (acl)", err))
		}
	}

	if opts.readOnlyMode == ReadOnlyModeIgnore {
		return nil
	}

	fi, err := os.Stat(name)
	if err != nil {
		return chmodError(name, fmt.Errorf("%w (stat)", err))
	}

	// Set or clear Windows' read-only attribute
	want := perm&windows.S_IWRITE != 0 // 0x80 (0o200)

	got := fi.Mode().Perm()&windows.S_IWRITE != 0

	if opts.readOnlyMode == ReadOnlyModeClear {
		want = false
	}

	if want == got {
		return nil
	}

	if want {
		perm |= windows.S_IWRITE
	} else {
		perm &^= os.FileMode(windows.S_IWRITE)
	}

	err = os.Chmod(name, perm)
	if err != nil {
		return chmodError(name, fmt.Errorf("%w (chmod)", err))
	}

	return nil
}

func create(name string, opts options) (*os.File, error) {
	if opts.readOnlyMode != ReadOnlyModeFromPermissions {
		opts.openFlags |= golang.O_FILE_FLAG_NO_RO_ATTR
	}

	if opts.deleteOnClose {
		opts.openFlags |= golang.FILE_FLAG_DELETE_ON_CLOSE
	}

	sa, err := saFromPerm(opts.fileMode, true)
	if err != nil {
		return nil, createError(name, err)
	}

	return golang.OpenFileNolog(name, opts.openFlags, opts.fileMode, sa)
}

func createTemp(dir, pattern string, opts options) (*os.File, error) {
	if opts.deleteOnClose {
		opts.openFlags |= golang.FILE_FLAG_DELETE_ON_CLOSE
	}

	if opts.readOnlyMode != ReadOnlyModeFromPermissions {
		opts.openFlags |= golang.O_FILE_FLAG_NO_RO_ATTR
	}

	sa, err := saFromPerm(opts.fileMode, true)
	if err != nil {
		return nil, createTempError(dir, err)
	}

	fp, err := golang.CreateTemp(dir, pattern, opts.openFlags, opts.fileMode, sa)
	if err != nil {
		return nil, createTempError(dir, err)
	}

	return fp, nil
}

func fchmod(fp *os.File, mode os.FileMode, opts options) error {
	if fp == nil {
		return chmodError("", os.ErrInvalid)
	}

	path, err := golang.Filepath(fp)
	if err == nil {
		err = chmod(path, mode, opts)
		if err == nil {
			return nil
		}
	}

	return chmodError(path, err)
}

func mkdir(name string, perm os.FileMode) error {
	sa, err := saFromPerm(perm, true)
	if err == nil {
		err = golang.Mkdir(name, perm, sa)
		if err == nil {
			return nil
		}
	}

	return mkdirError(name, err)
}

func mkdirAll(name string, perm os.FileMode) error {
	sa, err := saFromPerm(perm, true)
	if err == nil {
		err = golang.MkdirAll(name, perm, sa)
		if err == nil {
			return nil
		}
	}

	return mkdirTempError(name, err)
}

func mkdirTemp(dir, pattern string, opts options) (string, error) {
	sa, err := saFromPerm(opts.fileMode, true)

	var path string
	if err == nil {
		path, err = golang.MkdirTemp(dir, pattern, sa)
		if err == nil {
			return path, nil
		}
	}

	prefix, suffix, _ := golang.PrefixAndSuffix(pattern)

	return "", mkdirError(dir+string(os.PathSeparator)+prefix+"*"+suffix, err)
}

func openFile(name string, opts options) (*os.File, error) {
	if opts.deleteOnClose {
		opts.openFlags |= golang.FILE_FLAG_DELETE_ON_CLOSE
	}

	sa, err := saFromPerm(opts.fileMode, (opts.openFlags&os.O_CREATE) == os.O_CREATE)
	if err != nil {
		return nil, openError(name, err)
	}

	return golang.OpenFileNolog(name, opts.openFlags, opts.fileMode, sa)
}

func remove(path string, opts options) error {
	if opts.retryTimeout == 0 {
		return golang.Remove(path)
	}

	return robustio.Retry(func() (err error, mayRetry bool) {
		err = golang.Remove(path)

		return err, robustio.IsEphemeralError(err)
	}, opts.retryTimeout.Seconds())
}

func removeAll(path string, opts options) error {
	if opts.retryTimeout == 0 {
		return golang.RemoveAll(path)
	}

	return robustio.Retry(func() (err error, mayRetry bool) {
		err = golang.RemoveAll(path)

		return err, robustio.IsEphemeralError(err)
	}, opts.retryTimeout.Seconds())
}

func symlink(oldname, newname string, opts options) error {
	err := os.Symlink(oldname, newname)
	if err != nil {
		return err
	}

	if opts.setSymlinkOwner {
		err = setOwnerToCurrentUser(newname)
		if err != nil {
			return symlinkError(oldname, newname, err)
		}
	}

	return nil
}

// saFromPerm converts a perm (FileMode) to an *sa (*syscall.SecurityAttributes).
// @TODO return a *windows.SecurityAttributes.
func saFromPerm(perm os.FileMode, create bool) (*syscall.SecurityAttributes, error) {
	var sa syscall.SecurityAttributes

	sa.Length = uint32(unsafe.Sizeof(sa))

	if !create {
		return &sa, nil
	}

	umask, err := GetUmask()
	if err != nil {
		return nil, err
	}
	perm &^= os.FileMode(umask) //nolint:gosec

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
	umask, err := GetUmask()
	if err != nil {
		return nil, err
	}
	perm &^= os.FileMode(umask) //nolint:gosec

	ownerSid, groupSid, worldSid, err := getSIDs()
	if err != nil {
		return nil, err
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

var getSIDsOnce struct {
	sync.Once

	ownerSID *windows.SID
	groupSID *windows.SID
	worldSID *windows.SID
	err      error
}

func getSIDs() (*windows.SID, *windows.SID, *windows.SID, error) {
	getSIDsOnce.Do(func() {
		getSIDsOnce.ownerSID, getSIDsOnce.groupSID, getSIDsOnce.worldSID, getSIDsOnce.err = _getSIDs()
	})

	return getSIDsOnce.ownerSID, getSIDsOnce.groupSID, getSIDsOnce.worldSID, getSIDsOnce.err
}

func _getSIDs() (*windows.SID, *windows.SID, *windows.SID, error) {
	// // Get current user's SID
	token := windows.Token(0)

	err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get process token for %s: %w", currentUsername(), err)
	}
	defer token.Close()

	ownerSID, err := getOwnerSID(token) // works: getOwnerSID(token)
	if err != nil {
		return nil, nil, nil, err
	}

	groupSID, err := getPrimaryGroupSID(token) // works: getGroupSID(token)
	if err != nil {
		return nil, nil, nil, err
	}

	worldSID, err := windows.CreateWellKnownSid(windows.WinWorldSid)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create world SID: %w", err)
	}

	return ownerSID, groupSID, worldSID, nil
}

func getOwnerSID(token windows.Token) (*windows.SID, error) {
	var err error

	bufSize := initialBufSize

	var buf0 *byte

	for {
		var newBufSize uint32

		buf := make([]byte, bufSize)
		if bufSize > 0 {
			buf0 = &buf[0]
		}

		err = windows.GetTokenInformation(
			token,
			windows.TokenUser,
			buf0,
			bufSize,
			&newBufSize,
		)
		if err == nil {
			tu := (*windows.Tokenuser)(unsafe.Pointer(&buf[0]))

			return tu.User.Sid, nil
		}

		if !errors.Is(err, windows.ERROR_INSUFFICIENT_BUFFER) {
			return nil, fmt.Errorf("failed to get token information: %w", err)
		}

		if newBufSize > bufSize {
			bufSize = newBufSize
		} else {
			bufSize *= 2
		}
	}
}

// @TODO(rasa) improve this logic per
// https://github.com/golang/go/blob/cc8a6780/src/os/user/lookup_windows.go#L351
func getPrimaryGroupSID(token windows.Token) (*windows.SID, error) {
	// @TODO TEST IF  windows.GetTokenPrimaryGroup() can replace.
	var err error

	bufSize := initialBufSize

	var buf0 *byte

	for {
		var newBufSize uint32

		buf := make([]byte, bufSize)
		if bufSize > 0 {
			buf0 = &buf[0]
		}

		err = windows.GetTokenInformation(
			token,
			windows.TokenPrimaryGroup,
			buf0,
			bufSize,
			&newBufSize,
		)
		if err == nil {
			pg := (*windows.Tokenprimarygroup)(unsafe.Pointer(&buf[0]))

			return pg.PrimaryGroup, nil
		}

		if !errors.Is(err, windows.ERROR_INSUFFICIENT_BUFFER) {
			return nil, fmt.Errorf("failed to get token information: %w", err)
		}

		if newBufSize > bufSize {
			bufSize = newBufSize
		} else {
			bufSize *= 2
		}
	}
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

	err = sd.SetControl(
		windows.SE_DACL_PROTECTED,
		windows.SE_DACL_PROTECTED,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to set control on security descriptor: %w", err)
	}

	return sd, nil
}

func currentUsername() string {
	usr, err := user.Current()
	if err != nil {
		return UnknownUsername
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

var (
	seTakeOwnershipPrivilegeW, _ = windows.UTF16PtrFromString("SeTakeOwnershipPrivilege")
	seRestorePrivilegeW, _       = windows.UTF16PtrFromString("SeRestorePrivilege")
)

func setOwnerToCurrentUser(path string) error {
	var tok windows.Token

	err := windows.OpenProcessToken(
		windows.CurrentProcess(),
		windows.TOKEN_ADJUST_PRIVILEGES|windows.TOKEN_QUERY,
		&tok,
	)
	if err != nil {
		return fmt.Errorf("failed to open process token: %w", err)
	}
	defer tok.Close()

	// Current user SID (needs TOKEN_QUERY)
	tu, err := tok.GetTokenUser()
	if err != nil {
		return fmt.Errorf("failed to get token user: %w", err)
	}

	userSID := tu.User.Sid

	// Enable SeTakeOwnershipPrivilege (required to take ownership when you don't own it)
	err = enablePrivilege(tok, seTakeOwnershipPrivilegeW)
	if err != nil {
		return fmt.Errorf("failed to take ownership privilege: %w", err)
	}
	// Optional, sometimes helpful
	err = enablePrivilege(tok, seRestorePrivilegeW)
	if err != nil {
		return fmt.Errorf("failed to restore privileges: %w", err)
	}

	// Set owner by name (affects target if path is a symlink)
	err = windows.SetNamedSecurityInfo(
		path,
		windows.SE_FILE_OBJECT,
		windows.OWNER_SECURITY_INFORMATION,
		userSID, nil, nil, nil,
	)
	if err != nil {
		return fmt.Errorf("failed to set named security info: %w", err)
	}

	return nil
}

func enablePrivilege(tok windows.Token, name *uint16) error {
	var luid windows.LUID

	err := windows.LookupPrivilegeValue(nil, name, &luid)
	if err != nil {
		return fmt.Errorf("failed to lookup privilege: %w", err)
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
		return fmt.Errorf("failed to adjust token privileges: %w", err)
	}
	// AdjustTokenPrivileges can "succeed" but not assign; check last error.
	le := windows.GetLastError()
	if errors.Is(le, windows.ERROR_NOT_ALL_ASSIGNED) {
		return fmt.Errorf("failed to hold privilege: %w", le)
	}

	return nil
}
