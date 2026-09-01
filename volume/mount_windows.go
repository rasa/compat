// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package volume

import (
	"cmp"
	"slices"
	"syscall"

	"golang.org/x/sys/windows"
)

func Mounts() ([]Mount, error) {
	mounts := []Mount{}

	const bufSize = 1024

	var volNameBuf [bufSize]uint16

	// Start volume enumeration
	h, err := windows.FindFirstVolume(&volNameBuf[0], bufSize)
	if err != nil {
		return mounts, err
	}
	defer windows.FindVolumeClose(h)

	for {
		volName := syscall.UTF16ToString(volNameBuf[:])

		// Collect mount points for this volume
		paths, err := getVolumePaths(volName)
		if err == nil {
			for _, path := range paths {
				mounts = append(mounts, Mount{Device: volName, Mountpoint: path})
			}
		}

		// Next volume
		err = windows.FindNextVolume(h, &volNameBuf[0], bufSize)
		if err != nil {
			break // no more volumes
		}
	}

	// Sort by mount path
	slices.SortFunc(mounts, func(a, b Mount) int {
		return cmp.Compare(a.Mountpoint, b.Mountpoint)
	})

	return mounts, nil
}

func getVolumePaths(volName string) ([]string, error) {
	const bufSize = 1024

	var (
		mountBuf  [bufSize]uint16
		returnLen uint32
		paths     []string
	)

	volPtr, err := syscall.UTF16PtrFromString(volName)
	if err != nil {
		return nil, err
	}

	err = windows.GetVolumePathNamesForVolumeName(
		volPtr,
		&mountBuf[0],
		bufSize,
		&returnLen,
	)
	if err != nil {
		return nil, err
	}

	// Parse MULTI_SZ (series of null-terminated UTF-16 strings)
	start := 0

	for i, c := range mountBuf {
		if c == 0 {
			if i > start {
				paths = append(paths, syscall.UTF16ToString(mountBuf[start:i]))
			}

			start = i + 1
		}
	}

	return paths, nil
}
