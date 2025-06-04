// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: Portions copyright (c) 2015 Nate Finch (@natefinch)
// SPDX-FileCopyrightText: Portions copyright (c) 2022 Simon Dassow (@sdassow)

package compat

import (
	"bytes"
)

// WriteFileAtomic atomically writes the contents of data to the specified filepath. If
// an error occurs, the target file is guaranteed to be either fully written, or
// not written at all. WriteFileAtomic overwrites any file that exists at the
// location (but only if the write fully succeeds, otherwise the existing file
// is unmodified). Additional option arguments can be used to change the
// default configuration for the target file.
func WriteFileAtomic(filename string, data []byte, opts ...Option) (err error) {
	return WriteReaderAtomic(filename, bytes.NewReader(data), opts...)
}
