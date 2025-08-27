// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows && debug

package compat

import (
	"flag"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
	"testing"

	"golang.org/x/sys/windows"
)

func init() {
	testing.Init()
	flag.Parse()
}

func dumpMasks(perm os.FileMode, ownerMask uint32, groupMask uint32, worldMask uint32) { //nolint:unused // quiet linter
	if !strings.Contains(compatDebug, "DUMP") {
		return
	}
	omask := aMask(ownerMask)
	gmask := aMask(groupMask)
	wmask := aMask(worldMask)

	fmt.Printf("perm=%04o ownerMask=%v groupMask=%v worldMask=%v\n", perm, omask, gmask, wmask)
}

// https://github.com/golang/sys/blob/3d9a6b80/windows/security_windows.go#L992
var maskMap = map[uint32]string{ //nolint:unused // quiet linter
	windows.DELETE:                 "D",    // 0x00010000
	windows.READ_CONTROL:           "RC",   // 0x00020000
	windows.WRITE_DAC:              "WDAC", // 0x00040000
	windows.WRITE_OWNER:            "WO",   // 0x00080000
	windows.SYNCHRONIZE:            "S",    // 0x00100000
	windows.ACCESS_SYSTEM_SECURITY: "AS",   // 0x01000000
	windows.MAXIMUM_ALLOWED:        "MA",   // 0x02000000
	windows.GENERIC_READ:           "GR",   // 0x80000000
	windows.GENERIC_WRITE:          "GW",   // 0x40000000
	windows.GENERIC_EXECUTE:        "GE",   // 0x20000000
	windows.GENERIC_ALL:            "GA",   // 0x10000000
	// windows.STANDARD_RIGHTS_REQUIRED = 0x000F0000
	// windows.STANDARD_RIGHTS_READ     = READ_CONTROL
	// windows.STANDARD_RIGHTS_WRITE    = READ_CONTROL
	// windows.STANDARD_RIGHTS_EXECUTE  = READ_CONTROL
	// windows.STANDARD_RIGHTS_ALL      = 0x001F0000
	// windows.SPECIFIC_RIGHTS_ALL      = 0x0000FFFF
}

type aMask uint32 //nolint:unused // quiet linter

func (a aMask) String() string { //nolint:unused // quiet linter
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
		return "N"
	}
	keys := slices.Collect(maps.Keys(rights))
	slices.Sort(keys)
	rv += strings.Join(keys, ",")

	if mask != 0 {
		rv += "," + fmt.Sprintf("0x%x", mask)
	}

	return rv
}

func lookupAccountSid(sidString string) (name, domain string, use uint32, err error) {
	sid, err := windows.StringToSid(sidString)
	if err != nil {
		return "", "", 0, err
	}

	var nLen, dLen uint32
	_ = windows.LookupAccountSid(nil, sid, nil, &nLen, nil, &dLen, &use) // size query

	nameBuf := make([]uint16, nLen)
	domBuf := make([]uint16, dLen)
	err = windows.LookupAccountSid(nil, sid, &nameBuf[0], &nLen, &domBuf[0], &dLen, &use)
	if err != nil {
		return "", "", 0, err
	}

	return windows.UTF16ToString(nameBuf), windows.UTF16ToString(domBuf), use, nil
}
