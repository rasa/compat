// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9

package compat

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
)

func rename(source, destination string, _ ...Option) error {
	sourceAbs, err := filepath.Abs(source)
	if err != nil {
		return &os.LinkError{
			Op:  "rename",
			Old: source,
			New: destination,
			Err: err,
		}
	}

	destinationAbs, err := filepath.Abs(destination)
	if err != nil {
		return &os.LinkError{
			Op:  "rename",
			Old: source,
			New: destination,
			Err: err,
		}
	}

	// Plan 9 wstat can rename only within the same directory.
	if filepath.Dir(sourceAbs) != filepath.Dir(destinationAbs) {
		return &os.LinkError{
			Op:  "rename",
			Old: source,
			New: destination,
			Err: os.ErrInvalid,
		}
	}

	// Atomic replacement of an existing destination is unavailable.
	//
	// Do not use os.Rename here: its Plan 9 implementation removes an
	// existing destination before attempting the wstat rename.
	if _, err := os.Stat(destinationAbs); err == nil {
		return &os.LinkError{
			Op:  "rename",
			Old: source,
			New: destination,
			Err: errors.ErrUnsupported,
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return &os.LinkError{
			Op:  "rename",
			Old: source,
			New: destination,
			Err: err,
		}
	}

	var dir syscall.Dir
	dir.Null()
	dir.Name = filepath.Base(destinationAbs)

	buf := make([]byte, syscall.STATFIXLEN+len(dir.Name))

	n, err := dir.Marshal(buf)
	if err != nil {
		return &os.LinkError{
			Op:  "rename",
			Old: source,
			New: destination,
			Err: err,
		}
	}

	// This is a single wstat transaction. If another process creates
	// destination after the Stat above, Wstat fails without removing it.
	if err := syscall.Wstat(sourceAbs, buf[:n]); err != nil {
		return &os.LinkError{
			Op:  "rename",
			Old: source,
			New: destination,
			Err: err,
		}
	}

	return nil
}
