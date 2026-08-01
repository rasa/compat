// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/rasa/compat"
)

func TestGetOptions(t *testing.T) {
	opts := make([]compat.Option, 0)
	opts = append(opts, compat.WithAllowNonAtomicReplace(true))
	opts = append(opts, compat.WithAtomicity(true))
	opts = append(opts, compat.WithDefaultFileMode(perm777))
	opts = append(opts, compat.WithFileMode(perm777))
	opts = append(opts, compat.WithFlags(os.O_CREATE))
	opts = append(opts, compat.WithKeepFileMode(true))
	opts = append(opts, compat.WithReadOnlyMode(compat.ReadOnlyModeSet))
	opts = append(opts, compat.WithRetrySeconds(1))
	opts = append(opts, compat.WithSetSymlinkOwner(true))

	compat.SetOptions(opts...)
	o := compat.GetOptions()
	got := len(o)
	want := reflect.TypeFor[compat.Options]().NumField()
	if got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func TestBuildOptions(t *testing.T) {
	opts := make([]compat.Option, 0)
	fopts := compat.BuildOptions(opts...)
	got := fopts.String()
	want :=
		`allowNonAtomicReplace:      true
atomically:      true
defaultFileMode: 0o777 (-rwxrwxrwx)
fileMode:        0o777 (-rwxrwxrwx)
flags:           0x40
keepFileMode:    true
readOnlyMode:    1
retrySeconds:    1
setSymlinkOwner: true
`
	if got != want {
		t.Fatalf("got:\n---\n%v\n---\nwant:\n---\n%v\n---\n", got, want)
	}
}
