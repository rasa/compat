// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package debug

import (
	"flag"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
	"testing"

	"golang.org/x/sys/windows"

	compat_consts "github.com/rasa/compat/consts"
	"github.com/rasa/compat/debug/consts"
)

func init() {
	testing.Init()
	flag.Parse()
}

type FileFlags int

func (f FileFlags) String() string {
	type Flag struct {
		k int
		v string
	}

	flags := []Flag{
		// {consts.O_RDONLY: "O_RDONLY"}, // 0x00000
		{consts.O_WRONLY, "O_WRONLY"},             // 0x00001
		{consts.O_RDWR, "O_RDWR"},                 // 0x00002
		{consts.O_CREAT, "O_CREAT"},               // 0x00040
		{consts.O_EXCL, "O_EXCL"},                 // 0x00080
		{consts.O_NOCTTY, "O_NOCTTY"},             // 0x00100
		{consts.O_TRUNC, "O_TRUNC"},               // 0x00200
		{consts.O_NONBLOCK, "O_NONBLOCK"},         // 0x00800
		{consts.O_APPEND, "O_APPEND"},             // 0x00400
		{consts.O_SYNC, "O_SYNC"},                 // 0x01000
		{consts.O_ASYNC, "O_ASYNC"},               // 0x02000
		{consts.O_DIRECTORY, "o_DIRECTORY"},       // 0x04000
		{consts.O_CLOEXEC, "O_CLOEXEC"},           // 0x80000
		{consts.O_NOFOLLOW_ANY, "o_NOFOLLOW_ANY"}, // 0x200000000
		{consts.O_OPEN_REPARSE, "o_OPEN_REPARSE"}, // 0x400000000
		{consts.O_WRITE_ATTRS, "o_WRITE_ATTRS"},   // 0x800000000

		{consts.O_FILE_FLAG_DELETE_ON_CLOSE, "O_FILE_FLAG_DELETE_ON_CLOSE"}, // 0x04000000
		{compat_consts.O_FILE_FLAG_NO_RO_ATTR, "O_FILE_FLAG_NO_RO_ATTR"},    // 0x00010000
	}

	s := ""
	rdonly := true
	for _, flag := range flags {
		k := flag.k
		v := flag.v
		if int(f)&k == k {
			if k == os.O_WRONLY || k == os.O_RDWR {
				rdonly = false
			}
			if s != "" {
				s += ","
			}
			s += v
			f &^= FileFlags(k)
		}
	}
	if rdonly {
		s = "O_RDONLY," + s
	}
	if f != 0 {
		s += fmt.Sprintf(",%x", int(f))
	}
	return s
}

type FileAttrs int

func (f FileAttrs) String() string {
	type Flag struct {
		k int
		v string
	}

	flags := []Flag{
		{consts.FILE_ATTRIBUTE_ARCHIVE, "FILE_ATTRIBUTE_ARCHIVE"},             // 0x00000020 (32)
		{consts.FILE_ATTRIBUTE_ENCRYPTED, "FILE_ATTRIBUTE_ENCRYPTED"},         // 0x00004000 (16384)
		{consts.FILE_ATTRIBUTE_HIDDEN, "FILE_ATTRIBUTE_HIDDEN"},               // 0x00000000 (2)
		{consts.FILE_ATTRIBUTE_NORMAL, "FILE_ATTRIBUTE_NORMAL"},               // 0x00000080 (128)
		{consts.FILE_ATTRIBUTE_OFFLINE, "FILE_ATTRIBUTE_OFFLINE"},             // 0x00001000 (4096)
		{consts.FILE_ATTRIBUTE_READONLY, "FILE_ATTRIBUTE_READONLY"},           // 0x00000001 (1)
		{consts.FILE_ATTRIBUTE_SYSTEM, "FILE_ATTRIBUTE_SYSTEM"},               // 0x00000004 (4)
		{consts.FILE_ATTRIBUTE_TEMPORARY, "FILE_ATTRIBUTE_TEMPORARY"},         // 0x00000100 (256)
		{consts.FILE_FLAG_BACKUP_SEMANTICS, "FILE_FLAG_BACKUP_SEMANTICS"},     // 0x02000000
		{consts.FILE_FLAG_DELETE_ON_CLOSE, "FILE_FLAG_DELETE_ON_CLOSE"},       // 0x04000000
		{consts.FILE_FLAG_NO_BUFFERING, "FILE_FLAG_NO_BUFFERING"},             // 0x20000000
		{consts.FILE_FLAG_OPEN_NO_RECALL, "FILE_FLAG_OPEN_NO_RECALL"},         // 0x00100000
		{consts.FILE_FLAG_OPEN_REPARSE_POINT, "FILE_FLAG_OPEN_REPARSE_POINT"}, // 0x00200000
		{consts.FILE_FLAG_OVERLAPPED, "FILE_FLAG_OVERLAPPED"},                 // 0x40000000
		{consts.FILE_FLAG_POSIX_SEMANTICS, "FILE_FLAG_POSIX_SEMANTICS"},       // 0x01000000
		{consts.FILE_FLAG_RANDOM_ACCESS, "FILE_FLAG_RANDOM_ACCESS"},           // 0x10000000
		{consts.FILE_FLAG_SESSION_AWARE, "FILE_FLAG_SESSION_AWARE"},           // 0x00800000
		{consts.FILE_FLAG_SEQUENTIAL_SCAN, "FILE_FLAG_SEQUENTIAL_SCAN"},       // 0x08000000
		{consts.FILE_FLAG_WRITE_THROUGH, "FILE_FLAG_WRITE_THROUGH"},           // 0x80000000
	}

	s := ""
	for _, flag := range flags {
		k := flag.k
		v := flag.v
		if int(f)&k == k {
			if s != "" {
				s += ","
			}
			s += v
			f &^= FileAttrs(k)
		}
	}
	if f != 0 {
		s += fmt.Sprintf(",%x", int(f))
	}
	return s
}

const (
	FILE_ADD_FILE             = 2
	FILE_ADD_SUBDIRECTORY     = 4
	FILE_ALL_ACCESS           = 0x001F01FF
	FILE_APPEND_DATA          = 4
	FILE_CREATE_PIPE_INSTANCE = 4 // dup of FILE_APPEND_DATA
	FILE_DELETE_CHILD         = 0x40

	// windows.STANDARD_RIGHTS_REQUIRED = 0x000F0000
	// windows.STANDARD_RIGHTS_READ     = READ_CONTROL
	// windows.STANDARD_RIGHTS_WRITE    = READ_CONTROL
	// windows.STANDARD_RIGHTS_EXECUTE  = READ_CONTROL
	// windows.STANDARD_RIGHTS_ALL      = 0x001F0000
	// windows.SPECIFIC_RIGHTS_ALL      = 0x0000FFFF.
)

/*
	DE - delete
	RC - read control
	WDAC - write DAC
	WO - write owner
	S - synchronize
	AS - access system security
	MA - maximum allowed
	GR - generic read
	GW - generic write
	GE - generic execute
	GA - generic all
	RD - read data/list directory
 	WD - write data/add file
	AD - append data/add subdirectory
	REA - read extended attributes
	WEA - write extended attributes
 	X - execute/traverse
 	DC - delete child
 	RA - read attributes
 	WA - write attributes
*/

var icalcsMap = map[uint32]string{ //nolint:unused
	// https://learn.microsoft.com/en-us/windows/win32/secauthz/access-mask?source=recommendations
	windows.DELETE:       "DE",   // 0x00010000
	windows.READ_CONTROL: "RC",   // 0x00020000
	windows.WRITE_DAC:    "WDAC", // 0x00040000
	windows.WRITE_OWNER:  "WO",   // 0x00080000
	windows.SYNCHRONIZE:  "S",    // 0x00100000

	// https://learn.microsoft.com/en-us/windows/win32/SecAuthZ/generic-access-rights
	windows.GENERIC_READ:    "GR", // 0x80000000
	windows.GENERIC_WRITE:   "GW", // 0x40000000
	windows.GENERIC_EXECUTE: "GE", // 0x20000000
	windows.GENERIC_ALL:     "GA", // 0x10000000

	windows.ACCESS_SYSTEM_SECURITY: "AS", // 0x01000000
	windows.MAXIMUM_ALLOWED:        "MA", // 0x02000000

	// FILE_ADD_FILE: "WD", // 2 dup of FILE_WRITE_DATA
	// FILE_ADD_SUBDIRECTORY: "AD", // 4 (dup)
	// FILE_ALL_ACCESS:  "FILE_ALL_ACCESS", // 0x001F01FF
	FILE_APPEND_DATA: "AD", // 4
	// FILE_CREATE_PIPE_INSTANCE: "FILE_CREATE_PIPE_INSTANCE", // 4 (dup)
	FILE_DELETE_CHILD:            "DC", // 64 (0x40)
	windows.FILE_EXECUTE:         "X",  // 32 (0x20)
	windows.FILE_LIST_DIRECTORY:  "RD", // 1
	windows.FILE_READ_ATTRIBUTES: "RA", // 128 (0x80)
	// FILE_READ_DATA: "X", // 1 // dup of FILE_LIST_DIRECTORY
	windows.FILE_READ_EA: "REA", // 8
	// FILE_TRAVERSE: "FILE_TRAVERSE", // 32 (0x20) // dup of FILE_EXECUTE
	windows.FILE_WRITE_ATTRIBUTES: "WA",  // 256 (0x100)
	windows.FILE_WRITE_DATA:       "WD",  // 2
	windows.FILE_WRITE_EA:         "WEA", // 16 (0x10)
}

// https://github.com/golang/sys/blob/3d9a6b80/windows/security_windows.go#L992
var maskMap = map[uint32]string{ //nolint:unused
	// https://learn.microsoft.com/en-us/windows/win32/secauthz/access-mask?source=recommendations
	windows.DELETE:       "DELETE",       // 0x00010000
	windows.READ_CONTROL: "READ_CONTROL", // 0x00020000
	windows.WRITE_DAC:    "WRITE_DAC",    // 0x00040000
	windows.WRITE_OWNER:  "WRITE_OWNER",  // 0x00080000
	windows.SYNCHRONIZE:  "SYNCHRONIZE",  // 0x00100000

	// https://learn.microsoft.com/en-us/windows/win32/SecAuthZ/generic-access-rights
	windows.GENERIC_READ:    "GENERIC_READ",    // 0x80000000
	windows.GENERIC_WRITE:   "GENERIC_WRITE",   // 0x40000000
	windows.GENERIC_EXECUTE: "GENERIC_EXECUTE", // 0x20000000
	windows.GENERIC_ALL:     "GENERIC_ALL",     // 0x10000000

	windows.ACCESS_SYSTEM_SECURITY: "ACCESS_SYSTEM_SECURITY", // 0x01000000
	windows.MAXIMUM_ALLOWED:        "MAXIMUM_ALLOWED",        // 0x02000000

	// https://learn.microsoft.com/en-us/windows/win32/fileio/file-access-rights-constants

	// FILE_ADD_FILE: "FILE_ADD_FILE", // 2 (dup of FILE_WRITE_DATA)
	// FILE_ADD_SUBDIRECTORY: "FILE_ADD_SUBDIRECTORY", // 4 (dup)
	// FILE_ALL_ACCESS:  "FILE_ALL_ACCESS",  // 0x001F01FF
	FILE_APPEND_DATA: "FILE_APPEND_DATA", // 4 (dup)
	// FILE_CREATE_PIPE_INSTANCE: "FILE_CREATE_PIPE_INSTANCE", // 4 (dup)
	FILE_DELETE_CHILD:            "FILE_DELETE_CHILD",    // 64 (0x40)
	windows.FILE_EXECUTE:         "FILE_EXECUTE",         // 32 (0x20)
	windows.FILE_LIST_DIRECTORY:  "FILE_LIST_DIRECTORY",  // 1
	windows.FILE_READ_ATTRIBUTES: "FILE_READ_ATTRIBUTES", // 128 (0x80)
	// FILE_READ_DATA: "FILE_READ_DATA", // 1 // dup of FILE_LIST_DIRECTORY
	windows.FILE_READ_EA: "FILE_READ_EA", // 8
	// FILE_TRAVERSE: "FILE_TRAVERSE", // 32 (0x20) // dup of FILE_EXECUTE
	windows.FILE_WRITE_ATTRIBUTES: "FILE_WRITE_ATTRIBUTES", // 256 (0x100)
	windows.FILE_WRITE_DATA:       "FILE_WRITE_DATA",       // 2
	windows.FILE_WRITE_EA:         "FILE_WRITE_EA",         // 16 (0x10)
}

type AccessMask uint32 //nolint:unused

func (a AccessMask) String() string { //nolint:unused
	mask := uint32(a)
	rv := ""
	rights := map[string]uint32{}
	for k, v := range maskMap {
		if mask&k == k {
			rights[v] = k
			mask &^= k
		}
	}
	if len(rights) == 0 {
		return "NO ACCESS"
	}
	keys := slices.Collect(maps.Keys(rights))
	slices.Sort(keys)
	rv += strings.Join(keys, ",")

	if mask != 0 {
		rv += "," + fmt.Sprintf("0x%x", mask)
	}

	return rv
}

func DumpMasks(perm os.FileMode, ownerMask uint32, groupMask uint32, worldMask uint32) { //nolint:unused
	if !strings.Contains(os.Getenv("COMPAT_DEBUG"), "DUMP") {
		return
	}
	omask := AccessMask(ownerMask)
	gmask := AccessMask(groupMask)
	wmask := AccessMask(worldMask)

	fmt.Printf("perm=%04o ownerMask=%v groupMask=%v worldMask=%v\n", perm, omask, gmask, wmask)
}

// https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-createfilea

var createModeMap = map[uint32]string{ //nolint:unused
	windows.CREATE_ALWAYS:     "CREATE_ALWAYS",     // 2
	windows.CREATE_NEW:        "CREATE_NEW",        // 1
	windows.OPEN_ALWAYS:       "OPEN_ALWAYS",       // 4
	windows.OPEN_EXISTING:     "OPEN_EXISTING",     // 3
	windows.TRUNCATE_EXISTING: "TRUNCATE_EXISTING", // 5
}

type CreateMode uint32 //nolint:unused

func (c CreateMode) String() string { //nolint:unused
	u := uint32(c)
	v, ok := createModeMap[u]
	if ok {
		return v
	}

	return fmt.Sprintf("Unknown create mode 0x%x", u)
}

var shareModeMap = map[uint32]string{ //nolint:unused
	// windows.FILE_SHARE_NONE = 0
	windows.FILE_SHARE_READ:   "FILE_SHARE_READ",   // 0x00000001
	windows.FILE_SHARE_WRITE:  "FILE_SHARE_WRITE",  // 0x00000002
	windows.FILE_SHARE_DELETE: "FILE_SHARE_DELETE", // 0x00000004
}

type ShareMode uint32 //nolint:unused

func (s ShareMode) String() string { //nolint:unused
	mask := uint32(s)
	if mask == 0 {
		return "FILE_SHARE_NONE"
	}

	rv := ""
	rights := map[string]uint32{}
	for k, v := range shareModeMap {
		if mask&k == k {
			rights[v] = k
			mask &^= k
		}
	}
	keys := slices.Collect(maps.Keys(rights))
	slices.Sort(keys)
	rv += strings.Join(keys, ",")

	if mask != 0 {
		rv += "," + fmt.Sprintf("0x%x", mask)
	}

	return rv
}
