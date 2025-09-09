// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// demo is a sample application that runs functions provided by this library.
// This also helps with our code coverage report.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rasa/compat"
)

const mode = os.FileMode(0o654) // Something other than the default.

func main() {
	name := "hello.txt"
	err := compat.WriteFile(name, []byte("Hello World"), mode)
	if err != nil {
		log.Fatal(err)
	}

	fi, err := compat.Stat(name)
	if err != nil {
		log.Fatal(err)
	}

	print(fi.String())

	samePartition()
	samePartitionl()
	samePartitions()
	sameFile()
	sameFilel()
	sameFiles()
	_ = os.Remove(name)
}

func samePartition() {
	exe, _ := os.Executable()

	fi1, err := compat.Stat(exe)
	if err != nil {
		log.Fatal(err)
	}

	fi2, err := compat.Stat(exe)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("SamePartition(): %v\n", compat.SamePartition(fi1, fi2))
}

func samePartitionl() {
	exe, _ := os.Executable()
	link := "link" + filepath.Ext(exe)

	_ = os.Link(exe, link)

	defer os.Remove(link)

	fi1, err := compat.Lstat(link)
	if err != nil {
		log.Fatal(err) //nolint:gocritic
	}

	fi2, err := compat.Lstat(link)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("SamePartition() (link): %v\n", compat.SamePartition(fi1, fi2))
}

func samePartitions() {
	exe, _ := os.Executable()

	fmt.Printf("SamePartitions(): %v\n", compat.SamePartitions(exe, exe))
}

func sameFile() {
	exe, _ := os.Executable()

	fi1, err := compat.Stat(exe)
	if err != nil {
		log.Fatal(err)
	}

	fi2, err := compat.Stat(exe)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("SameFile(): %v\n", compat.SameFile(fi1, fi2))
}

func sameFilel() {
	exe, _ := os.Executable()
	link := "link" + filepath.Ext(exe)

	_ = os.Link(exe, link)

	defer os.Remove(link)

	fi1, err := compat.Lstat(link)
	if err != nil {
		log.Fatal(err) //nolint:gocritic
	}

	fi2, err := compat.Lstat(link)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("SameFile() (link): %v\n", compat.SameFile(fi1, fi2))
}

func sameFiles() {
	exe, _ := os.Executable()

	fmt.Printf("SameFiles(): %v\n", compat.SameFiles(exe, exe))
}
