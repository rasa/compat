// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"bytes"
	"os"
)

// WriteFile writes data to the named file, creating it if necessary.
// If the file does not exist, WriteFile creates it using perm's permissions
// bits (before umask); otherwise WriteFile truncates it before writing,
// without changing permissions. Since WriteFile requires multiple system
// calls to complete, a failure mid-operation can leave the file in a partially
// written state. Use WriteFile() with the WithAtomicity(true) option,
// if this is a concern.
//
// When WithAtomicity(true) is passed, WriteFile atomically writes the contents
// of data to the specified filename. The target file is guaranteed to be either
// fully written, or not written at all. WriteFile overwrites any file that
// exists at the location (but only if the write fully succeeds, otherwise the
// existing file is unmodified).
//
// If perm is zero, then 0o666 is used, as this is what the os.Create() function
// uses. If both perm, and WithFileMode(perm) are provided, WithFileMode(perm)
// takes precedence.
//
// Additional option arguments can be used to change the default configuration
// for the target file.
//
// On Plan 9, atomic creation of a new file is supported, but atomic replacement
// of an existing file is not. If the destination exists, WriteFile returns an
// error matching errors.ErrUnsupported and leaves the destination unchanged.
// To work around this issue, use the WithNonAtomicReplace option.
func WriteFile(name string, data []byte, perm os.FileMode, opts ...Option) error {
	return WriteReader(name, bytes.NewReader(data), perm, opts...)
}
