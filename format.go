// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

// Source: https://github.com/golang/go/blob/ac803b59/src/io/fs/format.go#L60-L76

// FormatDirEntry returns a formatted version of dir for human readability.
// Implementations of [DirEntry] can call this from a String method.
// The outputs for a directory named subdir and a file named hello.go are:
//
//	d subdir/
//	- hello.go
func FormatDirEntry(dir DirEntry) string {
	name := dir.Name()
	b := make([]byte, 0, 5+len(name)) //nolint:mnd

	// The Type method does not return any permission bits,
	// so strip them from the string.
	mode := dir.Type().String()
	mode = mode[:len(mode)-9]
	if len(mode) > 1 {
		mode = mode[:1]
	}
	b = append(b, mode...)
	b = append(b, ' ')
	b = append(b, name...)
	if dir.IsDir() {
		b = append(b, '/')
	}
	return string(b)
}
