// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package main

import "testing"

func Test_SameFile(t *testing.T) {
	want := "true"
	if got := sameFile(); got != want {
		t.Errorf("sameFile(): got %v, want %v", got, want)
	}
}

func Test_SameFiles(t *testing.T) {
	want := "true"
	if got := sameFiles(); got != want {
		t.Errorf("sameFiles(): got %v, want %v", got, want)
	}
}
