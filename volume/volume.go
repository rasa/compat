// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package volume

import (
	"fmt"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
// Volume
///////////////////////////////////////////////////////////////////////////////

type Volume struct {
	Device     string // /dev/sda1
	Mountpoint string // /mnt/sda1
	ID         uint64 // unique device ID
	Label      string
	// UUID      string
	SerialNumber string
	Type         Type
	// FilesystemName string     // ext4, etc.
	Options []string

	//
	Total           uint64
	Used            uint64
	Free            uint64
	InodesTotal     uint64
	InodesUsed      uint64
	InodesFree      uint64
	KnownFilesystem bool

	// // @TODO Move to Filesystem
	// Features map[Feature]Availability
	// // @TODO Move to Filesystem
	// OSFeatures map[OSFeature]Availability

	// Default for the filesystem found on the volume/partition
	Filesystem Filesystem
}

///////////////////////////////////////////////////////////////////////////////
// Type
///////////////////////////////////////////////////////////////////////////////

// Type defines the type of volume.
type Type uint

const (
	TypeUnknown Type = 1 << iota
	TypeUnavailable
	TypeFixed
	TypeRemovable
	TypeNetwork
	TypeOptical
	TypeRamdisk
	TypeLoop
	TypePseudo
)

var volTypeMap = map[Type]string{
	TypeUnknown:     "Unknown",     // 0: The drive type cannot be determined.
	TypeUnavailable: "Unavailable", // 1: The root path is invalid; for example, there is no volume mounted at the specified path.
	TypeFixed:       "Fixed",       // 2: The drive has removable media; for example, a floppy drive, thumb drive, or flash card reader.
	TypeRemovable:   "Removable",   // 3: The drive has fixed media; for example, a hard disk drive or flash drive.
	TypeNetwork:     "Network",     // 4: The drive is a remote (network) drive.
	TypeOptical:     "Optical",     // 5: The drive is a CD-ROM drive.
	TypeRamdisk:     "Ramdisk",     // 6: The drive is a RAM disk.
}

func (v Volume) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "Device:   %v\n", v.Device)
	fmt.Fprintf(&b, "Mount:    %v\n", v.Mountpoint)
	fmt.Fprintf(&b, "ID:       %v\n", v.ID)
	fmt.Fprintf(&b, "Label:    %v\n", v.Label)
	fmt.Fprintf(&b, "Serial#:  %v\n", v.SerialNumber)
	// fmt.Fprintf(&b, "GUID:     %v\n", v.guid)
	// fmt.Fprintf(&b, "FS:       %v\n", v.Filesystem.Name)
	fmt.Fprintf(&b, "Type:     %v\n", v.Type)

	fmt.Fprintf(&b, "Total:    %10v\n", si(float64(v.Total)))
	fmt.Fprintf(&b, "Used:     %10v\n", si(float64(v.Used)))
	fmt.Fprintf(&b, "Free:     %10v\n", si(float64(v.Free)))

	if v.InodesTotal != 0 {
		fmt.Fprintf(&b, "ITotal:   %10v\n", v.InodesTotal)
		fmt.Fprintf(&b, "IUsed:    %19v\n", v.InodesUsed)
		fmt.Fprintf(&b, "IFree:    %10v\n", v.InodesFree)
	}
	// fmt.Fprintf(&b, "Features:   %v\n", v.Features)
	// fmt.Fprintf(&b, "OSFeatures: %v\n", v.OSFeatures)
	fmt.Fprintf(&b, "Options:    %v\n", v.Options)
	fmt.Fprintf(&b, "Filesystem: %v\n", v.Filesystem)

	return b.String()
}

/*

func (v Volume) ID() uint64 {
	return v.Id
}

func (v Volume) MountPoint() string {
	return v.MountPoint
}

func (v Volume) Device() string {
	return v.Device
}

func (v Volume) Label() string {
	return v.Label
}

func (v Volume) SerialNumber() string {
	return v.SerialNumber
}

// func (v Volume) GUID() string {
// 	return v.Guid
// }

func (v Volume) FilesystemName() string {
	return v.FilesystemName
}

func (v Volume) Filesystem() Filesystem {
	return v.Filesystem
}

func (v Volume) Type() Type {
	return v.Type
}

func (v Volume) Total() uint64 {
	return v.Total
}

func (v Volume) Used() uint64 {
	return v.Used
}

func (v Volume) Free() uint64 {
	return v.Free
}

func (v Volume) Supports(supports OSFeature) bool {
	_, ok := v.Supports[supports]
	return ok
}

func (v Volume) SupportsMap() map[OSFeature]bool {
	return v.Supports
}

func (v Volume) Options() []string {
	return v.options
}
*/

func (v Type) String() string {
	s, ok := volTypeMap[v]
	if ok {
		return s
	}

	return fmt.Sprintf("unknown Type value: %d", v)
}
