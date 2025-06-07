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

func main() {
	exe, _ := os.Executable()
	fi, err := compat.Stat(exe)
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Printf("Name()   =%v\n", fi.Name())
	fmt.Printf("Size()   =%v\n", fi.Size())
	fmt.Printf("Mode()   =0o%o\n", fi.Mode())
	fmt.Printf("ModTime()=%v\n", fi.ModTime())
	fmt.Printf("IsDir()  =%v\n", fi.IsDir())
	fmt.Printf("Sys()    =%+v\n", fi.Sys())
	fmt.Printf("PartID() =%v\n", fi.PartitionID())
	fmt.Printf("FileID() =%v\n", fi.FileID())
	fmt.Printf("Links()  =%v\n", fi.Links())
	fmt.Printf("ATime()  =%v\n", fi.ATime())
	fmt.Printf("BTime()  =%v\n", fi.BTime())
	fmt.Printf("CTime()  =%v\n", fi.CTime())
	fmt.Printf("MTime()  =%v\n", fi.MTime())
	fmt.Printf("UID()    =%v\n", fi.UID())
	fmt.Printf("GID()    =%v\n", fi.GID())

	samePartition()
	samePartitionl()
	samePartitions()
	sameFile()
	sameFilel()
	sameFiles()
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
		log.Fatal(err) //nolint:gocritic // quiet linter
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
		log.Fatal(err) //nolint:gocritic // quiet linter
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
