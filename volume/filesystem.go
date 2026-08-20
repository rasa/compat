// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package volume

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/dustin/go-humanize"
	"golang.org/x/text/unicode/norm"
)

///////////////////////////////////////////////////////////////////////////////
// Filesystem
///////////////////////////////////////////////////////////////////////////////

type Filesystem struct {
	Name          string // Filesystem name, in all lowercase.
	MaxNameLength uint32 // maximum length (in CharSize units) of filename component*
	MaxPathLength uint32 // maximum length (in CharSize units) of an entire path*
	CharSize      uint8  // 1 for bytes, 2 for UTF-16/UCS-2, 4 for UTF-8/Unicode
	MaxFileSize   Size   // maxiumum file size*
	MaxVolumeSize Size   // maxiumum file size*
	MaxFiles      uint64 // maximum number of files*

	ATimeGranularity time.Duration // Last Accessed timestamp granularity (0=not supported)
	BTimeGranularity time.Duration // Creation (Birthed) timestamp granularity (0=not supported)
	CTimeGranularity time.Duration // Last metadata Change timestamp granularity (0=not supported)
	MTimeGranularity time.Duration // Last Modified timestamp granularity (0=not supported)

	CaseMapping          CaseMapping
	UnicodeNormalization UnicodeNormalization

	// Mostly related to FAT/NTFS/ReFS filesystems:
	DisallowedRunes   DisallowedRunes   // characters disallowed by the filesystem.
	ReservedFilenames ReservedFilenames // regexes of filenames that are reserved by the filesystem.

	Features   map[Feature]Availability
	OSFeatures map[OSFeature]Availability
}

const (
	// Binary (IEC) powers of two exponents.

	// KiBPow2 defines 10.
	KiBPow2 = (iota + 1) * 10 // 10
	MiBPow2                   // 20
	GiBPow2                   // 30
	TiBPow2                   // 40
	PiBPow2                   // 50
	EiBPow2                   // 60
	ZiBPow2                   // 70
	YiBPow2                   // 80

	// Binary (IEC) byte sizes
	// Does not compile.
	// _  = 1 << (10 * iota) // ignore iota=0
	// KiB                    // 1 << 10 (2^10 bytes)
	// MiB                    // 1 << 20 (2^20 bytes)
	// GiB                    // 1 << 30 (2^30 bytes)
	// TiB                    // 1 << 40 (2^40 bytes)
	// PiB                    // 1 << 50 (2^50 bytes)
	// EiB                    // 1 << 60 (2^60 bytes)
	// ZiB                    // 1 << 70 (2^70 bytes) (overflow if typed as uint64)
	// YiB                    // 1 << 80 (2^80 bytes).

	// Binary (IEC) byte sizes.

	KiB = 1 << 10 // 2^10 bytes
	MiB = 1 << 20 // 2^20 bytes
	GiB = 1 << 30 // 2^30 bytes
	TiB = 1 << 40 // 2^40 bytes
	PiB = 1 << 50 // 2^50 bytes
	EiB = 1 << 60 // 2^60 bytes
	ZiB = 1 << 70 // 2^70 bytes (overflows uint64, use math/big or float64)
	YiB = 1 << 80 // 2^80 bytes "

	_  = 1e3 * iota
	KB // 1e3 (10^3)
	MB // 1e6 (10^6)
	GB // 1e9 (10^9)
	TB // 1e12 (10^12)
	PB // 1e15 (10^15)
	EB // 1e18 (10^18)
	ZB // 1e21 (10^21)
	YB // 1e24 (10^24)

	MaxChars255    = 255
	MaxChars32767  = 32767
	CharSizeByte   = 1
	CharSizeUTF16  = 2
	CharSizeUTF8   = 4
	MaxUint32      = (1 << 32) - 1 //nolint:mnd
	MaxUint64      = (1 << 64) - 1 //nolint:mnd
	moLimitDefined = 0
)

///////////////////////////////////////////////////////////////////////////////
// Size
///////////////////////////////////////////////////////////////////////////////

type Size struct {
	Float    float64
	Base     uint64
	Exponent uint64
}

///////////////////////////////////////////////////////////////////////////////
// CaseMapping
///////////////////////////////////////////////////////////////////////////////

type CaseMapping uint

const (
	CaseMappingUnknown = 0 + iota
	CaseMappingASCII
	CaseMappingUnicode
)

///////////////////////////////////////////////////////////////////////////////
// UnicodeNormalization
///////////////////////////////////////////////////////////////////////////////

type UnicodeNormalization uint

const (
	UnicodeNormalizationUnknown     = 0 + iota
	UnicodeNormalizationUnsupported // exFAT on Linux 6.8
	UnicodeNormalizationNone
	UnicodeNormalizationNFC
	UnicodeNormalizationNFD
)

func normalizeNFC(s string) string { //nolint:unused
	return norm.NFC.String(s)
}

func normalizeNFD(s string) string { //nolint:unused
	return norm.NFD.String(s)
}

///////////////////////////////////////////////////////////////////////////////
// Feature
///////////////////////////////////////////////////////////////////////////////

type Feature uint

const (
	// Metadata
	// https://en.wikipedia.org/wiki/Comparison_of_file_systems#Metadata

	FeatureStoresFileOwner Feature = 1 << iota
	FeaturePOSIXPermissions
	FeatureCreationTimestamps
	FeatureLastAccessTimestamps
	FeatureLastMetadataChangeTimestamps
	FeatureLastArchiveTimestamps
	FeatureACLs
	FeatureMACLabels
	FeatureExtendedFilesystem
	FeatureMetadataChecksums

	// File capabilities
	// https://en.wikipedia.org/wiki/Comparison_of_file_systems#File_capabilities

	FeatureHardLinks
	FeatureSymbolicLinks
	FeatureBlockJournaling
	FeatureMetadataOnlyJournaling
	FeatureCaseSensitive
	FeatureCasePreserving
	FeatureXIP

	// Block capabilities
	// https://en.wikipedia.org/wiki/Comparison_of_file_systems#Block_capabilities

	FeatureSnapshots
	FeatureEncryption
	FeatureDeduplication
	FeatureDataChecksums
	FeaturePersistentCache
	FeatureMultipleDevices
	FeatureCompression
	FeatureSelfHealing

	// Resize capabilities
	// https://en.wikipedia.org/wiki/Comparison_of_file_systems#Resize_capabilities

	FeatureOfflineGrow
	FeatureOnlineGrow
	FeatureOfflineShrink
	FeatureOnlineShrink
	FeatureAddRemoveVolumes

	// Allocation and layout policies
	// https://en.wikipedia.org/wiki/Comparison_of_file_systems#Allocation_and_layout_policies

	FeatureSparseFiles
	FeatureBlockSuballocation
	FeatureTailPacking
	FeatureExtents
	FeatureVariableFileBlockSize
	FeatureAllocateOnFlush
	FeatureCopyOnWrite
	FeatureTrimSupport

	// Not listed on wikipedia.

	FeatureIsCompressed
	FeatureIsEncrypted
	FeatureIsReadOnly

	FeatureCreationTimestampsUpdateable
	FeatureLastMetadataChangeTimestampsUpdateable
	FeatureLastArchiveTimestampsUpdateable
)

type Availability uint

const (
	AvailabilityNever Availability = 0 + iota
	AvailabilityOptional
	AvailabilityAlways
)

var FeaturesNone = map[Feature]Availability{}

///////////////////////////////////////////////////////////////////////////////
// OSFeature
///////////////////////////////////////////////////////////////////////////////

type OSFeature uint64

var OSFeaturesNone = map[OSFeature]Availability{}

///////////////////////////////////////////////////////////////////////////////
// DisallowedRunes
///////////////////////////////////////////////////////////////////////////////

type RuneRange struct {
	Start rune
	End   rune
}

type DisallowedRunes []RuneRange

// Source: https://en.wikipedia.org/wiki/Comparison_of_file_systems#Limits

var (
	DisallowedRunesNone = []RuneRange{}
	DisallowedRunesNUL  = []RuneRange{
		{0x00, 0x00},
	}
	DisallowedRunesNULSlash = []RuneRange{
		{0x00, 0x00},
		{'/', '/'},
	}
	DisallowedRunesNULColon = []RuneRange{
		{0x00, 0x00},
		{':', ':'},
	}
	DisallowedRunesFAT = []RuneRange{
		// 0x00-0x1f, " * / : < > ? \ |.
		{0x00, 0x1f},
		{'"', '"'},
		{'*', '*'},
		{'/', '/'},
		{':', ':'},
		{'<', '<'},
		{'>', '>'},
		{'?', '?'},
		{'\\', '\\'},
		{'|', '|'},
	}
	DisallowedRunesAndroid = []RuneRange{
		// 0x00-0x1f, " * / : < > ? \ | DEL.
		{0x00, 0x1f},
		{'"', '"'},
		{'*', '*'},
		{'/', '/'},
		{':', ':'},
		{'<', '<'},
		{'>', '>'},
		{'?', '?'},
		{'\\', '\\'},
		{'|', '|'},
		{0x7f, 0x7f}, // DEL
	}
)

func (v DisallowedRunes) String() string {
	var builder strings.Builder

	for i, runeRange := range v {
		if i > 0 {
			fmt.Fprint(&builder, ",")
		}

		fmt.Fprint(&builder, addRune(runeRange.Start))

		if runeRange.Start != runeRange.End {
			fmt.Fprint(&builder, "-")
			fmt.Fprint(&builder, addRune(runeRange.End))
		}
	}

	return builder.String()
}

func addRune(r rune) string {
	if unicode.IsPrint(r) {
		return string(r)
	}

	return fmt.Sprintf("0x%02x", r)
}

///////////////////////////////////////////////////////////////////////////////
// ReservedFilenames
///////////////////////////////////////////////////////////////////////////////

type ReservedFilenames []string

var (
	ReservedFilenamesNone ReservedFilenames = []string{}
	// Source: https://learn.microsoft.com/en-us/windows/win32/fileio/naming-a-file
	// See also https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-createfilea#consoles

	ReservedFilenamesNTFS ReservedFilenames = []string{
		`^((CON|PRN|AUX|NUL|COM[1-9¹²³]|LPT[1-9¹²³])(\..*)?|CONIN\$|CONOUT\$)$`,
	}
	ReservedFilenamesNTFSSuffixes ReservedFilenames = []string{
		`(\. )$`,
	}
)

func (v Feature) String() string {
	s, ok := featuresMap[v]
	if ok {
		return s
	}

	return fmt.Sprintf("unknown Feature value: %d", v)
}

func (v OSFeature) String() string {
	s, ok := osFeatureMap[v]
	if ok {
		return s
	}

	return fmt.Sprintf("unknown OSFeature value: %d", v)
}

func si(f float64) string {
	return strings.ReplaceAll(humanize.SIWithDigits(f, 2, ""), " ", "") //nolint:mnd
}

func NewFilesystem(name string) Filesystem {
	var fs Filesystem

	fs.Name = name
	fs.Features = make(map[Feature]Availability)
	fs.OSFeatures = make(map[OSFeature]Availability)

	return fs
}

func (a Filesystem) String() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "MaxNameLength       : %v\n", a.MaxNameLength)
	fmt.Fprintf(&builder, "MaxPathLength       : %v\n", a.MaxPathLength)
	fmt.Fprintf(&builder, "CharSize            : %v\n", a.CharSize)
	fmt.Fprintf(&builder, "MaxFileSize         : %10v\n", si(a.MaxFileSize.Float))
	fmt.Fprintf(&builder, "MaxVolumeSize       : %10v\n", si(a.MaxVolumeSize.Float))
	fmt.Fprintf(&builder, "MaxFiles            : %10v\n", si(float64(a.MaxFiles)))

	fmt.Fprintf(&builder, "ATimeGranularity    : %v\n", a.ATimeGranularity)
	fmt.Fprintf(&builder, "BTimeGranularity    : %v\n", a.BTimeGranularity)
	fmt.Fprintf(&builder, "CTimeGranularity    : %v\n", a.CTimeGranularity)
	fmt.Fprintf(&builder, "MTimeGranularity    : %v\n", a.MTimeGranularity)

	fmt.Fprintf(&builder, "CaseMapping         : %v\n", a.CaseMapping)
	fmt.Fprintf(&builder, "UnicodeNormalization: %v\n", a.UnicodeNormalization)

	// Sort of FAT/NTFS/ReFS specific:
	fmt.Fprintf(&builder, "DisallowedRunes     : %v\n", a.DisallowedRunes)
	fmt.Fprintf(&builder, "ReservedFilenames   : %v\n", a.ReservedFilenames)

	fmt.Fprintf(&builder, "Features            : %v\n", a.Features)
	fmt.Fprintf(&builder, "OSFeatures          : %v\n", a.OSFeatures)

	return builder.String()
}

// @TODO use this mapping
// fatLikeFS marks Filesystems that disallow Windows-reserved characters
// such as " * / : < > ? \ | and DEL (0x7F).
var fatLikeFS = map[string]bool{ //nolint:unused
	"cifs":   true, // Windows SMB/CIFS network Filesystem (Windows-compatible)
	"exfat":  true, // Microsoft exFAT for flash storage (FAT32 successor)
	"fat":    true, // Classic FAT12/16/32 Filesystem
	"fatx":   true, // Xbox FAT variant with stricter naming limits (no Linux magic)
	"msdos":  true, // FAT12/16 using 8.3 short filenames (no VFAT)
	"novell": true, // NetWare Core Protocol (NCP) network Filesystem (Windows/DOS semantics)
	"ntfs":   true, // Windows NT Filesystem with Unicode long names
	"smb":    true, // SMB network mount (Windows file sharing)
	"vfat":   true, // Classic FAT12/16/32 Filesystem
}

// @TODO use this mapping
// maybeFatLikeFS marks Filesystems that might disallow Windows-reserved
// characters (" * / : < > ? \ | and DEL) depending on backend or configuration.
var maybeFatLikeFS = map[string]bool{ //nolint:unused
	"fuseblk": false, // Backend may be ntfs-3g, exfat-fuse, etc. — restrictions depend on module
	"fusectl": false, // Rarely used for Filesystems, but could host Windows-limited FUSE daemons
	"gpfs":    false, // IBM GPFS can enforce Windows-safe names when shared with SMB clients
	"gfs":     false, // Cluster FS used in mixed Windows/Linux environments; optional enforcement
	"panfs":   false, // Panasas FS; Windows interop mode applies FAT/NT-like restrictions
	"zfs":     false, // When shared via SMB or vfs_fruit, may enforce Windows-safe charset
}

var featuresMap = map[Feature]string{
	FeatureStoresFileOwner:              "StoresFileOwner",
	FeaturePOSIXPermissions:             "POSIXPermissions",
	FeatureCreationTimestamps:           "CreationTimestamps",
	FeatureLastAccessTimestamps:         "LastAccessTimestamps",
	FeatureLastMetadataChangeTimestamps: "LastMetadataChangeTimestamps",
	FeatureLastArchiveTimestamps:        "LastArchiveTimestamps",
	FeatureACLs:                         "ACLs",
	FeatureMACLabels:                    "MACLabels",
	FeatureExtendedFilesystem:           "ExtendedFilesystem",
	FeatureMetadataChecksums:            "MetadataChecksums",
	FeatureHardLinks:                    "HardLinks",
	FeatureSymbolicLinks:                "SymbolicLinks",
	FeatureBlockJournaling:              "BlockJournaling",
	FeatureMetadataOnlyJournaling:       "MetadataOnlyJournaling",
	FeatureCaseSensitive:                "CaseSensitive",
	FeatureCasePreserving:               "CasePreserving",
	FeatureXIP:                          "XIP",
	FeatureSnapshots:                    "Snapshots",
	FeatureEncryption:                   "Encryption",
	FeatureDeduplication:                "Deduplication",
	FeatureDataChecksums:                "DataChecksums",
	FeaturePersistentCache:              "PersistentCache",
	FeatureMultipleDevices:              "MultipleDevices",
	FeatureCompression:                  "Compression",
	FeatureSelfHealing:                  "SelfHealing",
	FeatureOfflineGrow:                  "OfflineGrow",
	FeatureOnlineGrow:                   "OnlineGrow",
	FeatureOfflineShrink:                "OfflineShrink",
	FeatureOnlineShrink:                 "OnlineShrink",
	FeatureAddRemoveVolumes:             "AddRemoveVolumes",
	FeatureSparseFiles:                  "SparseFiles",
	FeatureBlockSuballocation:           "BlockSuballocation",
	FeatureTailPacking:                  "TailPacking",
	FeatureExtents:                      "Extents",
	FeatureVariableFileBlockSize:        "VariableFileBlockSize",
	FeatureAllocateOnFlush:              "AllocateOnFlush",
	FeatureCopyOnWrite:                  "CopyOnWrite",
	FeatureTrimSupport:                  "TrimSupport",
	FeatureIsCompressed:                 "IsCompressed",
	FeatureIsEncrypted:                  "IsEncrypted",
	FeatureIsReadOnly:                   "IsReadOnly",
}

const (
	FilesystemAPFS     = "apfs"  // APFS
	FilesystemBtrfs    = "btrfs" // Btrfs
	FilesystemExFAT    = "exfat" // exFAT
	FilesystemExt2     = "ext2"
	FilesystemExt3     = "ext3"
	FilesystemExt4     = "ext4"
	FilesystemFAT32    = "fat32"    // FAT32
	FilesystemF2FS     = "f2fs"     // F2FS
	FilesystemHFSPlus  = "hfs+"     // HFS+
	FilesystemNTFS     = "ntfs"     // NTFS
	FilesystemReFS     = "refs"     // ReFS
	FilesystemReiserFS = "reiserfs" // ReiserFS
	FilesystemSquashFS = "squashfs" // SquashFS
	FilesystemUDF      = "udf"      // UDF
	FilesystemXFS      = "xfs"      // XFS
	FilesystemZFS      = "zfs"      // ZFS
)

var filesystemMap = map[string]Filesystem{}

func init() { //nolint:funlen
	// @TODO complete this table
	featuresExt4 := map[Feature]Availability{
		FeatureStoresFileOwner:              AvailabilityAlways,
		FeaturePOSIXPermissions:             AvailabilityAlways,
		FeatureCreationTimestamps:           AvailabilityAlways,
		FeatureLastAccessTimestamps:         AvailabilityAlways,
		FeatureLastMetadataChangeTimestamps: AvailabilityAlways,
	}

	// Source: https://en.wikipedia.org/wiki/Comparison_of_file_systems#Limits

	var (
		// https://en.wikipedia.org/wiki/Apple_File_System
		filesystemAPFS = Filesystem{
			MaxNameLength:        MaxChars255,
			MaxPathLength:        0, // No limit defined
			CharSize:             CharSizeUTF8,
			MaxFileSize:          Size{0, 0, 0},                  // ?
			MaxVolumeSize:        Size{EiB << 3, 2, EiBPow2 + 3}, // 8 EiB
			MaxFiles:             1 << 63,                        //nolint:mnd
			UnicodeNormalization: UnicodeNormalizationNFD,
			DisallowedRunes:      DisallowedRunesNULSlash,
		}
		//
		filesystemBtrfs = Filesystem{
			MaxNameLength:   MaxChars255,
			MaxPathLength:   0, // No limit defined
			CharSize:        CharSizeByte,
			MaxFileSize:     Size{TiB << 4, 2, TiBPow2 + 4}, // 16 TiB
			MaxVolumeSize:   Size{EiB << 4, 2, EiBPow2 + 4}, // 16 EiB
			MaxFiles:        MaxUint64,
			DisallowedRunes: DisallowedRunesNULSlash,
		}
		filesystemExFAT = Filesystem{
			MaxNameLength:   MaxChars255,
			MaxPathLength:   32760, //nolint:mnd
			CharSize:        CharSizeUTF16,
			MaxFileSize:     Size{EiB, 2, EiB},              // 1 EiB
			MaxVolumeSize:   Size{ZiB << 6, 2, ZiBPow2 + 6}, // 64 ZiB
			MaxFiles:        0,                              // ?
			DisallowedRunes: DisallowedRunesFAT,
		}
		filesystemExt2 = Filesystem{
			MaxNameLength:   MaxChars255,
			MaxPathLength:   0, // No limit defined
			CharSize:        CharSizeByte,
			MaxFileSize:     Size{TiB << 1, 2, TiBPow2 + 1}, // 2 TiB
			MaxVolumeSize:   Size{TiB << 5, 2, TiBPow2 + 5}, // 32 TiB
			MaxFiles:        0,                              // ?
			DisallowedRunes: DisallowedRunesNULSlash,
		}
		filesystemExt3 = Filesystem{
			MaxNameLength:   MaxChars255,
			MaxPathLength:   0, // No limit defined
			CharSize:        CharSizeByte,
			MaxFileSize:     Size{TiB << 1, 2, TiBPow2 + 1}, // 2 TiB
			MaxVolumeSize:   Size{TiB << 5, 2, TiBPow2 + 5}, // 32 TiB
			MaxFiles:        0,                              // ?
			DisallowedRunes: DisallowedRunesNULSlash,
		}
		filesystemExt4 = Filesystem{
			MaxNameLength: MaxChars255,
			MaxPathLength: 0, // No limit defined
			CharSize:      CharSizeByte,
			MaxFileSize:   Size{TiB << 4, 2, TiBPow2 + 4}, // 16 TiB
			MaxVolumeSize: Size{EiB, 2, EiBPow2},          // 1 EiB
			MaxFiles:      MaxUint32 + 1,

			ATimeGranularity: time.Duration(1), // nanoseconds
			BTimeGranularity: time.Duration(1),
			CTimeGranularity: time.Duration(1),
			MTimeGranularity: time.Duration(1),

			CaseMapping:          CaseMappingUnknown,
			UnicodeNormalization: UnicodeNormalizationNone,
			Features:             featuresExt4,
			OSFeatures:           OSFeaturesNone,
			DisallowedRunes:      DisallowedRunesNULSlash,
			ReservedFilenames:    ReservedFilenamesNone,
		}
		// @TODO FAT.
		filesystemFAT32 = Filesystem{
			MaxNameLength:   MaxChars255,
			MaxPathLength:   32760, //nolint:mnd
			CharSize:        CharSizeUTF16,
			MaxFileSize:     Size{GiB << 2, 2, GiBPow2 + 2}, // 4 GiB
			MaxVolumeSize:   Size{TiB << 4, 2, TiBPow2 + 4}, // 16 TiB
			MaxFiles:        0,                              // ?
			DisallowedRunes: DisallowedRunesFAT,
		}
		filesystemF2FS = Filesystem{
			MaxNameLength:   MaxChars255,
			MaxPathLength:   0, // No limit defined
			CharSize:        CharSizeByte,
			MaxVolumeSize:   Size{TiB << 6, 2, TiBPow2 + 6}, // 64 TiB
			MaxFileSize:     Size{1059 << 42, 1059, 42},     // 4,228,213,756 KiB
			MaxFiles:        0,                              // ?
			DisallowedRunes: DisallowedRunesNULSlash,
		}
		filesystemHFSPlus = Filesystem{
			MaxNameLength:   MaxChars255,
			MaxPathLength:   0, // No limit defined
			CharSize:        CharSizeUTF16,
			MaxVolumeSize:   Size{9.223 * EB, 0, 0}, // slightly less than 8 EiB
			MaxFileSize:     Size{9.223 * EB, 0, 0}, // slightly less than 8 EiB
			MaxFiles:        MaxUint32,              // https://en.wikipedia.org/wiki/HFS_Plus
			DisallowedRunes: DisallowedRunesNULColon,
		}
		// @TODO HFS+J.
		// @TODO HFSX.
		// @TODO JHFS+.
		// @TODO JHFS+X.

		// @TODO ISO9660.
		// @TODO Joliet.
		filesystemNTFS = Filesystem{
			MaxNameLength:   MaxChars255,
			MaxPathLength:   MaxChars32767,
			CharSize:        CharSizeUTF16,
			MaxVolumeSize:   Size{PiB << 3, 2, PiBPow2 + 3}, // 8 PiB
			MaxFileSize:     Size{PiB << 3, 2, PiBPow2 + 3}, // 8 PiB
			MaxFiles:        MaxUint32 + 1,
			DisallowedRunes: DisallowedRunesFAT,
		}
		filesystemReFS = Filesystem{
			MaxNameLength:   MaxChars255,
			MaxPathLength:   MaxChars32767,
			CharSize:        CharSizeUTF16,
			MaxVolumeSize:   Size{EiB << 4, 2, EiBPow2 + 4}, // 16 EiB
			MaxFileSize:     Size{YiB, 2, YiBPow2},          // 1 YiB
			MaxFiles:        0,                              // ?
			DisallowedRunes: DisallowedRunesFAT,
		}
		filesystemReiserFS = Filesystem{
			MaxNameLength:   4032, //nolint:mnd
			MaxPathLength:   0,    // no limit defined
			CharSize:        CharSizeByte,
			MaxVolumeSize:   Size{TiB << 3, 2, TiBPow2 + 3}, // 8 TiB
			MaxFileSize:     Size{TiB << 4, 2, TiBPow2 + 4}, // 16 YiB
			MaxFiles:        0,                              // ?
			DisallowedRunes: DisallowedRunesNULSlash,
		}
		// @TODO RockRidge.
		filesystemSquashFS = Filesystem{
			MaxNameLength:   MaxChars255 + 1, // 256
			MaxPathLength:   0,               // no limit defined
			CharSize:        CharSizeByte,
			MaxFileSize:     Size{EiB << 4, 2, EiBPow2 + 4}, // 16 EiB
			MaxVolumeSize:   Size{EiB << 4, 2, EiBPow2 + 4}, // 16 EiB
			MaxFiles:        0,                              // ?
			DisallowedRunes: DisallowedRunesNULSlash,
		}
		filesystemUDF = Filesystem{
			MaxNameLength:   MaxChars255,
			MaxPathLength:   1023, //nolint:mnd
			CharSize:        CharSizeByte,
			MaxFileSize:     Size{EiB << 4, 2, EiBPow2 + 4}, // 16 EiB (with sparse files)
			MaxVolumeSize:   Size{TiB << 3, 2, TiBPow2 + 3}, // 16 TiB
			MaxFiles:        0,                              // ?
			DisallowedRunes: DisallowedRunesNUL,
		}
		filesystemXFS = Filesystem{
			MaxNameLength:   MaxChars255,
			MaxPathLength:   0, // no limit defined
			CharSize:        CharSizeByte,
			MaxFileSize:     Size{EiB << 3, 2, EiBPow2 + 3}, // 8 EiB
			MaxVolumeSize:   Size{EiB << 3, 2, EiBPow2 + 3}, // 8 EiB
			MaxFiles:        0,                              // ?
			DisallowedRunes: DisallowedRunesNULSlash,
		}
		filesystemZFS = Filesystem{
			MaxNameLength:   1023, //nolint:mnd
			MaxPathLength:   0,    // no limit defined
			CharSize:        CharSizeByte,
			MaxFileSize:     Size{EiB << 4, 2, EiBPow2 + 4}, // 16 EiB
			MaxVolumeSize:   Size{YiB << 48, 2, 128},        // 281,474,976,710,656 YiB (2^128)
			MaxFiles:        0,                              // ?
			DisallowedRunes: DisallowedRunesNULSlash,
		}
	)

	filesystemMap = map[string]Filesystem{
		FilesystemAPFS:     filesystemAPFS,
		FilesystemBtrfs:    filesystemBtrfs,
		FilesystemExFAT:    filesystemExFAT,
		FilesystemExt2:     filesystemExt2,
		FilesystemExt3:     filesystemExt3,
		FilesystemExt4:     filesystemExt4,
		FilesystemFAT32:    filesystemFAT32,
		FilesystemF2FS:     filesystemF2FS,
		FilesystemHFSPlus:  filesystemHFSPlus,
		FilesystemNTFS:     filesystemNTFS,
		FilesystemReFS:     filesystemReFS,
		FilesystemReiserFS: filesystemReiserFS,
		FilesystemSquashFS: filesystemSquashFS,
		FilesystemUDF:      filesystemUDF,
		FilesystemXFS:      filesystemXFS,
		FilesystemZFS:      filesystemZFS,
	}
}
