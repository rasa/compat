// SPDX-FileCopyrightText: Copyright © 2026 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package golang

import (
	"strconv"
	"testing"
)

func TestLastIndexByteString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		s    string
		c    byte
		want int
	}{
		{name: "found", s: "abca", c: 'a', want: 3},
		{name: "first", s: "abc", c: 'a', want: 0},
		{name: "missing", s: "abc", c: 'z', want: -1},
		{name: "empty", s: "", c: 'a', want: -1},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := lastIndexByteString(tt.s, tt.c); got != tt.want {
				t.Fatalf("lastIndexByteString(%q, %q): got %v, want %v", tt.s, tt.c, got, tt.want)
			}
		})
	}
}

func TestUitoa(t *testing.T) {
	t.Parallel()

	tests := []uint{0, 1, 9, 10, 12345, ^uint(0)}
	for _, val := range tests {
		val := val
		t.Run(strconv.FormatUint(uint64(val), 10), func(t *testing.T) {
			t.Parallel()
			want := strconv.FormatUint(uint64(val), 10)
			if got := Uitoa(val); got != want {
				t.Fatalf("Uitoa(%v): got %q, want %q", val, got, want)
			}
		})
	}
}

func TestPrefixAndSuffix(t *testing.T) {
	t.Parallel()

	sepPattern := "bad" + string(PathSeparator) + "pattern"
	tests := []struct {
		name       string
		pattern    string
		wantPrefix string
		wantSuffix string
		wantErr    error
	}{
		{name: "no wildcard", pattern: "prefix", wantPrefix: "prefix"},
		{name: "wildcard", pattern: "pre*suf", wantPrefix: "pre", wantSuffix: "suf"},
		{name: "last wildcard wins", pattern: "a*b*c", wantPrefix: "a*b", wantSuffix: "c"},
		{name: "separator error", pattern: sepPattern, wantErr: errPatternHasSeparator},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotPrefix, gotSuffix, err := prefixAndSuffix(tt.pattern)
			if err != tt.wantErr {
				t.Fatalf("prefixAndSuffix(%q): got err %v, want %v", tt.pattern, err, tt.wantErr)
			}
			if gotPrefix != tt.wantPrefix || gotSuffix != tt.wantSuffix {
				t.Fatalf("prefixAndSuffix(%q): got (%q, %q), want (%q, %q)", tt.pattern, gotPrefix, gotSuffix, tt.wantPrefix, tt.wantSuffix)
			}
		})
	}
}

func TestJoinPath(t *testing.T) {
	t.Parallel()

	sep := string(PathSeparator)
	if got := joinPath("base", "name"); got != "base"+sep+"name" {
		t.Fatalf("joinPath without trailing separator: got %q", got)
	}
	if got := joinPath("base"+sep, "name"); got != "base"+sep+"name" {
		t.Fatalf("joinPath with trailing separator: got %q", got)
	}
}

func TestEndsWithDot(t *testing.T) {
	t.Parallel()

	sep := string(PathSeparator)
	tests := []struct {
		path string
		want bool
	}{
		{path: ".", want: true},
		{path: "abc" + sep + ".", want: true},
		{path: "abc.", want: false},
		{path: "abc" + sep + "x", want: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()
			if got := endsWithDot(tt.path); got != tt.want {
				t.Fatalf("endsWithDot(%q): got %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
