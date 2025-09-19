// SPDX-FileCopyrightText: Copyright (c) 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/rasa/compat"
)

const (
	perm000     = os.FileMode(0)
	perm100     = os.FileMode(0o100)
	perm200     = os.FileMode(0o200)
	perm400     = os.FileMode(0o400)
	perm555     = os.FileMode(0o555)
	perm644     = os.FileMode(0o644)
	perm600     = os.FileMode(0o600)
	perm700     = os.FileMode(0o700)
	perm777     = os.FileMode(0o777)
	invalidName = "\x00/a/name/with/an/embedded/\x00/byte"
)

var (
	compatDebug = strings.ToUpper(os.Getenv("COMPAT_DEBUG"))
	helloBytes  = []byte("hello")
	helloBuf    = bytes.NewBuffer(helloBytes)
)

func init() {
	// Needed for testing.Verbose() and testing.Short() to be available.
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

func compareTimes(a, b time.Time, granularity int) bool {
	if granularity < 0 {
		return a.IsZero()
	}
	// add 1 second for fractional seconds
	granularity += 1

	return a.Sub(b).Abs() < time.Duration(granularity)*time.Second
}

func debugln(t *testing.T, msg string) { //nolint:unused
	t.Helper()

	if testing.Verbose() && strings.Contains(compatDebug, "DEBUG") {
		fmt.Println(msg)
	}
}

func debugf(t *testing.T, format string, a ...any) { //nolint:unused
	t.Helper()

	debugln(t, fmt.Sprintf(format, a...))
}

func fatal(t *testing.T, msg any) { //nolint:unused
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

func fatalf(t *testing.T, format string, a ...any) { //nolint:unused
	t.Helper()

	fatal(t, fmt.Sprintf(format, a...))
}

func fatalTimes(t *testing.T, prefix string, got, want time.Time, granularity int) { //nolint:unused
	t.Helper()

	diff := got.Sub(want).Abs().Seconds()

	t.Fatalf("%v: got %.2fs difference, want <%ds (%v vs %v)", prefix, diff, granularity, got, want)
}

func fclose(f *os.File) {
	if f != nil {
		_ = f.Close()
	}
}

func fixPerms(perm os.FileMode, isDir bool) os.FileMode {
	if compat.IsWasip1 {
		if compat.IsTinygo {
			return perm600
		}
		if isDir {
			return perm700
		}

		return perm600
	}

	if !testEnv.noACLs {
		return perm
	}

	if isDir {
		switch {
		case compat.IsWindows:
			return compat.DefaultWindowsDirPerm
		case compat.IsApple:
			return compat.DefaultAppleDirPerm
		default:
			return compat.DefaultUnixDirPerm
		}
	}

	switch {
	case compat.IsWindows:
		return compat.DefaultWindowsFilePerm
	case compat.IsApple:
		return compat.DefaultAppleFilePerm
	default:
		return compat.DefaultUnixFilePerm
	}
}

func fixPosixPerms(perm os.FileMode, isDir bool) os.FileMode {
	if compat.IsWasip1 {
		if compat.IsTinygo {
			return perm000
		}
		if isDir {
			return perm700
		}

		return perm600
	}

	if compat.IsWindows {
		if isDir {
			return compat.DefaultWindowsDirPerm
		}

		return compat.DefaultWindowsFilePerm
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

func must(err error) { // nolint:unused
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func normalizeSize(s string) string { //nolint:unused
	r := strings.ToUpper(strings.TrimSpace(s))
	r = strings.ReplaceAll(r, "BYTES", "B")
	r = strings.ReplaceAll(r, "IB", "I")
	r = strings.ReplaceAll(r, "KIB", "K")
	r = strings.ReplaceAll(r, "MIB", "M")
	r = strings.ReplaceAll(r, "GIB", "G")
	r = strings.ReplaceAll(r, "TIB", "T")
	r = strings.ReplaceAll(r, "KB", "K")
	r = strings.ReplaceAll(r, "MB", "M")
	r = strings.ReplaceAll(r, "GB", "G")
	r = strings.ReplaceAll(r, "TB", "T")

	return r
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

func partitionType(name string) string {
	partType, err := compat.PartitionType(context.Background(), name)
	if err != nil {
		return "n/a"
	}
	return partType
}

func randomBase36String(n int) string { //nolint:unparam,unused
	const base36 = "0123456789abcdefghijklmnopqrstuvwxyz"
	out := make([]byte, n)
	for i := range out {
		out[i] = base36[rand.IntN(len(base36))] //nolint:gosec
	}
	return string(out)
}

func removeIt(name string) { //nolint:unused
	if os.IsPermission(os.Remove(name)) {
		_ = compat.Chmod(name, perm600)
		_ = compat.Remove(name)
	}
}

func run(name string, args ...string) error { //nolint:unparam,unused
	log("Executing: " + name + " " + strings.Join(args, " "))
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, name, args...) //nolint:gosec
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = io.NopCloser(bytes.NewReader(nil))
	return cmd.Run()
}

func runCapture(name string, args ...string) (string, error) { //nolint:unused
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

	if compat.IsTinygo {
		skip(t, "Skipping test: hard links are not supported on tinygo")

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

	if compat.IsTinygo {
		skip(t, "Skipping test: symlinks are not supported on tinygo")

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
