// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// main is a sample application that runs functions provided by this library.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rasa/compat"
)

func main() {
	fmt.Printf("sameDevice() returned  %v\n", sameDevice())
	fmt.Printf("sameDevices() returned %v\n", sameDevices())
	fmt.Printf("sameFile() returned    %v\n", sameFile())
	fmt.Printf("sameFiles() returned   %v\n", sameFiles())
}

func sameDevice() string {
	exe, _ := os.Executable()
	fi1, _ := compat.Stat(exe)
	fi2, _ := compat.Stat(exe)

	return strconv.FormatBool(compat.SameDevice(fi1, fi2))
}

func sameDevices() string {
	exe, _ := os.Executable()

	return strconv.FormatBool(compat.SameDevices(exe, exe))
}

func sameFile() string {
	exe, _ := os.Executable()
	fi1, _ := compat.Stat(exe)
	fi2, _ := compat.Stat(exe)

	return strconv.FormatBool(compat.SameFile(fi1, fi2))
}

func sameFiles() string {
	exe, _ := os.Executable()

	return strconv.FormatBool(compat.SameFiles(exe, exe))
}
