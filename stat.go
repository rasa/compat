// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"
)

// UnknownID is returned when the UID or GID could not be determined.
const UnknownID = int(-1)

// UserIDSourceType defines if the underlying source for the user's ID, is an
// int (UID() function), or a string (User() function).
type UserIDSourceType uint

const (
	// UserIDSourceUnknown is unused.
	UserIDSourceUnknown UserIDSourceType = iota
	// UserIDSourceIsInteger defines if the OS uses an int to identify the user.
	UserIDSourceIsInteger
	// UserIDSourceIsString defines if the OS uses a string to identify the
	// user.
	UserIDSourceIsString
	// UserIDSourceIsSID defines if the OS uses a SID to identify the user.
	UserIDSourceIsSID
	// UserIDSourceIsNone defines if the OS does not provide user's IDs, so
	// we provide sane defaults instead.
	UserIDSourceIsNone
)

// A FileInfo describes a file and is returned by [Stat].
// See https://github.com/golang/go/blob/ad7a6f81/src/io/fs/fs.go#L158
type FileInfo interface { //nolang:interfacebloat
	fs.FileInfo

	ATime() time.Time // last accessed time, or 0 if unsupported
	BTime() time.Time // created (birthed) time, or 0 if unsupported
	CTime() time.Time // status/metadata changed time, or 0 if unsupported
	// MTime() time.Time    // last modified time (alias)
	Links() uint         // number of hard links, or 0 if unsupported
	UID() int            // user ID, or -1 if an error or unsupported
	GID() int            // group ID, or -1 if an error or unsupported
	User() string        // user name, or "" if an error or unsupported
	Group() string       // group name, or "" if an error or unsupported
	PartitionID() uint64 // unique disk partition ID
	FileID() uint64      // unique file ID (on a specific partition)
	Error() error        // error result of the last system call that failed
	String() string
}

func (fs *fileStat) Name() string       { return fs.name }
func (fs *fileStat) Size() int64        { return fs.size }
func (fs *fileStat) Mode() os.FileMode  { return fs.mode }
func (fs *fileStat) ModTime() time.Time { return fs.mtime }
func (fs *fileStat) IsDir() bool        { return fs.mode.IsDir() }
func (fs *fileStat) Sys() any           { return &fs.sys }

func (fs *fileStat) ATime() time.Time    { return fs.atime }
func (fs *fileStat) Links() uint         { return fs.links }
func (fs *fileStat) PartitionID() uint64 { return fs.partID }
func (fs *fileStat) FileID() uint64      { return fs.fileID }
func (fs *fileStat) Error() error        { return fs.err }

func (fs *fileStat) String() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "Name:   %v\n", fs.Name())
	fmt.Fprintf(&builder, "Size:   %v\n", fs.Size())
	fmt.Fprintf(&builder, "Mode:   0o%o (%v)\n", fs.Mode(), fs.Mode())
	fmt.Fprintf(&builder, "ModTime:%v\n", fs.ModTime())
	fmt.Fprintf(&builder, "ATime:  %v\n", fs.ATime())
	fmt.Fprintf(&builder, "BTime:  %v\n", fs.BTime())
	fmt.Fprintf(&builder, "CTime:  %v\n", fs.CTime())
	fmt.Fprintf(&builder, "IsDir:  %v\n", fs.IsDir())
	fmt.Fprintf(&builder, "Links:  %v\n", fs.Links())
	fmt.Fprintf(&builder, "UID:    %v (%v)\n", fs.UID(), fs.User())
	fmt.Fprintf(&builder, "GID:    %v (%v)\n", fs.GID(), fs.Group())
	fmt.Fprintf(&builder, "PartID: %v\n", fs.PartitionID())
	fmt.Fprintf(&builder, "FileID: %v\n", fs.FileID())

	return builder.String()
}

// Fstat returns a [FileInfo] structure describing file.
// If there is an error, it will be of type *PathError.
func Fstat(f *os.File) (FileInfo, error) {
	return fstat(f)
}

// Lstat returns a [FileInfo] describing the named file.
// If the file is a symbolic link, the returned FileInfo
// describes the symbolic link. Lstat makes no attempt to follow the link.
// If there is an error, it will be of type [*os.PathError].
//
// On Windows, if the file is a reparse point that is a surrogate for another
// named entity (such as a symbolic link or mounted folder), the returned
// FileInfo describes the reparse point, and makes no attempt to resolve it.
func Lstat(name string) (FileInfo, error) {
	fi, err := os.Lstat(name)
	if err != nil {
		return nil, err
	}

	return stat(fi, name, false)
}

// SameFile reports whether fi1 and fi2 describe the same file. For example,
// on Unix this means that the partition (device) and inode fields of the two
// underlying structures are identical; on other systems the decision may be
// based on the path names.
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
func SameFiles(name1, name2 string) (bool, error) {
	fi1, err := Stat(name1)
	if err != nil {
		return false, err
	}

	fi2, err := Stat(name2)
	if err != nil {
		return false, err
	}

	return SameFile(fi1, fi2), nil
}

// SamePartition reports whether fi1 and fi2 describe files on the same disk
// partition. For example, on Unix this means that the partition (device) fields
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

// SamePartitions reports whether name1 and name2 are files on the same disk
// partition.
// The function follow symlinks.
func SamePartitions(name1, name2 string) (bool, error) {
	fi1, err := Stat(name1)
	if err != nil {
		return false, err
	}

	fi2, err := Stat(name2)
	if err != nil {
		return false, err
	}

	return SamePartition(fi1, fi2), nil
}

// Stat returns a [FileInfo] describing the named file.
// If there is an error, it will be of type [*os.PathError].
func Stat(name string) (FileInfo, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return nil, err
	}

	return stat(fi, name, true)
}
