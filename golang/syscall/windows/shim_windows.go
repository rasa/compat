// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package windows

/*
golang\syscall\windows\zsyscall_windows.go:553:26: undefined: _OSVERSIONINFOEXW
golang\syscall\windows\zsyscall_windows.go:570:63: undefined: PROCESS_MEMORY_COUNTERS
*/

// Snippet: https://github.com/golang/go/blob/8ad27fb6/src/internal/runtime/syscall/windows/defs_windows.go#L139-L149

// https://learn.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-memory_basic_information
type MemoryBasicInformation struct {
	BaseAddress       uintptr
	AllocationBase    uintptr
	AllocationProtect uint32
	PartitionId       uint16
	RegionSize        uintptr
	State             uint32
	Protect           uint32
	Type              uint32
}

// windows.GetProfilesDirectory
// windows.GetSidIdentifierAuthority
// windows.GetSidSubAuthority
// windows.GetSidSubAuthorityCount
// windows.IsValidSid

// const (
// 	SECURITY_LOCAL_SERVICE_RID = 19
// 	SID_REVISION = 1
// )
