// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build freebsd || netbsd

package compat

// Not supported: Fstat.
const supports supportsType = supportsATime | supportsBTime | supportsCTime | supportsLinks | supportsNice | supportsSymlinks

const userIDSource UserIDSourceType = UserIDSourceIsInt
