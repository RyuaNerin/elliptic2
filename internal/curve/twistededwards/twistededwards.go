package twistededwards

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
		P       *field.Modulus // prime
		A, D    *field.GFp     // a, d
		N       *field.GFp     // order
		Gx, Gy  *big.Int       // generator point
	}
	Curve struct {
		CurveParams
		newOperator   NewOpFunc
		withGenerator bool
	}
	NewOpFunc func(c *Curve) curve.GFpOperator
)

var _ curve.CurveArithmeticBase = (*Curve)(nil)

var defaultOpFuncSelector = SelectExtended

func Build(params *CurveParams) *Curve {
	return BuildOp(params, defaultOpFuncSelector(params))
}

func BuildOp(params *CurveParams, fnNewOp func(c *Curve) curve.GFpOperator) *Curve {
	c := &Curve{
		CurveParams: *params,
		newOperator: fnNewOp,
	}

	c.A.SetModulus(params.P)
	c.D.SetModulus(params.P)

	if params.Gx != nil && params.Gx.Sign() != 0 && params.Gy != nil && params.Gy.Sign() != 0 {
		if !c.IsOnCurve(params.Gx, params.Gy) {
			panic("twistededwards: invalid generator point")
		}

		c.withGenerator = true
	}
	return c
}

func (c *Curve) RawParams() any                      { return &c.CurveParams }
func (c *Curve) Modulus() *field.Modulus             { return c.P }
func (Curve) FieldType() curve.FieldType             { return curve.FieldTypeGFp }
func (c *Curve) Generator() (x, y *big.Int, ok bool) { return c.Gx, c.Gy, c.withGenerator }
func (c *Curve) NewOperator() curve.GFpOperator      { return c.newOperator(c) }

func (c *Curve) Params() *elliptic.CurveParams {
	p := &elliptic.CurveParams{
		P:       c.P.ToBigInt(nil),
		N:       c.N.ToBigInt(nil),
		BitSize: c.BitSize,
		Name:    c.Name,
	}
	if c.withGenerator {
		p.Gx, p.Gy = c.Gx, c.Gy
	}
	return p
}

func (c *Curve) Params2() *elliptic2.CurveParams {
	var gx, gy *big.Int
	if c.withGenerator {
		if c.Gx != nil {
			gx = new(big.Int).Set(c.Gx)
		}
		if c.Gy != nil {
			gy = new(big.Int).Set(c.Gy)
		}
	}

	return &elliptic2.CurveParams{
		Type:    elliptic2.CurveTypeTwistedEdwards,
		Name:    strings.Clone(c.Name),
		BitSize: c.BitSize,
		P:       c.P.ToBigInt(nil),
		N:       c.N.ToBigInt(nil),
		A:       c.A.ToBigInt(nil),
		D:       c.D.ToBigInt(nil),
		Gx:      gx,
		Gy:      gy,
	}
}

func (c *Curve) IsOnCurve(x, y *big.Int) bool {
	/**
	axx + yy = 1 + dxxyy
	*/

	A, D, P := c.A, c.D, c.P

	var X, Y field.GFp
	X.SetModulus(P).SetBigInt(x)
	Y.SetModulus(P).SetBigInt(y)

	var xx, yy, lhs, rhs field.GFp
	xx.SetModulus(P)
	yy.SetModulus(P)
	lhs.SetModulus(P)
	rhs.SetModulus(P)

	xx.Sqr(&X).ReduceSoft()
	yy.Sqr(&Y).ReduceSoft()

	// lhs = axx + yy
	lhs.Mul(A, &xx)    // lhs = axx
	lhs.Add(&lhs, &yy) // lhs = axx + yy
	lhs.Reduce()

	// rhs = 1 + dxxyy
	rhs.Mul(D, &xx).ReduceSoft() // rhs =     dxx
	rhs.Mul(&rhs, &yy)           // rhs =     dxxyy
	rhs.Add(field.GFpOne, &rhs)  // rhs = 1 + dxxyy
	rhs.Reduce()

	return lhs.Cmp(&rhs) == 0
}

func (c *Curve) ComputeY(x *big.Int, largeY bool) *big.Int {
	/**
	axx + yy = 1 + dxxyy

	yy - dxxyy = 1 - axx
	yy(1 - dxx) = 1 - axx
	yy = (1 - axx) / (1 - dxx)
	y = sqrt( (1 - axx) / (1 - dxx) )
	*/

	A, D, P := c.A, c.D, c.P

	var X field.GFp
	X.SetModulus(P).SetBigInt(x)

	var y1, y2 field.GFp
	y1.SetModulus(P)
	y2.SetModulus(P)

	y2.Sqr(&X).ReduceSoft()            // y2 = xx
	y1.Mul(A, &y2).ReduceSoft()        // y1 = axx
	y2.Mul(D, &y2).ReduceSoft()        // y2 = dxx
	y1.Sub(field.GFpOne, &y1).Reduce() // y1 = 1 - axx
	y2.Sub(field.GFpOne, &y2)          // y2 =      1 - dxx
	if y2.Inv(&y2) == nil {            // y2 = 1 / (1 - dxx)
		return nil
	}
	y1.Mul(&y1, &y2).ReduceSoft() // y1 =       (1 - axx) / (1 - dxx)
	if y1.Sqrt(&y1) == nil {      // y1 = sqrt( (1 - axx) / (1 - dxx) )
		return nil
	}

	y2.Neg(&y1).Reduce() // y2 = -y1

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
