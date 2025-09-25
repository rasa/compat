// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build dragonfly

package compat

// Not supported: BTime.
const supports supportsType = supportsLinks | supportsATime | supportsCTime | supportsFstat | supportsNice | supportsSymlinks
