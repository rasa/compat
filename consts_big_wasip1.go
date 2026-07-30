// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build wasip1 && !tinygo

package compat

// Not supported: BTime | Fstat | Nice | Symlinks
const supports supportsType = supportsATime |
	supportsATimeSetting |
	supportsAtomicReplace |
	supportsCTime |
	supportsLinks

const userIDSource UserIDSourceType = UserIDSourceIsNone
