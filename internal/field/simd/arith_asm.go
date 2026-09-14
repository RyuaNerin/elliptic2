// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build gc && (amd64 || 386 || arm || arm64) && !purego

package simd

import "math/big"

//go:noescape
func addVV(z, x, y []big.Word) (c big.Word)

//go:noescape
func subVV(z, x, y []big.Word) (c big.Word)

//go:noescape
func lshVU(z, x []big.Word, s uint) (c big.Word)

//go:noescape
func rshVU(z, x []big.Word, s uint) (c big.Word)

//go:noescape
func mulAddVWW(z, x []big.Word, m, a big.Word) (c big.Word)

//go:noescape
func addMulVVWW(z, x, y []big.Word, m, a big.Word) (c big.Word)
