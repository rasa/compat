// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package robustio

// Retry retries ephemeral errors from f up to an arbitrary timeout
// to work around filesystem flakiness on Windows and Darwin.
func Retry(f func() (err error, mayRetry bool), retrySeconds float64) error {
	return retry(f, retrySeconds)
}

// IsEphemeralError returns true if err may be resolved by waiting.
func IsEphemeralError(err error) bool {
	return isEphemeralError(err)
}
