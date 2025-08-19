// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"io/fs"
	"os"
	pathpkg "path"
	"path/filepath"
	"slices"
	"testing"
	"testing/fstest"

	"github.com/rasa/compat"
)

// Source: https://github.com/golang/go/blob/77f911e3/src/io/fs/walk_test.go#L17-L183

// The following code is:
// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

type Node struct {
	name    string
	entries []*Node // nil if the entry is a file
	mark    int
}

var tree = &Node{
	"testdata",
	[]*Node{
		{"a", nil, 0},
		{"b", []*Node{}, 0},
		{"c", nil, 0},
		{
			"d",
			[]*Node{
				{"x", nil, 0},
				{"y", []*Node{}, 0},
				{
					"z",
					[]*Node{
						{"u", nil, 0},
						{"v", nil, 0},
					},
					0,
				},
			},
			0,
		},
	},
	0,
}

func walkTree(n *Node, path string, f func(path string, n *Node)) {
	f(path, n)
	for _, e := range n.entries {
		walkTree(e, pathpkg.Join(path, e.name), f)
	}
}

func makeTree() fs.FS {
	fsys := fstest.MapFS{}
	walkTree(tree, tree.name, func(path string, n *Node) {
		if n.entries == nil {
			fsys[path] = &fstest.MapFile{}
		} else {
			fsys[path] = &fstest.MapFile{Mode: os.ModeDir}
		}
	})
	return fsys
}

// Assumes that each node name is unique. Good enough for a test.
// If clear is true, any incoming error is cleared before return. The errors
// are always accumulated, though.
func mark(entry compat.DirEntry, err error, errors *[]error, clear bool) error {
	name := entry.Name()
	walkTree(tree, tree.name, func(path string, n *Node) {
		if n.name == name {
			n.mark++
		}
	})
	if err != nil {
		*errors = append(*errors, err)
		if clear {
			return nil
		}
		return err
	}
	return nil
}

func TestWalkDir(t *testing.T) {
	t.Chdir(t.TempDir())

	fsys := makeTree()
	errors := make([]error, 0, 10)
	clear := true
	markFn := func(path string, entry compat.DirEntry, err error) error {
		return mark(entry, err, &errors, clear)
	}
	// Expect no errors.
	err := compat.WalkDir(fsys, ".", markFn)
	if err != nil {
		t.Fatalf("no error expected, found: %s", err)
	}
	if len(errors) != 0 {
		t.Fatalf("unexpected errors: %s", errors)
	}
	walkTree(tree, tree.name, func(path string, n *Node) {
		if n.mark != 1 {
			t.Errorf("node %s mark = %d; expected 1", path, n.mark)
		}
		n.mark = 0
	})
}

func TestWalkDirSymlink(t *testing.T) {
	if !compat.SupportsSymlinks() {
		skip(t, "Skipping test: symlinks are not supported on "+runtime.GOOS)

		return // tinygo doesn't support t.Skip
	}

	fsys := fstest.MapFS{
		"link":    {Data: []byte("dir"), Mode: os.ModeSymlink},
		"dir/a":   {},
		"dir/b/c": {},
		"dir/d":   {Data: []byte("b"), Mode: os.ModeSymlink},
	}

	wantTypes := map[string]os.FileMode{
		"link":     os.ModeDir,
		"link/a":   0,
		"link/b":   os.ModeDir,
		"link/b/c": 0,
		"link/d":   os.ModeSymlink,
	}
	marks := make(map[string]int)
	walkFn := func(path string, entry compat.DirEntry, err error) error {
		marks[path]++
		if want, ok := wantTypes[path]; !ok {
			t.Errorf("Unexpected path %q in walk", path)
		} else if got := entry.Type(); got != want {
			t.Errorf("%s entry type = %o (%v); want %o (%v)", path, got, got, want, want)
		}
		if err != nil {
			t.Errorf("Walking %s: %v", path, err)
		}
		return nil
	}

	// Expect no errors.
	err := compat.WalkDir(fsys, "link", walkFn)
	if err != nil {
		t.Fatalf("no error expected, found: %s", err)
	}
	for path := range wantTypes {
		if got := marks[path]; got != 1 {
			t.Errorf("%s visited %d times; expected 1", path, got)
		}
	}
}

func TestIssue51617(t *testing.T) {
	if compat.IsWindows {
		skip(t, "@TODO reenable this test once the underlying issue is fixed")
		return
	}

	dir := t.TempDir()
	for _, sub := range []string{"a", filepath.Join("a", "bad"), filepath.Join("a", "next")} {
		if err := os.Mkdir(filepath.Join(dir, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	bad := filepath.Join(dir, "a", "bad")
	if err := os.Chmod(bad, 0); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chmod(bad, 0o700) }() // avoid errors on cleanup
	var saw []string
	err := compat.WalkDir(os.DirFS(dir), ".", func(path string, d compat.DirEntry, err error) error {
		if err != nil {
			return filepath.SkipDir
		}
		if d.IsDir() {
			saw = append(saw, path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{".", "a", "a/bad", "a/next"}
	if !slices.Equal(saw, want) {
		t.Errorf("got directories %v, want %v", saw, want)
	}
}
