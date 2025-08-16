// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compat_test

import (
	"os"
	"testing"
	"time"

	"github.com/rasa/compat"
)

// Source: https://github.com/golang/go/blob/77f911e3/src/io/fs/format_test.go#L13

// formatTest implements FileInfo to test FormatFileInfo,
// and implements DirEntry to test FormatDirEntry.
type formatTest struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (fs *formatTest) Name() string {
	return fs.name
}

func (fs *formatTest) Size() int64 {
	return fs.size
}

func (fs *formatTest) Mode() os.FileMode {
	return fs.mode
}

func (fs *formatTest) ModTime() time.Time {
	return fs.modTime
}

func (fs *formatTest) IsDir() bool {
	return fs.isDir
}

func (fs *formatTest) Sys() any {
	return nil
}

func (fs *formatTest) Type() os.FileMode {
	return fs.mode.Type()
}

func (fs *formatTest) Info() (os.FileInfo, error) {
	return fs, nil
}

func (fs *formatTest) PartitionID() uint64 { return 0 }
func (fs *formatTest) FileID() uint64      { return 0 }
func (fs *formatTest) Links() uint64       { return 1 }
func (fs *formatTest) ATime() time.Time    { return time.Time{} }
func (fs *formatTest) BTime() time.Time    { return time.Time{} }
func (fs *formatTest) CTime() time.Time    { return time.Time{} }
func (fs *formatTest) MTime() time.Time    { return fs.modTime }
func (fs *formatTest) UID() int            { return 0 }
func (fs *formatTest) GID() int            { return 0 }
func (fs *formatTest) User() string        { return "" }
func (fs *formatTest) Group() string       { return "" }
func (fs *formatTest) Error() error        { return nil }
func (fs *formatTest) String() string      { return "" }

var formatTests = []struct {
	input        formatTest
	wantFileInfo string
	wantDirEntry string
}{
	{
		formatTest{
			name:    "hello.go",
			size:    100,
			mode:    0o644,
			modTime: time.Date(1970, time.January, 1, 12, 0, 0, 0, time.UTC),
			isDir:   false,
		},
		"-rw-r--r-- 100 1970-01-01 12:00:00 hello.go",
		"- hello.go",
	},
	{
		formatTest{
			name:    "home/gopher",
			size:    0,
			mode:    os.ModeDir | 0o755,
			modTime: time.Date(1970, time.January, 1, 12, 0, 0, 0, time.UTC),
			isDir:   true,
		},
		"drwxr-xr-x 0 1970-01-01 12:00:00 home/gopher/",
		"d home/gopher/",
	},
	{
		formatTest{
			name:    "big",
			size:    0x7fffffffffffffff,
			mode:    os.ModeIrregular | 0o644,
			modTime: time.Date(1970, time.January, 1, 12, 0, 0, 0, time.UTC),
			isDir:   false,
		},
		"?rw-r--r-- 9223372036854775807 1970-01-01 12:00:00 big",
		"? big",
	},
	{
		formatTest{
			name:    "small",
			size:    -0x8000000000000000,
			mode:    os.ModeSocket | os.ModeSetuid | 0o644,
			modTime: time.Date(1970, time.January, 1, 12, 0, 0, 0, time.UTC),
			isDir:   false,
		},
		"Surw-r--r-- -9223372036854775808 1970-01-01 12:00:00 small",
		"S small",
	},
}

func TestFormatDirEntry(t *testing.T) {
	for i, test := range formatTests {
		di := compat.FileInfoToDirEntry(&test.input)
		got := compat.FormatDirEntry(di)
		if got != test.wantDirEntry {
			t.Errorf("%d: FormatDirEntry(%#v) = %q, want %q", i, test.input, got, test.wantDirEntry)
		}
	}

}
