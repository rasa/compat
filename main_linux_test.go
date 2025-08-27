// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build linux && !android && !tinygo

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
	defaultMountBase  = "/mnt"
	defaultMountPoint = defaultMountBase + "/"
)

var fsTests = []fsTest{
	{nativeFS, testVars{}},
	{"btrfs", testVars{}},
	{"ext2", testVars{}},
	{"ext3", testVars{}},
	{"ext4", testVars{}},
	{"f2fs", testVars{false, false, false, 0, -1, 0, 0, -1, ""}},
	{"reiserfs", testVars{false, false, false, 0, -1, 0, 0, -1, ""}},
	{"xfs", testVars{}},
	//
	{"exFAT", testVars{true, true, true, 2, 2, 2, 2, -1, ""}},
	{"FAT32", testVars{true, true, true, 86400, 2, 2, 2, -1, ""}},
	{"FAT", testVars{true, true, true, 86400, 2, 2, 2, -1, ""}}, // aka FAT16
	{"NTFS", testVars{true, false, false, 0, 0, 0, 0, -1, ""}}, // requires ntfs-3g/ntfsprogs
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

		if fsTest.fsName == "Native" {
			code = m.Run()
			if code != 0 {
				return code
			}

			continue
		}

		spec, ok := mkfsSpecFor(fsTest.fsName)
		if !ok {
			fmt.Printf("Skipping testing on %v: unsupported on Linux\n", fsTest.fsName)
			code = 0
			continue
		}
		_, err := exec.LookPath(spec.tool)
		if err != nil {
			fmt.Printf("Skipping testing on %v: missing tool %q\n", fsTest.fsName, spec.tool)
			code = 0
			continue
		}
		if os.Geteuid() != 0 {
			fmt.Printf("Skipping testing on %v: must run as root for loop/mount\n", fsTest.fsName)
			code = 0
			continue
		}

		imgPath := filepath.Join(workdir, "img-"+strings.ToLower(fsTest.fsName)+".bin")
		err = allocateImage(imgPath, normalizeSize(tempSize))
		if err != nil {
			fmt.Printf("Skipping testing on %v: %v\n", fsTest.fsName, err)
			code = 0
			continue
		}

		loopDev, err := runCapture("losetup", "-f", "--show", imgPath)
		if err != nil {
			fmt.Printf("Skipping testing on %v: losetup: %v\n", fsTest.fsName, err)
			_ = os.Remove(imgPath)
			code = 0
			continue
		}
		loopDev = strings.TrimSpace(loopDev)

		_, err = runCapture(spec.tool, append(spec.args, loopDev)...)
		if err != nil {
			fmt.Printf("Skipping testing on %v: mkfs failed: %v\n", fsTest.fsName, err)
			_ = run("losetup", "-d", loopDev)
			_ = os.Remove(imgPath)
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
			fmt.Printf("Skipping testing on %v: mkdir: %v\n", fsTest.fsName, err)
			_ = run("losetup", "-d", loopDev)
			_ = os.Remove(imgPath)
			code = 0
			continue
		}

		_, mountErr := runCapture("mount", "-t", spec.fstype, loopDev, mnt)
		if mountErr != nil && spec.fstype == "ntfs3" {
			_, mountErr = runCapture("mount", "-t", "ntfs", loopDev, mnt)
		}
		if mountErr != nil {
			fmt.Printf("Skipping testing on %v: mount: %v\n", fsTest.fsName, mountErr)
			_ = os.RemoveAll(mnt)
			_ = run("losetup", "-d", loopDev)
			_ = os.Remove(imgPath)
			code = 0
			continue
		}

		fsTest.vars.fsType = fsTest.fsName
		tempPath = mnt
		testEnv = fsTest.vars

		runCode := m.Run()

		_ = run("umount", mnt)
		_ = run("losetup", "-d", loopDev)
		_ = os.RemoveAll(mnt)
		_ = os.Remove(imgPath)

		if runCode != 0 {
			return runCode
		}
		code = 0
	}

	if code == 0 {
		return 0
	}
	fmt.Printf("Unsupported filesystem: %q; use one of %v\n", fsToTest, strings.Join(supported, ","))

	return 1
}

func getMountPoint() string {
	base := defaultMountBase
	_, err := os.Stat(base)
	if err != nil {
		base = os.TempDir()
	}

	return filepath.Join(base, "compat-fs-"+randomBase36String(8)) + "/"
}

func randomBase36String(n int) string {
	const base36 = "0123456789abcdefghijklmnopqrstuvwxyz"
	out := make([]byte, n)
	for i := range out {
		out[i] = base36[rand.IntN(len(base36))] //nolint: gosec
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

func allocateImage(path, size string) error {
	if _, err := runCapture("fallocate", "-l", size, path); err == nil {
		return nil
	}

	return run("truncate", "-s", size, path)
}

type mkSpec struct {
	tool   string
	args   []string
	fstype string // passed to mount -t
}

func mkfsSpecFor(fsName string) (mkSpec, bool) {
	switch strings.ToLower(fsName) {
	case "btrfs":
		return mkSpec{tool: "mkfs.btrfs", args: []string{"-f"}, fstype: "btrfs"}, true
	case "ext2":
		return mkSpec{tool: "mkfs.ext2", args: []string{"-F"}, fstype: "ext2"}, true
	case "ext3":
		return mkSpec{tool: "mkfs.ext3", args: []string{"-F"}, fstype: "ext3"}, true
	case "ext4":
		return mkSpec{tool: "mkfs.ext4", args: []string{"-F"}, fstype: "ext4"}, true
	case "f2fs":
		return mkSpec{tool: "mkfs.f2fs", args: []string{"-f"}, fstype: "f2fs"}, true
	case "reiserfs":
		return mkSpec{tool: "mkfs.reiserfs", args: []string{"-f"}, fstype: "reiserfs"}, true
	case "xfs":
		return mkSpec{tool: "mkfs.xfs", args: []string{"-f"}, fstype: "xfs"}, true

	case "exfat":
		return mkSpec{tool: "mkfs.exfat", args: nil, fstype: "exfat"}, true
	case "fat":
		return mkSpec{tool: "mkfs.vfat", args: []string{"-F", "16"}, fstype: "vfat"}, true
	case "fat32":
		return mkSpec{tool: "mkfs.vfat", args: []string{"-F", "32"}, fstype: "vfat"}, true
	case "ntfs":
		// Requires ntfs-3g/ntfsprogs; kernel driver type "ntfs3" (fallback to "ntfs" if needed).
		return mkSpec{tool: "mkfs.ntfs", args: []string{"-F"}, fstype: "ntfs3"}, true
	default:
		fmt.Printf("Unsupported filesystem: %q\n", fsName)
	}

	return mkSpec{}, false
}
