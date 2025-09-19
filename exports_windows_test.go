// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

// file_windows.go

var CurrentUsername = currentUsername

// stat_posix_windows.go

var CopySid = copySid

var EqualDomainSid = equalDomainSid

var GetFileOwnerAndGroupSIDs = getFileOwnerAndGroupSIDs

var GetPrimaryDomainSID = getPrimaryDomainSID

var GetRID = getRID

var GetUserGroup = getUserGroup

var IsValidSid = isValidSid

var LSAOpenPolicy = lsaOpenPolicy //nolint:unused

var NameFromSID = nameFromSID

var SIDToPOSIXID = sidToPOSIXID

// stat_acls_windows.go

var SupportsACLs = supportsACLs

var SupportsACLsCached = supportsACLsCached

var SupportsACLsHandle = supportsACLsHandle

var OpenForQuery = openForQuery

var GetFinalPathNameByHandleGUID = getFinalPathNameByHandleGUID

var GetVolumePathNamesForVolumeName = getVolumePathNamesForVolumeName

var GetVolumeInfoByHandle = getVolumeInfoByHandle

var ResolveCanonicalRootFromHandle = resolveCanonicalRootFromHandle

var MultiSZToStrings = multiSZToStrings //nolint:unused

var IsDriveLetterRoot = isDriveLetterRoot //nolint:unused

var NormalizeRoot = normalizeRoot
