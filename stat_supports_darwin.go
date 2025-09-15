// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build darwin && !ios

// The darwin build flag includes ios

package compat

const supports supportsType = supportsATime | supportsBTime | supportsCTime | supportsLinks | supportsNice | supportsSymlinks
