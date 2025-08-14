// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build wasip1 && tinygo

package compat

import "time"

// Not supported: ATime | BTime | CTime | Symlinks
const supports supportsType = supportsLinks

const userIDSource UserIDSourceType = UserIDSourceIsNone

func (fs *fileStat) times() {}

func (fs *fileStat) BTime() time.Time { return fs.btime }
func (fs *fileStat) CTime() time.Time { return fs.ctime }

func (fs *fileStat) UID() int { return fs.uid }
func (fs *fileStat) GID() int { return fs.gid }

func (fs *fileStat) User() string  { return fs.user }
func (fs *fileStat) Group() string { return fs.group }
