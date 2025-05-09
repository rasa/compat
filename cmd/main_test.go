// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package main

import "testing"

func Test_SameDevice(t *testing.T) {
	want := "true"
	if got := sameDevice(); got != want {
		t.Errorf("sameDevice(): got %v, want %v", got, want)
	}
}

func Test_SameDevicel(t *testing.T) {
	want := "true"
	if got := sameDevicel(); got != want {
		t.Errorf("sameDevicel(): got %v, want %v", got, want)
	}
}

func Test_SameDevices(t *testing.T) {
	want := "true"
	if got := sameDevices(); got != want {
		t.Errorf("sameDevices(): got %v, want %v", got, want)
	}
}

func Test_SameFile(t *testing.T) {
	want := "true"
	if got := sameFile(); got != want {
		t.Errorf("sameFile(): got %v, want %v", got, want)
	}
}

func Test_SameFilel(t *testing.T) {
	want := "true"
	if got := sameFilel(); got != want {
		t.Errorf("sameFilel(): got %v, want %v", got, want)
	}
}

func Test_SameFiles(t *testing.T) {
	want := "true"
	if got := sameFiles(); got != want {
		t.Errorf("sameFiles(): got %v, want %v", got, want)
	}
}
