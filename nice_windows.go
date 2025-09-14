// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build windows

package compat

import (
	"fmt"

	"golang.org/x/sys/windows"
)

var niceMap = map[uint32]int{
	windows.REALTIME_PRIORITY_CLASS:     -20,
	windows.HIGH_PRIORITY_CLASS:         -15,
	windows.ABOVE_NORMAL_PRIORITY_CLASS: -5,
	windows.NORMAL_PRIORITY_CLASS:       0,
	windows.BELOW_NORMAL_PRIORITY_CLASS: 5,  //nolint:mnd
	windows.IDLE_PRIORITY_CLASS:         19, //nolint:mnd
}

// Nice gets the CPU process priority. The return value is in a range from
// -20 (least nice), to 19 (most nice), even on non-Unix systems such as
// Windows, plan9, etc. If not supported by the operating system, an error is
// returned.
func Nice() (int, error) {
	handle := windows.CurrentProcess()

	priorityClass, err := windows.GetPriorityClass(handle)
	if err != nil {
		return 0, &NiceError{err}
	}

	nice, ok := niceMap[priorityClass]
	if !ok {
		panic(fmt.Sprintf("nice: unknown priority class %v", priorityClass))
	}

	return nice, nil
}

// reniceMap maps unix's nice levels (-20 to 19), to Windows' levels:
//
//	input value      Windows priority
//
// ---------------   ---------------------
//
//		       -20 : Realtime priority (use with caution)
//	    -19 to -10 : High priority
//		 -9 to  -1 : Above normal priority
//		         0 : Normal priority (Windows' default priority)
//		  1 to   9 : Below normal priority
//	     10 to  19 : Idle priority
var reniceMap = map[int]uint32{
	-20: windows.REALTIME_PRIORITY_CLASS,
	-19: windows.HIGH_PRIORITY_CLASS,
	-18: windows.HIGH_PRIORITY_CLASS,
	-17: windows.HIGH_PRIORITY_CLASS,
	-16: windows.HIGH_PRIORITY_CLASS,
	-15: windows.HIGH_PRIORITY_CLASS,
	-14: windows.HIGH_PRIORITY_CLASS,
	-13: windows.HIGH_PRIORITY_CLASS,
	-12: windows.HIGH_PRIORITY_CLASS,
	-11: windows.HIGH_PRIORITY_CLASS,
	-10: windows.HIGH_PRIORITY_CLASS,
	-9:  windows.ABOVE_NORMAL_PRIORITY_CLASS,
	-8:  windows.ABOVE_NORMAL_PRIORITY_CLASS,
	-7:  windows.ABOVE_NORMAL_PRIORITY_CLASS,
	-6:  windows.ABOVE_NORMAL_PRIORITY_CLASS,
	-5:  windows.ABOVE_NORMAL_PRIORITY_CLASS,
	-4:  windows.ABOVE_NORMAL_PRIORITY_CLASS,
	-3:  windows.ABOVE_NORMAL_PRIORITY_CLASS,
	-2:  windows.ABOVE_NORMAL_PRIORITY_CLASS,
	-1:  windows.ABOVE_NORMAL_PRIORITY_CLASS,
	0:   windows.NORMAL_PRIORITY_CLASS,
	1:   windows.BELOW_NORMAL_PRIORITY_CLASS,
	2:   windows.BELOW_NORMAL_PRIORITY_CLASS,
	3:   windows.BELOW_NORMAL_PRIORITY_CLASS,
	4:   windows.BELOW_NORMAL_PRIORITY_CLASS,
	5:   windows.BELOW_NORMAL_PRIORITY_CLASS,
	6:   windows.BELOW_NORMAL_PRIORITY_CLASS,
	7:   windows.BELOW_NORMAL_PRIORITY_CLASS,
	8:   windows.BELOW_NORMAL_PRIORITY_CLASS,
	9:   windows.BELOW_NORMAL_PRIORITY_CLASS,
	10:  windows.IDLE_PRIORITY_CLASS,
	11:  windows.IDLE_PRIORITY_CLASS,
	12:  windows.IDLE_PRIORITY_CLASS,
	13:  windows.IDLE_PRIORITY_CLASS,
	14:  windows.IDLE_PRIORITY_CLASS,
	15:  windows.IDLE_PRIORITY_CLASS,
	16:  windows.IDLE_PRIORITY_CLASS,
	17:  windows.IDLE_PRIORITY_CLASS,
	18:  windows.IDLE_PRIORITY_CLASS,
	19:  windows.IDLE_PRIORITY_CLASS,
}

// Renice sets the CPU process priority. The nice parameter can range from
// -20 (least nice), to 19 (most nice), even on non-Unix systems such as
// Windows, plan9, etc.
func Renice(nice int) error {
	priorityClass, ok := reniceMap[nice]
	if !ok {
		return &InvalidNiceError{nice}
	}

	handle := windows.CurrentProcess()

	err := windows.SetPriorityClass(handle, priorityClass)
	if err != nil {
		return &ReniceError{nice, err}
	}

	return nil
}
