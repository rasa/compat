// SPDX-FileCopyrightText: Copyright © 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build android || (openbsd && fstat)

package compat

const (
	supportsATime         = true
	supportsATimeSetting  = true
	supportsAtomicReplace = true
	supportsBTime         = false
	supportsCTime         = true
	supportsFstat         = true
	supportsLinks         = true
	supportsNice          = true
	supportsRelativeFstat = true
	supportsSymlinks      = true
	supportsUmask         = true
)

const userIDSource UserIDSourceType = UserIDSourceIsInt
