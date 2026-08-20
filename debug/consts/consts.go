// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package consts

const (
	// https://github.com/golang/sys/blob/5e63aa5e0fdbc13e970e0b19c47af41dd3c96f45/windows/types_windows.go#L17

	O_RDONLY    = 0x00000
	O_WRONLY    = 0x00001
	O_RDWR      = 0x00002
	O_CREAT     = 0x00040
	O_EXCL      = 0x00080
	O_NOCTTY    = 0x00100
	O_TRUNC     = 0x00200
	O_NONBLOCK  = 0x00800
	O_APPEND    = 0x00400
	O_SYNC      = 0x01000
	O_ASYNC     = 0x02000
	O_DIRECTORY = 0x04000
	O_CLOEXEC   = 0x80000

	// https://github.com/golang/go/blob/fc88e18b4a781a8751799a123cdac8b29a92409d/src/syscall/types_windows.go#L52

	O_NOFOLLOW_ANY = 0x200000000
	O_OPEN_REPARSE = 0x400000000
	O_WRITE_ATTRS  = 0x800000000

	// https://github.com/golang/sys/blob/5e63aa5e0fdbc13e970e0b19c47af41dd3c96f45/windows/types_windows.go#L105C1-L123C51

	FILE_ATTRIBUTE_READONLY              = 0x00000001
	FILE_ATTRIBUTE_HIDDEN                = 0x00000002
	FILE_ATTRIBUTE_SYSTEM                = 0x00000004
	FILE_ATTRIBUTE_DIRECTORY             = 0x00000010
	FILE_ATTRIBUTE_ARCHIVE               = 0x00000020
	FILE_ATTRIBUTE_DEVICE                = 0x00000040
	FILE_ATTRIBUTE_NORMAL                = 0x00000080
	FILE_ATTRIBUTE_TEMPORARY             = 0x00000100
	FILE_ATTRIBUTE_SPARSE_FILE           = 0x00000200
	FILE_ATTRIBUTE_REPARSE_POINT         = 0x00000400
	FILE_ATTRIBUTE_COMPRESSED            = 0x00000800
	FILE_ATTRIBUTE_OFFLINE               = 0x00001000
	FILE_ATTRIBUTE_NOT_CONTENT_INDEXED   = 0x00002000
	FILE_ATTRIBUTE_ENCRYPTED             = 0x00004000
	FILE_ATTRIBUTE_INTEGRITY_STREAM      = 0x00008000
	FILE_ATTRIBUTE_VIRTUAL               = 0x00010000
	FILE_ATTRIBUTE_NO_SCRUB_DATA         = 0x00020000
	FILE_ATTRIBUTE_RECALL_ON_OPEN        = 0x00040000
	FILE_ATTRIBUTE_RECALL_ON_DATA_ACCESS = 0x00400000

	// https://github.com/golang/sys/blob/5e63aa5e0fdbc13e970e0b19c47af41dd3c96f45/windows/types_windows.go#L133C1-L145C46

	FILE_FLAG_OPEN_REQUIRING_OPLOCK = 0x00040000
	FILE_FLAG_FIRST_PIPE_INSTANCE   = 0x00080000
	FILE_FLAG_OPEN_NO_RECALL        = 0x00100000
	FILE_FLAG_OPEN_REPARSE_POINT    = 0x00200000
	FILE_FLAG_SESSION_AWARE         = 0x00800000
	FILE_FLAG_POSIX_SEMANTICS       = 0x01000000
	FILE_FLAG_BACKUP_SEMANTICS      = 0x02000000
	FILE_FLAG_DELETE_ON_CLOSE       = 0x04000000
	FILE_FLAG_SEQUENTIAL_SCAN       = 0x08000000
	FILE_FLAG_RANDOM_ACCESS         = 0x10000000
	FILE_FLAG_NO_BUFFERING          = 0x20000000
	FILE_FLAG_OVERLAPPED            = 0x40000000
	FILE_FLAG_WRITE_THROUGH         = 0x80000000

	// See https://github.com/golang/sys/blob/master/windows/types_windows.go#L68

	O_FILE_FLAG_OPEN_NO_RECALL     = FILE_FLAG_OPEN_NO_RECALL
	O_FILE_FLAG_OPEN_REPARSE_POINT = FILE_FLAG_OPEN_REPARSE_POINT
	O_FILE_FLAG_SESSION_AWARE      = FILE_FLAG_SESSION_AWARE
	O_FILE_FLAG_POSIX_SEMANTICS    = FILE_FLAG_POSIX_SEMANTICS
	O_FILE_FLAG_BACKUP_SEMANTICS   = FILE_FLAG_BACKUP_SEMANTICS
	O_FILE_FLAG_DELETE_ON_CLOSE    = FILE_FLAG_DELETE_ON_CLOSE
	O_FILE_FLAG_SEQUENTIAL_SCAN    = FILE_FLAG_SEQUENTIAL_SCAN
	O_FILE_FLAG_RANDOM_ACCESS      = FILE_FLAG_RANDOM_ACCESS
	O_FILE_FLAG_NO_BUFFERING       = FILE_FLAG_NO_BUFFERING
	O_FILE_FLAG_OVERLAPPED         = FILE_FLAG_OVERLAPPED
	O_FILE_FLAG_WRITE_THROUGH      = FILE_FLAG_WRITE_THROUGH
)

// https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-getdrivetypew

// https://github.com/golang/sys/blob/5e63aa5e0fdbc13e970e0b19c47af41dd3c96f45/windows/syscall_windows.go#L35
const (
	DRIVE_UNKNOWN     = 0
	DRIVE_NO_ROOT_DIR = 1
	DRIVE_REMOVABLE   = 2
	DRIVE_FIXED       = 3
	DRIVE_REMOTE      = 4
	DRIVE_CDROM       = 5
	DRIVE_RAMDISK     = 6
)

// https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-getvolumeinformationw

const (
	// The specified volume supports case-sensitive file names.
	FILE_CASE_SENSITIVE_SEARCH = 0x00000001
	// The specified volume supports preserved case of file names when it places a name on disk.
	FILE_CASE_PRESERVED_NAMES = 0x00000002
	// The specified volume supports Unicode in file names as they appear on disk.
	FILE_UNICODE_ON_DISK = 0x00000004
	// The specified volume preserves and enforces access control lists (ACL).
	// For example, the NTFS file system preserves and enforces ACLs, and the FAT file system does not.
	FILE_PERSISTENT_ACLS = 0x00000008
	// The specified volume supports file-based compression.
	FILE_FILE_COMPRESSION = 0x00000010
	// The specified volume supports disk quotas.
	FILE_VOLUME_QUOTAS = 0x00000020
	// The specified volume supports sparse files.
	FILE_SUPPORTS_SPARSE_FILES = 0x00000040
	// The specified volume supports reparse points.
	// ReFS: ReFS supports reparse points but does not index them so FindFirstVolumeMountPoint
	// and FindNextVolumeMountPoint will not function as expected.
	FILE_SUPPORTS_REPARSE_POINTS = 0x00000080
	// The file system supports remote storage.
	FILE_SUPPORTS_REMOTE_STORAGE = 0x00000100
	// On a successful cleanup operation, the file system returns information that describes
	// additional actions taken during cleanup, such as deleting the file.
	// File system filters can examine this information in their post-cleanup callback.
	FILE_RETURNS_CLEANUP_RESULT_INFO = 0x00000200
	// The file system supports POSIX-style delete and rename operations.
	FILE_SUPPORTS_POSIX_UNLINK_RENAME = 0x00000400
	// The specified volume is a compressed volume, for example, a DoubleSpace volume.
	FILE_VOLUME_IS_COMPRESSED = 0x00008000
	// The specified volume supports object identifiers.
	FILE_SUPPORTS_OBJECT_IDS = 0x00010000
	// The specified volume supports the Encrypted File System (EFS).
	// For more information, see File Encryption.
	FILE_SUPPORTS_ENCRYPTION = 0x00020000
	// The specified volume supports named streams.
	FILE_NAMED_STREAMS = 0x00040000
	// The specified volume is read-only.
	FILE_READ_ONLY_VOLUME = 0x00080000
	// The specified volume supports a single sequential write.
	FILE_SEQUENTIAL_WRITE_ONCE = 0x00100000
	// The specified volume supports transactions. For more information, see About KTM.
	FILE_SUPPORTS_TRANSACTIONS = 0x00200000
	// The specified volume supports hard links. For more information, see Hard Links and Junctions.
	// Windows Server 2008, Windows Vista, Windows Server 2003 and Windows XP:
	// This value is not supported until Windows Server 2008 R2 and Windows 7.
	FILE_SUPPORTS_HARD_LINKS = 0x00400000
	// The specified volume supports extended attributes. An extended attribute is a
	// piece of application-specific metadata that an application can associate with
	// a file and is not part of the file's data.
	// Windows Server 2008, Windows Vista, Windows Server 2003 and Windows XP:
	// This value is not supported until Windows Server 2008 R2 and Windows 7.
	FILE_SUPPORTS_EXTENDED_ATTRIBUTES = 0x00800000
	// The file system supports open by FileID. For more information, see FILE_ID_BOTH_DIR_INFO.
	// Windows Server 2008, Windows Vista, Windows Server 2003 and Windows XP:
	// This value is not supported until Windows Server 2008 R2 and Windows 7.
	FILE_SUPPORTS_OPEN_BY_FILE_ID = 0x01000000
	// The specified volume supports update sequence number (USN) journals.
	// For more information, see Change Journal Records.
	// Windows Server 2008, Windows Vista, Windows Server 2003 and Windows XP:
	// This value is not supported until Windows Server 2008 R2 and Windows 7.
	FILE_SUPPORTS_USN_JOURNAL = 0x02000000
	// The file system supports integrity streams.
	FILE_SUPPORTS_INTEGRITY_STREAMS = 0x04000000
	// The specified volume supports sharing logical clusters between files on the same volume.
	// The file system reallocates on writes to shared clusters.
	// Indicates that FSCTL_DUPLICATE_EXTENTS_TO_FILE is a supported operation.
	FILE_SUPPORTS_BLOCK_REFCOUNTING = 0x08000000
	// The file system tracks whether each cluster of a file contains valid data
	// (either from explicit file writes or automatic zeros) or invalid data
	// (has not yet been written to or zeroed). File systems that use sparse valid
	// data length (VDL) do not store a valid data length and do not require that
	// valid data be contiguous within a file.
	FILE_SUPPORTS_SPARSE_VDL = 0x10000000
	// The specified volume is a direct access (DAX) volume.
	// Note: This flag was introduced in Windows 10, version 1607.
	FILE_DAX_VOLUME = 0x20000000
	// The file system supports ghosting.
	FILE_SUPPORTS_GHOSTING = 0x40000000
)
