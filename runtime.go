// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import (
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
