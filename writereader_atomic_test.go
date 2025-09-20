// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: Portions copyright (c) 2015 Nate Finch (@natefinch)
// SPDX-FileCopyrightText: Portions copyright (c) 2022 Simon Dassow (@sdassow)

package compat_test

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/rasa/compat"
)

func TestWriteReaderAtomic(t *testing.T) {
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	err = compat.WriteReaderAtomic(file, helloBuf)
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

func TestWriteReaderAtomicCurrentDir(t *testing.T) {
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	dir, base := filepath.Split(file)
	t.Chdir(dir)

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	err = compat.WriteReaderAtomic(base, helloBuf)
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

func TestWriteReaderAtomicDefaultFileMode(t *testing.T) {
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	err = compat.WriteReaderAtomic(file, helloBuf, compat.WithDefaultFileMode(perm644))
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

	err = compat.WriteReaderAtomic(file, helloBuf, compat.WithDefaultFileMode(perm644))
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

func TestWriteReaderAtomicKeepFileMode(t *testing.T) { //nolint:dupl
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

	err = compat.WriteReaderAtomic(file, helloBuf, compat.KeepFileMode(true))
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
}

func TestWriteReaderAtomicKeepFileModeFalse(t *testing.T) { //nolint:dupl
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

	err = compat.WriteReaderAtomic(file, helloBuf, compat.KeepFileMode(false))
	if err != nil {
		t.Fatalf("Failed to write file: %q: %v", file, err)
	}

	fi, err := compat.Stat(file)
	if err != nil {
		t.Fatalf("Failed to stat file: %q: %v", file, err)
	}

	want := fixPerms(perm, false)
	got := fi.Mode().Perm()
	if got == want {
		if perm != want {
			partType := partitionType(file)
			t.Logf("got %v, want %v (ignoring: %v on %v)", got, want, partType, runtime.GOOS)
			return
		}
		t.Fatalf("got %04o, want !%04o (2)", got, want)
	}
}

func TestWriteReaderAtomicWithFileMode(t *testing.T) { //nolint:dupl
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	err = compat.WriteReaderAtomic(file, helloBuf, compat.WithFileMode(perm644))
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

	err = compat.WriteReaderAtomic(file, helloBuf, compat.WithFileMode(perm644))
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

func TestWriteReaderAtomicReadOnlyModeReset(t *testing.T) {
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

	err = compat.WriteReaderAtomic(file, helloBuf, compat.WithFileMode(perm400), compat.WithReadOnlyMode(compat.ReadOnlyModeReset))
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

func TestWriteReaderAtomicInvalid(t *testing.T) {
	err := compat.WriteReaderAtomic(invalidName, helloBuf)
	if err == nil {
		t.Fatalf("got nil, want an error")
	}
}

func TestWriteReaderAtomicCantRead(t *testing.T) {
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

	err = compat.WriteReaderAtomic(file, helloBuf, compat.KeepFileMode(true))
	if err != nil {
		fatalf(t, "WriteReaderAtomic: %v", err)

		return // Tinygo doesn't support T.Fatal
	}
}

type errReader struct{}

func (errReader) Read(p []byte) (int, error) {
	return 0, errors.New("boom: read failure")
}

func TestWriteReaderAtomicError(t *testing.T) {
	file, err := tempName(t)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(file)
	})

	err = compat.WriteReaderAtomic(file, errReader{})
	if err == nil {
		fatal(t, "got nil, want an error")

		return // Tinygo doesn't support T.Fatal
	}
}
