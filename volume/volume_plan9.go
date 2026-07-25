// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9

package volume

import (
	"errors"
)

var osFeatureMap = map[OSFeature]string{}

func Volumes(mounts []Mount) ([]Volume, error) {
	volumes := []Volume{}

	return volumes, errors.ErrUnsupported
}
