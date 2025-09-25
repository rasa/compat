// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build darwin || dragonfly || freebsd || netbsd || openbsd

package compat

import (
	"bytes"
	"os"
	"syscall"
	"unsafe"
)

const (
	// See https://github.com/apple-oss-distributions/xnu/blob/e3723e1f/bsd/sys/syslimits.h#L111
	MAXPATHLEN = 1024
	// See https://github.com/apple-oss-distributions/xnu/blob/e3723e1f/bsd/sys/fcntl.h#L303
	F_GETPATH = 50
)

func fstat(f *os.File) (FileInfo, error) {
	if f == nil {
		return nil, &os.PathError{Op: "stat", Path: "", Err: os.ErrInvalid}
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, &os.PathError{Op: "stat", Path: f.Name(), Err: err}
	}

	fd := int(f.Fd())

	var buf [MAXPATHLEN]byte
	_, _, errno := syscall.Syscall(syscall.SYS_FCNTL, uintptr(fd), F_GETPATH, uintptr(unsafe.Pointer(&buf[0])))
	if errno != 0 {
		return nil, &os.PathError{Op: "stat", Path: f.Name(), Err: errno}
	}
	i := bytes.IndexByte(buf[:], 0)
	if i < 0 {
		i = len(buf)
	}
	path := string(buf[:i])

	return stat(fi, path, false)
}
