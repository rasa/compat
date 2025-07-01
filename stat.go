// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"os"
	"time"
)

// SupportedType defines a bitmask that identifies if the OS supports specific
// fields, or not.
type SupportedType uint

const (
	// Links defines if FileInfo's Links() function is supported by the OS.
	// Links() returns the number of hard links to the file.
	Links SupportedType = 1 << iota
	// ATime defines if FileInfo's ATime() function is supported by the OS.
	// ATime() returns the time the file was last accessed.
	ATime
	// BTime defines if FileInfo's BTime() function is supported by the OS.
	// BTime() returns the time the file was created (or "birthed").
	BTime
	// CTime defines if FileInfo's CTime() function is supported by the OS.
	// CTime() returns the time the file's status/metadata was last changed.
	CTime
	// UID defines if FileInfo's UID() function is supported by the OS.
	// UID() returns the user ID of the file's owner.
	UID
	// GID defines if FileInfo's GID() function is supported by the OS.
	// GID() returns the group ID of the file's group.
	GID
	// UnknownID is returned when the UID or GID could not be determined.
	UnknownID = ^uint64(0)
)

// A FileInfo describes a file and is returned by [Stat].
// See https://github.com/golang/go/blob/ad7a6f81/src/io/fs/fs.go#L158
type FileInfo interface {
	Name() string        // base name of the file
	Size() int64         // length in bytes for regular files; system-dependent for others
	Mode() os.FileMode   // file mode bits
	ModTime() time.Time  // last modified time
	IsDir() bool         // abbreviation for Mode().IsDir()
	Sys() any            // underlying data source
	PartitionID() uint64 // unique disk partition ID
	FileID() uint64      // unique file ID (on a specific partition)
	Links() uint64       // number of hard links, or 0 if unsupported
	ATime() time.Time    // last accessed time, or 0 if unsupported
	BTime() time.Time    // created (birthed) time, or 0 if unsupported
	CTime() time.Time    // status/metadata changed time, or 0 if unsupported
	MTime() time.Time    // last modified time (alias)
	UID() uint64         // user ID, or 0 if unsupported
	GID() uint64         // group ID, or 0 if unsupported
}

func (fs *fileStat) Name() string        { return fs.name }
func (fs *fileStat) Size() int64         { return fs.size }
func (fs *fileStat) Mode() os.FileMode   { return fs.mode }
func (fs *fileStat) ModTime() time.Time  { return fs.mtime }
func (fs *fileStat) IsDir() bool         { return fs.mode.IsDir() }
func (fs *fileStat) Sys() any            { return &fs.sys }
func (fs *fileStat) PartitionID() uint64 { return fs.partID }
func (fs *fileStat) FileID() uint64      { return fs.fileID }
func (fs *fileStat) Links() uint64       { return fs.links }
func (fs *fileStat) ATime() time.Time    { return fs.atime }
func (fs *fileStat) BTime() time.Time    { return fs.btime }
func (fs *fileStat) CTime() time.Time    { return fs.ctime }
func (fs *fileStat) MTime() time.Time    { return fs.mtime } // duplicates ModTime
func (fs *fileStat) UID() uint64         { return fs.uid }
func (fs *fileStat) GID() uint64         { return fs.gid }

// Supports returns whether function is supported by the operating system.
func Supports(function SupportedType) bool {
	return supported&function == function
}

// Stat returns a [FileInfo] describing the named file.
// If there is an error, it will be of type [*PathError].
func Stat(name string) (FileInfo, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return nil, err
	}

	return stat(fi, name)
}

// Lstat returns a [FileInfo] describing the named file.
// If the file is a symbolic link, the returned FileInfo
// describes the symbolic link. Lstat makes no attempt to follow the link.
// If there is an error, it will be of type [*PathError].
//
// On Windows, if the file is a reparse point that is a surrogate for another
// named entity (such as a symbolic link or mounted folder), the returned
// FileInfo describes the reparse point, and makes no attempt to resolve it.
func Lstat(name string) (FileInfo, error) {
	fi, err := os.Lstat(name)
	if err != nil {
		return nil, err
	}

	return stat(fi, name)
}

// SamePartition reports whether fi1 and fi2 describe files on the same Partition.
// For example, on Unix this means that the Partition fields
// of the two underlying structures are identical; on other systems
// the decision may be based on the path names.
// SamePartition only applies to results returned by this package's [Stat].
// It returns false in other cases.
func SamePartition(fi1, fi2 FileInfo) bool {
	fs1, ok1 := fi1.(*fileStat)
	fs2, ok2 := fi2.(*fileStat)
	if !ok1 || !ok2 {
		return false
	}

	return fs1.partID == fs2.partID
}

// SamePartitions reports whether name1 and name2 are files on the same Partition.
// The function follow symlinks.
func SamePartitions(name1, name2 string) bool {
	fi1, err := Stat(name1)
	if err != nil {
		return false
	}
	fi2, err := Stat(name2)
	if err != nil {
		return false
	}

	return SamePartition(fi1, fi2)
}

// SameFile reports whether fi1 and fi2 describe the same file. For example,
// on Unix this means that the Partition and inode fields of the two underlying
// structures are identical; on other systems the decision may be based on the
// path names.
// SamePartition only applies to results returned by this package's [Stat].
// It returns false in other cases.
func SameFile(fi1, fi2 FileInfo) bool {
	fs1, ok1 := fi1.(*fileStat)
	fs2, ok2 := fi2.(*fileStat)
	if !ok1 || !ok2 {
		return false
	}

	return fs1.partID == fs2.partID && fs1.fileID == fs2.fileID
}

// SameFiles reports whether name1 and name2 are the same file.
// The function follow symlinks.
func SameFiles(name1, name2 string) bool {
	fi1, err := Stat(name1)
	if err != nil {
		return false
	}
	fi2, err := Stat(name2)
	if err != nil {
		return false
	}

	return SameFile(fi1, fi2)
}
