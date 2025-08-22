// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package compat_test

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/rasa/compat"
)

func TestMain(m *testing.M) {
	code := m.Run()
	if code == 0 {
		os.Exit(0)
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
