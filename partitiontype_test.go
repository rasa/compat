// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"

	"github.com/rasa/compat"
)

func TestPartitionType(t *testing.T) {
	f, err := os.CreateTemp(tempDir(t), "")
	if err != nil {
		t.Fatal(err)
	}

	name := f.Name()
	_ = f.Close()

	testPartitionType(t, name)
}

func TestPartitionTypeRel(t *testing.T) {
	dir := tempDir(t)

	f, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Fatal(err)
	}

	t.Chdir(dir)

	name := filepath.Base(f.Name())
	_ = f.Close()

	testPartitionType(t, name)
}

func TestPartitionTypePrefix(t *testing.T) {
	if !compat.IsWindows {
		skip(t, "Skipping test: requires Windows")

		return
	}

	f, err := os.CreateTemp(tempDir(t), "")
	if err != nil {
		t.Fatal(err)
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

	usr, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}

	dir := tempDir(t)
	ctx := context.Background()
	sharename := randomBase36String(8)
	args := []string{"share", sharename + "=" + dir, "/grant:" + usr.Username + ",READ"}

	err = exec.CommandContext(ctx, "net.exe", args...).Run()
	if err != nil {
		t.Fatal(err)
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

func TestPartitionTypeInvalid(t *testing.T) {
	_, err := compat.PartitionType(context.Background(), invalidName)
	if err == nil {
		t.Fatal("got nil, want an error")
	}
}

func testPartitionType(t *testing.T, name string) {
	t.Helper()

	ctx := context.Background()

	partitionType, err := compat.PartitionType(ctx, name)
	if err != nil {
		var errno syscall.Errno
		errorNo := ""
		if errors.As(err, &errno) {
			errorNo = fmt.Sprintf("(error 0x%x) ", uint32(errno))
		}
		if errors.Is(err, errors.ErrUnsupported) {
			t.Skipf("Skipping test on %v/%v: %v %v(ErrUnsupported)", runtime.GOOS, runtime.GOARCH, err, errorNo)
		}
		if errors.Is(err, syscall.EACCES) {
			t.Skipf("Skipping test on %v/%v: %v %v(EACCES)", runtime.GOOS, runtime.GOARCH, err, errorNo)
		}
		if strings.Contains(err.Error(), "not implemented") {
			t.Skipf("Skipping test on %v/%v: %v %v(not implemented)", runtime.GOOS, runtime.GOARCH, err, errorNo)
		}
		if strings.Contains(err.Error(), "permission denied") {
			t.Skipf("Skipping test on %v/%v: %v %v(permission denied)", runtime.GOOS, runtime.GOARCH, err, errorNo)
		}

		t.Fatal(err)
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
