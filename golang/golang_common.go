//nolint:all
// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package golang

import (
	"errors"
	"sync"
	_ "unsafe" // for go:linkname
)

// SPDX-FileCopyrightText: Copyright 2012 The Go Authors. All rights reserved.
// SPDX-License-Identifier: BSD-3

// The following code is:
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Snippet: https://github.com/golang/go/blob/77f911e3/src/os/getwd.go#L13-L16

var getwdCache struct {
	sync.Mutex
	dir string
}

// Snippet: https://github.com/golang/go/blob/77f911e3/src/os/tempfile.go#L22-L24

func nextRandom() string {
	return Uitoa(uint(uint32(runtime_rand())))
}

// Snippet: https://github.com/golang/go/blob/77f911e3/src/internal/bytealg/lastindexbyte_generic.go#L16-L23

func lastIndexByteString(s string, c byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			return i
		}
	}
	return -1
}

// Snippet: https://github.com/golang/go/blob/77f911e3/src/internal/itoa/itoa.go#L18-L33

func Uitoa(val uint) string {
	if val == 0 { // avoid string allocation
		return "0"
	}
	var buf [20]byte // big enough for 64bit value base 10
	i := len(buf) - 1
	for val >= 10 {
		q := val / 10
		buf[i] = byte('0' + val - q*10)
		i--
		val = q
	}
	// val < 10
	buf[i] = byte('0' + val)
	return string(buf[i:])
}

// Snippet: https://github.com/golang/go/blob/77f911e3/src/os/tempfile.go#L60-L60

var errPatternHasSeparator = errors.New("pattern contains path separator")

// Snippet: https://github.com/golang/go/blob/77f911e3/src/os/tempfile.go#L64-L76

func prefixAndSuffix(pattern string) (prefix, suffix string, err error) {
	for i := 0; i < len(pattern); i++ {
		if IsPathSeparator(pattern[i]) {
			return "", "", errPatternHasSeparator
		}
	}
	if pos := lastIndexByteString(pattern, '*'); pos != -1 { // removed bytealg
		prefix, suffix = pattern[:pos], pattern[pos+1:]
	} else {
		prefix = pattern
	}
	return prefix, suffix, nil
}

// Snippet: https://github.com/golang/go/blob/77f911e3/src/os/tempfile.go#L119-L124

func joinPath(dir, name string) string {
	if len(dir) > 0 && IsPathSeparator(dir[len(dir)-1]) {
		return dir + name
	}
	return dir + string(PathSeparator) + name
}

// Snippet: https://github.com/golang/go/blob/77f911e3/src/os/path.go#L77-L86

func endsWithDot(path string) bool {
	if path == "." {
		return true
	}
	if len(path) >= 2 && path[len(path)-1] == '.' && IsPathSeparator(path[len(path)-2]) {
		return true
	}
	return false
}
