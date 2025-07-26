// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/rasa/compat"
)

func skip(t *testing.T, msg any) {
	t.Helper()
	s := fmt.Sprint(msg)
	if compat.IsWasip1 {
		goos := runtime.GOOS
		if compat.IsTinygo {
			goos += "/tinygo"
		}
		s += " (" + goos + ")"
		t.Log(s)
		return
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
	if compat.IsWasip1 {
		s = "Skipping test: fatal error: " + s
		goos := runtime.GOOS
		if compat.IsTinygo {
			goos += "/tinygo"
		}
		s += " (" + goos + ")"
		t.Log(s)
		return
	}
	t.Fatal(s)
}

func fatalf(t *testing.T, format string, a ...any) { //nolint:unused // quiet linter
	t.Helper()
	fatal(t, fmt.Sprintf(format, a...))
}
