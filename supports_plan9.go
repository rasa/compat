// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9

package compat

const (
	supportsATime         = true
	supportsATimeSetting  = false
	supportsAtomicReplace = false
	supportsBTime         = false
	supportsCTime         = false
	supportsFstat         = false
	supportsLinks         = false
	supportsNice          = true
	supportsRelativeFstat = true
	supportsSymlinks      = false
	supportsUmask         = false
)

const userIDSource UserIDSourceType = UserIDSourceIsString
