// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

const nativeFS = "Native"

var (
	tempPath string
	tempSize string
)

type testVars struct {
	noACLs                  bool
	noSymlinks              bool
	noHardLinks             bool
	atimeGranularity        int // seconds
	btimeGranularity        int
	ctimeGranularity        int
	mtimeGranularity        int
	btimeSymlinkGranularity int
	fsType                  string
}

var testEnv = testVars{}
