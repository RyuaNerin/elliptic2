// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build gc && amd64 && !purego

package simd

import "golang.org/x/sys/cpu"

var hasADX = cpu.X86.HasADX && cpu.X86.HasBMI2
