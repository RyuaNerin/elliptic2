package curve

import (
	"crypto/elliptic"
	"math/big"

	"github.com/RyuaNerin/elliptic2"
	"github.com/RyuaNerin/elliptic2/internal/field/simd"
)

type (
	curveGF2mMadd struct {
		curveGF2m
	}
)

var (
	_ elliptic2.CurveExtended = (*curveGF2mMadd)(nil)
	_ elliptic.Curve          = (*curveGF2mMadd)(nil)
)

func (c *curveGF2mMadd) ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
	c.panicIfNotOnCurve(x1, y1)

	return c.scalarMultWithMAdd(x1, y1, k)
}

func (c *curveGF2mMadd) ScalarBaseMult(k []byte) (x, y *big.Int) {
	gx, gy, ok := c.base.Generator()
	if !ok {
		panic("elliptic2: curve has no generator")
	}
	return c.scalarMultWithMAdd(gx, gy, k)
}

func (c *curveGF2mMadd) scalarMultWithMAdd(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
	if len(k) == 0 {
		return new(big.Int), new(big.Int)
	}

	for len(k) > 0 && k[len(k)-1] == 0 {
		k = k[:len(k)-1]
	}
	if len(k) == 0 {
		return new(big.Int), new(big.Int)
	}

	op := c.base.NewOperator().(GF2mMaddOperator)

	modulus := c.base.Modulus()

	// table[0] = P (Z=1)
	// table[1] = 3P (Z=1)
	// table[2] = 5P (Z=1)
	// ...
	var table, negTable [maddTableSize]GF2mCoordinate
	for idx := range table {
		table[idx].SetModulus(modulus)
		negTable[idx].SetModulus(modulus)
	}

	op.ToCoordinate(&table[0], x1, y1)
	op.Neg(&negTable[0], &table[0])

	// 2P 계산
	var p2 GF2mCoordinate
	p2.SetModulus(modulus)
	op.Double(&p2, &table[0])
	op.ScaleZ(&p2)

	for idx := 1; idx < maddTableSize; idx++ {
		op.Madd(&table[idx], &table[idx-1], &p2)
		op.ScaleZ(&table[idx])
		op.Neg(&negTable[idx], &table[idx])
	}

	// wNAF (windowed Non-Adjacent Form)
	var nafBuf [simd.WordByteSize * simd.Words * 8]int8
	naf := computeWNAF(nafBuf[:0], k)

	var resultValue, tmpValue GF2mCoordinate
	result, tmp := &resultValue, &tmpValue
	resultValue.SetModulus(modulus)
	tmpValue.SetModulus(modulus)

	op.SetInfinity(result)

	for idx := len(naf) - 1; idx >= 0; idx-- {
		if !op.IsInfinity(result) {
			op.Double(tmp, result)
			result, tmp = tmp, result
		}

		digit := naf[idx]
		if digit > 0 {
			// result += table[(digit-1)/2]
			idx := (digit - 1) / 2
			if op.IsInfinity(result) {
				result.Set(&table[idx])
			} else {
				op.Madd(tmp, result, &table[idx])
				result, tmp = tmp, result
			}
		} else if digit < 0 {
			// result -= table[(-digit-1)/2]
			idx := (-digit - 1) / 2
			if op.IsInfinity(result) {
				result.Set(&negTable[idx])
			} else {
				op.Madd(tmp, result, &negTable[idx])
				result, tmp = tmp, result
			}
		}
	}

	x, y = new(big.Int), new(big.Int)
	op.ToAffinePoint(x, y, result)
	return x, y
}
