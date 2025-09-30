// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build linux

package compat

import (
	"time"

	"golang.org/x/sys/unix"
)

const supports supportsType = supportsATime | supportsBTime | supportsCTime | supportsFstat | supportsLinks | supportsNice | supportsSymlinks

const userIDSource UserIDSourceType = UserIDSourceIsInt
