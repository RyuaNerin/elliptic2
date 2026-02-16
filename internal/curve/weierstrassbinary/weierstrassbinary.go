package weierstrassbinary

import (
	"crypto/elliptic"
	"math/big"
	"strings"

	"github.com/RyuaNerin/elliptic2"
	"github.com/RyuaNerin/elliptic2/internal/curve"
	"github.com/RyuaNerin/elliptic2/internal/field"
)

type (
	CurveParams struct {
		Name    string
		BitSize int
		Poly    *field.Modulus // Polynomials
		A2, A6  *field.GF2m    // a, b
		N       *field.GF2m    // order
		Gx, Gy  *big.Int       // generator point
	}
	Curve struct {
		CurveParams
		newOperator   NewOpFunc
		withGenerator bool
	}
	NewOpFunc func(c *Curve) curve.GF2mOperator
)

var _ curve.CurveArithmeticBase = (*Curve)(nil)

var defaultOpFuncSelector = SelectJacobianOp

func Build(params *CurveParams) *Curve {
	return BuildOp(params, defaultOpFuncSelector(params))
}

func BuildOp(params *CurveParams, fnNewOp func(c *Curve) curve.GF2mOperator) *Curve {
	c := &Curve{
		CurveParams: *params,
		newOperator: fnNewOp,
	}

	c.A2.SetModulus(c.Poly)
	c.A6.SetModulus(c.Poly)
	c.N.SetModulus(c.Poly)

	if params.Gx != nil && params.Gx.Sign() != 0 && params.Gy != nil && params.Gy.Sign() != 0 {
		if !c.IsOnCurve(params.Gx, params.Gy) {
			panic("weierstrassbinary: invalid generator point")
		}

		c.withGenerator = true
	}
	return c
}

func (c *Curve) RawParams() any                      { return &c.CurveParams }
func (c *Curve) Modulus() *field.Modulus             { return c.Poly }
func (Curve) FieldType() curve.FieldType             { return curve.FieldTypeGF2m }
func (c *Curve) Generator() (x, y *big.Int, ok bool) { return c.Gx, c.Gy, c.withGenerator }
func (c *Curve) NewOperator() curve.GF2mOperator     { return c.newOperator(c) }

func (c *Curve) Params() *elliptic.CurveParams {
	p := &elliptic.CurveParams{
		Name:    c.Name,
		BitSize: c.BitSize,
		N:       c.N.ToBigInt(nil),
	}
	if c.withGenerator {
		p.Gx, p.Gy = c.Gx, c.Gy
	}
	return p
}

func (c *Curve) Params2() *elliptic2.CurveParams {
	return &elliptic2.CurveParams{
		Type:    elliptic2.CurveTypeWeierstrassBinary,
		Name:    strings.Clone(c.Name),
		BitSize: c.BitSize,
		P:       c.Poly.ToBigInt(nil),
		N:       c.N.ToBigInt(nil),
		A:       c.A2.ToBigInt(nil),
		B:       c.A6.ToBigInt(nil),
		Gx:      new(big.Int).Set(c.Gx),
		Gy:      new(big.Int).Set(c.Gy),
	}
}

func (c *Curve) IsOnCurve(x, y *big.Int) bool {
	var X, Y field.GF2m
	var tmp, lhs, rhs field.GF2m

	X.SetModulus(c.Poly)
	Y.SetModulus(c.Poly)
	tmp.SetModulus(c.Poly)
	lhs.SetModulus(c.Poly)
	rhs.SetModulus(c.Poly)

	X.SetBigInt(x)
	Y.SetBigInt(y)

	A, B := c.A2, c.A6

	// yy + xy = xxx + axx + b

	lhs.Sqr(&Y)                    // lhs = yy
	lhs.Add(&lhs, tmp.Mul(&X, &Y)) // lhs = yy + xy

	tmp.Sqr(&X)       //       xx
	rhs.Mul(&tmp, &X) // rhs = xxx
	rhs.Add(&rhs, B)  // rhs = xxx + b
	if !A.IsZero() {
		tmp.Mul(&tmp, A)    //                 axx
		rhs.Add(&rhs, &tmp) // rhs = xxx + b + axx
	}

	return lhs.Cmp(&rhs) == 0
}

func (c *Curve) ComputeY(x *big.Int, largeY bool) *big.Int {
	/**
	yy + xy = xxx + axx + b

	yy + xy = xxx + axx + b
	z = y/x
	-> zz + z = x + a + b/xx

	z = HalfTrace(x + a + b/xx)
	-> y = x * z
	*/

	var X, y1, y2, z field.GF2m
	X.SetModulus(c.Poly)
	y1.SetModulus(c.Poly)
	y2.SetModulus(c.Poly)
	z.SetModulus(c.Poly)

	X.SetBigInt(x)

	// x = 0인 경우: y = sqrt(b)
	if X.IsZero() {
		y1.Sqrt(c.A6)
		return y1.ToBigInt(nil)
	}

	A, B := c.A2, c.A6

	// w = x + a + b/(xx)
	y2.Sqr(&X)      //                xx
	y2.Inv(&y2)     //             1/(xx)
	y2.Mul(B, &y2)  // w =         b/(xx)
	y2.Add(&y2, &X) // w = x +     b/(xx)
	y2.Add(&y2, A)  // w = x + a + b/(xx)

	// zz + z = w
	if y2.Trace() != 0 {
		return nil
	}
	// z = HalfTrace(w)
	if c.Poly.Degree()%2 == 1 {
		z.HalfTrace(&y2)
	} else {
		gamma := c.Poly.FindGF2mGamma()
		z.SolveQuadraticEven(&y2, gamma)
	}

	y1.Mul(&X, &z)  // y1 = xz
	y2.Add(&y1, &X) // y2 = xz + x = x(z + 1)

	if y1.Cmp(&y2) < 0 { // y1 < y2
		if largeY {
			return y2.ToBigInt(nil)
		}
		return y1.ToBigInt(nil)
	}
	// y1 > y2
	if largeY {
		return y1.ToBigInt(nil)
	}
	return y2.ToBigInt(nil)
}
