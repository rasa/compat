// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import "runtime"

const (
	aIX       = "aix"
	android   = "android"
	darwin    = "darwin"
	dragonfly = "dragonfly"
	freeBSD   = "freebsd"
	illumos   = "illumos"
	iOS       = "ios"
	jS        = "js"
	linux     = "linux"
	netBSD    = "netbsd"
	openBSD   = "openbsd"
	plan9     = "plan9"
	solaris   = "solaris"
	wasip1    = "wasip1"
	windows   = "windows"
)

const (
	IsAIX       = runtime.GOOS == aIX
	IsAndroid   = runtime.GOOS == android
	IsDarwin    = runtime.GOOS == darwin
	IsDragonfly = runtime.GOOS == dragonfly
	IsFreeBSD   = runtime.GOOS == freeBSD
	IsIllumos   = runtime.GOOS == illumos
	IsIOS       = runtime.GOOS == iOS
	IsJS        = runtime.GOOS == jS
	IsLinux     = runtime.GOOS == linux
	IsNetBSD    = runtime.GOOS == netBSD
	IsOpenBSD   = runtime.GOOS == openBSD
	IsPlan9     = runtime.GOOS == plan9
	IsSolaris   = runtime.GOOS == solaris
	IsWasip1    = runtime.GOOS == wasip1
	IsWindows   = runtime.GOOS == windows
)

const (
	i386     = "386"
	amd64    = "amd64"
	arm      = "arm"
	arm64    = "arm64"
	loong64  = "loong64"
	mips     = "mips"
	mips64   = "mips64"
	mips64le = "mips64le"
	mipsle   = "mipsle"
	ppc64    = "ppc64"
	ppc64le  = "ppc64le"
	riscv64  = "riscv64"
	s390x    = "s390x"
	wasm     = "wasm"
)

const (
	Is386      = runtime.GOARCH == i386
	IsAmd64    = runtime.GOARCH == amd64
	IsArm      = runtime.GOARCH == arm
	IsArm64    = runtime.GOARCH == arm64
	IsLoong64  = runtime.GOARCH == loong64
	IsMips     = runtime.GOARCH == mips
	IsMips64   = runtime.GOARCH == mips64
	IsMips64le = runtime.GOARCH == mips64le
	IsMipsle   = runtime.GOARCH == mipsle
	IsPpc64    = runtime.GOARCH == ppc64
	IsPpc64le  = runtime.GOARCH == ppc64le
	IsRiscv64  = runtime.GOARCH == riscv64
	IsS390x    = runtime.GOARCH == s390x
	IsWasm     = runtime.GOARCH == wasm
)
