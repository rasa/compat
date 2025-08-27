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
			skip(t, "Skipping test on "+runtime.GOOS+"/"+runtime.GOARCH+": "+err.Error())

			return
		}

		t.Error(err)

		return
	}
	t.Logf("partitionType=%v (path=%v)", partitionType, name)
}
