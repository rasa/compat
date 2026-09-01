// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: Portions copyright (c) 2015 Nate Finch (@natefinch)
// SPDX-FileCopyrightText: Portions copyright (c) 2022 Simon Dassow (@sdassow)

package compat

import (
	"fmt"
	"os"
	"strings"
)

// Options define the behavior of `WriteFile()`, etc.
type options struct {
	nonAtomicReplace bool         // default false
	atomically       bool         // default false
	defaultFileMode  os.FileMode  // default 0
	deleteOnClose    bool         // default 0
	fileMode         os.FileMode  // default 0
	flags            int          // default 0
	keepFileMode     bool         // default false
	readOnlyMode     ReadOnlyMode // default 0
	retrySeconds     float64      // default 0.0
	setSymlinkOwner  bool         // default false
	skipACLs         bool         // default false
}

// Option functions modify Options.
type Option func(*options)

// WithAtomicity creates or renames a file atomically.
// Used by the WriteFile and WriteReader functions.
func WithAtomicity(atomically bool) Option {
	return func(opts *options) {
		opts.atomically = atomically
	}
}

// WithDefaultFileMode sets the default file mode instead of using the
// `os.CreateTemp()` default of `0600`.
// Used by the WriteFile and WriteReader functions.
func WithDefaultFileMode(mode os.FileMode) Option {
	return func(opts *options) {
		opts.defaultFileMode = mode
	}
}

// WithDeleteOnClose deletes the file when the file is closed.
// Used by the Create, CreateTemp, Open, OpenFile, WriteFile and WriteReader functions.
func WithDeleteOnClose(deleteOnClose bool) Option {
	return func(opts *options) {
		opts.deleteOnClose = deleteOnClose
	}
}

// WithFileMode sets the file mode to the desired value and has precedence over
// all other options.
// Used by the Create, CreateTemp, MkdirTemp, Open, OpenFile, WriteFile, and
// WriteReader functions.
func WithFileMode(mode os.FileMode) Option {
	return func(opts *options) {
		opts.fileMode = mode
	}
}

// WithFlags sets the flag option.
// Used by the Create, CreateTemp, Open, OpenFile, WriteFile, and WriteReader
// functions.
func WithFlags(flags int) Option {
	return func(opts *options) {
		opts.flags |= flags
	}
}

// WithKeepFileMode preserves the file mode of an existing file instead of using the
// default value.
// Used by the WriteFile and WriteReader functions.
func WithKeepFileMode(keep bool) Option {
	return func(opts *options) {
		opts.keepFileMode = keep
	}
}

// WithNonAtomicReplace permits a non-atomic replacement when the
// operating system cannot atomically replace an existing destination.
//
// On Plan 9, this may temporarily leave the destination absent.
// The default is false.
func WithNonAtomicReplace(nonAtomicReplace bool) Option {
	return func(opts *options) {
		opts.nonAtomicReplace = nonAtomicReplace
	}
}

// WithReadOnlyMode is used to determine if/when to set a file's read-only
// (RO) attribute. The following values are supported:
// ReadOnlyModeIgnore do not set a file's RO attribute, and ignore if it's set.
// ReadOnlyModeSet set a file's RO attribute if the file's FileMode has the
// user writable bit set.
// ReadOnlyModeReset  do not set a file's RO attribute, and if it's set, reset it.
// The option is functional on Windows only. On other OSes, it is ignored.
// Used by the Chmod, Create, CreateTemp, Fchmod, Open, OpenFile, WriteFile,
// and WriteReader functions.
func WithReadOnlyMode(mode ReadOnlyMode) Option {
	return func(opts *options) {
		opts.readOnlyMode = mode
	}
}

// WithRetrySeconds sets the retry timeout option in seconds. The default is 0
// which means to not retry at all.
// Used by the Rename and RemoveAll functions.
func WithRetrySeconds(seconds float64) Option {
	return func(opts *options) {
		opts.retrySeconds = seconds
	}
}

// WithSetSymlinkOwner sets the symlink's owner to be the current user.
// Otherwise, the symlink will have a default owner assigned by the system,
// such as BUILTIN\Administrator.
// The option is functional on Windows only. On other OSes, it is ignored.
// Used by the Symlink function.
func WithSetSymlinkOwner(setSymlinkOwner bool) Option {
	return func(opts *options) {
		opts.setSymlinkOwner = setSymlinkOwner
	}
}

// WithSkipACLs skips setting ACLs (Access Control Lists) on Windows.
// The option is functional on Windows only. On other OSes, it is ignored.
// Used by the Chmod function.
func WithSkipACLs(skipACLs bool) Option {
	return func(opts *options) {
		opts.skipACLs = skipACLs
	}
}

func (o options) String() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "atomically:       %v\n", o.atomically)
	fmt.Fprintf(&builder, "defaultFileMode:  0o%03o (%v)\n", o.defaultFileMode, o.defaultFileMode)
	fmt.Fprintf(&builder, "deleteOnClose:    %v\n", o.deleteOnClose)
	fmt.Fprintf(&builder, "fileMode:         0o%03o (%v)\n", o.fileMode, o.fileMode)
	fmt.Fprintf(&builder, "flags:            0x%x\n", o.flags)
	fmt.Fprintf(&builder, "keepFileMode:     %v\n", o.keepFileMode)
	fmt.Fprintf(&builder, "nonAtomicReplace: %v\n", o.nonAtomicReplace)
	fmt.Fprintf(&builder, "readOnlyMode:     %v\n", o.readOnlyMode)
	fmt.Fprintf(&builder, "retrySeconds:     %v\n", o.retrySeconds)
	fmt.Fprintf(&builder, "setSymlinkOwner:  %v\n", o.setSymlinkOwner)
	fmt.Fprintf(&builder, "skipACLs:         %v\n", o.skipACLs)

	return builder.String()
}
