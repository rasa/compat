// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"bytes"
	"os/exec"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/rasa/compat"
)

var oses = []string{
	"aix",
	"android",
	"darwin",
	"dragonfly",
	"freebsd",
	"illumos",
	"ios",
	"js",
	"linux",
	"netbsd",
	"openbsd",
	"plan9",
	"solaris",
	"wasip1",
	"windows",
}

var arches = []string{
	"386",
	"amd64",
	"arm",
	"arm64",
	"loong64",
	"mips",
	"mips64",
	"mips64le",
	"mipsle",
	"ppc64",
	"ppc64le",
	"riscv64",
	"s390x",
	"wasm",
}

func TestRuntimeConsts(t *testing.T) { //nolint:funlen,gocyclo
	goExe, err := exec.LookPath("go")
	if err != nil {
		if compat.IsTinygo || compat.IsWasm {
			skipf(t, "Skipping test: %v (1)", err)
			return
		}

		t.Fatal(err)
	}

	out, err := exec.Command(goExe, "tool", "dist", "list").Output() //nolint:noctx
	if err != nil {
		if compat.IsTinygo || compat.IsWasm {
			skipf(t, "Skipping test: %v (2)", err)
			return
		}

		t.Fatal(err)
	}

	gooses := make(map[string]struct{})
	goarches := make(map[string]struct{})

	lines := bytes.SplitSeq(out, []byte{'\n'})
	for line := range lines {
		trimmed := strings.TrimSpace(string(line))

		before, after, found := strings.Cut(trimmed, "/")
		if !found {
			continue
		}

		gooses[before] = struct{}{}
		goarches[after] = struct{}{}
	}

	if len(gooses) == 0 {
		t.Fatal("failed to parse output of: go tool dist list")
	}

	for goos := range gooses {
		if !slices.Contains(oses, goos) {
			t.Errorf("found new GOOS: %q", goos)
		}
	}

	for goarch := range goarches {
		if !slices.Contains(arches, goarch) {
			t.Errorf("found new GOARCH: %q", goarch)
		}
	}

	for _, goos := range oses {
		_, ok := gooses[goos]
		if !ok {
			t.Logf("go no longer supports GOOS: %q", goos)
		}
	}

	for _, goarch := range arches {
		_, ok := goarches[goarch]
		if !ok {
			t.Logf("go no longer supports GOARCH: %q", goarch)
		}
	}
}

func TestRuntimeGoVersion(t *testing.T) {
	if compat.IsTinygo {
		skip(t, "Skipping test: fails on tinygo")
		return
	}

	want := runtime.Version()
	got := compat.GoVersion()
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestRuntimeGetGoVersionFalse(t *testing.T) {
	want := runtime.Version()
	got := compat.ExportedGoVersion(want, false)
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestRuntimeGetGoVersionTrue(t *testing.T) {
	want := "go1.25"
	got := compat.ExportedGoVersion("0.39.1", true)
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}
