// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build darwin && !ios

// The darwin build flag includes ios (which doesn't support Nice())

package compat

// Not supported: none.
const supports supportsType = supportsATime |
	supportsATimeSetting |
	supportsAtomicReplace |
	supportsBTime |
	supportsCTime |
	supportsFstat |
	supportsLinks |
	supportsNice |
	supportsSymlinks

const userIDSource UserIDSourceType = UserIDSourceIsInt
