package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/rasa/compat"
)

func supportsChmod(path string) (bool, error) {
	// Create a temp file in the target path
	tmp := path + "/.permtest.tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return false, err
	}
	_ = f.Close()
	defer os.Remove(tmp) //nolint:errcheck

	// Try to chmod
	perm := os.FileMode(0o765)
	fi, err := compat.Stat(tmp)
	if err == nil {
		log.Printf("stat:\n%v", fi)
	}

	if time.Since(fi.CTime()) <= 10*time.Millisecond {
		time.Sleep(10 * time.Millisecond)
	}

	err = os.Chmod(tmp, perm)
	if err != nil {
		// Check for ENOTSUP / EOPNOTSUPP
		errno := &os.PathError{}
		ok := errors.As(err, &errno)
		if ok {
			if IsUnsupportedError(err) {
				return false, nil
			}
		}

		return false, err
	}

	// Some filesystems just silently ignore chmod but don’t fail
	// If chmod succeeds but stat doesn’t reflect change, also false
	fi2, err := compat.Stat(tmp)
	if err != nil {
		return false, err
	}
	log.Printf("stat:\n%v", fi2)
	if fi.ModTime() != fi2.ModTime() {
		return false, errors.New("file updated externally")
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
