// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

// Rename atomically replaces the destination file or directory with the
// source. It is guaranteed to either replace the target file entirely, or not
// change either file.
//
// On Plan 9, atomic replacement of an existing file is not supported. If the
// destination exists, WriteReader returns an error matching
// errors.ErrUnsupported and leaves the destination unchanged.
// To work around this issue, use the WithNonAtomicReplace option.
func Rename(source, destination string, fns ...Option) error {
	return rename(source, destination, fns...)
}
