package montgomery

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
		A, B    *field.GFp     // a, b
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

var defaultOpFuncSelector = SelectAffine

func Build(params *CurveParams) *Curve {
	return BuildOp(params, defaultOpFuncSelector(params))
}

func BuildOp(params *CurveParams, fnNewOp func(c *Curve) curve.GFpOperator) *Curve {
	c := &Curve{
		CurveParams: *params,
		newOperator: fnNewOp,
	}

	c.A.SetModulus(params.P)
	c.B.SetModulus(params.P)

	if params.Gx != nil && params.Gx.Sign() != 0 && params.Gy != nil && params.Gy.Sign() != 0 {
		if !c.IsOnCurve(params.Gx, params.Gy) {
			panic("montgomery: invalid generator point")
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
		Type:    elliptic2.CurveTypeMontgomery,
		Name:    strings.Clone(c.Name),
		BitSize: c.BitSize,
		P:       c.P.ToBigInt(nil),
		N:       c.N.ToBigInt(nil),
		A:       c.A.ToBigInt(nil),
		B:       c.B.ToBigInt(nil),
		Gx:      gx,
		Gy:      gy,
	}
}

func (c *Curve) IsOnCurve(x, y *big.Int) bool {
	/**
	byy = xxx + axx + x
	*/

	A, B, P := c.A, c.B, c.P

	var X, Y field.GFp
	X.SetModulus(P).SetBigInt(x)
	Y.SetModulus(P).SetBigInt(y)

	var lhs, rhs field.GFp
	lhs.SetModulus(P)
	rhs.SetModulus(P)

	// rhs = xxx + axx + x
	lhs.Sqr(&X).ReduceSoft()       //   t = xx
	rhs.Mul(&X, &lhs).ReduceSoft() // rhs = xxx
	lhs.Mul(A, &lhs).ReduceSoft()  //         t = axx
	rhs.Add(&rhs, &lhs)            // rhs = xxx + axx
	rhs.Add(&rhs, &X)              // rhs = xxx + axx + x
	rhs.Reduce()

	// lhs = byy
	lhs.Sqr(&Y).ReduceSoft() // lhs =  yy
	lhs.Mul(B, &lhs)         // lhs = byy
	lhs.Reduce()

	return lhs.Cmp(&rhs) == 0
}

func (c *Curve) ComputeY(x *big.Int, largeY bool) *big.Int {
	/**
	byy = xxx + axx + x
	yy = (xxx + axx + x) / b
	y = sqrt((xxx + axx + x) / b)
	*/

	A, B, P := c.A, c.B, c.P

	var X field.GFp
	X.SetModulus(P).SetBigInt(x)

	var y1, y2, bInv field.GFp
	y1.SetModulus(P)
	y2.SetModulus(P)
	bInv.SetModulus(P)

	// bInv = 1/b
	if bInv.Inv(B) == nil {
		return nil
	}

	y2.Sqr(&X).ReduceSoft()      // y2 =        xx
	y1.Mul(&X, &y2).ReduceSoft() // y1 =       xxx
	y2.Mul(A, &y2).ReduceSoft()  //             y2 = axx
	y1.Add(&y1, &y2)             // y1 =       xxx + axx
	y1.Add(&y1, &X)              // y1 =       xxx + axx + x
	y1.Mul(&y1, &bInv).Reduce()  // y1 =      (xxx + axx + x) / b
	if y1.Sqrt(&y1) == nil {     // y1 = sqrt((xxx + axx + x) / b)
		return nil
	}

	y2.Neg(&y1).Reduce()

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
