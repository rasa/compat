// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

/*
AreShortNamesEnabled
https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-areshortnamesenabled

GetDiskSpaceInformationW
https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-getdiskspaceinformationw
*/

package volume

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"golang.org/x/sys/windows"

	"github.com/rasa/compat"
	"github.com/rasa/compat/debug/consts"
)

func Volumes(Mounts []Mount) ([]Volume, error) {
	var volumes []Volume

	for _, mount := range Mounts {
		// fmt.Printf("%s: \n", mount.Mountpoint)
		volume, err := getVolume(mount)
		if err != nil {
			log.Printf("getVolume(%v): %v", mount.Mountpoint, err)

			continue
		}

		volumes = append(volumes, volume)
	}

	return volumes, nil
}

func getVolume(mount Mount) (Volume, error) {
	var volume Volume

	volume.Device = mount.Device
	volume.Mountpoint = mount.Mountpoint
	volume.Options = mount.Opts

	volume.Type, _ = typeOf(mount)

	root := mount.Mountpoint
	if len(root) < 3 { //nolint:mnd
		root += `\`
	}

	// Open a handle to the root of the drive (e.g. C:\)
	// Must use FILE_FLAG_BACKUP_SEMANTICS to open directories
	handle, err := windows.CreateFile(
		windows.StringToUTF16Ptr(root),
		windows.GENERIC_READ,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_FLAG_BACKUP_SEMANTICS,
		0,
	)
	if err != nil {
		return volume, err
	}
	defer windows.CloseHandle(handle) //nolint:errcheck

	var (
		volNameBuf      [windows.MAX_PATH + 1]uint16
		fsNameBuf       [windows.MAX_PATH + 1]uint16
		serialNumber    uint32
		maxComponentLen uint32
		fileSystemFlags uint32
	)

	err = windows.GetVolumeInformationByHandle(
		handle,
		&volNameBuf[0],
		uint32(len(volNameBuf)),
		&serialNumber,
		&maxComponentLen,
		&fileSystemFlags,
		&fsNameBuf[0],
		uint32(len(fsNameBuf)),
	)
	if err != nil {
		err := fmt.Errorf("GetVolumeInformationByHandleW failed: %w", err)

		return volume, err
	}

	volume.Label = windows.UTF16ToString(volNameBuf[:])
	volume.SerialNumber = strconv.Itoa(int(serialNumber))

	fi, err := compat.Stat(root)
	if err != nil {
		log.Printf("Stat(%v): %v", root, err)
	} else {
		volume.ID = fi.PartitionID()
	}

	features := map[Feature]Availability{}
	osFeatures := map[OSFeature]Availability{}

	for k := range osFeatureMap {
		u := uint32(k)

		f := OSFeature(k)
		if fileSystemFlags&u != u {
			osFeatures[f] = AvailabilityNever

			continue
		}

		osFeatures[f] = AvailabilityAlways

		feature, ok := featureMap[f]
		if ok {
			features[feature] = AvailabilityAlways
		}
	}

	fsName := strings.ToLower(windows.UTF16ToString(fsNameBuf[:]))

	filesystem, ok := filesystemMap[fsName]
	if ok {
		volume.Filesystem = filesystem
		volume.Filesystem.MaxNameLength = maxComponentLen

		log.Printf("Unknown filesystem %s", fsName)
	} else {
		volume.Filesystem = NewFilesystem(fsName)
	}

	volume.KnownFilesystem = ok

	volume.Filesystem.Name = fsName
	volume.Filesystem.Features = features
	volume.Filesystem.OSFeatures = osFeatures

	return volume, nil
}

// https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-getvolumeinformationw
var osFeatureMap = map[OSFeature]string{
	// @TODO reformat these labels
	consts.FILE_CASE_SENSITIVE_SEARCH:        "case_sensitive_search",        // 0x00000001: The specified volume supports case-sensitive file names.
	consts.FILE_CASE_PRESERVED_NAMES:         "case_preserved_names",         // 0x00000002: The specified volume supports preserved case of file names when it places a name on disk.
	consts.FILE_UNICODE_ON_DISK:              "unicode_on_disk",              // 0x00000004: The specified volume supports Unicode in file names as they appear on disk.
	consts.FILE_PERSISTENT_ACLS:              "persistent_acls",              // 0x00000008: The specified volume preserves and enforces access control lists (ACL). For example, the NTFS file system preserves and enforces ACLs, and the FAT file system does not.
	consts.FILE_FILE_COMPRESSION:             "file_compression",             // 0x00000010: The specified volume supports file-based compression.
	consts.FILE_VOLUME_QUOTAS:                "volume_quotas",                // 0x00000020: The specified volume supports disk quotas.
	consts.FILE_SUPPORTS_SPARSE_FILES:        "supports_sparse_files",        // 0x00000040: The specified volume supports sparse files.
	consts.FILE_SUPPORTS_REPARSE_POINTS:      "supports_reparse_points",      // 0x00000080: The specified volume supports reparse points.
	consts.FILE_SUPPORTS_REMOTE_STORAGE:      "supports_remote_storage",      // 0x00000100: The file system supports remote storage.
	consts.FILE_RETURNS_CLEANUP_RESULT_INFO:  "returns_cleanup_result_info",  // 0x00000200: On a successful cleanup operation, the file system returns information that describes additional actions taken during cleanup, such as deleting the file. File system filters can examine this information in their post-cleanup callback.
	consts.FILE_SUPPORTS_POSIX_UNLINK_RENAME: "supports_posix_unlink_rename", // 0x00000400: The file system supports POSIX-style delete and rename operations.
	consts.FILE_VOLUME_IS_COMPRESSED:         "volume_is_compressed",         // 0x00008000: The specified volume is a compressed volume, for example, a DoubleSpace volume.
	consts.FILE_SUPPORTS_OBJECT_IDS:          "supports_object_ids",          // 0x00010000: The specified volume supports object identifiers.
	consts.FILE_SUPPORTS_ENCRYPTION:          "supports_encryption",          // 0x00020000: The specified volume supports the Encrypted File System (EFS). For more information, see File Encryption.
	consts.FILE_NAMED_STREAMS:                "named_streams",                // 0x00040000: The specified volume supports named streams.
	consts.FILE_READ_ONLY_VOLUME:             "read_only_volume",             // 0x00080000: The specified volume is read-only.
	consts.FILE_SEQUENTIAL_WRITE_ONCE:        "sequential_write_once",        // 0x00100000: The specified volume supports a single sequential write.
	consts.FILE_SUPPORTS_TRANSACTIONS:        "supports_transactions",        // 0x00200000: The specified volume supports transactions. For more information, see About KTM.
	consts.FILE_SUPPORTS_HARD_LINKS:          "supports_hard_links",          // 0x00400000: The specified volume supports hard links. For more information, see Hard Links and Junctions.
	consts.FILE_SUPPORTS_EXTENDED_ATTRIBUTES: "supports_extended_filesystem", // 0x00800000: The specified volume supports extended filesystem. An extended attribute is a piece of application-specific metadata that an application can associate with a file and is not part of the file's data.
	consts.FILE_SUPPORTS_OPEN_BY_FILE_ID:     "supports_open_by_file_id",     // 0x01000000: The file system supports open by FileID. For more information, see FILE_ID_BOTH_DIR_INFO.
	consts.FILE_SUPPORTS_USN_JOURNAL:         "supports_usn_journal",         // 0x02000000: The specified volume supports update sequence number (USN) journals. For more information, see Change Journal Records.
	consts.FILE_SUPPORTS_INTEGRITY_STREAMS:   "supports_integrity_streams",   // 0x04000000: The file system supports integrity streams.
	consts.FILE_SUPPORTS_BLOCK_REFCOUNTING:   "supports_block_refcounting",   // 0x08000000: The specified volume supports sharing logical clusters between files on the same volume. The file system reallocates on writes to shared clusters. Indicates that FSCTL_DUPLICATE_EXTENTS_TO_FILE is a supported operation.
	consts.FILE_SUPPORTS_SPARSE_VDL:          "supports_sparse_vdl",          // 0x10000000: The file system tracks whether each cluster of a file contains valid data (either from explicit file writes or automatic zeros) or invalid data (has not yet been written to or zeroed). File systems that use sparse valid data length (VDL) do not store a valid data length and do not require that valid data be contiguous within a file.
	consts.FILE_DAX_VOLUME:                   "dax_volume",                   // 0x20000000: The specified volume is a direct access (DAX) volume.
	consts.FILE_SUPPORTS_GHOSTING:            "supports_ghosting",            // 0x40000000: The file system supports ghosting.
}

var featureMap = map[OSFeature]Feature{
	consts.FILE_CASE_PRESERVED_NAMES:         FeatureCasePreserving,
	consts.FILE_FILE_COMPRESSION:             FeatureCompression,
	consts.FILE_PERSISTENT_ACLS:              FeatureACLs,
	consts.FILE_READ_ONLY_VOLUME:             FeatureIsReadOnly,
	consts.FILE_SUPPORTS_BLOCK_REFCOUNTING:   FeatureDeduplication,
	consts.FILE_SUPPORTS_ENCRYPTION:          FeatureEncryption,
	consts.FILE_SUPPORTS_EXTENDED_ATTRIBUTES: FeatureExtendedFilesystem,
	consts.FILE_SUPPORTS_HARD_LINKS:          FeatureHardLinks,
	consts.FILE_SUPPORTS_REPARSE_POINTS:      FeatureSymbolicLinks,
	consts.FILE_SUPPORTS_SPARSE_FILES:        FeatureSparseFiles,
	consts.FILE_VOLUME_IS_COMPRESSED:         FeatureIsCompressed,

	// consts.FILE_SUPPORTS_TRANSACTIONS:        FeatureTransactions,
	// consts.FILE_UNICODE_ON_DISK:              FeatureUnicode,

	// consts.FILE_CASE_SENSITIVE_SEARCH
	// consts.FILE_DAX_VOLUME:                   "dax_volume",                   // 0x20000000: The specified volume is a direct access (DAX) volume.
	// consts.FILE_NAMED_STREAMS:                "named_streams",                // 0x00040000: The specified volume supports named streams.
	// consts.FILE_RETURNS_CLEANUP_RESULT_INFO:  "returns_cleanup_result_info",  // 0x00000200: On a successful cleanup operation, the file system returns information that describes additional actions taken during cleanup, such as deleting the file. File system filters can examine this information in their post-cleanup callback.
	// consts.FILE_SEQUENTIAL_WRITE_ONCE:        "sequential_write_once",        // 0x00100000: The specified volume supports a single sequential write.
	// consts.FILE_SUPPORTS_GHOSTING:            "supports_ghosting",            // 0x40000000: The file system supports ghosting.
	// consts.FILE_SUPPORTS_INTEGRITY_STREAMS:   "supports_integrity_streams",   // 0x04000000: The file system supports integrity streams.
	// consts.FILE_SUPPORTS_OBJECT_IDS:          "supports_object_ids",          // 0x00010000: The specified volume supports object identifiers.
	// consts.FILE_SUPPORTS_OPEN_BY_FILE_ID:     "supports_open_by_file_id",     // 0x01000000: The file system supports open by FileID. For more information, see FILE_ID_BOTH_DIR_INFO.
	// consts.FILE_SUPPORTS_POSIX_UNLINK_RENAME: "supports_posix_unlink_rename", // 0x00000400: The file system supports POSIX-style delete and rename operations.
	// consts.FILE_SUPPORTS_REMOTE_STORAGE:      "supports_remote_storage",      // 0x00000100: The file system supports remote storage.
	// consts.FILE_SUPPORTS_SPARSE_VDL:          "supports_sparse_vdl",          // 0x10000000: The file system tracks whether each cluster of a file contains valid data (either from explicit file writes or automatic zeros) or invalid data (has not yet been written to or zeroed). File systems that use sparse valid data length (VDL) do not store a valid data length and do not require that valid data be contiguous within a file.
	// consts.FILE_SUPPORTS_USN_JOURNAL:         "supports_usn_journal",         // 0x02000000: The specified volume supports update sequence number (USN) journals. For more information, see Change Journal Records.
	// consts.FILE_VOLUME_QUOTAS:                "volume_quotas",                // 0x00000020: The specified volume supports disk quotas.
}
