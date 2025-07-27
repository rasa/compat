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

// WriteReaderAtomic atomically writes the contents of r to the specified filename.
// The target file is guaranteed to be either fully written, or not written at all.
// WriteReaderAtomic overwrites any file that exists at the location (but only if
// the write fully succeeds, otherwise the existing file is unmodified).
// Additional option arguments can be used to change the default configuration
// for the target file.
func WriteReaderAtomic(filename string, r io.Reader, opts ...Option) (err error) { //nolint:funlen,gocyclo // quiet linter
	// original behavior is to preserve the mode of an existing file.
	fopts := &FileOptions{
		keepFileMode: true,
	}
	for _, opt := range opts {
		opt(fopts)
	}

	// write to a temp file first, then we'll atomically replace the target file
	// with the temp file.
	dir, file := filepath.Split(filename)
	if dir == "" {
		dir = "."
	}

	f, err := CreateTemp(dir, file)
	if err != nil {
		return fmt.Errorf("cannot create temp file: %w", err)
	}
	name := f.Name()
	defer func() {
		if err != nil {
			// Don't leave the temp file lying around on error.
			_ = os.Remove(name) // yes, ignore the error, not much we can do about it.
		}
	}()
	// ensure we always close f. Note that this does not conflict with the
	// close below, as close is idempotent.
	defer f.Close()
	_, err = io.Copy(f, r)
	if err != nil {
		return fmt.Errorf("cannot write data to tempfile %q: %w", name, err)
	}
	// fsync is important, otherwise os.Rename could rename a zero-length file
	err = f.Sync()
	if err != nil {
		return fmt.Errorf("cannot flush tempfile %q: %w", name, err)
	}
	err = f.Close()
	if err != nil {
		return fmt.Errorf("cannot close tempfile %q: %w", name, err)
	}
	sourceInfo, err := Stat(name)
	if err != nil {
		return err
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
	// apply possible file mode change
	if fileMode != 0 && fileMode != sourceInfo.Mode() {
		err = Chmod(name, fileMode)
		if err != nil {
			return fmt.Errorf("cannot set permissions on tempfile %q: %w", name, err)
		}
	}
	err = Rename(name, filename)
	if err != nil {
		return fmt.Errorf("cannot replace %q with tempfile %q: %w", filename, name, err)
	}

	return nil
}
