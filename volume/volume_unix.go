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

		label, _ := disk.LabelWithContext(ctx, mount.Device)
		volume.Label = label

		serialNumber, _ := disk.SerialNumberWithContext(ctx, mount.Device)
		volume.SerialNumber = serialNumber

		usage, _ := disk.UsageWithContext(ctx, mount.Mountpoint)
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
		if err == nil {
			volume.ID = fi.PartitionID()
		}

		volType, _ := typeOf(mount)
		volume.Type = volType

		fsName := strings.ToLower(mount.Fstype)
		filesystem, ok := filesystemMap[fsName]
		if ok {
			volume.Filesystem = filesystem
		} else {
			log.Printf("No filesystem defined for %v filesystem", fsName)
			volume.Filesystem = NewFilesystem(fsName)
		}
		volume.KnownFilesystem = ok
		volume.Filesystem.Name = fsName
		volumes = append(volumes, volume)
	}

	return volumes, nil
}
