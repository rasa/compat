// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// volinfo retrieves information about volumes (partitions).
package main

import (
	"cmp"
	"fmt"
	"log"
	"maps"
	"slices"
	"strings"

	"github.com/rasa/compat/volume"
)

var IgnoreMounts = []string{"/snap/"}

func main() {
	mounts, err := volume.Mounts()
	if err != nil {
		log.Fatalf("Mounts: %v", err)
	}

	// Sort by mount path
	slices.SortFunc(mounts, func(a, b volume.Mount) int {
		return cmp.Compare(a.Mountpoint, b.Mountpoint)
	})

nextMount:
	for i, mnt := range mounts {
		for _, ignore := range IgnoreMounts {
			if strings.HasPrefix(mnt.Mountpoint, ignore) {
				continue nextMount
			}
		}

		fmt.Printf("%d: %v: %v\n", i, mnt.Mountpoint, mnt.Device)
	}

	volumes, err := volume.Volumes(mounts)
	if err != nil {
		log.Fatalf("Volumes: %v", err)
	}

nextVolume:
	for i, vol := range volumes {
		for _, ignore := range IgnoreMounts {
			if strings.HasPrefix(vol.Mountpoint, ignore) {
				continue nextVolume
			}
		}

		fmt.Printf("%d: %v: \n", i, vol.Mountpoint)
		fmt.Printf("%v\n", vol)
		// entries, err := compat.ReadDir(volume.MountPoint())
		// if err != nil {
		// 	fmt.Printf("ReadDir(%v): %v\n", volume.MountPoint(), err)
		// }
		// for _, entry := range entries {
		// 	name := entry.Name()
		// 	path := filepath.Join(volume.MountPoint(), name)
		// 	fi, err := compat.Stat(path)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	fmt.Printf("%s", fi.String())
		// 	break
		// }
	}

	fmt.Println("")

nextVolume2:
	for i, vol := range volumes {
		for _, ignore := range IgnoreMounts {
			if strings.HasPrefix(vol.Mountpoint, ignore) {
				continue nextVolume2
			}
		}

		fmt.Printf("%d: %v: \n", i, vol.Mountpoint)
		keys := slices.Collect(maps.Keys(vol.Filesystem.OSFeatures))

		// Sort keys by comparing their values in the map
		slices.SortFunc(keys, func(a, b volume.OSFeature) int {
			// Compare values: return <0 if a<b, 0 if equal, >0 if a>b
			return cmp.Compare(a.String(), b.String())
		})

		fmt.Println("Supported:")

		for _, k := range keys {
			fmt.Printf("  %s\n", k)
		}

		fmt.Println("")
	}
}
