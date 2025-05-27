// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"bytes"
	"os/exec"
	"slices"
	"strings"
	"testing"
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

func TestRuntime(t *testing.T) {
	out, err := exec.Command("go", "tool", "dist", "list").Output()
	if err != nil {
		t.Fatal(err)
	}

	gooses := make(map[string]struct{})
	goarches := make(map[string]struct{})

	lines := bytes.Split(out, []byte{'\n'})
	for _, line := range lines {
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
