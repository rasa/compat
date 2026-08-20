// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

// file_windows.go

var CurrentUsername = currentUsername

var EnablePrivilege = enablePrivilege

var SaFromPerm = saFromPerm

var SetOwnerToCurrentUser = setOwnerToCurrentUser

// stat_acls_windows.go

var SupportsACLs = supportsACLs

var SupportsACLsCached = supportsACLsCached

var SupportsACLsHandle = supportsACLsHandle

var OpenForQuery = openForQuery

var GetFinalPathNameByHandleGUID = getFinalPathNameByHandleGUID

var GetVolumePathNamesForVolumeName = getVolumePathNamesForVolumeName

var GetVolumeInfoByHandle = getVolumeInfoByHandle

var ResolveCanonicalRootFromHandle = resolveCanonicalRootFromHandle

var MultiSZToStrings = multiSZToStrings

var IsDriveLetterRoot = isDriveLetterRoot

var NormalizeRoot = normalizeRoot

// stat_posix_windows.go

var CopySid = copySid

var EqualDomainSid = equalDomainSid

var GetFileOwnerAndGroupSIDs = getFileOwnerAndGroupSIDs

var GetPrimaryDomainSID = getPrimaryDomainSID

var GetRID = getRID

var GetUserGroup = getUserGroup

var IsValidSid = isValidSid

var LSAOpenPolicy = lsaOpenPolicy

var NameFromSID = nameFromSID

var SIDToPOSIXID = sidToPOSIXID
