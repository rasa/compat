// SPDX-FileCopyrightText: Copyright ?? 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

// Not supported: none.
const supports supportsType = supportsLinks | supportsATime | supportsBTime | supportsCTime | supportsFstat | supportsNice | supportsSymlinks

const userIDSource UserIDSourceType = UserIDSourceIsSID
