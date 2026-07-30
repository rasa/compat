//go:build plan9

package compat

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
)

func rename(source, destination string, opts ...Option) error {
	var fopts Options

	for _, opt := range opts {
		opt(&fopts)
	}

	sourceAbs, err := filepath.Abs(source)
	if err != nil {
		return renameError(source, destination, err)
	}

	destinationAbs, err := filepath.Abs(destination)
	if err != nil {
		return renameError(source, destination, err)
	}

	// Plan 9 wstat can rename only within the same directory.
	if filepath.Dir(sourceAbs) != filepath.Dir(destinationAbs) {
		return renameError(source, destination, os.ErrInvalid)
	}

	_, err = os.Stat(destinationAbs)

	switch {
	case err == nil:
		// Replacing an existing destination cannot be atomic on Plan 9.
		if !fopts.allowNonAtomicReplace {
			return renameError(
				source,
				destination,
				errors.ErrUnsupported,
			)
		}

		// Explicitly permitted non-atomic fallback.
		if err := os.Remove(destinationAbs); err != nil {
			return renameError(source, destination, err)
		}

	case errors.Is(err, os.ErrNotExist):
		// Normal atomic rename to a nonexistent destination.
		// Continue to Wstat below.

	default:
		return renameError(source, destination, err)
	}

	var dir syscall.Dir
	dir.Null()
	dir.Name = filepath.Base(destinationAbs)

	buf := make([]byte, syscall.STATFIXLEN+len(dir.Name))

	n, err := dir.Marshal(buf)
	if err != nil {
		return renameError(source, destination, err)
	}

	if err := syscall.Wstat(sourceAbs, buf[:n]); err != nil {
		return renameError(source, destination, err)
	}

	return nil
}

func renameError(source, destination string, err error) error {
	return &os.LinkError{
		Op:  "rename",
		Old: source,
		New: destination,
		Err: err,
	}
}

