// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build wasip1 && tinygo

package compat

// Not supported: ATime | BTime | CTime.
const supported SupportedType = Links | UID | GID

func (fs *fileStat) times() {
}
