// SPDX-FileCopyrightText: Copyright (c) 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build darwin

package compat

import (
	"bytes"
	"os"
	"syscall"
	"unsafe"
)

const (
	// See https://github.com/apple-oss-distributions/xnu/blob/e3723e1f/bsd/sys/syslimits.h#L111
	maxPathLen = 1024
	// See https://github.com/apple-oss-distributions/xnu/blob/e3723e1f/bsd/sys/fcntl.h#L303
	fGetPath = 50
)

func fstat(fp *os.File) (FileInfo, error) {
	if fp == nil {
		return nil, statError("", os.ErrInvalid)
	}

	fi, err := fp.Stat()
	if err != nil {
		return nil, statError(fp.Name(), err)
	}

	fd := int(fp.Fd())

	var buf [maxPathLen]byte

	_, _, errno := syscall.Syscall(syscall.SYS_FCNTL, uintptr(fd), fGetPath, uintptr(unsafe.Pointer(&buf[0])))
	if errno != 0 {
		return nil, statError(fp.Name(), errno)
	}

	i := bytes.IndexByte(buf[:], 0)
	if i < 0 {
		i = len(buf)
	}

	path := string(buf[:i])

	return stat(fi, path, false)
}
