// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build unix && !(openbsd && ppc64) && !(netbsd && 386) && !(freebsd && riscv64) && !(cgo && aix && ppc64)

package volume

import (
	"context"
	"log"
	"strings"

	"github.com/shirou/gopsutil/v4/disk"

	"github.com/rasa/compat"
)

var osFeatureMap = map[OSFeature]string{}

func Volumes(mounts []Mount) ([]Volume, error) {
	volumes := make([]Volume, 0, len(mounts))

	ctx := context.Background()

	for _, mount := range mounts {
		volume := Volume{}

		volume.Mountpoint = mount.Mountpoint
		volume.Device = mount.Device
		volume.Options = mount.Opts

		label, err := disk.LabelWithContext(ctx, mount.Device)
		if err != nil {
			return []Volume{}, err
		}

		volume.Label = label

		serialNumber, err := disk.SerialNumberWithContext(ctx, mount.Device)
		if err != nil {
			return []Volume{}, err
		}

		volume.SerialNumber = serialNumber

		usage, err := disk.UsageWithContext(ctx, mount.Mountpoint)
		if err != nil {
			return []Volume{}, err
		}

		volume.Total = usage.Total
		volume.Used = usage.Used
		volume.Free = usage.Free
		volume.InodesTotal = usage.InodesTotal
		volume.InodesUsed = usage.InodesUsed
		volume.InodesFree = usage.InodesFree

		// volume.guid         = "guid"
		// volume.FilesystemName = strings.ToLower(mount.Fstype)

		// volume.Features = make(map[Feature]Availability)
		// volume.OSFeatures = make(map[OSFeature]Availability)

		fi, err := compat.Stat(mount.Mountpoint)
		if err != nil {
			return []Volume{}, err
		}

		volume.ID = fi.PartitionID()

		volType, _ := typeOf(mount)
		volume.Type = volType

		fsName := strings.ToLower(mount.Fstype)

		filesystem, found := filesystemMap[fsName]
		if found {
			volume.Filesystem = filesystem
		} else {
			log.Printf("No filesystem defined for %v filesystem", fsName)
			volume.Filesystem = NewFilesystem(fsName)
		}

		volume.KnownFilesystem = found
		volume.Filesystem.Name = fsName
		volumes = append(volumes, volume)
	}

	return volumes, nil
}
