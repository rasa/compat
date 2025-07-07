// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

// IsWSL returns true if run instead a Windows Subsystem for Linux (WSL)
// environment, otherwise false.
//
// It's counter-intuitive that IsWSL() returns false in Windows, but here's why:
// WSL can run executables built to run on Linux, and those built to run on
// Windows. For example, executing `whoami` will run WSL's `/usr/bin/whoami`,
// but append a `.exe`, and execute `whoami.exe`, then WSL will instead run
// Windows' whoami, in `C:/Windows/System32/whoami.exe`. But it doesn't
// appear to me that executables built to run on Windows can't tell they were
// started from inside a WSL environment. For example, the program doesn't see
// the `WSL_DISTRO_NAME“ environment variable that other programs run inside
// WSL see. Hence, this function must return false.
func IsWSL() bool {
	return false
}
