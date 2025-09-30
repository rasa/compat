// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build js

package compat

// Not supported: BTime | Fstat | Nice.
const supports supportsType = supportsATime | supportsCTime | supportsLinks | supportsSymlinks

const userIDSource UserIDSourceType = UserIDSourceIsNone
