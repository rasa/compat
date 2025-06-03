// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows && debug

package compat

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
	"testing"

	"golang.org/x/sys/windows"
)

func dumpMasks(perm os.FileMode, ownerMask uint32, groupMask uint32, worldMask uint32) {
	if !testing.Verbose() {
		return
	}
	if os.Getenv("COMPAT_DEBUG") == "" {
		return
	}
	omask := aMask(ownerMask)
	gmask := aMask(groupMask)
	wmask := aMask(worldMask)

	fmt.Printf("perm=%04o ownerMask=%v groupMask=%v worldMask=%v\n", perm, omask, gmask, wmask)
}

// https://github.com/golang/sys/blob/3d9a6b80/windows/security_windows.go#L992
var maskMap = map[uint32]string{
	windows.GENERIC_READ:    "GR", // 0x80000000
	windows.GENERIC_WRITE:   "GW", // 0x40000000
	windows.GENERIC_EXECUTE: "GE", // 0x20000000
	windows.GENERIC_ALL:     "GA", // 0x10000000
	windows.DELETE:          "D",  // 0x00010000
}

type aMask uint32

func (a aMask) String() string {
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
