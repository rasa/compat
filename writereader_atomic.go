// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: Portions copyright (c) 2015 Nate Finch (@natefinch)
// SPDX-FileCopyrightText: Portions copyright (c) 2022 Simon Dassow (@sdassow)

package compat

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// WriteReader writes r to the named file, creating it if necessary.
// If the file does not exist, WriteReader creates it using perm's permissions
// bits (before umask); otherwise WriteReader truncates it before writing,
// without changing permissions. Since WriteReader requires multiple system
// calls to complete, a failure mid-operation can leave the file in a partially
// written state. Use WriteReader() with the WithAtomicity(true) options,
// if this is a concern.
//
// When WithAtomicity(true) is passed, WriteReader atomically writes the
// contents of r to the specified filename. The target file is guaranteed to be
// either fully written, or not written at all. WriteReader overwrites any file
// that exists at the location (but only if the write fully succeeds, otherwise
// the existing file is unmodified).
//
// If perm is zero, then 0o666 is used, as this is what the os.Create() function
// uses. If both perm, and WithFileMode(perm) are provided, WithFileMode(perm)
// takes precedence.
//
// Additional option arguments can be used to change the default configuration
// for the target file.
func WriteReader(name string, r io.Reader, perm os.FileMode, opts ...Option) error {
	if perm.Perm() == 0 {
		perm |= CreatePerm // 0o666
	}

	fopts := Options{
		flags:    os.O_CREATE | os.O_WRONLY | os.O_TRUNC,
		fileMode: perm,
	}

	for _, opt := range opts {
		opt(&fopts)
	}

	if !fopts.atomically {
		if IsWindows {
			if fopts.readOnlyMode != ReadOnlyModeSet {
				fopts.flags |= O_FILE_FLAG_NO_RO_ATTR
			}
		}

		return writeReader(name, r, fopts.flags, fopts.fileMode)
	}

	return writeReaderAtomic(name, r, opts...)
}

func writeReaderAtomic(filename string, r io.Reader, opts ...Option) (err error) { //nolint:funlen,gocyclo
	// original behavior is to preserve the mode of an existing file.
	fopts := Options{
		keepFileMode: true,
	}

	for _, opt := range opts {
		opt(&fopts)
	}

	var fileMode os.FileMode
	// change default file mode for when file does not exist yet.
	if fopts.defaultFileMode != 0 {
		fileMode = fopts.defaultFileMode
	}

	// get the file mode from the original file and use that for the replacement
	// file, too.
	if fopts.keepFileMode {
		var destInfo os.FileInfo

		destInfo, err = Stat(filename)
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		if destInfo != nil {
			fileMode = destInfo.Mode()
		}
	}
	// given file mode always takes precedence
	if fopts.fileMode != 0 {
		fileMode = fopts.fileMode
	}

	if fileMode.Perm() == 0 {
		fileMode = CreatePerm
	}

	if IsWindows {
		if fopts.readOnlyMode != ReadOnlyModeSet {
			fopts.flags |= O_FILE_FLAG_NO_RO_ATTR
		}
	}

	// write to a temp file first, then we'll atomically replace the target file
	// with the temp file.
	dir, _ := filepath.Split(filename)
	if dir == "" {
		dir = "."
	}

	f, err := createTemp(dir, "~*.tmp", fileMode, fopts.flags)
	if err != nil {
		err = fmt.Errorf("cannot create temp file: %w", err)
		return &os.PathError{Op: "write", Path: filename, Err: err}
	}

	name := f.Name()

	defer func() {
		if err != nil {
			// Don't leave the temp file lying around on error.
			_ = Chmod(name, CreateTempPerm) // 0o600
			_ = Remove(name)
		}
	}()
	// ensure we always close f. Note that this does not conflict with the
	// close below, as close is idempotent.
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		err = fmt.Errorf("cannot write data to tempfile '%v': %w", name, err)
		return &os.PathError{Op: "write", Path: filename, Err: err}
	}
	// fsync is important, otherwise os.Rename could rename a zero-length file
	err = f.Sync()
	if err != nil {
		err = fmt.Errorf("cannot flush tempfile '%v': %w", name, err)
		return &os.PathError{Op: "write", Path: filename, Err: err}
	}

	err = f.Close()
	if err != nil {
		err = fmt.Errorf("cannot close tempfile '%v': %w", name, err)
		return &os.PathError{Op: "write", Path: filename, Err: err}
	}

	err = Rename(name, filename)
	if err != nil {
		err = fmt.Errorf("cannot replace '%v' with tempfile '%v': %w", filename, name, err)
		return &os.PathError{Op: "write", Path: filename, Err: err}
	}

	return nil
}

func writeReader(filename string, r io.Reader, flag int, perm os.FileMode) error {
	f, err := openFile(filename, flag, perm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return &os.PathError{Op: "write", Path: filename, Err: err}
	}

	err = f.Sync()
	if err != nil {
		return &os.PathError{Op: "write", Path: filename, Err: err}
	}

	return nil
}
