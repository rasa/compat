// SPDX-FileCopyrightText: Copyright (c) 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
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

type errReader struct{}

func (errReader) Read(_ []byte) (int, error) {
	return 0, errors.New("simulated read failure")
}

type nilReader struct{}

func (nilReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 'x'
	}

	return len(p), nil
}

func init() {
	// Needed for testing.Verbose() and testing.Short() to be available.
	testing.Init()
	flag.Parse()

	// @TODO(rasa): test different umask settings
	compat.Umask(0)
}

func cleanup(t *testing.T, names ...string) {
	t.Helper()

	for _, name := range names {
		t.Cleanup(removeItFunc(name))
	}
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
	granularity++

	return a.Sub(b).Abs() <= time.Duration(granularity)*time.Second
}

func debugln(t *testing.T, msg string) {
	t.Helper()

	if testing.Verbose() && strings.Contains(compatDebug, "DEBUG") {
		fmt.Println(msg)
	}
}

func debugf(t *testing.T, format string, a ...any) {
	t.Helper()

	debugln(t, fmt.Sprintf(format, a...))
}

func fatal(t *testing.T, msg any) {
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

func fatalf(t *testing.T, format string, a ...any) {
	t.Helper()

	fatal(t, fmt.Sprintf(format, a...))
}

func fatalTimes(t *testing.T, prefix string, got, want time.Time, granularity int) {
	t.Helper()

	diff := got.Sub(want).Abs().Seconds()

	t.Fatalf("%v: got %fs difference, want <%ds (%v vs %v)", prefix, diff, granularity, got, want)
}

func fclose(f *os.File) {
	if f != nil {
		_ = f.Close()
	}
}

type semanticVersion struct {
	major int
	minor int
	patch int
}

var osVersion semanticVersion

func init() {
	osVersion, _ = getOSVersion()
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
			if osVersion.major != 13 {
				return compat.DefaultAppleDirPerm
			}

			fallthrough
		default:
			return compat.DefaultUnixDirPerm
		}
	}

	switch {
	case compat.IsWindows:
		return compat.DefaultWindowsFilePerm
	case compat.IsApple:
		if osVersion.major != 13 {
			return compat.DefaultAppleFilePerm
		}

		fallthrough
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

func log(msg string) {
	if testing.Verbose() {
		fmt.Println(msg)
	}
}

func logf(format string, a ...any) {
	if testing.Verbose() {
		fmt.Printf(format, a...)
	}
}

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func normalizeSize(s string) string {
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

func randomBase36String(n int) string {
	const base36 = "0123456789abcdefghijklmnopqrstuvwxyz"

	out := make([]byte, n)
	for i := range out {
		out[i] = base36[rand.IntN(len(base36))]
	}

	return string(out)
}

func removeIt(name string) {
	fi, err := os.Stat(name)
	if errors.Is(err, os.ErrNotExist) {
		return
	}

	if strings.Contains(compatDebug, "NODEL") {
		return
	}

	if fi.IsDir() {
		_ = filepath.WalkDir(name, func(path string, _ fs.DirEntry, err error) error {
			if err != nil {
				// ignore errors
				return nil //nolint:nilerr
			}

			_ = compat.Chmod(path, 0o777, compat.WithReadOnlyMode(compat.ReadOnlyModeClear))

			return nil
		})
	}

	_ = compat.RemoveAll(name)

	_, err = os.Stat(name)
	if errors.Is(err, os.ErrNotExist) {
		return
	}

	if compat.IsWindows {
		args := []string{name, "/q", "/t", "/c", "/grant", os.Getenv("USERNAME") + ":F"}
		_ = exec.CommandContext(context.Background(), "icacls.exe", args...).Run()
		_ = compat.RemoveAll(name)
	}
}

func removeItFunc(name string) func() {
	return func() {
		removeIt(name)
	}
}

func run(name string, args ...string) error {
	log("Executing: " + name + " " + strings.Join(args, " "))

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = io.NopCloser(bytes.NewReader(nil))

	return cmd.Run()
}

func runCapture(name string, args ...string) (string, error) {
	log("Executing: " + name + " " + strings.Join(args, " "))

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, name, args...)

	var out, errb bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &errb
	cmd.Stdin = io.NopCloser(bytes.NewReader(nil))

	err := cmd.Run()
	if err != nil {
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

	if os.Getenv("ACT") != "" {
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

	if compat.IsTinygo {
		skip(t, "Skipping test: hard links are not supported on tinygo")

		return false // tinygo doesn't support t.Skip
	}

	if !compat.SupportsLinks() {
		t.Skipf("Skipping test: Links() not supported on %v", runtime.GOOS)

		return false // tinygo doesn't support t.Skip
	}

	if testEnv.noHardLinks {
		t.Skipf("Skipping test: hard links are not supported on a %v filesystem", testEnv.fsType)

		return false // tinygo doesn't support t.Skip
	}

	// probe Android SELinux restriction introduced in Android 7/API 24
	dir := t.TempDir()
	src := filepath.Join(dir, "source")
	dst := filepath.Join(dir, "link")

	err := os.WriteFile(src, nil, 0o600)
	if err != nil {
		t.Fatalf("probe hard-link support on %v: %v", runtime.GOOS, err)

		return false // tinygo doesn't support t.Fatalf
	}

	err = os.Link(src, dst)
	if err == nil {
		return true
	}

	err = compat.NormalizeUnsupportedError(err)

	if errors.Is(err, syscall.EPERM) ||
		errors.Is(err, syscall.EACCES) ||
		errors.Is(err, errors.ErrUnsupported) {
		t.Logf("hard links unavailable on %v: %v", runtime.GOOS, err)

		return false
	}

	t.Fatalf("unexpected hard-link probe failure: %v", err)

	return false
}

func supportsSymlinks(t *testing.T) bool {
	t.Helper()

	if compat.IsTinygo {
		skip(t, "Skipping test: symlinks are not supported on tinygo")

		return false // tinygo doesn't support t.Skip
	}

	if !compat.SupportsSymlinks() {
		t.Skipf("Skipping test: symlinks are not supported on %v", runtime.GOOS)

		return false // tinygo doesn't support t.Skip
	}

	if testEnv.noSymlinks {
		t.Skipf("Skipping test: symlinks are not supported on a %v filesystem", testEnv.fsType)

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

	err = f.Close()
	if err != nil {
		return "", err
	}

	name := f.Name()

	_, err = os.Stat(name)
	if errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	return name, nil
}

func tempName(t *testing.T) string {
	t.Helper()

	return filepath.Join(tempDir(t), randomBase36String(8)+".tmp")
}

func tempDir(t *testing.T) string {
	t.Helper()

	tempDir := t.TempDir()

	if tempPath != "" {
		parts := strings.Split(t.TempDir(), string(os.PathSeparator))

		idx := -1

		for i, p := range parts {
			if strings.HasPrefix(p, "Test") {
				idx = i

				break
			}
		}

		if idx == -1 {
			t.Fatalf("no component starting with 'Test' found in %v", t.TempDir())
		}

		tempDir = filepath.Join(append([]string{tempPath, "tmp"}, parts[idx:]...)...)

		err := compat.MkdirAll(tempDir, perm777)
		if err != nil {
			t.Fatal(err)
		}

		_, err = os.Stat(tempDir)
		if errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		}
	}

	if compat.IsPlan9 {
		err := os.Chmod(tempDir, 0o777)
		if err != nil {
			t.Fatal(err)
		}
	}

	return tempDir
}
