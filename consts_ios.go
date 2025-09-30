// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build ios

package compat

// Not supported: Nice.
const supports supportsType = supportsATime | supportsBTime | supportsCTime | supportsFstat | supportsLinks | supportsSymlinks

const userIDSource UserIDSourceType = UserIDSourceIsInt
