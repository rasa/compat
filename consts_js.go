// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build js

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
	supportsSymlinks      = true
	supportsUmask         = true
)

const userIDSource UserIDSourceType = UserIDSourceIsNone
