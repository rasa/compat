// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !unix

package compat

// IsUnix is true if the unix build tag has been set.
// Currently, it's equivalent to: IsAIX || IsAndroid || IsApple || IsBSD || IsLinux || IsSolaria.
const IsUnix = false
