// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: Portions copyright (c) 2015 Nate Finch (@natefinch)
// SPDX-FileCopyrightText: Portions copyright (c) 2022 Simon Dassow (@sdassow)

package compat_test

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/rasa/compat"
)

func TestWriteFileAtomic(t *testing.T) {
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	err = compat.WriteFileAtomic(file, helloBytes)
	if err != nil {
		fatalf(t, "Failed to write file: %q: %v", file, err)

		return // Tinygo doesn't support T.Fatal
	}

	fi, err := compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	perm := compat.CreateTempPerm // 0o600
	want := fixPerms(perm, false)
	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got %04o, want %04o", got, want)
	}
}

func TestWriteFileAtomicCurrentDir(t *testing.T) {
	file := randomBase36String(8) + ".tmp"

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	err := compat.WriteFileAtomic(file, helloBytes)
	if err != nil {
		fatalf(t, "Failed to write file: %q: %v", file, err)

		return // Tinygo doesn't support T.Fatal
	}

	fi, err := compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	perm := compat.CreateTempPerm // 0o600
	want := fixPerms(perm, false)
	got := fi.Mode().Perm()
	if got != want {
		dir, _ := os.Getwd()
		partType := partitionType(dir)
		if strings.Contains(partType, "fat") || strings.Contains(partType, "dos") {
			skipf(t, "got %04o, want %04o (ignoring: on %v filesystem", got, want, partType)

			return
		}
	
		t.Fatalf("got %04o, want %04o", got, want)
	}
}

func TestWriteFileAtomicDefaultFileMode(t *testing.T) {
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	err = compat.WriteFileAtomic(file, helloBytes, compat.WithDefaultFileMode(perm644))
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	var fi os.FileInfo

	fi, err = compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	want := fixPerms(perm644, false)

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got %04o, want %04o (1)", got, want)
	}
	// check if file mode is preserved
	err = compat.Chmod(file, perm600)
	if err != nil {
		t.Fatalf("Failed to change file mode: %q: %v", file, err)
	}

	err = compat.WriteFileAtomic(file, helloBytes, compat.WithDefaultFileMode(perm644))
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err = compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	want = fixPerms(perm600, false)

	got = fi.Mode().Perm()

	if got != want {
		t.Fatalf("got %04o, want %04o (2)", got, want)
	}
}

func TestWriteFileAtomicKeepFileMode(t *testing.T) { //nolint:dupl
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	perm := perm555

	err = compat.WriteFile(file, helloBytes, perm)
	if err != nil {
		t.Fatalf("Failed to create file: %q: %v", file, err)
	}

	err = compat.WriteFileAtomic(file, helloBytes, compat.KeepFileMode(true))
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err := compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	want := fixPerms(perm, false)
	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got %04o, want %04o: perm=%3o (%v) (1)", got, want, perm, perm)
	}

	err = compat.WriteFileAtomic(file, helloBytes, compat.KeepFileMode(false))
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err = compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	got = fi.Mode().Perm()
	if got == want {
		if perm != want {
			partType := partitionType(file)
			t.Logf("got %v, want %v (ignoring: %v on %v)", got, want, partType, runtime.GOOS)
			return
		}
		t.Fatalf("got %04o, want !%04o (2)", got, want)
	}
}

func TestWriteFileAtomicWithFileMode(t *testing.T) {
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	err = compat.WriteFileAtomic(file, helloBytes, compat.WithFileMode(perm644))
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err := compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	want := fixPerms(perm644, false)

	got := fi.Mode().Perm()
	if got != want {
		t.Fatalf("got %04o, want %04o (1)", got, want)
	}
	// ensure previous file mode is ignored
	err = compat.Chmod(file, perm600)
	if err != nil {
		t.Fatalf("Failed to change file mode: %q: %v", file, err)
	}

	err = compat.WriteFileAtomic(file, helloBytes, compat.WithFileMode(perm644))
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err = compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	got = fi.Mode().Perm()
	if got != want {
		t.Fatalf("got %04o, want %04o (2)", got, want)
	}
}

func TestWriteFileAtomicReadOnlyModeReset(t *testing.T) {
	if !compat.IsWindows {
		skip(t, "Skipping test: requires Windows")

		return
	}

	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		_ = compat.Chmod(file, perm600)
		_ = os.Remove(file)
	})

	err = compat.WriteFileAtomic(file, helloBytes, compat.WithFileMode(perm400), compat.WithReadOnlyMode(compat.ReadOnlyModeReset))
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err := os.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	want := true // user-writable bit is set.
	got := fi.Mode().Perm()&perm200 == perm200
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestWriteFileAtomicInvalid(t *testing.T) {
	err := compat.WriteFileAtomic(invalidName, helloBytes)
	if err == nil {
		t.Fatalf("got nil, want an error")
	}
}

func TestWriteFileAtomicCantRead(t *testing.T) {
	file, err := tempFile(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		_ = compat.Chmod(file, perm600)
		_ = os.Remove(file)
	})

	perm := fixPerms(perm100, false)
	if perm != perm100 {
		partType := partitionType(file)
		skipf(t, "Skipping test: ACLs are not supported on a %v filesystem", partType)

		return
	}
	err = compat.Chmod(file, perm)
	if err != nil {
		t.Fatalf("Chmod: %v", err)
	}

	err = compat.WriteFileAtomic(file, helloBytes, compat.KeepFileMode(true))
	if err != nil {
		fatalf(t, "WriteFileAtomic: %v", err)

		return // Tinygo doesn't support T.Fatal
	}
}
