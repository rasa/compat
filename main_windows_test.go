// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows && !tinygo

package compat_test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"golang.org/x/sys/windows"
)

// "github.com/davecgh/go-spew/spew"

const defaultMountPoint = `Z:\`

var fsTests = []fsTest{
	{nativeFS, testVars{}}, // assume we are on NTFS
	{"exFAT", testVars{
		noACLs:                  true,
		noSymlinks:              true,
		noHardLinks:             true,
		atimeGranularity:        2,
		btimeGranularity:        2,
		ctimeGranularity:        -1,
		mtimeGranularity:        2,
		btimeSymlinkGranularity: -1,
	}},
	{"FAT32", testVars{
		noACLs:                  true,
		noSymlinks:              true,
		noHardLinks:             true,
		atimeGranularity:        86400,
		btimeGranularity:        2,
		ctimeGranularity:        -1,
		mtimeGranularity:        2,
		btimeSymlinkGranularity: -1,
	}},
	// @TODO(rasa) determine why FAT fails more often than it succeeds.
	// {"FAT", testVars{true, true, 86400, 2, -1, 2, -1, ""}},
	{"NTFS", testVars{}},
	{"ReFS", testVars{}},
}

func testMain(m *testing.M, fsToTest, nativeFSType, fsPath string) int { //nolint:gocyclo
	if tempPath != "" {
		if len(tempPath) < 2 {
			tempPath += ":"
		}
		if len(tempPath) < 3 {
			tempPath += `\`
		}
	}
	fsSize := os.Getenv("COMPAT_DEBUG_FS_SIZE")
	if fsSize != "" {
		tempSize = fsSize
	}

	code := -1
	supported := []string{allFS}

	fsToTestUpper := strings.ToUpper(fsToTest)
	testsToRun := 0
	for _, fsTest := range fsTests {
		supported = append(supported, fsTest.fsName)
		fsNameUpper := strings.ToUpper(fsTest.fsName)
		if fsToTest != "" && fsToTestUpper != strings.ToUpper(allFS) && fsToTestUpper != fsNameUpper {
			continue
		}
		testsToRun++
	}

	n := 0
	for _, fsTest := range fsTests {
		fsNameUpper := strings.ToUpper(fsTest.fsName)
		if fsToTest != "" && fsToTestUpper != strings.ToUpper(allFS) && fsToTestUpper != fsNameUpper {
			continue
		}
		n++

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

		fmt.Printf("%d/%d: Testing on %v filesystem mounted on %v\n", n, testsToRun, fsName, mountPath)

		if fsTest.fsName == "Native" {
			code = m.Run()
			if code != 0 {
				return code
			}

			continue
		}

		tempDrive := string(tempPath[0])

		exe := "powershell.exe"
		exe, _ = exec.LookPath(exe)
		if exe == "" {
			exe, _ = exec.LookPath("pwsh.exe")
		}
		if exe == "" {
			fmt.Printf("Cannot find powershell or pwsh in the PATH\n")
			os.Exit(1)
		}
		args := []string{
			"-file",
			"create-vhdx.ps1",
			tempDrive,
			fsTest.fsName,
			tempSize,
		}
		out, err := runCapture(exe, args...)
		log(out)
		if err == nil {
			code = m.Run()
		}
		if !strings.Contains(compatDebug, "NODEL") {
			args = []string{
				"-file",
				"remove-vhdx.ps1",
			}
			out2, _ := runCapture(exe, args...)
			log(out2)
		}
		if err != nil {
			fmt.Printf("Skipping testing on %v: %v\n", fsTest.fsName, err)
			if !testing.Verbose() {
				fmt.Println(out)
			}
			code = 0
			continue
		}

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
	drives, err := windows.GetLogicalDrives()
	if err != nil {
		fmt.Printf("cannot get list of drives letters: %v, using %v\n", err, defaultMountPoint)
	}

	for i := 25; i >= 0; i-- {
		var mask uint32 = 1 << i
		if drives&mask == 0 {
			return string(rune('A'+i)) + `:\`
		}
	}

	fmt.Printf("no unused drive letter was found, using %v\n", defaultMountPoint)

	return defaultMountPoint
}
