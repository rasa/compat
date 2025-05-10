// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat

import "runtime"

const (
	_aix       = "aix"
	_android   = "android"
	_darwin    = "darwin"
	_dragonfly = "dragonfly"
	_freebsd   = "freebsd"
	_illumos   = "illumos"
	_ios       = "ios"
	_js        = "js"
	_linux     = "linux"
	_netbsd    = "netbsd"
	_openbsd   = "openbsd"
	_plan9     = "plan9"
	_solaris   = "solaris"
	_wasip1    = "wasip1"
	_windows   = "windows"
)

const (
	IsAIX       = runtime.GOOS == _aix
	IsAndroid   = runtime.GOOS == _android
	IsDarwin    = runtime.GOOS == _darwin
	IsDragonfly = runtime.GOOS == _dragonfly
	IsFreeBSD   = runtime.GOOS == _freebsd
	IsIllumos   = runtime.GOOS == _illumos
	IsIOS       = runtime.GOOS == _ios
	IsJS        = runtime.GOOS == _js
	IsLinux     = runtime.GOOS == _linux
	IsNetBSD    = runtime.GOOS == _netbsd
	IsOpenBSD   = runtime.GOOS == _openbsd
	IsPlan9     = runtime.GOOS == _plan9
	IsSolaris   = runtime.GOOS == _solaris
	IsWasip1    = runtime.GOOS == _wasip1
	IsWindows   = runtime.GOOS == _windows
)

const (
	_386      = "386"
	_amd64    = "amd64"
	_arm      = "arm"
	_arm64    = "arm64"
	_loong64  = "loong64"
	_mips     = "mips"
	_mips64   = "mips64"
	_mips64le = "mips64le"
	_mipsle   = "mipsle"
	_ppc64    = "ppc64"
	_ppc64le  = "ppc64le"
	_riscv64  = "riscv64"
	_s390x    = "s390x"
	_wasm     = "wasm"
)

const (
	Is386      = runtime.GOARCH == _386
	IsAmd64    = runtime.GOARCH == _amd64
	IsArm      = runtime.GOARCH == _arm
	IsArm64    = runtime.GOARCH == _arm64
	IsLoong64  = runtime.GOARCH == _loong64
	IsMips     = runtime.GOARCH == _mips
	IsMips64   = runtime.GOARCH == _mips64
	IsMips64le = runtime.GOARCH == _mips64le
	IsMipsle   = runtime.GOARCH == _mipsle
	IsPpc64    = runtime.GOARCH == _ppc64
	IsPpc64le  = runtime.GOARCH == _ppc64le
	IsRiscv64  = runtime.GOARCH == _riscv64
	IsS390x    = runtime.GOARCH == _s390x
	IsWasm     = runtime.GOARCH == _wasm
)
