// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build (!unix && !windows) || (openbsd && ppc64) || (netbsd && 386) || (freebsd && riscv64) || (cgo && aix && ppc64)

package volume

import (
	"errors"

	"github.com/rasa/compat"
)

var osFeatureMap = map[OSFeature]string{}

func Volumes(mounts []Mount) ([]Volume, error) {
	volumes := []Volume{}

	return volumes, &compat.NotYetImplementedError{Op: "volumes"}
}
