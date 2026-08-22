package curve

import (
	"crypto/elliptic"
	"encoding/asn1"
	"fmt"
	"math/big"

	"github.com/RyuaNerin/elliptic2"
	"github.com/RyuaNerin/elliptic2/internal/field"
)

type (
	GFpCurveArithmetic CurveArithmetic[GFpCoordinate, *GFpCoordinate, GFpOperator]
	GFpOperator        Operator[GFpCoordinate, *GFpCoordinate]
	GFpMaddOperator    MaddOperator[GFpCoordinate, *GFpCoordinate]

	curveGFp struct {
		curveBase
		base GFpCurveArithmetic
	}
)

var (
	_ elliptic2.CurveExtended = (*curveGFp)(nil)
	_ elliptic.Curve          = (*curveGFp)(nil)
)

func (c *curveGFp) Params() *elliptic.CurveParams   { return c.base.Params() }
func (c *curveGFp) Params2() *elliptic2.CurveParams { return c.base.Params2() }

func (c *curveGFp) IsOnCurve(x, y *big.Int) bool { return c.base.IsOnCurve(x, y) }

func (c *curveGFp) IsEquivalentCurve(other elliptic2.Curve) bool {
	eq, _ := c.isEquivalent(other, false)
	return eq
}

var (
	oidNamedCurveP224 = asn1.ObjectIdentifier{1, 3, 132, 0, 33}
	oidNamedCurveP256 = asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7}
	oidNamedCurveP384 = asn1.ObjectIdentifier{1, 3, 132, 0, 34}
	oidNamedCurveP521 = asn1.ObjectIdentifier{1, 3, 132, 0, 35}
)

func (c *curveGFp) IsEquivalentCurveWithOID(other elliptic2.Curve) bool {
	eq, otherIsStd := c.isEquivalent(other, true)
	if !eq {
		return false
	}

	if otherIsStd {
		switch {
		case other == elliptic.P224() && c.oid.Equal(oidNamedCurveP224):
			return true
		case other == elliptic.P256() && c.oid.Equal(oidNamedCurveP256):
			return true
		case other == elliptic.P384() && c.oid.Equal(oidNamedCurveP384):
			return true
		case other == elliptic.P521() && c.oid.Equal(oidNamedCurveP521):
			return true
		}

		return false
	}

	return true
}

func (c *curveGFp) isEquivalent(other elliptic2.Curve, withOID bool) (eq bool, otherIsStd bool) {
	if isNilCurve(c) || isNilCurve(other) {
		return false, false
	}
	if isSameCurveInstance(c, other) {
		return true, false
	}

	var C *curveGFp

	if c2, ok := other.(*curveGFp); ok {
		C = c2
	}
	if c2, ok := other.(*curveGFpMadd); ok {
		C = &c2.curveGFp
	}
	if C == nil {
		// compare with standard library curve
		if c.base.CurveType() != elliptic2.CurveTypeWeierstrassPrime {
			return false, false
		}

		// (A + 3) % P == 0  <=>  A == -3 mod P
		params := c.base.Params2()
		var tmp field.GFp
		var a field.GFp
		a.SetModulus(c.base.Modulus()).SetBigInt(params.A)
		tmp.SetModulus(c.base.Modulus()).Add(&a, field.GFpThree).Reduce()
		if !tmp.IsZero() {
			return false, false
		}
		if C, ok := other.(elliptic.Curve); ok {
			return isEquivalentParams(params.Params(), C.Params()), true
		}
		return false, false
	}

	if withOID {
		if c.hashOID != C.hashOID {
			return false, false
		}
		// full comparison of curve parameters
		return isEquivalentParams2(c.base.Params2(), C.base.Params2()) && c.oid.Equal(C.oid), false
	}

	if c.hash != C.hash {
		return false, false
	}
	// full comparison of curve parameters
	return isEquivalentParams2(c.base.Params2(), C.base.Params2()), false
}

func (c *curveGFp) HasGenerator() bool {
	_, _, ok := c.base.Generator()
	return ok
}

func (c *curveGFp) ComputeY(x *big.Int, largeY bool) *big.Int { return c.base.ComputeY(x, largeY) }

func (c *curveGFp) Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	c.panicIfNotOnCurve(x1, y1)
	c.panicIfNotOnCurve(x2, y2)

	if c.base.IsInfinity(x1, y1) {
		return new(big.Int).Set(x2), new(big.Int).Set(y2)
	}
	if c.base.IsInfinity(x2, y2) {
		return new(big.Int).Set(x1), new(big.Int).Set(y1)
	}

	op := c.base.NewOperator()

	var dst, p1 GFpCoordinate
	dst.SetModulus(c.base.Modulus())
	p1.SetModulus(c.base.Modulus())

	op.ToCoordinate(&p1, x1, y1)

	if x1.Cmp(x2) == 0 && y1.Cmp(y2) == 0 {
		op.Double(&dst, &p1)
	} else {
		var p2 GFpCoordinate
		p2.SetModulus(c.base.Modulus())
		op.ToCoordinate(&p2, x2, y2)
		op.Add(&dst, &p1, &p2)
	}

	x, y = new(big.Int), new(big.Int)
	op.ToAffinePoint(x, y, &dst)
	return x, y
}

func (c *curveGFp) Double(x1, y1 *big.Int) (x, y *big.Int) {
	c.panicIfNotOnCurve(x1, y1)

	if c.base.IsInfinity(x1, y1) {
		return c.base.InfinityPoint()
	}

	op := c.base.NewOperator()

	var dst, p1 GFpCoordinate
	dst.SetModulus(c.base.Modulus())
	p1.SetModulus(c.base.Modulus())

	op.ToCoordinate(&p1, x1, y1)

	op.Double(&dst, &p1)

	x, y = new(big.Int), new(big.Int)
	op.ToAffinePoint(x, y, &dst)
	return x, y
}

func (c *curveGFp) ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
	c.panicIfNotOnCurve(x1, y1)

	return c.scalarMult(x1, y1, k)
}

func (c *curveGFp) ScalarBaseMult(k []byte) (x, y *big.Int) {
	x, y, ok := c.base.Generator()
	if !ok {
		panic("elliptic2: curve has no generator")
	}
	return c.scalarMult(x, y, k)
}

func (c *curveGFp) scalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
	if len(k) == 0 {
		return c.base.InfinityPoint()
	}

	var num big.Int
	num.SetBytes(k)
	if num.Sign() == 0 {
		return c.base.InfinityPoint()
	}

	op := c.base.NewOperator()

	var r0Value, r1Value, tmValue GFpCoordinate
	r0, r1, tm := &r0Value, &r1Value, &tmValue

	r0.SetModulus(c.base.Modulus())
	r1.SetModulus(c.base.Modulus())
	tm.SetModulus(c.base.Modulus())

	op.ToCoordinate(r0, x1, y1)

	// Montgomery Ladder
	op.Double(r1, r0)
	for i := num.BitLen() - 2; i >= 0; i-- {
		if num.Bit(i) == 1 {
			op.Add(tm, r0, r1)
			r0, tm = tm, r0

			op.Double(tm, r1)
			r1, tm = tm, r1
		} else {
			op.Add(tm, r1, r0)
			r1, tm = tm, r1

			op.Double(tm, r0)
			r0, tm = tm, r0
		}
	}

	x, y = new(big.Int), new(big.Int)
	op.ToAffinePoint(x, y, r0)
	return x, y
}

func (c *curveGFp) panicIfNotOnCurve(x, y *big.Int) {
	if x.Sign() == 0 && y.Sign() == 0 {
		return
	}

	if !c.base.IsOnCurve(x, y) {
		panic(fmt.Sprintf("elliptic2: point (%s, %s) is not on curve %s", x.Text(16), y.Text(16), c.base.Params().Name))
	}
}
