// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"context"
	"os"
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

func TestPartitionTypeBad(t *testing.T) {
	name := "/an/invalid/filename"
	ctx := context.Background()
	_, err := compat.PartitionType(ctx, name)
	if err == nil {
		t.Fatalf("got not error for invalid file %q", name)
	}
}
