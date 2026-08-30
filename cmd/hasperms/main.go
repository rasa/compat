// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rasa/compat"
)

const mode765 = 0o765

var ErrFileUpdatedExternally = errors.New("file updated externally")

func supportsChmod(path string) (bool, error) {
	// Create a temp file in the target path
	tmp := path + "/.permtest.tmp"

	f, err := os.Create(tmp) //nolint:gosec
	if err != nil {
		return false, fmt.Errorf("cannot create %q: %w", tmp, err)
	}

	_ = f.Close()

	defer os.Remove(tmp)

	// Try to chmod
	perm := os.FileMode(mode765)

	fi, err := compat.Stat(tmp)
	if err == nil {
		log.Printf("stat:\n%v", fi)
	}

	if time.Since(fi.CTime()) <= 10*time.Millisecond {
		time.Sleep(10 * time.Millisecond) //nolint:mnd
	}

	err = os.Chmod(tmp, perm) //nolint:gosec
	if err != nil {
		// Check for ENOTSUP / EOPNOTSUPP
		errno := &os.PathError{}

		ok := errors.As(err, &errno)
		if ok {
			if errors.Is(err, errors.ErrUnsupported) {
				return false, nil
			}
		}

		return false, fmt.Errorf("cannot chmod %q: %w", tmp, err)
	}

	// Some filesystems just silently ignore chmod but don’t fail
	// If chmod succeeds but stat doesn’t reflect change, also false
	fi2, err := compat.Stat(tmp)
	if err != nil {
		return false, fmt.Errorf("cannot stat %q: %w", tmp, err)
	}

	log.Printf("stat:\n%v", fi2)

	if fi.ModTime() != fi2.ModTime() {
		return false, ErrFileUpdatedExternally
	}

	if fi2.Mode().Perm() != perm {
		return false, nil
	}

	return true, nil
}

func main() {
	arg := "/tmp"
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	b, err := supportsChmod(arg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Supports chmod: %v (%v)\n", b, arg)
}
