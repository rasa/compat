// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
	"go/version"
	"os"
	"runtime"
)

// IsTinygo is true if the go compiler is tinygo.
const IsTinygo = runtime.Compiler == "tinygo"

// IsAct is true when running github actions locally using the act command.
var IsAct = os.Getenv("ACT") == "true"

const (
	IsAIX       = runtime.GOOS == "aix"
	IsAndroid   = runtime.GOOS == "android"
	IsDarwin    = runtime.GOOS == "darwin"
	IsDragonfly = runtime.GOOS == "dragonfly"
	IsFreeBSD   = runtime.GOOS == "freebsd"
	IsIllumos   = runtime.GOOS == "illumos"
	IsIOS       = runtime.GOOS == "ios"
	IsJS        = runtime.GOOS == "js"
	IsLinux     = runtime.GOOS == "linux"
	IsNetBSD    = runtime.GOOS == "netbsd"
	IsOpenBSD   = runtime.GOOS == "openbsd"
	IsPlan9     = runtime.GOOS == "plan9"
	IsSolaris   = runtime.GOOS == "solaris"
	IsWasip1    = runtime.GOOS == "wasip1"
	IsWindows   = runtime.GOOS == "windows"
)

const (
	IsApple   = IsDarwin || IsIOS
	IsBSD     = IsDragonfly || IsFreeBSD || IsNetBSD || IsOpenBSD
	IsBSDLike = IsApple || IsBSD
	IsSolaria = IsIllumos || IsSolaris
)

const (
	Is386      = runtime.GOARCH == "386"
	IsAmd64    = runtime.GOARCH == "amd64"
	IsArm      = runtime.GOARCH == "arm"
	IsArm64    = runtime.GOARCH == "arm64"
	IsLoong64  = runtime.GOARCH == "loong64"
	IsMips     = runtime.GOARCH == "mips"
	IsMips64   = runtime.GOARCH == "mips64"
	IsMips64le = runtime.GOARCH == "mips64le"
	IsMipsle   = runtime.GOARCH == "mipsle"
	IsPpc64    = runtime.GOARCH == "ppc64"
	IsPpc64le  = runtime.GOARCH == "ppc64le"
	IsRiscv64  = runtime.GOARCH == "riscv64"
	IsS390x    = runtime.GOARCH == "s390x"
	IsWasm     = runtime.GOARCH == "wasm"
)

const (
	IsX86CPU  = Is386 || IsAmd64
	IsArmCPU  = IsArm || IsArm64
	IsMipsCPU = IsMips || IsMips64 || IsMips64le || IsMipsle
	IsPpcCPU  = IsPpc64 || IsPpc64le
)

var tinygoThresholds = []struct {
	tinygo   string
	goMinVer string
	goMaxVer string
}{
	// https://github.com/tinygo-org/tinygo/blob/v0.12.0/builder/config.go#L28
	{"0.12.0", "go1.11", "go1.13"},
	// https://github.com/tinygo-org/tinygo/blob/v0.14.0/builder/config.go#L28
	{"0.14.0", "go1.11", "go1.14"},
	// https://github.com/tinygo-org/tinygo/blob/v0.16.0/builder/config.go#L28
	{"0.16.0", "go1.11", "go1.15"},
	// https://github.com/tinygo-org/tinygo/blob/v0.17.0/builder/config.go#L28
	{"0.17.0", "go1.11", "go1.16"},
	// https://github.com/tinygo-org/tinygo/blob/v0.19.0/builder/config.go#L36
	{"0.19.0", "go1.13", "go1.16"},
	// https://github.com/tinygo-org/tinygo/blob/v0.22.0/builder/config.go#L36
	{"0.22.0", "go1.15", "go1.17"},
	// https://github.com/tinygo-org/tinygo/blob/v0.23.0/builder/config.go#L36
	{"0.23.0", "go1.15", "go1.18"},
	// https://github.com/tinygo-org/tinygo/blob/v0.25.0/builder/config.go#L36
	{"0.25.0", "go1.16", "go1.19"},
	// https://github.com/tinygo-org/tinygo/blob/v0.26.0/builder/config.go#L36
	{"0.26.0", "go1.18", "go1.19"},
	// https://github.com/tinygo-org/tinygo/blob/v0.27.0/builder/config.go#L36
	{"0.27.0", "go1.18", "go1.20"},
	// https://github.com/tinygo-org/tinygo/blob/v0.29.0/builder/config.go#L30
	{"0.29.0", "go1.18", "go1.21"},
	// https://github.com/tinygo-org/tinygo/blob/v0.31.0/builder/config.go#L30
	{"0.31.0", "go1.18", "go1.22"},
	// https://github.com/tinygo-org/tinygo/blob/v0.33.0/builder/config.go#L30
	{"0.33.0", "go1.19", "go1.23"},
	// https://github.com/tinygo-org/tinygo/blob/v0.36.0/builder/config.go#L28
	{"0.36.0", "go1.19", "go1.24"},
	// https://github.com/tinygo-org/tinygo/blob/v0.39.0/builder/config.go#L28
	{"0.39.0", "go1.19", "go1.25"},
}

// UnderlyingGoVersion returns the effective Go toolchain version string ("goX.Y")
// for the current environment.
// - On standard Go: returns runtime.Version() (already "go1.xx").
// - On TinyGo: picks the highest Go version supported based on thresholds.
func UnderlyingGoVersion() string {
	v := runtime.Version()

	if !IsTinygo {
		return v
	}

	v = "go" + v

	// TinyGo: runtime.Version() is like "0.39.1"
	best := ""
	for _, th := range tinygoThresholds {
		if version.Compare(v, "go"+th.tinygo) >= 0 {
			best = th.goMaxVer
		}
	}

	return best
}
