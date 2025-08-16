// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"testing"

	"github.com/rasa/compat"
)

// Source: https://github.com/golang/go/blob/77f911e3/src/os/read_test.go#L102

func TestReadDir(t *testing.T) {
	dirname := "rumpelstilzchen"
	_, err := compat.ReadDir(dirname)
	if err == nil {
		t.Fatalf("ReadDir %s: error expected, none found", dirname)
	}

	dirname = "."
	list, err := compat.ReadDir(dirname)
	if err != nil {
		t.Fatalf("ReadDir %s: %v", dirname, err)
	}

	foundFile := false
	foundSubDir := false
	for _, dir := range list {
		switch {
		case !dir.IsDir() && dir.Name() == "dir_test.go":
			foundFile = true
		case dir.IsDir() && dir.Name() == "golang":
			foundSubDir = true
		}
	}
	if !foundFile {
		t.Fatalf("ReadDir %s: dir_test.go file not found", dirname)
	}
	if !foundSubDir {
		t.Fatalf("ReadDir %s: golang directory not found", dirname)
	}
}
