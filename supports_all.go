// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build (darwin && !ios) || freebsd || (linux && !android) || (netbsd && fstat) || windows

// The darwin build flag includes ios (which doesn't support Nice())
// The linux build flag includes android (which doesn't support BTime())

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
	supportsRelativeFstat = true
	supportsSymlinks      = true
	supportsUmask         = true
)

const userIDSource UserIDSourceType = UserIDSourceIsInteger
