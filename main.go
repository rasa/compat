// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

// main is a sample application that runs many of the functions provided by this
// library.
package main

import "fmt"

func main() {
	fmt.Println(greet())
}

func greet() string {
	return "Hi!"
}
