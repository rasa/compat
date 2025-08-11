// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"context"
	"os"
	"testing"

	"github.com/rasa/compat"
)

func TestPartitionType(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "")
	if err != nil {
		t.Error(err)

		return
	}
	name := f.Name()
	_ = f.Close()
	ctx := context.Background()
	partitionType, err := compat.PartitionType(ctx, name)
	if err != nil {
		t.Error(err)

		return
	}
	t.Logf("partitionType=%v (path=%v)", partitionType, name)
}
