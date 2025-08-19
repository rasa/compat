// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"errors"
	"fmt"
	"runtime"
	"syscall"
	"testing"

	"github.com/rasa/compat"
)

func skip(t *testing.T, msg any) {
	t.Helper()

	s := fmt.Sprint(msg)
	if compat.IsTinygo {
		s += " (" + runtime.GOOS + "/tinygo" + ")"
		t.Log(s)

		return
	}
	if compat.IsAct {
		s += " (" + runtime.GOOS + "/act" + ")"
	}

	t.Skip(s)
}

func skipf(t *testing.T, format string, a ...any) {
	t.Helper()
	skip(t, fmt.Sprintf(format, a...))
}

func fatal(t *testing.T, msg any) { //nolint:unused // quiet linter
	t.Helper()

	s := fmt.Sprint(msg)
	if compat.IsTinygo {
		s = "Skipping test: fatal error: " + s
		s += " (" + runtime.GOOS + "/tinygo" + ")"
		t.Log(s)

		return
	}

	t.Fatal(s)
}

func fatalf(t *testing.T, format string, a ...any) { //nolint:unused // quiet linter
	t.Helper()
	fatal(t, fmt.Sprintf(format, a...))
}

func errno(err error) uint32 {
	if err == nil {
		return 0
	}
	var errno syscall.Errno
	if errors.As(err, &errno) {
		return uint32(errno)
	}

	return ^uint32(0)
}
