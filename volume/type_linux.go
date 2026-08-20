// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build linux

package volume

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shirou/gopsutil/v4/disk"
	"golang.org/x/sys/unix"
)

var magicMap = map[int64]Type{
	disk.AFS_SUPER_MAGIC:   TypeNetwork, // Andrew File System (AFS)
	disk.CEPH_SUPER_MAGIC:  TypeNetwork, // Ceph distributed FS
	disk.CIFS_MAGIC_NUMBER: TypeNetwork, // CIFS/SMB (Windows shares, Samba)
	disk.CODA_SUPER_MAGIC:  TypeNetwork, // Coda FS (experimental distributed FS)
	disk.FHGFS_SUPER_MAGIC: TypeNetwork, // Fraunhofer FS (now BeeGFS)
	disk.FUSE_SUPER_MAGIC:  TypeNetwork, // FUSE can be backed by network FS (sshfs, gcsfuse, s3fs, etc.)
	// dup: disk.FUSEBLK_SUPER_MAGIC: TypeNetwork, // "
	disk.GPFS_SUPER_MAGIC:   TypeNetwork, // IBM Spectrum Scale (GPFS)
	disk.HOSTFS_SUPER_MAGIC: TypeNetwork, // hostfs (UML, can map into host FS including remote)
	disk.KAFS_SUPER_MAGIC:   TypeNetwork, // newer AFS (“kafs” implementation)
	disk.LUSTRE_SUPER_MAGIC: TypeNetwork, //nolint:misspell // Lustre parallel distributed FS
	disk.NCP_SUPER_MAGIC:    TypeNetwork, // Novell NetWare (NCP protocol)
	disk.NFS_SUPER_MAGIC:    TypeNetwork, // NFS (Network File System)
	disk.NFSD_SUPER_MAGIC:   TypeNetwork, // NFS daemon’s pseudo-FS (not data, but still network related)
	disk.PANFS_SUPER_MAGIC:  TypeNetwork, // Panasas PanFS (cluster FS)
	disk.SMB_SUPER_MAGIC:    TypeNetwork, // older SMB
	disk.VMHGFS_SUPER_MAGIC: TypeNetwork, // VMware host/guest shared folder FS
	disk.VZFS_SUPER_MAGIC:   TypeNetwork, // Virtuozzo / OpenVZ shared FS

	disk.ISOFS_R_WIN_SUPER_MAGIC: TypeOptical, // ISO-9660 (Rock Ridge / Windows extension variant)
	disk.ISOFS_SUPER_MAGIC:       TypeOptical, // ISO-9660 (standard CD-ROM filesystem)
	disk.ISOFS_WIN_SUPER_MAGIC:   TypeOptical, // ISO-9660 (Windows extension variant)
	disk.UDF_SUPER_MAGIC:         TypeOptical, // UDF (Universal Disk Format, common on DVD/Blu-Ray)

	disk.ANON_INODE_FS_SUPER_MAGIC: TypePseudo, // anonymous inodes fs
	disk.BINFMTFS_MAGIC:            TypePseudo, // Binary format handler fs
	disk.CGROUP_SUPER_MAGIC:        TypePseudo, // Control groups (cgroups v1)
	disk.CGROUP2_SUPER_MAGIC:       TypePseudo, // Control groups v2
	disk.CONFIGFS_MAGIC:            TypePseudo, // configfs
	disk.DEBUGFS_MAGIC:             TypePseudo, // debugfs
	disk.DEVPTS_SUPER_MAGIC:        TypePseudo, // Pseudo-terminal devices
	disk.EFIVARFS_MAGIC:            TypePseudo, // EFI variables
	// dup: disk.HOSTFS_SUPER_MAGIC:     TypePseudo, // UML hostfs (pass-through pseudo fs)
	disk.HUGETLBFS_MAGIC:        TypePseudo, // HugeTLB pages FS (memory, pseudo)
	disk.INOTIFYFS_SUPER_MAGIC:  TypePseudo, // Inotify fs
	disk.MQUEUE_MAGIC:           TypePseudo, // POSIX message queues (memory, pseudo)
	disk.NSFS_MAGIC:             TypePseudo, // Namespace file system
	disk.PIPEFS_MAGIC:           TypePseudo, // Pipefs (internal pipes representation)
	disk.PROC_SUPER_MAGIC:       TypePseudo, // /proc procfs
	disk.PSTOREFS_MAGIC:         TypePseudo, // pstore (persistent storage for logs)
	disk.RPC_PIPEFS_SUPER_MAGIC: TypePseudo, // RPC pipefs (used by NFS)
	disk.SECURITYFS_SUPER_MAGIC: TypePseudo, // securityfs
	disk.SELINUX_MAGIC:          TypePseudo, // SELinux virtual FS
	disk.SMACK_MAGIC:            TypePseudo, // SMACK security FS
	disk.SOCKFS_MAGIC:           TypePseudo, // Sockfs (internal sockets representation)
	disk.SYSFS_MAGIC:            TypePseudo, // /sys sysfs
	disk.TRACEFS_MAGIC:          TypePseudo, // tracefs (kernel tracing)
	disk.USBDEVICE_SUPER_MAGIC:  TypePseudo, // USB device pseudo fs
	disk.XENFS_SUPER_MAGIC:      TypePseudo, // Xen hypervisor fs

	// disk.HUGETLBFS_MAGIC: TypeRamdisk, // hugetlbfs (special RAM-backed fs for huge pages)
	// disk.MQUEUE_MAGIC:   TypeRamdisk, // POSIX message queues (tmpfs-like, memory-based)
	disk.RAMFS_MAGIC: TypeRamdisk, // ramfs (in-memory filesystem, non-persistent)
	// disk.SHM_MAGIC:      TypeRamdisk, // SysV shared memory (if exposed as fs)
	disk.TMPFS_MAGIC: TypeRamdisk, // tmpfs (RAM-backed, often used for /tmp, /dev/shm)
}

func typeOf(mount Mount) (Type, error) {
	device := mount.Device
	base := filepath.Base(device)

	// Handle special cases first
	if strings.HasPrefix(base, "loop") {
		return TypeLoop, nil
	}

	if strings.HasPrefix(base, "ram") || strings.HasPrefix(base, "zram") {
		return TypeRamdisk, nil
	}

	if strings.HasPrefix(base, "sr") {
		if readSys(filepath.Join("/sys/class/block", base, "device/type")) == "5" {
			return TypeOptical, nil
		}
	}

	// If it's not a block device, check filesystem type (network, tmpfs, etc.)
	var st unix.Statfs_t

	err := unix.Statfs(device, &st)
	if err != nil {
		return TypeUnknown, fmt.Errorf("statfs: %w", err)
	}

	magicID := st.Type

	typ, ok := magicMap[magicID]
	if ok {
		return typ, nil
	}

	// Try sysfs to detect fixed/removable
	sysPath := filepath.Join("/sys/class/block", base)

	_, err = os.Stat(sysPath)
	if os.IsNotExist(err) {
		return TypeUnavailable, nil
	}

	// Partitions: resolve parent
	path := filepath.Join(sysPath, "partition")

	_, err = os.Stat(path)
	if err != nil {
		return TypeUnknown, fmt.Errorf("os.Stat(%v): %w", path, err)
	}

	sysPath, _ = filepath.EvalSymlinks(filepath.Join(sysPath, ".."))

	if readSys(filepath.Join(sysPath, "removable")) == "1" {
		return TypeRemovable, nil
	}

	return TypeFixed, nil
}

func readSys(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}
