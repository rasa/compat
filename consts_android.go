// SPDX-FileCopyrightText: Copyright © 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build android

package compat

// Not supported: supportsBTime
const supports supportsType = supportsATime | supportsCTime | supportsFstat | supportsLinks | supportsNice | supportsSymlinks

const userIDSource UserIDSourceType = UserIDSourceIsInt
