// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build plan9

package compat

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Nice gets the CPU process priority. The return value is in a range from
// -20 (least nice), to 19 (most nice), even on non-Unix systems such as
// Windows, plan9, etc. If not supported by the operating system, 0 is
// returned.
func Nice() (int, error) {
	pid := os.Getpid()
	path := fmt.Sprintf("/proc/%d/status", pid)

	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close() //nolint:errcheck

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "pri ") {
			// format: "pri <value>"
			fields := strings.Fields(line)
			if len(fields) == 2 {
				val, err := strconv.Atoi(fields[1])
				if err == nil {
					return val, nil
				}
			}
		}
	}

	return 0, fmt.Errorf("nice: priority not found in %v", path)
}

// See https://9p.io/magic/man2html/3/proc

// priorityMap maps unix's nice levels (-20 to 19), to plan9's levels (0-19):
// input value   plan9 priority
// -----------   ---------------------
//
//		    -20 : 19 (highest priority)
//	   -19, -18 : 18
//	   -17, -16 : 17
//	   -15, -14 : 16
//	   -13, -12 : 15
//	   -11, -10 : 14
//	 -9, -8, -7 : 13 (plan9's default for kernel processes)
//	 -6, -5, -4 : 12
//	 -3, -2, -1 : 11
//	          0 : 10 (normal priority) (plan9's default for non-kernel processes)
//	      1,  2 :  9
//	      3,  4 :  8
//	      5,  6 :  7
//	      7,  8 :  6
//	      9, 10 :  5
//	     11, 12 :  4
//	     13, 14 :  3
//	     15, 16 :  2
//	     17, 18 :  1
//	         19 :  0 (lowest priority)
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

// Renice sets the CPU process priority. The nice parameter can range from
// -20 (least nice), to 19 (most nice), even on non-Unix systems such as
// Windows, plan9, etc. If not supported by the operating system, nil is returned.
func Renice(nice int) error {
	priority, ok := priorityMap[nice]
	if !ok {
		return &InvalidNiceError{nice}
	}
	filename := fmt.Sprintf("/proc/%d/ctl", os.Getpid())
	f, err := os.Open(filename)
	if err != nil {
		return &ReniceError{nice, err}
	}
	_, err = f.Write([]byte(fmt.Sprintf("pri %d", priority)))
	if err != nil {
		return &ReniceError{nice, err}
	}
	err = f.Close()
	if err != nil {
		return &ReniceError{nice, err}
	}

	return nil
}
