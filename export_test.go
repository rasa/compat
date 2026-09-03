// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

// dir.go

var (
	FSDirEntryToDirEntry = fsDirEntryToDirEntry
	FSFileInfoToDirEntry = fsFileInfoToDirEntry
	OSDirEntryToDirEntry = osDirEntryToDirEntry
)

// errors.go

var (
	ChmodError       = chmodError
	CreateError      = createError
	CreateTempError  = createTempError
	MkdirError       = mkdirError
	MkdirTempError   = mkdirTempError
	NiceError        = niceError
	OpenError        = openError
	RenameError      = renameError
	StatError        = statError
	SymlinkError     = symlinkError
	WriteError       = writeError
	UnsupportedError = unsupportedError
)

// errors_plan9.go/errors_other.go

var NormalizeUnsupportedError = normalizeUnsupportedError

// globals.go

var BuildOptions = buildOptions

// options.go

type Options = options

// partitiontype.go

var NormalizePath = normalizePath

// runtime.go

var ExportedGoVersion = goVersion

// stat_*.go

var ExportedStat = stat

// writereader.go

var ExportedWriteReader = writeReader

// umask.go

var InitUmask = initUmask
