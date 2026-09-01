// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"sync"
	"sync/atomic"
)

var (
	optionsPtr     atomic.Pointer[options]
	optionDefaults = &options{}
	optionsMux     sync.Mutex
)

func init() {
	optionsPtr.Store(optionDefaults)
}

func GetOptions() []Option {
	options := *optionsPtr.Load()
	opts := make([]Option, 0)

	if options.atomically != optionDefaults.atomically {
		opts = append(opts, WithAtomicity(options.atomically))
	}

	if options.defaultFileMode != optionDefaults.defaultFileMode {
		opts = append(opts, WithDefaultFileMode(options.defaultFileMode))
	}

	if options.deleteOnClose != optionDefaults.deleteOnClose {
		opts = append(opts, WithDeleteOnClose(options.deleteOnClose))
	}

	if options.fileMode != optionDefaults.fileMode {
		opts = append(opts, WithFileMode(options.fileMode))
	}

	if options.openFlags != optionDefaults.openFlags {
		opts = append(opts, WithOpenFlags(options.openFlags))
	}

	if options.keepFileMode != optionDefaults.keepFileMode {
		opts = append(opts, WithKeepFileMode(options.keepFileMode))
	}

	if options.nonAtomicReplace != optionDefaults.nonAtomicReplace {
		opts = append(opts, WithNonAtomicReplace(options.nonAtomicReplace))
	}

	if options.readOnlyMode != optionDefaults.readOnlyMode {
		opts = append(opts, WithReadOnlyMode(options.readOnlyMode))
	}

	if options.retryTimeout != optionDefaults.retryTimeout {
		opts = append(opts, WithRetryTimeout(options.retryTimeout))
	}

	if options.setSymlinkOwner != optionDefaults.setSymlinkOwner {
		opts = append(opts, WithSetSymlinkOwner(options.setSymlinkOwner))
	}

	if options.skipACLs != optionDefaults.skipACLs {
		opts = append(opts, WithSkipACLs(options.skipACLs))
	}

	return opts
}

func SetOptions(fns ...Option) {
	optionsMux.Lock()

	options := *optionsPtr.Load()
	for _, fn := range fns {
		fn(&options)
	}

	optionsPtr.Store(&options)
	optionsMux.Unlock()
}

func buildOptions(fns ...Option) options {
	options := *optionsPtr.Load()
	for _, fn := range fns {
		fn(&options)
	}

	return options
}
