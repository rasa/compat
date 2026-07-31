// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !linux && !windows

package volume

import (
	"errors"

	"github.com/rasa/compat"
)

func typeOf(_ Mount) (Type, error) {
	return TypeUnknown, &compat.NotYetImplementedError("typeOf")
}
