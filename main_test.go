// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !tinygo

package compat_test

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/rasa/compat"
)

const (
	nativeFS        = "Native"
	allFS           = "All"
	defaultTempSize = "2GB"
)

type fsTest struct { //nolint:unused
	fsName string
	vars   testVars
}

func TestMain(m *testing.M) {
	fsToTest := os.Getenv("COMPAT_DEBUG_FS")

	tempSize = defaultTempSize
	fsSize := os.Getenv("COMPAT_DEBUG_FS_SIZE")
	if fsSize != "" {
		tempSize = fsSize
	}

	nativeFSType, _ := compat.PartitionType(context.Background(), os.TempDir())
	if nativeFSType == "" {
		nativeFSType = "Unknown"
	}

	fsPath := os.Getenv("COMPAT_DEBUG_FS_PATH")

	code := testMain(m, fsToTest, nativeFSType, fsPath)

	if code == 0 {
		os.Exit(0)
	}

	if !strings.Contains(compatDebug, "DEBUG") {
		os.Exit(code)
	}

	fmt.Println("Tests failed")
	fmt.Printf("Error code: %v\n", code)
	fmt.Printf("Compiler:   %v\n", runtime.Compiler)
	v := runtime.Version()
	if compat.IsTinygo {
		fmt.Printf("Tinygo ver: %v\n", runtime.Version())
		v = compat.UnderlyingGoVersion()
	}
	fmt.Printf("Go version: %v\n", v)
	fmt.Printf("GOOS:       %v\n", runtime.GOOS)
	fmt.Printf("GOARCH:     %v\n", runtime.GOARCH)
	fmt.Printf("NumCPU:     %v\n", runtime.NumCPU())

	info, ok := debug.ReadBuildInfo()
	if ok {
		fmt.Println("Build info:")
		spew.Dump(info)
		fmt.Printf("%#v\n", info)
	}
	os.Exit(code)
}
