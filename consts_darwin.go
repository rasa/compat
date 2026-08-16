// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build darwin && !ios

// The darwin build flag includes ios (which doesn't support Nice())

package compat

const (
	supportsATime         = true
	supportsATimeSetting  = true
	supportsAtomicReplace = true
	supportsBTime         = true
	supportsCTime         = true
	supportsFstat         = true
	supportsLinks         = true
	supportsNice          = true
	supportsSymlinks      = true
	supportsUmask         = true
)

const userIDSource UserIDSourceType = UserIDSourceIsInt
