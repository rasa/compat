// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

// IsAdmin returns true if the user is root, or has Windows administrator rights.
//
// Deprecated: Use IsRoot() instead.
var IsAdmin = IsRoot

// ReplaceFile atomically replaces the destination file or directory with the
// source.  It is guaranteed to either replace the target file entirely, or not
// change either file.
//
// Deprecated: Use Rename() instead.
var ReplaceFile = Rename

// Deprecated: Use Supports*() functions instead.
type SupportedType = supportsType

const (
	// Links defines if FileInfo's Links() function is supported by the OS.
	// Deprecated: Use SupportsLinks() instead.
	Links = supportsLinks
	// ATime defines if FileInfo's ATime() function is supported by the OS.
	// Deprecated: No longer used or needed.
	ATime = supportsATime
	// BTime defines if FileInfo's BTime() function is supported by the OS.
	// Deprecated: Use SupportsBTime() instead.
	BTime = supportsBTime
	// CTime defines if FileInfo's CTime() function is supported by the OS.
	// Deprecated: Use SupportsCTime() instead.
	CTime = supportsCTime
	// UID defines if FileInfo's UID() function is supported by the OS.
	// Deprecated: No longer used or needed.
	UID = supportsUID
	// GID defines if FileInfo's GID() function is supported by the OS.
	// Deprecated: No longer used or needed.
	GID = supportsGID
)

// Supported returns whether supportedType is supported by the operating system.
// Deprecated: Use SupportsLinks(), SupportsBTime() and SupportsCTime() functions instead.
func Supported(supportedType SupportedType) bool {
	return supports&supportedType == supportedType
}
