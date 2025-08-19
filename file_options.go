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
	defaultFileMode os.FileMode
	useFileMode     os.FileMode
	keepFileMode    bool
	flag            int
}

// Option functions modify FileOptions.
type Option func(*Options)

// UseFileMode sets the file mode to the desired value and has precedence over all
// other options.
func UseFileMode(mode os.FileMode) Option {
	return func(opts *Options) {
		opts.useFileMode = mode
	}
}

// DefaultFileMode sets the default file mode instead of using the
// `os.CreateTemp()` default of `0600`.
func DefaultFileMode(mode os.FileMode) Option {
	return func(opts *Options) {
		opts.defaultFileMode = mode
	}
}

// KeepFileMode preserves the file mode of an existing file instead of using the
// default value.
func KeepFileMode(keep bool) Option {
	return func(opts *Options) {
		opts.keepFileMode = keep
	}
}

// Flag sets the flag option.
func Flag(flag int) Option {
	return func(opts *Options) {
		opts.flag = flag
	}
}
