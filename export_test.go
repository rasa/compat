// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

// dir.go

var FSDirEntryToDirEntry = fsDirEntryToDirEntry

var FSFileInfoToDirEntry = fsFileInfoToDirEntry

var OSDirEntryToDirEntry = osDirEntryToDirEntry

// errors.go

var (
	ExportedUnsupportedError   = unsupportedError
	ExportedUnimplementedError = unimplementedError
	ChmodError                 = chmodError
	CreateError                = createError
	CreateTempError            = createTempError
	MkdirError                 = mkdirError
	MkdirallError              = mkdirallError
	MkdirTempError             = mkdirTempError
	OpenError                  = openError
	RenameError                = renameError
	StatError                  = statError
	SymlinkError               = symlinkError
	WriteError                 = writeError
)

// globals.go

var BuildOptions = buildOptions //nolint:unused

// runtime.go

var ExportedGoVersion = goVersion //nolint:unused

// stat_*.go

var ExportedStat = stat //nolint:unused

// writereader.go

var ExportedWriteReader = writeReader //nolint:unused
