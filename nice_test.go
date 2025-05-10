// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"runtime"
	"testing"

	"github.com/rasa/compat"
)

func TestNice(t *testing.T) {
	_, err := compat.Nice()
	if err != nil {
		t.Error(err)
	}
}

func TestRenice(t *testing.T) {
	nice, err := compat.Nice()
	if err != nil {
		t.Error(err)
	}
	err = compat.Renice(nice)
	if err != nil {
		t.Error(err)
	}
}

func TestReniceWindows(t *testing.T) {
	isAdmin, _ := compat.IsAdmin()

	if !compat.IsWindows && !isAdmin {
		t.Skip("Skipping admin-only test on " + runtime.GOOS)
	}
	nice, err := compat.Nice()
	if err != nil {
		t.Error(err)
	}
	for n := compat.MinNice; n <= compat.MaxNice; n++ {
		err = compat.Renice(n)
		if err != nil {
			t.Error(err)
		}
	}
	err = compat.Renice(nice)
	if err != nil {
		t.Error(err)
	}
}
