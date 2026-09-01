// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"errors"
	"io"
	"time"
)

const (
	// ReadOnlyModeSet set a file's read-only attribute, if the specified
	// perm FileMode has the user writable bit (0o200) set. Otherwise, it will
	// resets (clears) it. (Windows only).
	// Deprecated.
	//
	// This constant will be removed in a future release.
	ReadOnlyModeSet = ReadOnlyModeFromPermissions
	// ReadOnlyModeReset does not set a file's read-only attribute, and if it's
	// set, it resets (clears) it. (Windows only).
	// Deprecated.
	//
	// This constant will be removed in a future release.
	ReadOnlyModeReset = ReadOnlyModeClear

	// UserIDSourceIsInt defines if the OS uses an int to identify the user.
	// Deprecated.
	//
	// This constant will be removed in a future release.
	UserIDSourceIsInt = UserIDSourceIsInteger
)

// IsUnsupportedError returns true if err indicates the function is unsupported.
// Deprecated.
//
// This function will be removed in a future release.
func IsUnsupportedError(err error) bool {
	err = normalizeUnsupportedError(err)

	return errors.Is(err, errors.ErrUnsupported)
}

// WithFlags sets the flag option.
// Used by the Create, Open, OpenFile, WriteFile, and WriteReader functions.
// Deprecated: Use WithOpenFlags() instead.
//
// This function will be removed in a future release.
func WithFlags(flags int) Option {
	return WithOpenFlags(flags)
}

// WithRetrySeconds sets the retry timeout option in seconds. The default is 0
// which means to not retry at all.
// Used by the Rename and RemoveAll functions.
// Deprecated: Use WithRetryTimeout() instead.
//
// This function will be removed in a future release.
func WithRetrySeconds(seconds float64) Option {
	return WithRetryTimeout(time.Duration(seconds * float64(time.Second)))
}

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
func WriteFileAtomic(filename string, data []byte, fns ...Option) error {
	fns = append(fns, WithAtomicity(true))

	return WriteFile(filename, data, CreatePerm, fns...)
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
func WriteReaderAtomic(filename string, r io.Reader, fns ...Option) error {
	fns = append(fns, WithAtomicity(true))

	return WriteReader(filename, r, CreatePerm, fns...)
}

// Deprecated: Use GoVersion() instead.
//
// This function will be removed in a future release.
var UnderlyingGoVersion = GoVersion
