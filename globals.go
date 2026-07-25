// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"fmt"
	"strings"
	"sync/atomic"
)

var (
	options        atomic.Pointer[Options]
	optionDefaults = &Options{}
)

func init() {
	options.Store(optionDefaults)
}

func GetOptions() []Option { //nolint:unused
	o := *options.Load()
	opts := make([]Option, 0)

	if o.atomically != optionDefaults.atomically {
		opts = append(opts, WithAtomicity(o.atomically))
	}
	if o.defaultFileMode != optionDefaults.defaultFileMode {
		opts = append(opts, WithDefaultFileMode(o.defaultFileMode))
	}
	if o.fileMode != optionDefaults.fileMode {
		opts = append(opts, WithFileMode(o.fileMode))
	}
	if o.flags != optionDefaults.flags {
		opts = append(opts, WithFlags(o.flags))
	}
	if o.keepFileMode != optionDefaults.keepFileMode {
		opts = append(opts, WithKeepFileMode(o.keepFileMode))
	}
	if o.readOnlyMode != optionDefaults.readOnlyMode {
		opts = append(opts, WithReadOnlyMode(o.readOnlyMode))
	}
	if o.retrySeconds != optionDefaults.retrySeconds {
		opts = append(opts, WithRetrySeconds(o.retrySeconds))
	}
	if o.setSymlinkOwner != optionDefaults.setSymlinkOwner {
		opts = append(opts, WithSetSymlinkOwner(o.setSymlinkOwner))
	}

	return opts
}

func SetOptions(opts ...Option) { //nolint:unused
	base := *options.Load()
	for _, fn := range opts {
		fn(&base)
	}
	options.Store(&base)
}

func buildOptions(opts ...Option) Options { //nolint:unused
	fopts := *options.Load()
	for _, fn := range opts {
		fn(&fopts)
	}
	return fopts
}

func (o Options) String() string { //nolint:unused
	var b strings.Builder
	fmt.Fprintf(&b, "atomically:      %v\n", o.atomically)
	fmt.Fprintf(&b, "defaultFileMode: 0o%03o (%v)\n", o.defaultFileMode, o.defaultFileMode)
	fmt.Fprintf(&b, "fileMode:        0o%03o (%v)\n", o.fileMode, o.fileMode)
	fmt.Fprintf(&b, "flags:           0x%x\n", o.flags)
	fmt.Fprintf(&b, "keepFileMode:    %v\n", o.keepFileMode)
	fmt.Fprintf(&b, "readOnlyMode:    %v\n", o.readOnlyMode)
	fmt.Fprintf(&b, "retrySeconds:    %v\n", o.retrySeconds)
	fmt.Fprintf(&b, "setSymlinkOwner: %v\n", o.setSymlinkOwner)

	return b.String()
}

func Example(opts ...Option) { //nolint:unused
	fopts := buildOptions(opts...)
	fmt.Printf("fopts:\n%v\n", fopts)
}
