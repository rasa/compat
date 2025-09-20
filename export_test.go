// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

// dir.go

var FSDirEntryToDirEntry = fsDirEntryToDirEntry

var FSFileInfoToDirEntry = fsFileInfoToDirEntry

var OSDirEntryToDirEntry = osDirEntryToDirEntry

// stat_*.go

var ExportedStat = stat //nolint:unused
