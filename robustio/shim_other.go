// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !windows && !darwin

package robustio

// retry retries ephemeral errors from f up to an arbitrary timeout
// to work around filesystem flakiness on Windows and Darwin.
func retry(f func() (_ error, _ bool), _ float64) error {
	err, _ := f()
	return err
}
