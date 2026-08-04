// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build dragonfly

package compat

// Not supported: BTime | Fstat.
const (
	supportsATime         = true
	supportsATimeSetting  = true
	supportsAtomicReplace = true
	supportsBTime         = false
	supportsCTime         = true
	supportsFstat         = false
	supportsLinks         = true
	supportsNice          = true
	supportsSymlinks      = true
)

const userIDSource UserIDSourceType = UserIDSourceIsInt
