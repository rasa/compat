// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build aix || dragonfly || illumos || (openbsd && !fstat) || solaris

package compat

const (
	supportsATime         = true
	supportsATimeSetting  = true
	supportsAtomicReplace = true
	supportsBTime         = false
	supportsCTime         = true
	supportsFstat         = false
	supportsLinks         = true
	supportsNice          = true
	supportsRelativeFstat = true
	supportsSymlinks      = true
	supportsUmask         = true
)

const userIDSource UserIDSourceType = UserIDSourceIsInt
