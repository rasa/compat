// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: Portions copyright (c) 2015 Nate Finch (@natefinch)
// SPDX-FileCopyrightText: Portions copyright (c) 2022 Simon Dassow (@sdassow)

package compat

import (
	"os"
)

// Options define the behavior of `WriteReaderAtomic()`, etc.
type Options struct {
	defaultFileMode os.FileMode  // default 0
	fileMode        os.FileMode  // default 0
	flags           int          // default 0
	keepFileMode    bool         // default false
	readOnlyMode    ReadOnlyMode // default 0
	setSymlinkOwner bool         // default false
}

// Option functions modify Options.
type Option func(*Options)

// WithDefaultFileMode sets the default file mode instead of using the
// `os.CreateTemp()` default of `0600`.
func WithDefaultFileMode(mode os.FileMode) Option {
	return func(opts *Options) {
		opts.defaultFileMode = mode
	}
}

// WithFileMode sets the file mode to the desired value and has precedence over all
// other options.
func WithFileMode(mode os.FileMode) Option {
	return func(opts *Options) {
		opts.fileMode = mode
	}
}

// WithFlags sets the flag option.
func WithFlags(flags int) Option {
	return func(opts *Options) {
		opts.flags = flags
	}
}

// WithKeepFileMode preserves the file mode of an existing file instead of using the
// default value.
func WithKeepFileMode(keep bool) Option {
	return func(opts *Options) {
		opts.keepFileMode = keep
	}
}

// WithReadOnlyMode is used to determine if/when to set a file's read-only
// (RO) attribute on Windows. The following values are supported:
// ReadOnlyModeIgnore do not set a file's RO attribute, and ignore if it's set.
// ReadOnlyMaskSet    set a file's RO attribute if the file's FileMode has the
//
//	user writable bit set.
//
// ReadOnlyMaskReset  do not set a file's RO attribute, and if it's set, reset it.
func WithReadOnlyMode(mode ReadOnlyMode) Option {
	return func(opts *Options) {
		opts.readOnlyMode = mode
	}
}

// WithSetSymlinkOwner sets the symlink's owner to be the current user.
// Otherwise, the symlink will have a default owner assigned by the system,
// such as BUILTIN\Administrator.
func WithSetSymlinkOwner(setSymlinkOwner bool) Option {
	return func(opts *Options) {
		opts.setSymlinkOwner = setSymlinkOwner
	}
}
