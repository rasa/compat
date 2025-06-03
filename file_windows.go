// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"syscall"
	"unsafe"

	"github.com/capnspacehook/go-acl"

	"golang.org/x/sys/windows"
)

type securityInfo struct {
	ownerSid *windows.SID
	groupSid *windows.SID
	acl      *windows.ACL
	perm     os.FileMode
}

func chmod(name string, perm os.FileMode) error {
	// set/reset syscall.FILE_ATTRIBUTE_READONLY/syscall.FILE_ATTRIBUTE_NORMAL
	err := os.Chmod(name, perm)
	if err != nil {
		return err
	}

	// set Windows' ACLs
	return acl.Chmod(name, perm)
}

func create(name string, perm os.FileMode, flag int) (*os.File, error) {
	flag |= O_CREATE
	sa, err := securityAttributes(perm, true)
	if err != nil {
		return nil, err
	}

	return openFileNolog(name, flag, perm, sa)
}

func createTemp(dir, pattern string, flag int) (*os.File, error) {
	sa, err := securityAttributes(CreateTempPerm, true)
	if err != nil {
		return nil, err
	}

	return _createTemp(dir, pattern, flag, sa)
}

func mkdir(name string, perm os.FileMode) error {
	si, err := _securityInfo(perm)
	if err != nil {
		return err
	}

	return _mkdir(name, si)
}

func mkdirAll(name string, perm os.FileMode) error {
	si, err := _securityInfo(perm)
	if err != nil {
		return err
	}

	return _mkdirAll(name, si)
}

func mkdirTemp(dir, pattern string) (string, error) {
	perm := MkdirTempPerm

	si, err := _securityInfo(perm)
	if err != nil {
		return "", err
	}

	return _mkdirTemp(dir, pattern, si)
}

func openFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	sa, err := securityAttributes(perm, flag|O_CREATE == O_CREATE)
	if err != nil {
		return nil, err
	}

	return openFileNolog(name, flag, perm, sa)
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

func securityAttributes(perm os.FileMode, create bool) (*syscall.SecurityAttributes, error) {
	var sa syscall.SecurityAttributes
	sa.Length = uint32(unsafe.Sizeof(sa))

	if !create {
		return &sa, nil
	}

	perm &^= os.FileMode(GetUmask()) //nolint:gosec // quiet linter
	sd, _, err := securityDescriptor(perm)
	if err != nil {
		return nil, err
	}

	// Convert security descriptor to SECURITY_ATTRIBUTES
	sa = syscall.SecurityAttributes{
		Length:             uint32(unsafe.Sizeof(syscall.SecurityAttributes{})),
		SecurityDescriptor: uintptr(unsafe.Pointer(sd)), // Directly pass the security descriptor pointer
		InheritHandle:      1,                           // No handle inheritance
	}
	return &sa, nil
}

func _securityInfo(perm os.FileMode) (*securityInfo, error) {
	perm &^= os.FileMode(GetUmask()) //nolint:gosec // quiet linter
	_, si, err := securityDescriptor(perm)

	return si, err
}

func securityDescriptor(perm os.FileMode) (*windows.SECURITY_DESCRIPTOR, *securityInfo, error) {
	// Get current user's SID
	usr, err := user.Current()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot get current user: %w", err)
	}
	ownerSid, _, _, err := windows.LookupSID("", usr.Username)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot lookup SID for %s: %w", usr.Username, err)
	}

	token := windows.Token(0)
	err = windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get process token for %s: %w", usr.Username, err)
	}
	defer token.Close()

	var size uint32
	// First call to get required buffer size
	_ = windows.GetTokenInformation(token, windows.TokenPrimaryGroup, nil, 0, &size)

	buf := make([]byte, size)
	err = windows.GetTokenInformation(token, windows.TokenPrimaryGroup, &buf[0], size, &size)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get token information for %s: %w", usr.Username, err)
	}

	group := (*tokenPrimaryGroup)(unsafe.Pointer(&buf[0]))
	groupSid := group.PrimaryGroup

	worldSid, err := windows.CreateWellKnownSid(windows.WinWorldSid)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create world SID: %w", err)
	}

	var ea [3]windows.EXPLICIT_ACCESS

	ownerMask := accessMask(perm, 6) //nolint:mnd // quiet linter
	setExplicitAccess(&ea[0], ownerSid, ownerMask, windows.TRUSTEE_IS_USER)

	groupMask := accessMask(perm, 3) //nolint:mnd // quiet linter
	setExplicitAccess(&ea[1], groupSid, groupMask, windows.TRUSTEE_IS_GROUP)

	worldMask := accessMask(perm, 0)
	setExplicitAccess(&ea[2], worldSid, worldMask, windows.TRUSTEE_IS_WELL_KNOWN_GROUP)

	// dumpMasks(perm, ownerMask, groupMask, worldMask).

	acl, err := windows.ACLFromEntries(ea[:], nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ACLs: %w", err)
	}

	sd, err := windows.NewSecurityDescriptor()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create security descriptor: %w", err)
	}

	err = sd.SetDACL(acl, true, false)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set ACL in security descriptor: %w", err)
	}

	err = sd.SetOwner(ownerSid, false)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set ACL owner in security descriptor: %w", err)
	}

	err = sd.SetGroup(groupSid, false)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set ACL group in security descriptor: %w", err)
	}

	si := securityInfo{ownerSid, groupSid, acl, perm}

	return sd, &si, nil
}

func accessMask(mode os.FileMode, shift int) uint32 {
	perm := uint32(mode.Perm())

	var mask uint32

	if perm&(0o4<<shift) == (0o4 << shift) { //nolint:mnd // quiet linter
		mask |= windows.GENERIC_READ
	}
	if perm&(0o2<<shift) == (0o2 << shift) { //nolint:mnd // quiet linter
		mask |= windows.GENERIC_WRITE | windows.DELETE
	}
	if perm&(0o1<<shift) == (0o1 << shift) { //nolint:mnd // quiet linter
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

func setNamedSecurityInfo(name string, si *securityInfo) error {
	return windows.SetNamedSecurityInfo(
		name,
		windows.SE_FILE_OBJECT,
		windows.OWNER_SECURITY_INFORMATION|
			windows.GROUP_SECURITY_INFORMATION|
			windows.DACL_SECURITY_INFORMATION,
		si.ownerSid,
		si.groupSid,
		si.acl,
		nil,
	)
}

type tokenPrimaryGroup struct {
	PrimaryGroup *windows.SID
}

const supportsCreateWithStickyBit = true

// The following code is:
// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

///////////////////////////////////////////////////////////////////////////////
// Copied from https://github.com/golang/go/blob/e282cbb1/src/os/file_windows.go#L137
///////////////////////////////////////////////////////////////////////////////

// openFileNolog is the Windows implementation of OpenFile.
func openFileNolog(name string, flag int, perm os.FileMode, sa *syscall.SecurityAttributes) (*os.File, error) {
	if name == "" {
		return nil, &os.PathError{Op: "open", Path: name, Err: syscall.ENOENT}
	}
	path := fixLongPath(name)
	r, err := open(path, flag|syscall.O_CLOEXEC, syscallMode(perm), sa)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: name, Err: err}
	}
	// syscall.Open always returns a non-blocking handle.
	// newFile() is private, so call NewFile() instead.
	// return newFile(r, name, "file", false), nil
	return os.NewFile(uintptr(r), name), nil
}

///////////////////////////////////////////////////////////////////////////////
// Copied from https://github.com/golang/go/blob/e282cbb1/src/syscall/syscall_windows.go#L365
///////////////////////////////////////////////////////////////////////////////

func open(name string, flag int, perm uint32, sa *syscall.SecurityAttributes) (fd syscall.Handle, err error) { //nolint:funlen,gocyclo // quiet linter
	if len(name) == 0 {
		return syscall.InvalidHandle, syscall.ERROR_FILE_NOT_FOUND
	}
	namep, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return syscall.InvalidHandle, err
	}
	var access uint32
	switch flag & (syscall.O_RDONLY | syscall.O_WRONLY | syscall.O_RDWR) {
	case syscall.O_RDONLY:
		access = syscall.GENERIC_READ
	case syscall.O_WRONLY:
		access = syscall.GENERIC_WRITE
	case syscall.O_RDWR:
		access = syscall.GENERIC_READ | syscall.GENERIC_WRITE
	}
	if flag&syscall.O_CREAT != 0 {
		access |= syscall.GENERIC_WRITE
	}
	if flag&syscall.O_APPEND != 0 {
		// Remove GENERIC_WRITE unless O_TRUNC is set, in which case we need it to truncate the file.
		// We can't just remove FILE_WRITE_DATA because GENERIC_WRITE without FILE_WRITE_DATA
		// starts appending at the beginning of the file rather than at the end.
		if flag&syscall.O_TRUNC == 0 {
			access &^= syscall.GENERIC_WRITE
		}
		// Set all access rights granted by GENERIC_WRITE except for FILE_WRITE_DATA.
		access |= syscall.FILE_APPEND_DATA | syscall.FILE_WRITE_ATTRIBUTES | windows.FILE_WRITE_EA | syscall.STANDARD_RIGHTS_WRITE | syscall.SYNCHRONIZE
	}
	sharemode := uint32(syscall.FILE_SHARE_READ | syscall.FILE_SHARE_WRITE)
	// <compat change>
	// var sa *syscall.SecurityAttributes
	// if flag&syscall.O_CLOEXEC == 0 {
	// 	sa = makeInheritSa()
	// }
	// </compat change>
	// We don't use CREATE_ALWAYS, because when opening a file with
	// FILE_ATTRIBUTE_READONLY these will replace an existing file
	// with a new, read-only one. See https://go.dev/issue/38225.
	//
	// Instead, we ftruncate the file after opening when O_TRUNC is set.
	var createmode uint32
	switch {
	case flag&(syscall.O_CREAT|syscall.O_EXCL) == (syscall.O_CREAT | syscall.O_EXCL):
		createmode = syscall.CREATE_NEW
	case flag&syscall.O_CREAT == syscall.O_CREAT:
		createmode = syscall.OPEN_ALWAYS
	default:
		createmode = syscall.OPEN_EXISTING
	}
	var attrs uint32 = syscall.FILE_ATTRIBUTE_NORMAL
	if perm&syscall.S_IWRITE == 0 {
		attrs = syscall.FILE_ATTRIBUTE_READONLY
	}
	if flag&syscall.O_WRONLY == 0 && flag&syscall.O_RDWR == 0 {
		// We might be opening or creating a directory.
		// CreateFile requires FILE_FLAG_BACKUP_SEMANTICS
		// to work with directories.
		attrs |= syscall.FILE_FLAG_BACKUP_SEMANTICS
	}
	if flag&syscall.O_SYNC != 0 {
		const _FILE_FLAG_WRITE_THROUGH = 0x80000000
		attrs |= _FILE_FLAG_WRITE_THROUGH
	}
	// <compat addition>
	if flag&O_DELETE == O_DELETE {
		attrs &^= uint32(windows.FILE_ATTRIBUTE_READONLY)
		attrs |= (windows.FILE_FLAG_DELETE_ON_CLOSE | windows.FILE_ATTRIBUTE_TEMPORARY)
		sharemode |= syscall.FILE_SHARE_DELETE
	}
	// </compat addition>
	h, err := syscall.CreateFile(namep, access, sharemode, sa, createmode, attrs, 0)
	if h == syscall.InvalidHandle {
		if err == syscall.ERROR_ACCESS_DENIED && (flag&syscall.O_WRONLY != 0 || flag&syscall.O_RDWR != 0) {
			// We should return EISDIR when we are trying to open a directory with write access.
			fa, e1 := syscall.GetFileAttributes(namep)
			if e1 == nil && fa&syscall.FILE_ATTRIBUTE_DIRECTORY != 0 {
				err = syscall.EISDIR
			}
		}
		return h, err
	}
	// Ignore O_TRUNC if the file has just been created.
	if flag&syscall.O_TRUNC == syscall.O_TRUNC &&
		(createmode == syscall.OPEN_EXISTING || (createmode == syscall.OPEN_ALWAYS && err == syscall.ERROR_ALREADY_EXISTS)) {
		err = syscall.Ftruncate(h, 0)
		if err != nil {
			syscall.CloseHandle(h) //nolint:errcheck // quiet linter
			return syscall.InvalidHandle, err
		}
	}
	return h, nil
}

///////////////////////////////////////////////////////////////////////////////
// Copied from https://github.com/golang/go/blob/e282cbb1/src/os/file_posix.go#L59C1-L73C2
///////////////////////////////////////////////////////////////////////////////

// syscallMode returns the syscall-specific mode bits from Go's portable mode bits.
func syscallMode(i os.FileMode) (o uint32) {
	o |= uint32(i.Perm())
	if i&os.ModeSetuid != 0 {
		o |= syscall.S_ISUID
	}
	if i&os.ModeSetgid != 0 {
		o |= syscall.S_ISGID
	}
	if i&os.ModeSticky != 0 {
		o |= syscall.S_ISVTX
	}
	// No mapping for Go's ModeTemporary (plan9 only).
	return
}

///////////////////////////////////////////////////////////////////////////////
// Copied from https://github.com/golang/go/blob/e282cbb1/src/os/tempfile.go#L14-L24
///////////////////////////////////////////////////////////////////////////////

// random number source provided by runtime.
// We generate random temporary file names so that there's a good
// chance the file doesn't exist yet - keeps the number of tries in
// TempFile to a minimum.
//
//go:linkname runtime_rand runtime.rand
func runtime_rand() uint64 //nolint:revive // quiet linter

func nextRandom() string {
	return uitoa(uint(uint32(runtime_rand()))) //nolint:gosec // quiet linter
}

///////////////////////////////////////////////////////////////////////////////
// Copied from https://github.com/golang/go/blob/e282cbb1/src/os/tempfile.go#L26-L123
///////////////////////////////////////////////////////////////////////////////

// CreateTemp creates a new temporary file in the directory dir,
// opens the file for reading and writing, and returns the resulting file.
// The filename is generated by taking pattern and adding a random string to the end.
// If pattern includes a "*", the random string replaces the last "*".
// The file is created with mode 0o600 (before umask).
// If dir is the empty string, CreateTemp uses the default directory for temporary files, as returned by [TempDir].
// Multiple programs or goroutines calling CreateTemp simultaneously will not choose the same file.
// The caller can use the file's Name method to find the pathname of the file.
// It is the caller's responsibility to remove the file when it is no longer needed.
func _createTemp(dir, pattern string, flag int, sa *syscall.SecurityAttributes) (*os.File, error) {
	if dir == "" {
		dir = os.TempDir()
	}

	prefix, suffix, err := prefixAndSuffix(pattern)
	if err != nil {
		return nil, &os.PathError{Op: "createtemp", Path: pattern, Err: err}
	}
	prefix = joinPath(dir, prefix)

	try := 0
	for {
		name := prefix + nextRandom() + suffix
		f, err := openFileNolog(name, O_RDWR|O_CREATE|O_EXCL|flag, CreateTempPerm, sa)
		if os.IsExist(err) {
			if try++; try < 10000 { //nolint:mnd // quiet linter
				continue
			}
			return nil, &os.PathError{Op: "createtemp", Path: prefix + "*" + suffix, Err: os.ErrExist}
		}
		return f, err
	}
}

var errPatternHasSeparator = errors.New("pattern contains path separator")

// prefixAndSuffix splits pattern by the last wildcard "*", if applicable,
// returning prefix as the part before "*" and suffix as the part after "*".
func prefixAndSuffix(pattern string) (prefix, suffix string, err error) {
	for i := 0; i < len(pattern); i++ {
		if os.IsPathSeparator(pattern[i]) {
			return "", "", errPatternHasSeparator
		}
	}
	if pos := lastIndexByteString(pattern, '*'); pos != -1 { // removed bytealg
		prefix, suffix = pattern[:pos], pattern[pos+1:]
	} else {
		prefix = pattern
	}
	return prefix, suffix, nil
}

// MkdirTemp creates a new temporary directory in the directory dir
// and returns the pathname of the new directory.
// The new directory's name is generated by adding a random string to the end of pattern.
// If pattern includes a "*", the random string replaces the last "*" instead.
// The directory is created with mode 0o700 (before umask).
// If dir is the empty string, MkdirTemp uses the default directory for temporary files, as returned by TempDir.
// Multiple programs or goroutines calling MkdirTemp simultaneously will not choose the same directory.
// It is the caller's responsibility to remove the directory when it is no longer needed.
func _mkdirTemp(dir, pattern string, si *securityInfo) (string, error) {
	if dir == "" {
		dir = os.TempDir()
	}

	prefix, suffix, err := prefixAndSuffix(pattern)
	if err != nil {
		return "", &os.PathError{Op: "mkdirtemp", Path: pattern, Err: err}
	}
	prefix = joinPath(dir, prefix)

	try := 0
	for {
		name := prefix + nextRandom() + suffix
		err := _mkdir(name, si)
		if err == nil {
			return name, nil
		}
		if os.IsExist(err) {
			if try++; try < 10000 { //nolint:mnd // quiet linter
				continue
			}
			return "", &os.PathError{Op: "mkdirtemp", Path: dir + string(os.PathSeparator) + prefix + "*" + suffix, Err: os.ErrExist}
		}
		if os.IsNotExist(err) {
			if _, err := os.Stat(dir); os.IsNotExist(err) { //nolint:govet // quiet linter
				return "", err
			}
		}
		return "", err
	}
}

func joinPath(dir, name string) string {
	if len(dir) > 0 && os.IsPathSeparator(dir[len(dir)-1]) {
		return dir + name
	}
	return dir + string(os.PathSeparator) + name
}

///////////////////////////////////////////////////////////////////////////////
// Copied from https://github.com/golang/go/blob/e282cbb1/src/internal/itoa/itoa.go#L17C1-L33C2
///////////////////////////////////////////////////////////////////////////////

// Uitoa converts val to a decimal string.
func uitoa(val uint) string {
	if val == 0 { // avoid string allocation
		return "0"
	}
	var buf [20]byte // big enough for 64bit value base 10
	i := len(buf) - 1
	for val >= 10 {
		q := val / 10 //nolint:mnd // quiet linter
		buf[i] = byte('0' + val - q*10)
		i--
		val = q
	}
	// val < 10
	buf[i] = byte('0' + val)
	return string(buf[i:])
}

///////////////////////////////////////////////////////////////////////////////
// Copied from https://github.com/golang/go/blob/e282cbb1/src/internal/bytealg/lastindexbyte_generic.go#L16-L23
///////////////////////////////////////////////////////////////////////////////

func lastIndexByteString(s string, c byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			return i
		}
	}
	return -1
}

///////////////////////////////////////////////////////////////////////////////
// Copied from https://github.com/golang/go/blob/e282cbb1/src/os/file.go#L324C1-L349C1
///////////////////////////////////////////////////////////////////////////////

// Mkdir creates a new directory with the specified name and permission
// bits (before umask).
// If there is an error, it will be of type [*PathError].
func _mkdir(name string, si *securityInfo) error {
	longName := fixLongPath(name)
	e := ignoringEINTR(func() error {
		name, err := syscall.UTF16PtrFromString(longName)
		if err != nil {
			return err
		}
		// return syscall.Mkdir(longName, syscallMode(perm))
		err = windows.CreateDirectory(name, nil)
		if err != nil {
			return err
		}
		err = setNamedSecurityInfo(longName, si)
		if err != nil {
			_ = os.Remove(longName)

			return err
		}

		return nil
	})

	if e != nil {
		return &os.PathError{Op: "mkdir", Path: name, Err: e}
	}

	// mkdir(2) itself won't handle the sticky bit on *BSD and Solaris
	if !supportsCreateWithStickyBit && si.perm&os.ModeSticky != 0 {
		e = setStickyBit(name)

		if e != nil {
			os.Remove(name)
			return e
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////
// Copied from https://github.com/golang/go/blob/e282cbb1/src/os/file_posix.go#L247C1-L262C1
///////////////////////////////////////////////////////////////////////////////

// ignoringEINTR makes a function call and repeats it if it returns an
// EINTR error. This appears to be required even though we install all
// signal handlers with SA_RESTART: see #22838, #38033, #38836, #40846.
// Also #20400 and #36644 are issues in which a signal handler is
// installed without setting SA_RESTART. None of these are the common case,
// but there are enough of them that it seems that we can't avoid
// an EINTR loop.
func ignoringEINTR(fn func() error) error {
	for {
		err := fn()
		if err != syscall.EINTR { //nolint:errorlint // quiet linter
			return err
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
// Copied from https://github.com/golang/go/blob/e282cbb1/src/os/file.go#L350C1-L357C2
///////////////////////////////////////////////////////////////////////////////

// setStickyBit adds ModeSticky to the permission bits of path, non atomic.
func setStickyBit(name string) error {
	fi, err := os.Stat(name)
	if err != nil {
		return err
	}
	return os.Chmod(name, fi.Mode()|os.ModeSticky)
}

///////////////////////////////////////////////////////////////////////////////
// Copied from https://github.com/golang/go/blob/e282cbb1/src/os/path.go#L12C1-L66C2
///////////////////////////////////////////////////////////////////////////////

// MkdirAll creates a directory named path,
// along with any necessary parents, and returns nil,
// or else returns an error.
// The permission bits perm (before umask) are used for all
// directories that MkdirAll creates.
// If path is already a directory, MkdirAll does nothing
// and returns nil.
func _mkdirAll(path string, si *securityInfo) error {
	// Fast path: if we can tell whether path is a directory or file, stop with success or error.
	dir, err := os.Stat(path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return &os.PathError{Op: "mkdir", Path: path, Err: syscall.ENOTDIR}
	}

	// Slow path: make sure parent exists and then call Mkdir for path.

	// Extract the parent folder from path by first removing any trailing
	// path separator and then scanning backward until finding a path
	// separator or reaching the beginning of the string.
	i := len(path) - 1
	for i >= 0 && os.IsPathSeparator(path[i]) {
		i--
	}
	for i >= 0 && !os.IsPathSeparator(path[i]) {
		i--
	}
	if i < 0 {
		i = 0
	}

	// If there is a parent directory, and it is not the volume name,
	// recurse to ensure parent directory exists.
	if parent := path[:i]; len(parent) > len(filepath.VolumeName(path)) {
		err = _mkdirAll(parent, si)
		if err != nil {
			return err
		}
	}

	// Parent now exists; invoke Mkdir and use its result.
	err = _mkdir(path, si)
	if err != nil {
		// Handle arguments like "foo/." by
		// double-checking that directory doesn't exist.
		dir, err1 := os.Lstat(path)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}
	return nil
}
