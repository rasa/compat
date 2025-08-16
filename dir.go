// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"os"
	"slices"
	"strings"
)

// Source: https://github.com/golang/go/blob/77f911e3/src/io/fs/fs.go#L91

// A DirEntry is an entry read from a directory
// (using the [ReadDir] function or a [ReadDirFile]'s ReadDir method).
type DirEntry interface {
	// Name returns the name of the file (or subdirectory) described by the entry.
	// This name is only the final element of the path (the base name), not the entire path.
	// For example, Name would return "hello.go" not "home/gopher/hello.go".
	Name() string

	// IsDir reports whether the entry describes a directory.
	IsDir() bool

	// Type returns the type bits for the entry.
	// The type bits are a subset of the usual FileMode bits, those returned by the FileMode.Type method.
	Type() os.FileMode

	// Info returns the FileInfo for the file or subdirectory described by the entry.
	// The returned FileInfo may be from the time of the original directory read
	// or from the time of the call to Info. If the file has been removed or renamed
	// since the directory read, Info may return an error satisfying errors.Is(err, ErrNotExist).
	// If the entry denotes a symbolic link, Info reports the information about the link itself,
	// not the link's target.
	Info() (FileInfo, error)
}

// Source: https://github.com/golang/go/blob/77f911e3/src/io/fs/readdir.go#L52

// dirInfo is a DirEntry based on a FileInfo.
type dirInfo struct {
	fileInfo FileInfo
}

func (di dirInfo) IsDir() bool {
	return di.fileInfo.IsDir()
}

func (di dirInfo) Type() os.FileMode {
	return di.fileInfo.Mode().Type()
}

func (di dirInfo) Info() (FileInfo, error) {
	return di.fileInfo, nil
}

func (di dirInfo) Name() string {
	return di.fileInfo.Name()
}

func (di dirInfo) String() string {
	return FormatDirEntry(di)
}

// Source: https://github.com/golang/go/blob/77f911e3/src/os/dir.go#L87

// ReadDir reads the named directory,
// returning all its directory entries sorted by filename.
// If an error occurs reading the directory,
// ReadDir returns the entries it was able to read before the error,
// along with the error.
func ReadDir(name string) ([]DirEntry, error) {
	dirs, err := os.ReadDir(name)
	if err != nil {
		return nil, err
	}
	var rv = make([]DirEntry, len(dirs))
	for i, dir := range dirs {
		fi, err := dir.Info()
		if err != nil {
			return nil, err
		}
		cfi, err := stat(fi, name, false)
		if err != nil {
			return nil, err
		}
		rv[i] = dirInfo{cfi}
	}
	slices.SortFunc(rv, func(a, b DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})

	return rv, nil
}

// Source: https://github.com/golang/go/blob/77f911e3/src/io/fs/readdir.go#L77

// FileInfoToDirEntry returns a [DirEntry] that returns information from info.
// If info is nil, FileInfoToDirEntry returns nil.
func FileInfoToDirEntry(info FileInfo) DirEntry {
	if info == nil {
		return nil
	}
	return dirInfo{fileInfo: info}
}
