// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build wasip1 && tinygo

package compat

// Not supported: ATimeSetting | ATime | BTime | CTime | Fstat | Links | Nice | Symlinks.
const supports supportsType = supportsAtomicReplace

const userIDSource UserIDSourceType = UserIDSourceIsNone
