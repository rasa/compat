// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build aix || illumos || openbsd || solaris

package compat

// Not supported: BTime | Fstat.
const supports supportsType = supportsLinks | supportsATime | supportsCTime | supportsNice | supportsSymlinks
