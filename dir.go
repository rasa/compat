// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Source: https://github.com/golang/go/blob/77f911e3/src/io/fs/fs.go#L91-L113

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

// Inspired by: https://github.com/golang/go/blob/77f911e3/src/os/file_unix.go#L446-L486

type dirEntry struct {
	parent string
	name   string
	typ    os.FileMode
	osInfo os.FileInfo
	info   FileInfo
	infoed bool
	err    error
}

func (d dirEntry) IsDir() bool {
	return d.typ.IsDir()
}

func (d dirEntry) Type() os.FileMode {
	if d.err != nil {
		return 0
	}

	return d.typ
}

func (d dirEntry) Info() (FileInfo, error) {
	if d.infoed {
		return d.info, d.err
	}
	d.infoed = true //nolint:staticcheck
	path := d.name
	if d.parent != "" {
		path = filepath.Join(d.parent, d.name)
	}
	if d.osInfo == nil {
		// WalkDir doesn't follow symlinks
		d.osInfo, d.err = os.Lstat(path)
		if d.err != nil {
			return nil, d.err
		}
	}
	d.info, d.err = stat(d.osInfo, path, false)
	if d.err != nil {
		return nil, d.err
	}
	d.typ = d.info.Mode().Type() //nolint:govet,staticcheck

	return d.info, nil
}

func (d dirEntry) Name() string {
	return d.name
}

func (d dirEntry) String() string {
	return FormatDirEntry(d)
}

// Source: https://github.com/golang/go/blob/77f911e3/src/os/dir.go#L87

// ReadDir reads the named directory,
// returning all its directory entries sorted by filename.
// If an error occurs reading the directory,
// ReadDir returns the entries it was able to read before the error,
// along with the error.
func ReadDir(name string) ([]DirEntry, error) {
	osDirs, err := os.ReadDir(name)
	if err != nil {
		return []DirEntry{}, err
	}
	dirs := make([]DirEntry, len(osDirs))
	for i, dir := range osDirs {
		dirs[i] = osDirEntryToDirEntry(dir, name)
	}
	slices.SortFunc(dirs, func(a, b DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})

	return dirs, nil
}

// Inspired by: https://github.com/golang/go/blob/77f911e3/src/io/fs/readdir.go#L77-L84

// FileInfoToDirEntry returns a [DirEntry] that returns information from info.
// If info is nil, FileInfoToDirEntry returns nil.
func FileInfoToDirEntry(info FileInfo, parent string) DirEntry {
	if info == nil {
		return nil
	}

	return dirEntry{
		parent: parent,
		name:   info.Name(),
		typ:    info.Mode().Type(),
		info:   info,
	}
}

func osDirEntryToDirEntry(entry os.DirEntry, parent string) DirEntry {
	if entry == nil {
		return nil
	}

	return dirEntry{
		parent: parent,
		name:   entry.Name(),
		typ:    entry.Type(),
	}
}

func fsDirEntryToDirEntry(entry fs.DirEntry, parent string) DirEntry {
	if entry == nil {
		return nil
	}

	return dirEntry{
		parent: parent,
		name:   entry.Name(),
		typ:    entry.Type(),
	}
}

func fsFileInfoToDirEntry(info fs.FileInfo, parent string) DirEntry {
	if info == nil {
		return nil
	}

	return dirEntry{
		parent: parent,
		name:   info.Name(),
		typ:    info.Mode().Type(),
		osInfo: info,
	}
}
