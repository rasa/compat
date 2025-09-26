// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build ignore

// was go:build freebsd

package compat

import (
	"os"
	"syscall"
	"unsafe"
)

const (
	// https://github.com/freebsd/freebsd-src/blob/7e8fb7756c3ed89a2141b923e6da1b6fd96f509c/sys/sys/fcntl.h#L286
	_F_KINFO = 22
	// https://github.com/freebsd/freebsd-src/blob/7e8fb7756c3ed89a2141b923e6da1b6fd96f509c/sys/sys/syslimits.h#L60
	_PATH_MAX = 1024
)

// https://github.com/freebsd/freebsd-src/blob/7e8fb7756c3ed89a2141b923e6da1b6fd96f509c/sys/sys/caprights.h#L43
type _CapRights struct {
	CrRights [2]uint64
}

// https://github.com/freebsd/freebsd-src/blob/7e8fb7756c3ed89a2141b923e6da1b6fd96f509c/sys/sys/_sockaddr_storage.h#L38-L52
type _SockaddrStorage struct {
	SsLen    uint8
	SsFamily uint8
	Pad1     [6]byte
	Align    int64
	Pad2     [112]byte
}

// https://github.com/freebsd/freebsd-src/blob/7e8fb7756c3ed89a2141b923e6da1b6fd96f509c/sys/sys/user.h#L344
type _KinfoFile struct {
	KfStructsize int32
	KfType       int32
	KfFd         int32
	KfRefCount   int32
	KfFlags      int32
	KfPad0       int32
	KfOffset     int64

	KfSock struct {
		KfSockSendq      uint32
		KfSockDomain0    int32
		KfSockType0      int32
		KfSockProtocol0  int32
		KfSaLocal        _SockaddrStorage
		KfSaPeer         _SockaddrStorage
		KfSockPcb        uint64
		KfSockInpcb      uint64
		KfSockUnpconn    uint64
		KfSockSndSbState uint16
		KfSockRcvSbState uint16
		KfSockRecvq      uint32
	}

	KfStatus    uint16
	KfPad1      uint16
	KfIspare0   int32
	KfCapRights _CapRights
	KfCapSpare  uint64
	KfPath      [_PATH_MAX]byte
}

func fstat(f *os.File) (FileInfo, error) {
	if f == nil {
		return nil, &os.PathError{Op: "stat", Path: "", Err: os.ErrInvalid}
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, &os.PathError{Op: "stat", Path: f.Name(), Err: err}
	}

	fd := int(f.Fd())

	var kif _KinfoFile

	_, _, errno := syscall.Syscall(
		syscall.SYS_FCNTL,
		uintptr(fd),
		uintptr(_F_KINFO),
		uintptr(unsafe.Pointer(&kif)),
	)
	if errno != 0 {
		return nil, &os.PathError{Op: "stat", Path: f.Name(), Err: errno}
	}

	n := 0
	for ; n < len(kif.KfPath); n++ {
		if kif.KfPath[n] == 0 {
			break
		}
	}
	if n == 0 {
		return nil, &os.PathError{Op: "stat", Path: f.Name(), Err: os.ErrInvalid}
	}
	path := string(kif.KfPath[:n])

	return stat(fi, path, false)
}
