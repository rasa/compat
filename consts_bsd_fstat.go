// SPDX-FileCopyrightText: Copyright © 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build fstat && (freebsd || netbsd)

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
