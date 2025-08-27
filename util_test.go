// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/rasa/compat"
)

var (
	compatDebug = strings.ToUpper(os.Getenv("COMPAT_DEBUG"))
	helloBytes  = []byte("hello")
)

func init() {
	testing.Init()
	flag.Parse()

	// @TODO(rasa): test different umask settings
	compat.Umask(0)
}

func compareNames(got string, want string) bool {
	if compat.IsWasip1 {
		if got == "" && want == "daemon" {
			return true
		}
	}

	if !compat.IsWindows {
		return got == want
	}

	if testEnv.noACLs {
		return true
	}

	if got == "" || want == "" {
		return false
	}
	gotDomain, gotName := parseName(got)
	wantDomain, wantName := parseName(want)
	if gotName == wantName {
		if gotDomain == wantDomain || gotDomain == "" || wantDomain == "" {
			return true
		}
	}

	return false
}

func parseName(name string) (string, string) {
	parts := strings.Split(name, `\`)
	switch {
	case len(parts) == 1:
		return "", strings.ToLower(parts[0])
	default:
		return strings.ToLower(parts[0]), strings.ToLower(parts[1])
	}
}

func compareTimes(a, b time.Time, granularity int) bool {
	if granularity < 0 {
		return a.IsZero()
	}
	if granularity == 0 {
		granularity = 1
	}

	return a.Sub(b).Abs() < time.Duration(granularity)*time.Second
}

func errno(err error) uint32 { //nolint:unused // quiet linter
	if err == nil {
		return 0
	}
	var errno syscall.Errno
	if errors.As(err, &errno) {
		return uint32(errno)
	}

	return ^uint32(0)
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

func fatalTimes(t *testing.T, prefix string, got, want time.Time, granularity int) { //nolint:unused // quiet linter
	t.Helper()

	diff := got.Sub(want).Abs().Seconds()

	t.Fatalf("%v: got %.2fs difference, want <%ds (%v vs %v)", prefix, diff, granularity, got, want)
}

func fixPerms(perm os.FileMode, isDir bool) os.FileMode {
	if testEnv.noACLs {
		if isDir {
			if compat.IsWindows {
				return compat.DefaultWindowsDirPerm
			} else {
				return compat.DefaultUnixDirPerm
			}
		} else {
			if compat.IsWindows {
				return compat.DefaultWindowsFilePerm
			} else {
				return compat.DefaultUnixFilePerm
			}
		}
	}

	if compat.IsWasip1 {
		if compat.IsTinygo {
			// Tinygo's os.Stat() returns mode 0o000
			return os.FileMode(0o000)
		} else {
			return perm & 0o700
		}
	}

	return perm
}

func fixPosixPerms(perm os.FileMode, isDir bool) os.FileMode {
	if compat.IsWindows {
		if isDir {
			return compat.DefaultWindowsDirPerm
		} else {
			return compat.DefaultWindowsFilePerm
		}
	}

	return fixPerms(perm, isDir)
}

func log(msg string) { //nolint:unused
	if testing.Verbose() {
		fmt.Println(msg)
	}
}

func logf(format string, a ...any) { //nolint:unused
	if testing.Verbose() {
		fmt.Printf(format, a...)
	}
}

func must(err error) { // nolint:unused // quiet linter
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(name string, args ...string) error { //nolint:unused
	log("Executing: " + name + " " + strings.Join(args, " "))
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, name, args...) //nolint:gosec
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = io.NopCloser(bytes.NewReader(nil))
	return cmd.Run()
}

func runCapture(name string, args ...string) (string, error) {
	log("Executing: " + name + " " + strings.Join(args, " "))
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, name, args...) //nolint:gosec
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	cmd.Stdin = io.NopCloser(bytes.NewReader(nil))
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s %v: %w\nstderr:\n%s", name, args, err, errb.String())
	}
	return out.String(), nil
}

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

func supportsHardLinks(t *testing.T) bool {
	t.Helper()

	if !compat.SupportsLinks() {
		skip(t, "Skipping test: Links() not supported on "+runtime.GOOS)

		return false // tinygo doesn't support t.Skip
	}

	if testEnv.noHardLinks {
		skipf(t, "Skipping test: hard links are not supported on a %v filesystem", testEnv.fsType)

		return false // tinygo doesn't support t.Skip
	}

	return true
}

func supportsSymlinks(t *testing.T) bool {
	t.Helper()

	if !compat.SupportsSymlinks() {
		skipf(t, "Skipping test: symlinks are not supported on %v", runtime.GOOS)

		return false // tinygo doesn't support t.Skip
	}

	if testEnv.noSymlinks {
		skipf(t, "Skipping test: symlinks are not supported on a %v filesystem", testEnv.fsType)

		return false // tinygo doesn't support t.Skip
	}

	return true
}

func tempFile(t *testing.T) (string, error) {
	t.Helper()

	f, err := compat.CreateTemp(tempDir(t), "")
	if err != nil {
		return "", err
	}

	name := f.Name()

	err = f.Close()
	if err != nil {
		return "", err
	}

	return name, nil
}

func tempName(t *testing.T) (string, error) {
	t.Helper()

	name, err := tempFile(t)
	if err != nil {
		return "", err
	}

	err = os.Remove(name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func tempDir(t *testing.T) string {
	t.Helper()

	if tempPath != "" {
		return tempPath
	}

	return t.TempDir()
}
