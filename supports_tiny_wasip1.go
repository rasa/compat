// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build wasip1 && tinygo

package compat

const (
	supportsATime         = false
	supportsATimeSetting  = false
	supportsAtomicReplace = true
	supportsBTime         = false
	supportsCTime         = false
	supportsFstat         = false
	supportsLinks         = false
	supportsNice          = false
	supportsRelativeFstat = true
	supportsSymlinks      = false
	supportsUmask         = false
)

const userIDSource UserIDSourceType = UserIDSourceIsNone
