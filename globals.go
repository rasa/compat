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

	if options.fileMode != optionDefaults.fileMode {
		opts = append(opts, WithFileMode(options.fileMode))
	}

	if options.flags != optionDefaults.flags {
		opts = append(opts, WithFlags(options.flags))
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

	if options.retrySeconds != optionDefaults.retrySeconds {
		opts = append(opts, WithRetrySeconds(options.retrySeconds))
	}

	if options.setSymlinkOwner != optionDefaults.setSymlinkOwner {
		opts = append(opts, WithSetSymlinkOwner(options.setSymlinkOwner))
	}

	if options.skipACLs != optionDefaults.skipACLs {
		opts = append(opts, WithSkipACLs(options.skipACLs))
	}

	return opts
}

func SetOptions(opts ...Option) {
	optionsMux.Lock()

	options := *optionsPtr.Load()
	for _, fn := range opts {
		fn(&options)
	}

	optionsPtr.Store(&options)
	optionsMux.Unlock()
}

func buildOptions(opts ...Option) options {
	options := *optionsPtr.Load()
	for _, fn := range opts {
		fn(&options)
	}

	return options
}
