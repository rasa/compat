// SPDX-FileCopyrightText: Copyright (c) 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build darwin && !tinygo

package compat_test

import (
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	defaultMountBase  = "/Volumes"
	defaultMountPoint = defaultMountBase + "/"
)

var fsTests = []fsTest{
	{nativeFS, testVars{}},
	{"APFS", testVars{}},
	{"HFS+", testVars{}},
	{"HFS+J", testVars{}},
	{"HFSX", testVars{}},
	{"JHFS+", testVars{}},
	{"JHFS+X", testVars{}},
	{"UDF", testVars{false, false, false, 86401, 0, 0, 0, 0, ""}},
	{"ExFAT", testVars{true, true, true, 2, 2, 2, 2, -1, ""}},
	{"FAT32", testVars{true, true, true, 86400, 2, 2, 2, -1, ""}},
}

func testMain(m *testing.M, fsToTest, nativeFSType, fsPath string) int { //nolint:gocyclo
	if tempPath != "" && !strings.HasSuffix(tempPath, "/") {
		tempPath += "/"
	}

	code := -1
	supported := []string{allFS}

	workdir, err := os.MkdirTemp("", "compat-fs-*")
	if err != nil {
		fmt.Printf("cannot create temp workdir: %v\n", err)
		return 1
	}
	defer os.RemoveAll(workdir)

	for i, fsTest := range fsTests {
		supported = append(supported, fsTest.fsName)
		fsNameUpper := strings.ToUpper(fsTest.fsName)
		fsToTestUpper := strings.ToUpper(fsToTest)
		if fsToTest != "" && fsToTest != allFS && fsToTestUpper != fsNameUpper {
			continue
		}

		if testing.Short() && code != -1 {
			break
		}

		fsName := fsTest.fsName
		if fsTest.fsName == nativeFS {
			fsTest.vars.fsType = testEnv.fsType
			fsName += " (" + nativeFSType + ")"
		}
		testEnv = fsTest.vars

		if fsTest.fsName != nativeFS {
			testEnv.fsType = fsTest.fsName
			tempPath = getMountPoint()
		}
		if fsPath != "" {
			tempPath = fsPath
		}

		mountPath := tempPath
		if mountPath == "" {
			mountPath = os.TempDir()
		}

		fmt.Printf("%d/%d: Testing on %v filesystem mounted on %v\n", i+1, len(fsTests), fsName, mountPath)

		if fsTest.fsName == nativeFS {
			code = m.Run()
			if code != 0 {
				return code
			}
			continue
		}

		spec, ok := mkfsSpecFor(fsTest.fsName)
		if !ok {
			fmt.Printf("Skipping testing on %v: unsupported on macOS\n", fsTest.fsName)
			code = 0
			continue
		}
		missing := 0
		for _, bin := range []string{"hdiutil"} {
			if _, err := exec.LookPath(bin); err != nil {
				fmt.Printf("Skipping testing on %v: missing tool %q\n", fsTest.fsName, bin)
				missing++
			}
		}
		if missing > 0 {
			code = 0
			continue
		}

		imgPath := filepath.Join(workdir, "img-"+strings.ToLower(fsTest.fsName)+".sparseimage")
		_, err := runCapture("hdiutil", "create",
			"-size", normalizeSize(tempSize),
			"-type", "SPARSE",
			"-fs", spec.fsArg,
			"-volname", spec.volname,
			"-ov",
			imgPath)
		if err != nil {
			fmt.Printf("Skipping testing on %v: hdiutil create: %v\n", fsTest.fsName, err)
			code = 0
			continue
		}

		mntBase := defaultMountBase
		_, err = os.Stat(mntBase)
		if err != nil {
			mntBase = workdir
		}
		mnt, err := os.MkdirTemp(mntBase, "mnt-*")
		if err != nil {
			_ = os.Remove(imgPath)
			fmt.Printf("Skipping testing on %v: mkdir: %v\n", fsTest.fsName, err)
			code = 0
			continue
		}

		out, err := runCapture("hdiutil", "attach", "-nobrowse", "-owners", "on", "-mountpoint", mnt, imgPath)
		if err != nil {
			_ = os.RemoveAll(mnt)
			_ = os.Remove(imgPath)
			fmt.Printf("Skipping testing on %v: hdiutil attach: %v\n", fsTest.fsName, err)
			code = 0
			continue
		}
		dev := parseDiskFromAttach(out)
		if dev == "" {
			_ = os.RemoveAll(mnt)
			_ = os.Remove(imgPath)
			fmt.Printf("Skipping testing on %v: could not parse disk node from attach output\n", fsTest.fsName)
			code = 0
			continue
		}

		fsTest.vars.fsType = fsTest.fsName
		tempPath = mnt
		testEnv = fsTest.vars

		code = m.Run()

		_, _ = runCapture("hdiutil", "detach", dev)
		_ = os.RemoveAll(mnt)
		_ = os.Remove(imgPath)

		if code != 0 {
			return code
		}
	}

	if code == 0 {
		return 0
	}
	fmt.Printf("Unsupported filesystem: %q; use one of %v\n", fsToTest, strings.Join(supported, ","))

	return 1
}

func getMountPoint() string {
	base := defaultMountBase
	if _, err := os.Stat(base); err != nil {
		base = os.TempDir()
	}
	return filepath.Join(base, "compat-fs-"+randomBase36String(8)) + "/"
}

func randomBase36String(n int) string {
	const base36 = "0123456789abcdefghijklmnopqrstuvwxyz"
	out := make([]byte, n)
	for i := range out {
		out[i] = base36[rand.IntN(len(base36))] //nolint:gosec
	}
	return string(out)
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

type mkSpec struct {
	fsArg   string // value to pass to: hdiutil create -fs <fsArg>
	volname string // value for -volname
}

func mkfsSpecFor(fsName string) (mkSpec, bool) {
	name := strings.ToUpper(fsName)
	switch name {
	case "APFS", "HFS+", "HFS+J", "HFSX", "JHFS+", "JHFS+X", "UDF":
		return mkSpec{fsArg: name, volname: "compatfs"}, true
	case "EXFAT":
		return mkSpec{fsArg: "ExFAT", volname: "compatfs"}, true
	case "FAT32":
		return mkSpec{fsArg: "MS-DOS FAT32", volname: "COMPATFS"}, true
	case "FAT":
		return mkSpec{fsArg: "MS-DOS FAT16", volname: "COMPATFS"}, true
	default:
		return mkSpec{}, false
	}
}

// Extract /dev/diskN from hdiutil attach output (works with/without -mountpoint).
func parseDiskFromAttach(out string) string {
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) > 0 && strings.HasPrefix(fields[0], "/dev/disk") {
			// Return base node (/dev/diskN), trim any slice suffix like /dev/diskNs1
			dev := fields[0]
			if i := strings.Index(dev, "s"); i > 0 {
				return dev[:i]
			}
			return dev
		}
	}
	return ""
}
