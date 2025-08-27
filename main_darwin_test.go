// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build darwin && !ios && !tinygo

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
	{"exFAT", testVars{true, true, true, 2, 2, 2, 2, -1, ""}},
	{"FAT32", testVars{true, true, true, 86400, 2, 2, 2, -1, ""}},
	{"FAT", testVars{true, true, true, 86400, 2, 2, 2, -1, ""}}, // aka FAT16
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
		if os.Geteuid() != 0 {
			fmt.Printf("Skipping testing on %v: must run as root on macOS (hdiutil/newfs/diskutil)\n", fsTest.fsName)
			code = 0
			continue
		}
		missing := 0
		for _, bin := range []string{"hdiutil", "diskutil", spec.tool} {
			_, err = exec.LookPath(bin)
			if err != nil {
				fmt.Printf("Skipping testing on %v: missing tool %q\n", fsTest.fsName, bin)
				missing++
			}
		}
		if missing > 0 {
			code = 0
			continue
		}

		// Create sparse image (fs=none), attach (nomount) to get /dev/diskN[/s1],
		// run newfs_*, then diskutil mount -mountPoint.
		imgPath := filepath.Join(workdir, "img-"+strings.ToLower(fsTest.fsName)+".sparseimage")

		// 1) Create image
		_, err = runCapture("hdiutil", "create",
			"-size", normalizeSize(tempSize),
			"-type", "SPARSE",
			"-fs", "none",
			"-ov",
			imgPath)
		if err != nil {
			fmt.Printf("Skipping testing on %v: hdiutil create: %v\n", fsTest.fsName, err)
			code = 0
			continue
		}

		// 2) Attach without auto-mount; capture device node
		out, err := runCapture("hdiutil", "attach", "-nomount", "-readwrite", imgPath)
		if err != nil {
			_ = os.Remove(imgPath)
			fmt.Printf("Skipping testing on %v: hdiutil attach: %v\n", fsTest.fsName, err)
			code = 0
			continue
		}
		dev := parseDiskFromAttach(out)
		if dev == "" {
			_ = run("hdiutil", "detach", imgPath) // best-effort
			_ = os.Remove(imgPath)
			fmt.Printf("Skipping testing on %v: cannot find disk node in attach output\n", fsTest.fsName)
			code = 0
			continue
		}
		raw := strings.Replace(dev, "/dev/disk", "/dev/rdisk", 1)

		// 3) Format via newfs_*
		args := append(append([]string{}, spec.args...), spec.label, raw)
		if _, err = runCapture(spec.tool, args...); err != nil {
			_ = run("hdiutil", "detach", dev)
			_ = os.Remove(imgPath)
			fmt.Printf("Skipping testing on %v: %s failed: %v\n", fsTest.fsName, spec.tool, err)
			code = 0
			continue
		}

		// 4) Mount at our mount point
		mntBase := defaultMountBase
		_, err = os.Stat(mntBase)
		if err != nil {
			mntBase = workdir
		}
		mnt, err := os.MkdirTemp(mntBase, "mnt-*")
		if err != nil {
			_ = run("hdiutil", "detach", dev)
			_ = os.Remove(imgPath)
			fmt.Printf("Skipping testing on %v: mkdir: %v\n", fsTest.fsName, err)
			code = 0
			continue
		}
		_, err = runCapture("diskutil", "mount", "-mountPoint", mnt, dev)
		if err != nil {
			_ = os.RemoveAll(mnt)
			_ = run("hdiutil", "detach", dev)
			_ = os.Remove(imgPath)
			fmt.Printf("Skipping testing on %v: diskutil mount: %v\n", fsTest.fsName, err)
			code = 0
			continue
		}

		// Run tests against mounted volume
		fsTest.vars.fsType = fsTest.fsName
		tempPath = mnt
		testEnv = fsTest.vars

		code = m.Run()

		// Teardown
		_, _ = runCapture("diskutil", "unmount", mnt)
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

// Keep mkSpec name/signature; map to macOS newfs_* tools. fstype is informational.
type mkSpec struct {
	tool   string
	args   []string
	fstype string
	label  string
}

func mkfsSpecFor(fsName string) (mkSpec, bool) {
	switch strings.ToUpper(fsName) {
	case "APFS":
		return mkSpec{tool: "newfs_apfs", args: []string{"-v"}, fstype: "apfs", label: "compatfs"}, true
	case "HFS+":
		return mkSpec{tool: "newfs_hfs", args: []string{"-v"}, fstype: "hfs", label: "compatfs"}, true
	case "EXFAT":
		return mkSpec{tool: "newfs_exfat", args: []string{"-v"}, fstype: "exfat", label: "compatfs"}, true
	case "FAT32":
		return mkSpec{tool: "newfs_msdos", args: []string{"-F", "32", "-v"}, fstype: "msdos", label: "COMPATFS"}, true
	case "FAT":
		return mkSpec{tool: "newfs_msdos", args: []string{"-F", "16", "-v"}, fstype: "msdos", label: "COMPATFS"}, true
	default:
		// Everything else (Btrfs, ext*, F2FS, ReiserFS, XFS, NTFS) is unsupported natively.
		return mkSpec{}, false
	}
}

// Parse /dev/diskN from `hdiutil attach -nomount` output.
func parseDiskFromAttach(out string) string {
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) > 0 && strings.HasPrefix(fields[0], "/dev/disk") {
			dev := fields[0]
			if i := strings.Index(dev, "s"); i > 0 { // strip any slice suffix
				return dev[:i]
			}
			return dev
		}
	}
	return ""
}
