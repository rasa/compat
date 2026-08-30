// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9

package compat

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Nice gets the CPU process priority. The return value is in a range from
// -20 (least nice), to 19 (most nice), even on non-Unix systems such as
// Windows, Plan 9, etc. If not supported by the operating system, 0 is
// returned.
//
// Plan 9 exposes the base and current scheduling priorities as the final two
// fields of /proc/pid/status. This function reports the base priority because
// that is the value changed by Renice.
func Nice() (int, error) {
	path := fmt.Sprintf("/proc/%d/status", os.Getpid())

	data, err := os.ReadFile(path)
	if err != nil {
		return 0, niceError("cannot read "+path, err)
	}

	fields := strings.Fields(string(data))
	if len(fields) < 2 {
		return 0, fmt.Errorf("nice: malformed process status in %q", path)
	}

	// The last two fields are the base and current scheduling priorities.
	pri := fields[len(fields)-2]
	basePriority, err := strconv.Atoi(pri)
	if err != nil {
		return 0, fmt.Errorf("nice: invalid base priority %q in %q: %w", pri, path, err)
	}

	nice, ok := niceMap[basePriority]
	if !ok {
		return 0, fmt.Errorf("nice: invalid priority %d in %q", basePriority, path)
	}

	return nice, nil
}

// See https://9p.io/magic/man2html/3/proc
//
// Plan 9 priorities range from 0 to 19, with larger values representing
// higher priorities. Unix nice values range from -20 to 19, with smaller
// values representing higher priorities.
//
// Because Unix has 40 nice values and Plan 9 has only 20 priority values,
// conversion is necessarily lossy.

// priorityMap maps Unix nice levels (-20 to 19) to Plan 9 priorities (0 to 19).
var priorityMap = map[int]uint32{
	-20: 19,
	-19: 18,
	-18: 18,
	-17: 17,
	-16: 17,
	-15: 16,
	-14: 16,
	-13: 15,
	-12: 15,
	-11: 14,
	-10: 14,
	-9:  13,
	-8:  13,
	-7:  13,
	-6:  12,
	-5:  12,
	-4:  12,
	-3:  11,
	-2:  11,
	-1:  11,
	0:   10,
	1:   9,
	2:   9,
	3:   8,
	4:   8,
	5:   7,
	6:   7,
	7:   6,
	8:   6,
	9:   5,
	10:  5,
	11:  4,
	12:  4,
	13:  3,
	14:  3,
	15:  2,
	16:  2,
	17:  1,
	18:  1,
	19:  0,
}

// niceMap maps each Plan 9 priority to a representative Unix nice value.
// Where multiple Unix nice values map to the same Plan 9 priority, one
// canonical value is returned.
var niceMap = map[int]int{
	19: -20,
	18: -18,
	17: -16,
	16: -14,
	15: -12,
	14: -10,
	13: -8,
	12: -5,
	11: -2,
	10: 0,
	9:  1,
	8:  3,
	7:  5,
	6:  7,
	5:  9,
	4:  11,
	3:  13,
	2:  15,
	1:  17,
	0:  19,
}

// Renice sets the CPU process priority. The nice parameter can range from
// -20 (least nice) to 19 (most nice), even on non-Unix systems such as
// Windows, Plan 9, etc. If not supported by the operating system, nil is
// returned.
func Renice(nice int) error {
	err := validateNice(nice)
	if err != nil {
		return err
	}

	priority, ok := priorityMap[nice]
	if !ok {
		return unexpectedNiceError(nice)
	}

	path := fmt.Sprintf("/proc/%d/ctl", os.Getpid())

	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return niceError("cannot open process control file", err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "pri %d", priority)
	if err != nil {
		return niceError("cannot write to process control file", err)
	}

	return nil
}
