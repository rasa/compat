// SPDX-FileCopyrightText: Copyright © 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"errors"
	"sync"
	"sync/atomic"
)

const (
	initialUmask = 0o022
	// These are the only bits used by umask.
	umaskMask = 0o777
)

var (
	// ErrUmaskChanged reports if the umask was changed outside compat.
	ErrUmaskChanged = errors.New("umask changed outside compat")
	// ErrInvalidUmask reports if the UMASK environment variable is > 0o777.
	ErrInvalidUmask = errors.New("invalid UMASK environment variable")

	savedUmask atomic.Uint64
	umaskMutex sync.Mutex
)

func init() {
	_ = initUmask()
}
