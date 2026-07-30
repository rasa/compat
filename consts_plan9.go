// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9

package compat

// Not supported: ATimeSetting | AtomicReplace | BTime | CTime | Fstat | Links | Symlinks
const supports supportsType = supportsATime |
	supportsNice

const userIDSource UserIDSourceType = UserIDSourceIsString
