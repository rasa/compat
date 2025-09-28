// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"bytes"
	"io"
)

// WriteFileAtomic atomically writes the contents of data to the specified filename.
// The target file is guaranteed to be either fully written, or not written at all.
// WriteFileAtomic overwrites any file that exists at the location (but only if
// the write fully succeeds, otherwise the existing file is unmodified).
// Additional option arguments can be used to change the default configuration
// for the target file.
//
// Deprecated: Use WriteFile() with WithAtomicity(true) instead.
//
// This function will be removed in a future release.
func WriteFileAtomic(filename string, data []byte, opts ...Option) (err error) {
	opts = append(opts, WithAtomicity(true))
	return writeReaderAtomic(filename, bytes.NewReader(data), opts...)
}

// WriteReaderAtomic atomically writes the contents of r to the specified filename.
// The target file is guaranteed to be either fully written, or not written at all.
// WriteReaderAtomic overwrites any file that exists at the location (but only if
// the write fully succeeds, otherwise the existing file is unmodified).
// Additional option arguments can be used to change the default configuration
// for the target file.
//
// Deprecated: Use WriteReader() with WithAtomicity(true) instead.
//
// This function will be removed in a future release.
func WriteReaderAtomic(filename string, r io.Reader, opts ...Option) (err error) { //nolint:funlen,gocyclo
	opts = append(opts, WithAtomicity(true))
	return writeReaderAtomic(filename, r, opts...)
}

// Deprecated: Use GoVersion() instead.
//
// This function will be removed in a future release.
var UnderlyingGoVersion = GoVersion
