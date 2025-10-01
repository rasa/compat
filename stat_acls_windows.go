// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

import (
	"errors"
	"strings"
	"sync"
	"testing"

	"golang.org/x/sys/windows"
)

var initialBufSize uint32 = 1024

func init() {
	if testing.Testing() {
		initialBufSize = 0
	}
}

type volumeACLCache struct {
	mu         sync.RWMutex
	byRoot     map[string]uint32 // root path "C:\", "\\server\share\", or "\\?\Volume{GUID}\"
	bySerial   map[uint32]uint32 // volume serial -> fsFlags
	guidToRoot map[string]string // "\\?\Volume{GUID}\" -> preferred root (drive letter if present)
}

var volCache = volumeACLCache{
	byRoot:     make(map[string]uint32),
	bySerial:   make(map[uint32]uint32),
	guidToRoot: make(map[string]string),
}

// supportsACLs checks the volume of 'path' and returns whether it supports
// persistent ACLs (e.g., NTFS/ReFS).
func supportsACLs(path string) (bool, error) {
	h, err := openForQuery(path)
	if err != nil {
		return false, err
	}
	defer func() { _ = windows.CloseHandle(h) }()
	return supportsACLsHandle(h)
}

func supportsACLsCached(fi FileInfo) (bool, error) {
	if fi == nil {
		return false, errors.New("fi is nil")
	}
	if fi.PartitionID() != 0 {
		volumeSerialNumber := uint32(fi.PartitionID()) //nolint:gosec,govet
		volCache.mu.RLock()
		fsFlags, ok := volCache.bySerial[volumeSerialNumber]
		volCache.mu.RUnlock()
		if ok {
			return flagsIndicatePersistentACLs(fsFlags), nil
		}
	}

	return false, errors.New("not cached")
}

// supportsACLsHandle checks the volume for an already-open file/dir handle.
func supportsACLsHandle(h windows.Handle) (bool, error) {
	// Cheap cache by volume serial first.
	var info windows.ByHandleFileInformation
	err := windows.GetFileInformationByHandle(h, &info)
	if err == nil {
		volCache.mu.RLock()
		fsFlags, ok := volCache.bySerial[info.VolumeSerialNumber]
		volCache.mu.RUnlock()
		if ok {
			return flagsIndicatePersistentACLs(fsFlags), nil
		}
	}

	// Ask Windows for serial + filesystem flags via handle.
	serial, fsFlags, err := getVolumeInfoByHandle(h)
	if err != nil {
		return false, err
	}

	// Cache by serial immediately.
	volCache.mu.Lock()
	volCache.bySerial[serial] = fsFlags
	volCache.mu.Unlock()

	// Also cache by a canonical root if we can resolve one.
	guidRoot, root, err := resolveCanonicalRootFromHandle(h)
	if err == nil {
		volCache.mu.Lock()
		if guidRoot != "" {
			guidRoot = strings.ToUpper(guidRoot)
			volCache.guidToRoot[guidRoot] = root
		}
		root = strings.ToUpper(root)
		volCache.byRoot[root] = fsFlags
		volCache.mu.Unlock()
	}

	return flagsIndicatePersistentACLs(fsFlags), nil
}

// openForQuery opens a file or directory with minimal rights for attribute/volume queries.
func openForQuery(path string) (windows.Handle, error) {
	p16, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}
	return windows.CreateFile(
		p16,
		windows.FILE_READ_ATTRIBUTES,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_FLAG_BACKUP_SEMANTICS,
		0,
	)
}

//------------------------------------------------------------------------------
// Windows API calls
//------------------------------------------------------------------------------

// getFinalPathNameByHandleGUID returns a path that starts with a volume GUID root
// ("\\?\Volume{GUID}\...") for local volumes, or "\\?\UNC\server\share\..." for UNC.
func getFinalPathNameByHandleGUID(h windows.Handle) (string, error) {
	// See https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-getfinalpathnamebyhandlea#parameters
	const FILE_NAME_NORMALIZED = 0x0
	const VOLUME_NAME_GUID = 0x1

	bufSize := initialBufSize
	for {
		buf := make([]uint16, bufSize)

		var n uint32
		var err error
		if bufSize == 0 {
			n, err = windows.GetFinalPathNameByHandle(h, nil, 0, FILE_NAME_NORMALIZED|VOLUME_NAME_GUID)
		} else {
			n, err = windows.GetFinalPathNameByHandle(h, &buf[0], uint32(len(buf)), FILE_NAME_NORMALIZED|VOLUME_NAME_GUID) //nolint:gosec
		}
		if n == 0 {
			if err != nil {
				return "", err
			}
			return "", errors.New("unexpected length with GetFinalPathNameByHandle()")
		}

		if n < uint32(len(buf)) { //nolint:gosec
			// success, truncate to returned length
			return windows.UTF16ToString(buf[:n]), nil
		}

		bufSize = n + 1
	}
}

// getVolumePathNamesForVolumeName fetches all mount points for a volume GUID.
func getVolumePathNamesForVolumeName(volGUID string) ([]string, error) {
	g16, err := windows.UTF16PtrFromString(volGUID)
	if err != nil {
		return nil, err
	}

	bufSize := initialBufSize
	for {
		var newBufSize uint32
		buf := make([]uint16, bufSize)
		if bufSize == 0 {
			err = windows.GetVolumePathNamesForVolumeName(g16, nil, bufSize, &newBufSize)
		} else {
			err = windows.GetVolumePathNamesForVolumeName(g16, &buf[0], bufSize, &newBufSize)
		}
		if err == nil {
			return multiSZToStrings(buf), nil
		}
		errno, ok := err.(windows.Errno) //nolint:errorlint
		if !ok || errno != windows.ERROR_MORE_DATA {
			return nil, err
		}
		if newBufSize > bufSize {
			bufSize = newBufSize
		} else {
			bufSize *= 2
		}
	}
}

// getVolumeInfoByHandle returns (serial, fsFlags) using GetVolumeInformationByHandle.
func getVolumeInfoByHandle(h windows.Handle) (uint32, uint32, error) {
	var (
		serial uint32
		flags  uint32
	)
	// We don't need names/lengths, pass nil/0 for those.
	err := windows.GetVolumeInformationByHandle(h, nil, 0, &serial, nil, &flags, nil, 0)
	if err != nil {
		return 0, 0, err
	}
	return serial, flags, nil
}

// (Optional) by-root variant if you want to seed/verify the byRoot cache directly.
/*
func getVolumeInfoByRoot(root string) (uint32, uint32, error) { //nolint:unused
	r16, err := windows.UTF16PtrFromString(root)
	if err != nil {
		return 0, 0, err
	}
	var (
		serial uint32
		flags  uint32
	)
	err = windows.GetVolumeInformation(r16, nil, 0, &serial, nil, &flags, nil, 0)
	if err != nil {
		return 0, 0, err
	}
	return serial, flags, nil
}
*/

//------------------------------------------------------------------------------
// Root resolution via GUID + GetVolumePathNamesForVolumeName
//------------------------------------------------------------------------------

func resolveCanonicalRootFromHandle(h windows.Handle) (guidRoot string, root string, err error) { //nolint:gocyclo
	full, err := getFinalPathNameByHandleGUID(h)
	if err != nil || full == "" {
		return "", "", err
	}

	// UNC: "\\?\UNC\server\share\..."
	if strings.HasPrefix(full, `\\?\UNC\`) {
		parts := strings.Split(full[len(`\\?\UNC\`):], `\`)
		if len(parts) >= 2 { //nolint:mnd
			root = `\\` + parts[0] + `\` + parts[1] + `\`
			return "", root, nil
		}
		return "", "", errors.New("unexpected UNC format from GetFinalPathNameByHandle()")
	}

	// Local volume GUID: "\\?\Volume{GUID}\..."
	if strings.HasPrefix(full, `\\?\Volume{`) {
		i := strings.Index(full, `}\`)
		if i <= 0 {
			return "", "", errors.New("unexpected GUID path from GetFinalPathNameByHandle()")
		}
		guidRoot = full[:i+2] // include trailing backslash

		guidRootUpper := strings.ToUpper(guidRoot)
		// Cache hit?
		volCache.mu.RLock()
		cached, ok := volCache.guidToRoot[guidRootUpper]
		volCache.mu.RUnlock()
		if ok {
			return guidRoot, cached, nil
		}

		mounts, err := getVolumePathNamesForVolumeName(guidRoot)
		if err != nil {
			return "", "", err
		}
		// Choose canonical root: prefer drive letter, else first mount, else the GUID itself.
		chosen := ""
		for _, m := range mounts {
			m = normalizeRoot(m)
			if isDriveLetterRoot(m) {
				chosen = m
				break
			}
			if chosen == "" && m != "" {
				chosen = m
			}
		}
		if chosen == "" {
			chosen = guidRoot
		}

		volCache.mu.Lock()
		volCache.guidToRoot[guidRootUpper] = chosen
		volCache.mu.Unlock()

		return guidRoot, chosen, nil
	}

	// Fallback to drive-root "C:\" style if present.
	if len(full) >= 3 && full[1] == ':' && (full[2] == '\\' || full[2] == '/') {
		return "", strings.ToUpper(full[:3]), nil
	}
	return "", "", errors.New("could not resolve canonical root")
}

func multiSZToStrings(buf []uint16) []string {
	var out []string
	start := 0
	for i, v := range buf {
		if v == 0 {
			if i == start {
				break // double-NUL terminator
			}
			out = append(out, windows.UTF16ToString(buf[start:i]))
			start = i + 1
		}
	}
	return out
}

func isDriveLetterRoot(s string) bool {
	return len(s) == 3 && s[1] == ':' && (s[2] == '\\' || s[2] == '/')
}

func normalizeRoot(root string) string {
	if root == "" {
		return root
	}
	root = strings.ReplaceAll(root, "/", `\`)
	if len(root) >= 2 && root[1] == ':' {
		root = strings.ToUpper(root[:1]) + root[1:]
	}
	if !strings.HasSuffix(root, `\`) {
		root += `\`
	}
	return root
}

func flagsIndicatePersistentACLs(fsFlags uint32) bool {
	return fsFlags&windows.FILE_PERSISTENT_ACLS == windows.FILE_PERSISTENT_ACLS
}
