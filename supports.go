// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

// SupportsATime returns true if FileInfo's ATime() function is supported by the OS.
func SupportsATime() bool {
	return supportsATime
}

// SupportsATimeSetting is defined in stat_other.go and stat_tiny.go

// SupportsAtomicReplace returns true if atomically replacing a file is supported by the OS.
func SupportsAtomicReplace() bool {
	return supportsAtomicReplace
}

// SupportsBTime returns true if FileInfo's BTime() function is supported by the OS.
func SupportsBTime() bool {
	return supportsBTime
}

// SupportsCTime returns true if FileInfo's CTime() function is supported by the OS.
func SupportsCTime() bool {
	return supportsCTime
}

// SupportsFstat returns true if the Fstat() function is supported by the OS.
func SupportsFstat() bool {
	return supportsFstat
}

// SupportsLinks returns true if FileInfo's Links() function is supported by the OS.
func SupportsLinks() bool {
	return supportsLinks
}

// SupportsNice returns true if the Nice() function is supported by the OS.
func SupportsNice() bool {
	return supportsNice
}

// SupportsRelativeFstat always returns true.
// Deprecated. Will be removed in a future version.
func SupportsRelativeFstat() bool {
	return supportsRelativeFstat
}

// SupportsSymlinks returns true if the os.Symlinks() function is supported by the OS.
// Note that the underlying filesystem may not allow symlinks.
func SupportsSymlinks() bool {
	return supportsSymlinks
}

// SupportsUmask returns true if the Umask() function is supported by the OS.
func SupportsUmask() bool {
	return supportsUmask
}

// UserIDSource returns the source of the user's ID: UserIDSourceIsInt,
// UserIDSourceIsString, or UserIDSourceIsNone.
func UserIDSource() UserIDSourceType {
	return userIDSource
}
