// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"os"
)

// Deprecated: Use IsRoot() instead. This function may be removed in the future.
var IsAdmin = IsRoot

// Deprecated: Use Rename() instead. This function may be removed in the future.
var ReplaceFile = Rename

// Deprecated: Use Options instead. This may be removed in the future.
type FileOptions = Options

// Deprecated: Use Supports*() functions instead. This may be removed in the future.
type SupportedType = supportsType

const (
	// Deprecated: Use SupportsLinks() instead. This may be removed in the future.
	Links = supportsLinks
	// Deprecated: Use SupportsATime() instead. This may be removed in the future.
	ATime = supportsATime
	// Deprecated: Use SupportsBTime() instead. This may be removed in the future.
	BTime = supportsBTime
	// Deprecated: Use SupportsCTime() instead. This may be removed in the future.
	CTime = supportsCTime
	// Deprecated: No longer used or needed. This may be removed in the future.
	UID = supportsUID
	// Deprecated: No longer used or needed. This may be removed in the future.
	GID = supportsGID
)

// Deprecated: Use SupportsLinks(), SupportsATime(), SupportsBTime() and
// SupportsCTime() functions instead. This function may be removed in the future.
var Supported = supported

func supported(supportedType SupportedType) bool {
	return supports&supportedType == supportedType
}

// Deprecated: Use UserIDSourceIsInt instead.
// This may be removed in the future.
const UserIDSourceIsNumeric = UserIDSourceIsInt

// Deprecated: Use Create() instead, and pass perm and flag via Option array.
// This function may be removed in the future.
var CreateEx = createex

func createex(name string, perm os.FileMode, flag int) (*os.File, error) {
	flag |= os.O_CREATE
	if flag&os.O_WRONLY == 0 {
		flag |= os.O_RDWR
	}

	return create(name, perm, flag)
}

// Deprecated: Use CreateTemp() instead, and pass flag via Option array.
// This function may be removed in the future.
var CreateTempEx = createTempEx

func createTempEx(dir, pattern string, flag int) (*os.File, error) {
	return createTemp(dir, pattern, CreateTempPerm, flag)
}

// Deprecated: Use WriteFile() instead, and pass flag via Option array.
// This function may be removed in the future.
var WriteFileEx = writeFileEx

func writeFileEx(name string, data []byte, perm os.FileMode, flag int) error {
	return writeFile(name, data, perm, flag)
}

// Deprecated: Use WithDefaultFileMode() instead.
// This function may be removed in the future.
var DefaultFileMode = WithDefaultFileMode

// Deprecated: Use WithFlags() instead.
// This function may be removed in the future.
var Flag = WithFlags

// Deprecated: Use WithKeepFileMode() instead.
// This function may be removed in the future.
var KeepFileMode = WithKeepFileMode

// Deprecated: Use WithFileMode() instead.
// This function may be removed in the future.
var UseFileMode = WithFileMode

// Deprecated: Use O_FILE_FLAG_DELETE_ON_CLOSE instead.
// This constant may be removed in the future.
const O_DELETE = O_FILE_FLAG_DELETE_ON_CLOSE

// Deprecated: Use O_FILE_FLAG_NO_RO_ATTR instead.
// This constant may be removed in the future.
const O_NOROATTR = O_FILE_FLAG_NO_RO_ATTR
