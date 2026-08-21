// SPDX-FileCopyrightText: Copyright © 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build wasip1 && !tinygo

package compat

const (
	supportsATime         = true
	supportsATimeSetting  = true
	supportsAtomicReplace = true
	supportsBTime         = false
	supportsCTime         = true
	supportsFstat         = false
	supportsLinks         = true
	supportsNice          = false
	supportsRelativeFstat = false
	supportsSymlinks      = false
	supportsUmask         = false
)

const userIDSource UserIDSourceType = UserIDSourceIsNone
