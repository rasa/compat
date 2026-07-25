// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package user

import (
	osuser "os/user"
)

var Current = osuser.Current

var (
	FindHomeDirInRegistry          = findHomeDirInRegistry
	GetCurrentToken                = getCurrentToken
	GetProfilesDirectory           = getProfilesDirectory
	IsDomainJoined                 = isDomainJoined
	IsServiceAccount               = isServiceAccount
	IsValidGroupAccountType        = isValidGroupAccountType
	IsValidUserAccountType         = isValidUserAccountType
	ListGroups                     = listGroups
	ListGroupsForUsernameAndDomain = listGroupsForUsernameAndDomain
	LookupFullName                 = lookupFullName
	LookupFullNameDomain           = lookupFullNameDomain
	LookupFullNameServer           = lookupFullNameServer
	LookupGroup                    = lookupGroup
	LookupGroupId                  = lookupGroupId
	LookupGroupName                = lookupGroupName
	LookupUser                     = lookupUser
	LookupUserId                   = lookupUserId
	LookupUsernameAndDomain        = lookupUsernameAndDomain
	LookupUserPrimaryGroup         = lookupUserPrimaryGroup
	NewUser                        = newUser
	NewUserFromSid                 = newUserFromSid
	RunAsProcessOwner              = runAsProcessOwner
)

type (
	Group = osuser.Group
	User  = osuser.User
)
