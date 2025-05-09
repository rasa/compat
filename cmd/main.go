// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// main is a sample application that runs functions provided by this library.
// This also helps with our code coverage report.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rasa/compat"
)

func main() {
	fmt.Printf("sameDevice() returned  %v\n", sameDevice())
	fmt.Printf("sameDevicel() returned %v\n", sameDevicel())
	fmt.Printf("sameDevices() returned %v\n", sameDevices())
	fmt.Printf("sameFile() returned    %v\n", sameFile())
	fmt.Printf("sameFilel() returned   %v\n", sameFilel())
	fmt.Printf("sameFiles() returned   %v\n", sameFiles())
	exe, _ := os.Executable()
	fi, err := compat.Stat(exe)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Name()    =%v\n", fi.Name())
	fmt.Printf("Size()    =%v\n", fi.Size())
	fmt.Printf("Mode()    =0o%o\n", fi.Mode())
	fmt.Printf("ModTime() =%v\n", fi.ModTime())
	fmt.Printf("IsDir()   =%v\n", fi.IsDir())
	fmt.Printf("Sys()     =%+v\n", fi.Sys())
	fmt.Printf("DeviceID()=%v\n", fi.DeviceID())
	fmt.Printf("FileID()  =%v\n", fi.FileID())
	fmt.Printf("Links()   =%v\n", fi.Links())
	fmt.Printf("ATime()   =%v\n", fi.ATime())
	fmt.Printf("BTime()   =%v\n", fi.BTime())
	fmt.Printf("CTime()   =%v\n", fi.CTime())
	fmt.Printf("MTime()   =%v\n", fi.MTime())
	fmt.Printf("UID()     =%v\n", fi.UID())
	fmt.Printf("GID()     =%v\n", fi.GID())
}

func sameDevice() string {
	exe, _ := os.Executable()
	fi1, err := compat.Stat(exe)
	if err != nil {
		log.Fatal(err)
	}
	fi2, err := compat.Stat(exe)
	if err != nil {
		log.Fatal(err)
	}

	return strconv.FormatBool(compat.SameDevice(fi1, fi2))
}

func sameDevicel() string {
	exe, _ := os.Executable()
	link := "link" + filepath.Ext(exe)
	_ = os.Link(exe, link)
	defer os.Remove(link)
	fi1, err := compat.Lstat(link)
	if err != nil {
		log.Fatal(err)
	}
	fi2, err := compat.Lstat(link)
	if err != nil {
		log.Fatal(err)
	}

	return strconv.FormatBool(compat.SameDevice(fi1, fi2))
}

func sameDevices() string {
	exe, _ := os.Executable()

	return strconv.FormatBool(compat.SameDevices(exe, exe))
}

func sameFile() string {
	exe, _ := os.Executable()
	fi1, err := compat.Stat(exe)
	if err != nil {
		log.Fatal(err)
	}
	fi2, err := compat.Stat(exe)
	if err != nil {
		log.Fatal(err)
	}

	return strconv.FormatBool(compat.SameFile(fi1, fi2))
}

func sameFilel() string {
	exe, _ := os.Executable()
	link := "link" + filepath.Ext(exe)
	_ = os.Link(exe, link)
	defer os.Remove(link)
	fi1, err := compat.Lstat(link)
	if err != nil {
		log.Fatal(err)
	}
	fi2, err := compat.Lstat(link)
	if err != nil {
		log.Fatal(err)
	}

	return strconv.FormatBool(compat.SameFile(fi1, fi2))
}

func sameFiles() string {
	exe, _ := os.Executable()

	return strconv.FormatBool(compat.SameFiles(exe, exe))
}
