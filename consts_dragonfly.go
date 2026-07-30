// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build dragonfly

package compat

// Not supported: BTime | Fstat.
const supports supportsType = supportsATime |
	supportsATimeSetting |
	supportsAtomicReplace |
	supportsCTime |
	supportsLinks |
	supportsNice |
	supportsSymlinks

const userIDSource UserIDSourceType = UserIDSourceIsInt
