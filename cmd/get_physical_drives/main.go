// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package main

import (
	"errors"
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	IOCTL_VOLUME_GET_VOLUME_DISK_EXTENTS = 0x00560000
	IOCTL_STORAGE_QUERY_PROPERTY         = 0x002D1400

	StorageDeviceProperty = 0
	PropertyStandardQuery = 0
)

type DISK_EXTENT struct {
	DiskNumber     uint32
	StartingOffset int64
	ExtentLength   int64
}

type VOLUME_DISK_EXTENTS struct {
	NumberOfExtents uint32
	Extents         [1]DISK_EXTENT
}

type STORAGE_PROPERTY_QUERY struct {
	PropertyId uint32
	QueryType  uint32
	Additional [1]byte
}

type STORAGE_DEVICE_DESCRIPTOR struct {
	Version               uint32
	Size                  uint32
	DeviceType            byte
	DeviceTypeModifier    byte
	RemovableMedia        byte
	CommandQueueing       byte
	VendorIdOffset        uint32
	ProductIdOffset       uint32
	ProductRevisionOffset uint32
	SerialNumberOffset    uint32
	BusType               uint32
	RawPropertiesLength   uint32
	// followed by variable-length RawDeviceProperties
}

func getPhysicalDrivesFromVolume(volume string) ([]uint32, error) {
	p16, err := windows.UTF16PtrFromString(volume)
	if err != nil {
		return nil, err
	}

	h, err := windows.CreateFile(
		p16,
		0,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil,
		windows.OPEN_EXISTING,
		0,
		0,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateFile volume: %w", err)
	}
	defer windows.CloseHandle(h)

	buf := make([]byte, 1024)

	var bytesReturned uint32

	err = windows.DeviceIoControl(
		h,
		IOCTL_VOLUME_GET_VOLUME_DISK_EXTENTS,
		nil,
		0,
		&buf[0],
		uint32(len(buf)),
		&bytesReturned,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("DeviceIoControl volume: %w", err)
	}

	extents := (*VOLUME_DISK_EXTENTS)(unsafe.Pointer(&buf[0]))

	drives := make([]uint32, extents.NumberOfExtents)
	for i := range uint32(extents.NumberOfExtents) {
		drives[i] = extents.Extents[i].DiskNumber
	}

	return drives, nil
}

func getDiskSerial(diskNum uint32) (string, error) {
	path := fmt.Sprintf(`\\.\PhysicalDrive%d`, diskNum)

	p16, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return "", err
	}

	h, err := windows.CreateFile(
		p16,
		0,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		0,
		0,
	)
	if err != nil {
		return "", fmt.Errorf("CreateFile disk: %w", err)
	}
	defer windows.CloseHandle(h)

	query := STORAGE_PROPERTY_QUERY{
		PropertyId: StorageDeviceProperty,
		QueryType:  PropertyStandardQuery,
	}

	buf := make([]byte, 1024)

	var returned uint32

	err = windows.DeviceIoControl(
		h,
		IOCTL_STORAGE_QUERY_PROPERTY,
		(*byte)(unsafe.Pointer(&query)),
		uint32(unsafe.Sizeof(query)),
		&buf[0],
		uint32(len(buf)),
		&returned,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("DeviceIoControl disk: %w", err)
	}

	desc := (*STORAGE_DEVICE_DESCRIPTOR)(unsafe.Pointer(&buf[0]))
	if desc.SerialNumberOffset != 0 && desc.SerialNumberOffset < uint32(len(buf)) {
		serial := buf[desc.SerialNumberOffset:]
		for i, b := range serial {
			if b == 0 {
				serial = serial[:i]

				break
			}
		}

		return string(serial), nil
	}

	return "", errors.New("serial number not available")
}

func main() {
	// Example: Start from C: volume
	vol := `\\.\C:`
	if len(os.Args) > 1 {
		vol = os.Args[1]
	}

	drives, err := getPhysicalDrivesFromVolume(vol)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Volume %s maps to PhysicalDrives: %v\n", vol, drives)

	for _, d := range drives {
		serial, err := getDiskSerial(d)
		if err != nil {
			fmt.Printf("PhysicalDrive%d: error %v\n", d, err)
		} else {
			fmt.Printf("PhysicalDrive%d serial: %s\n", d, serial)
		}
	}
}
