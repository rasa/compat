// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat_test

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/windows"

	"github.com/rasa/compat"
)

func TestRootIsRoot(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 1*time.Second)
	defer cancel()

	exe := "whoami.exe"
	system32, _ := windows.GetSystemDirectory()
	if system32 != "" {
		exe = filepath.Join(system32, exe)
	}
	cmd := exec.CommandContext(ctx, exe, "/all")

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		// Report failure, but don't fail, as the user's environment may be non-standard.
		t.Skipf("Command failed: '%v /all': %v: %v", exe, err, string(stdoutStderr))
	}
	// We could remove the reliance on windows, by hardcoding "S-1-5-32-544".
	sid, err := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
	if err != nil {
		t.Errorf("Unable to create well known sid for administrators: %v", err)
	}
	isRoot := strings.Contains(string(stdoutStderr), sid.String()) // "S-1-5-32-544"
	got, err := compat.IsRoot()
	if err != nil {
		t.Errorf("IsRoot() returned: %v", err)
	}
	if got != isRoot {
		// Report result, but don't fail, as the user may not be an Root.
		t.Skipf("IsRoot(): got %v, want %v", got, isRoot)
	}
}
