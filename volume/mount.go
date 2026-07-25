// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

/*

package mount

func New(deviceOrMountPath string) (Mount, error)

func Mounts(opts []Option...) ([]Mount, error)
	WithIgnore([]string["/snap/"])

*/

package volume

// Mount contains a volume's device name, mount point, filesystem type,
// and mount options, if any. It is the same structure as
// gopsutil's disk.PartitionType struct.
type Mount struct {
	// Device is the device name for disk partition. For example, on unix-like
	// systems, /dev/sda1. On Windows systems, \\.\Device{GUID}
	Device string
	// Mountpoint is the location in the filesystem, where the device is mounted.
	// For example, on unix-like systems, /mnt/sda1. On Windows systems, D:.
	Mountpoint string
	// Fstype is the name of the type of filesystem the disk partition is
	// formatted as, in lowercase ("ext4", "fat", etc.).
	Fstype string
	// Opts are the mount options used to mount the device, if any (e.g., "rw", "ro", etc.)
	Opts []string
}
