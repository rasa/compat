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

func TestNiceRenice(t *testing.T) {
	nice, err := compat.Nice()
	if err != nil {
		t.Error(err)
	}
	err = compat.Renice(nice)
	if err != nil {
		t.Error(err)
	}
}

func TestNiceReniceIfRoot(t *testing.T) {
	if compat.IsWasi {
		t.Log("Skipping test on wasi: operation not supported")
		return
	}

	isRoot, _ := compat.IsRoot()

	if !compat.IsWindows && !isRoot {
		t.Skip("Skipping root-only test on " + runtime.GOOS)
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
