// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package volume

import (
	"golang.org/x/sys/windows"
)

// https://github.com/golang/sys/blob/5e63aa5e0fdbc13e970e0b19c47af41dd3c96f45/windows/syscall_windows.go#L35
var driveTypeMap = map[uint32]Type{
	windows.DRIVE_UNKNOWN:     TypeUnknown,     // 0: The drive type cannot be determined.
	windows.DRIVE_NO_ROOT_DIR: TypeUnavailable, // 1: The root path is invalid; for example, there is no volume mounted at the specified path.
	windows.DRIVE_REMOVABLE:   TypeRemovable,   // 2: The drive has removable media; for example, a floppy drive, thumb drive, or flash card reader.
	windows.DRIVE_FIXED:       TypeFixed,       // 3: The drive has fixed media; for example, a hard disk drive or flash drive.
	windows.DRIVE_REMOTE:      TypeNetwork,     // 4: The drive is a remote (network) drive.
	windows.DRIVE_CDROM:       TypeOptical,     // 5: The drive is a CD-ROM drive.
	windows.DRIVE_RAMDISK:     TypeRamdisk,     // 6: The drive is a RAM disk.
}

func typeOf(mount Mount) (Type, error) {
	driveTypeId := windows.GetDriveType(windows.StringToUTF16Ptr(mount.Mountpoint))
	volType, ok := driveTypeMap[driveTypeId]
	if !ok {
		volType = driveTypeMap[windows.DRIVE_UNKNOWN]
	}
	return volType, nil
}
