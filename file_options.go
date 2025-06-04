// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: Portions copyright (c) 2015 Nate Finch (@natefinch)
// SPDX-FileCopyrightText: Portions copyright (c) 2022 Simon Dassow (@sdassow)

package compat

import (
	"os"
)

// FileOptions define the behavior of `WriteReaderAtomic()`, etc.
type FileOptions struct {
	defaultFileMode os.FileMode
	fileMode        os.FileMode
	keepFileMode    bool
}

// Option functions modify FileOptions.
type Option func(*FileOptions)

// FileMode sets the file mode to the desired value and has precedence over all
// other options.
func FileMode(mode os.FileMode) Option {
	return func(opts *FileOptions) {
		opts.fileMode = mode
	}
}

// DefaultFileMode sets the default file mode instead of using the
// `os.CreateTemp()` default of `0600`.
func DefaultFileMode(mode os.FileMode) Option {
	return func(opts *FileOptions) {
		opts.defaultFileMode = mode
	}
}

// KeepFileMode preserves the file mode of an existing file instead of using the
// default value.
func KeepFileMode(keep bool) Option {
	return func(opts *FileOptions) {
		opts.keepFileMode = keep
	}
}
