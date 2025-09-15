// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"context"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/rasa/compat"
)

func TestPartitionType(t *testing.T) {
	f, err := os.CreateTemp(tempDir(t), "")
	if err != nil {
		t.Error(err)

		return
	}
	name := f.Name()
	_ = f.Close()
	testPartitionType(t, name)
}

func TestPartitionTypeRel(t *testing.T) {
	dir := tempDir(t)
	f, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Error(err)

		return
	}
	t.Chdir(dir)

	name := filepath.Base(f.Name())
	_ = f.Close()
	testPartitionType(t, name)
}

func TestPartitionTypeBad(t *testing.T) {
	name := "/a/bad/filename/for/partitiontype"
	ctx := context.Background()
	_, err := compat.PartitionType(ctx, name)
	if err == nil {
		t.Fatalf("got not error for invalid file %q", name)
	}
}

func TestPartitionTypePrefix(t *testing.T) {
	if !compat.IsWindows {
		skip(t, "Skipping test: requires Windows")

		return
	}

	f, err := os.CreateTemp(tempDir(t), "")
	if err != nil {
		t.Fatal(err)

		return
	}
	name := `\\?\` + f.Name()
	_ = f.Close()
	testPartitionType(t, name)
}

func TestPartitionTypeUNC(t *testing.T) {
	if !compat.IsWindows {
		skip(t, "Skipping test: requires Windows")

		return
	}

	dir := tempDir(t)
	ctx := context.Background()
	sharename := randomBase36String(8)
	args := []string{"share", sharename + "=" + dir, "/grant:" + currentUsername() + ",READ"}
	err := exec.CommandContext(ctx, "net.exe", args...).Run()
	if err != nil {
		t.Fatal(err)

		return
	}

	defer func() {
		args := []string{"share", sharename, "/del", "/yes"}
		err = exec.CommandContext(ctx, "net.exe", args...).Run()
		if err != nil {
			t.Fatal(err)
		}
	}()

	f, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Fatal(err)

		return
	}

	name := `\\?\UNC\127.0.0.1\` + sharename + `\` + filepath.Base(f.Name())
	_ = f.Close()
	testPartitionType(t, name)
}

func TestPartitionTypeRoot(t *testing.T) {
	if !compat.IsWindows {
		skip(t, "Skipping test: requires Windows")

		return
	}

	systemDrive := os.Getenv("SystemDrive")
	if systemDrive == "" {
		systemDrive = "C:"
	}
	testPartitionType(t, systemDrive)
}

func testPartitionType(t *testing.T, name string) {
	t.Helper()

	ctx := context.Background()
	partitionType, err := compat.PartitionType(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), "not implemented") {
			skipf(t, "Skipping test on %v/%v: %v", runtime.GOOS, runtime.GOARCH, err)

			return
		}

		t.Error(err)

		return
	}
	if testEnv.fsType == "" || testEnv.fsType == nativeFS {
		return
	}
	fsType := strings.ToLower(testEnv.fsType)
	if !strings.Contains(partitionType, fsType) {
		// @TODO change this to Errorf eventually
		t.Logf("PartitionType(): got %v, want %v", partitionType, fsType)
	}
}

func currentUsername() string {
	usr, err := user.Current()
	if err != nil {
		return compat.UnknownUsername
	}

	return usr.Username
}
