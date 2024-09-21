package curve

import (
	"crypto/elliptic"
	"math/big"

	"github.com/RyuaNerin/elliptic2"
	"github.com/RyuaNerin/elliptic2/internal/field/simd"
)

type (
	curveGFpMadd struct {
		curveGFp
	}
)

var (
	_ elliptic2.CurveExtended = (*curveGFpMadd)(nil)
	_ elliptic.Curve          = (*curveGFpMadd)(nil)
)

func (c *curveGFpMadd) ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
	c.panicIfNotOnCurve(x1, y1)

	return c.scalarMultWithMAdd(x1, y1, k)
}

func (c *curveGFpMadd) ScalarBaseMult(k []byte) (x, y *big.Int) {
	gx, gy, ok := c.base.Generator()
	if !ok {
		panic("elliptic2: curve has no generator")
	}
	return c.scalarMultWithMAdd(gx, gy, k)
}

func (c *curveGFpMadd) scalarMultWithMAdd(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
	if len(k) == 0 {
		return new(big.Int), new(big.Int)
	}

	op := c.base.NewOperator().(GFpMaddOperator)

	modulus := c.base.Modulus()

	// table[0] = P (Z=1)
	// table[1] = 3P (Z=1)
	// table[2] = 5P (Z=1)
	// ...
	var table, negTable [maddTableSize]GFpCoordinate
	for idx := range table {
		table[idx].SetModulus(modulus)
		negTable[idx].SetModulus(modulus)
	}

	op.ToCoordinate(&table[0], x1, y1)
	op.Neg(&negTable[0], &table[0])

	// 2P 계산
	{
		var p2 GFpCoordinate
		p2.SetModulus(modulus)
		op.Double(&p2, &table[0])
		op.ScaleZ(&p2)

		for idx := 1; idx < maddTableSize; idx++ {
			op.Add(&table[idx], &table[idx-1], &p2)
			op.ScaleZ(&table[idx])
			op.Neg(&negTable[idx], &table[idx])
		}
	}

	// wNAF (windowed Non-Adjacent Form)
	var nafBuf [simd.WordByteSize * simd.Words * 8]int8
	naf := computeWNAF(nafBuf[:0], k)

	var resultValue, tmpValue GFpCoordinate
	result, tmp := &resultValue, &tmpValue
	resultValue.SetModulus(c.base.Modulus())
	tmpValue.SetModulus(c.base.Modulus())

	op.SetInfinity(result)

	for idx := len(naf) - 1; idx >= 0; idx-- {
		if !op.IsInfinity(result) {
			op.Double(tmp, result)
			result, tmp = tmp, result
		}

		digit := naf[idx]
		if digit > 0 {
			// result += table[(digit-1)/2]
			idxTable := (digit - 1) / 2
			if op.IsInfinity(result) {
				result.Set(&table[idxTable])
			} else {
				op.Madd(tmp, result, &table[idxTable])
				result, tmp = tmp, result
			}
		} else if digit < 0 {
			// result -= table[(-digit-1)/2]
			idxTable := (-digit - 1) / 2
			if op.IsInfinity(result) {
				result.Set(&negTable[idxTable])
			} else {
				op.Madd(tmp, result, &negTable[idxTable])
				result, tmp = tmp, result
			}
		}
	}

	x, y = new(big.Int), new(big.Int)
	op.ToAffinePoint(x, y, result)
	return x, y
}
