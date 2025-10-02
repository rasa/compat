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

type tokenPrimaryGroup struct {
	PrimaryGroup *windows.SID
}

type securityInfo struct {
	ownerSid *windows.SID
	groupSid *windows.SID
	acl      *windows.ACL
	perm     os.FileMode
}

func chmod(name string, perm os.FileMode, opts ...Option) error {
	fopts := Options{
		fileMode: perm,
	}

	for _, opt := range opts {
		opt(&fopts)
	}

	// acl.Chmod will panic otherwise
	_, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return &os.PathError{Op: "chmod", Path: name, Err: os.ErrInvalid}
	}

	perm = fopts.fileMode.Perm()

	// set Windows' ACLs
	err = acl.Chmod(name, perm)
	if err != nil {
		return &os.PathError{Op: "chmod", Path: name, Err: fmt.Errorf("%w (acl)", err)}
	}

	if fopts.readOnlyMode == ReadOnlyModeIgnore {
		return nil
	}

	fi, err := os.Stat(name)
	if err != nil {
		return &os.PathError{Op: "chmod", Path: name, Err: fmt.Errorf("%w (stat)", err)}
	}

	// Set or clear Windows' read-only attribute
	want := perm&windows.S_IWRITE != 0 // 0x80 (0o200)
	got := fi.Mode().Perm()&windows.S_IWRITE != 0
	if fopts.readOnlyMode == ReadOnlyModeReset {
		if !got {
			return nil
		}
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
		return &os.PathError{Op: "chmod", Path: name, Err: fmt.Errorf("%w (chmod)", err)}
	}

	return nil
}

func create(name string, opts ...Option) (*os.File, error) {
	fopts := Options{
		fileMode: CreatePerm,
		flags:    os.O_CREATE | os.O_TRUNC,
	}

	for _, opt := range opts {
		opt(&fopts)
	}

	if fopts.flags&os.O_WRONLY != os.O_WRONLY {
		fopts.flags |= os.O_RDWR
	}

	if fopts.readOnlyMode != ReadOnlyModeSet {
		fopts.flags |= O_FILE_FLAG_NO_RO_ATTR
	}

	sa, err := saFromPerm(fopts.fileMode, true)
	if err != nil {
		return nil, &os.PathError{Op: "create", Path: name, Err: err}
	}

	return golang.OpenFileNolog(name, fopts.flags, fopts.fileMode, sa)
}

func createTemp(dir, pattern string, perm os.FileMode, flag int) (*os.File, error) {
	if perm == 0 {
		perm = CreateTempPerm // 0o600
	}
	sa, err := saFromPerm(perm, true)
	if err != nil {
		return nil, &os.PathError{Op: "createtemp", Path: dir, Err: err}
	}

	f, err := golang.CreateTemp(dir, pattern, flag, perm, sa)
	if err != nil {
		return nil, &os.PathError{Op: "createtemp", Path: dir, Err: err}
	}

	return f, nil
}

func fchmod(f *os.File, mode os.FileMode, opts ...Option) error {
	if f == nil {
		return &os.PathError{Op: "chmod", Path: "", Err: os.ErrInvalid}
	}
	path, err := golang.Filepath(f)
	if err != nil {
		return &os.PathError{Op: "chmod", Path: f.Name(), Err: err}
	}

	err = chmod(path, mode, opts...)
	if err != nil {
		return &os.PathError{Op: "chmod", Path: path, Err: err}
	}

	return nil
}

func mkdir(name string, perm os.FileMode) error {
	sa, err := saFromPerm(perm, true)
	if err != nil {
		return &os.PathError{Op: "mkdir", Path: name, Err: err}
	}

	err = golang.Mkdir(name, MkdirTempPerm, sa)
	if err != nil {
		return &os.PathError{Op: "mkdir", Path: name, Err: err}
	}

	return nil
}

func mkdirAll(name string, perm os.FileMode) error {
	sa, err := saFromPerm(perm, true)
	if err != nil {
		return &os.PathError{Op: "mkdirall", Path: name, Err: err}
	}

	err = golang.MkdirAll(name, MkdirTempPerm, sa)
	if err != nil {
		return &os.PathError{Op: "mkdirall", Path: name, Err: err}
	}

	return nil
}

func mkdirTemp(dir, pattern string, opts ...Option) (string, error) {
	fopts := Options{
		fileMode: MkdirTempPerm,
	}

	for _, opt := range opts {
		opt(&fopts)
	}

	sa, err := saFromPerm(fopts.fileMode, true)
	if err != nil {
		prefix, suffix, _ := golang.PrefixAndSuffix(pattern)

		return "", &os.PathError{Op: "mkdirtemp", Path: dir + string(os.PathSeparator) + prefix + "*" + suffix, Err: err}
	}

	path, err := golang.MkdirTemp(dir, pattern, sa)
	if err != nil {
		prefix, suffix, _ := golang.PrefixAndSuffix(pattern)

		return "", &os.PathError{Op: "mkdirtemp", Path: dir + string(os.PathSeparator) + prefix + "*" + suffix, Err: err}
	}

	return path, nil
}

func openFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	sa, err := saFromPerm(perm, (flag&os.O_CREATE) == os.O_CREATE)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: name, Err: err}
	}

	return golang.OpenFileNolog(name, flag, perm, sa)
}

func remove(name string) error {
	return golang.Remove(name)
}

func removeAll(path string, opts ...Option) error {
	fopts := Options{}
	for _, opt := range opts {
		opt(&fopts)
	}

	if fopts.retrySeconds <= 0 {
		return golang.RemoveAll(path)
	}

	return robustio.Retry(func() (err error, mayRetry bool) {
		err = golang.RemoveAll(path)
		return err, robustio.IsEphemeralError(err)
	}, fopts.retrySeconds)
}

func symlink(oldname, newname string, opts ...Option) error {
	fopts := Options{}
	for _, opt := range opts {
		opt(&fopts)
	}

	err := os.Symlink(oldname, newname)
	if err != nil {
		return err
	}

	if fopts.setSymlinkOwner {
		err = setOwnerToCurrentUser(newname)
		if err != nil {
			return &os.LinkError{Op: "symlink", Old: oldname, New: newname, Err: err}
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
			&newBufSize)
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
			&newBufSize)
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
	err = sd.SetControl(windows.SE_DACL_PROTECTED,
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
		userSID, nil, nil, nil)
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
	if le := windows.GetLastError(); errors.Is(le, windows.ERROR_NOT_ALL_ASSIGNED) {
		return fmt.Errorf("failed to hold privilege: %w", le)
	}

	return nil
}
