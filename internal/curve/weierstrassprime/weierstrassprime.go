package weierstrassprime

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
		N       *big.Int       // order
		Gx, Gy  *big.Int       // generator point
	}
	Curve struct {
		CurveParams
		newOperator   NewOpFunc
		withGenerator bool
		nBytes        int
	}
	NewOpFunc func(c *Curve) curve.GFpOperator
)

var _ curve.CurveArithmeticBase = (*Curve)(nil)

var defaultOpFuncSelector = SelectJacobian

func Build(params *CurveParams) *Curve {
	return BuildOp(params, defaultOpFuncSelector(params))
}

func BuildOp(params *CurveParams, fnNewOp func(c *Curve) curve.GFpOperator) *Curve {
	c := &Curve{
		CurveParams: *params,
		newOperator: fnNewOp,
		nBytes:      (params.N.BitLen() + 7) / 8,
	}

	c.A.SetModulus(params.P)
	c.B.SetModulus(params.P)

	if params.Gx != nil && params.Gx.Sign() != 0 && params.Gy != nil && params.Gy.Sign() != 0 {
		if !c.IsOnCurve(params.Gx, params.Gy) {
			panic("weierstrassprime: invalid generator point")
		}

		c.withGenerator = true
	}
	return c
}

func (c *Curve) RawParams() any                      { return &c.CurveParams }
func (c *Curve) Modulus() *field.Modulus             { return c.P }
func (*Curve) FieldType() curve.FieldType            { return curve.FieldTypeGFp }
func (c *Curve) Generator() (x, y *big.Int, ok bool) { return c.Gx, c.Gy, c.withGenerator }
func (*Curve) IsInfinity(x, y *big.Int) bool         { return x.Sign() == 0 && y.Sign() == 0 }
func (*Curve) Identity() (x, y *big.Int)             { return new(big.Int), new(big.Int) }
func (c *Curve) NewOperator() curve.GFpOperator      { return c.newOperator(c) }
func (c *Curve) N() *big.Int                         { return c.CurveParams.N }

func (c *Curve) Params() *elliptic.CurveParams {
	p := &elliptic.CurveParams{
		P:       c.P.ToBigInt(nil),
		N:       new(big.Int).Set(c.CurveParams.N),
		B:       c.B.ToBigInt(nil),
		BitSize: c.BitSize,
		Name:    c.Name,
	}
	if c.withGenerator {
		p.Gx, p.Gy = new(big.Int).Set(c.Gx), new(big.Int).Set(c.Gy)
	}
	return p
}

func (c *Curve) Params2() *elliptic2.CurveParams {
	p := &elliptic2.CurveParams{
		Type:    elliptic2.CurveTypeWeierstrassPrime,
		Name:    strings.Clone(c.Name),
		BitSize: c.BitSize,
		P:       c.P.ToBigInt(nil),
		N:       new(big.Int).Set(c.CurveParams.N),
		A:       c.A.ToBigInt(nil),
		B:       c.B.ToBigInt(nil),
	}
	if c.withGenerator {
		p.Gx, p.Gy = new(big.Int).Set(c.Gx), new(big.Int).Set(c.Gy)
	}
	p.InfX, p.InfY = c.Identity()
	return p
}

func (c *Curve) IsOnCurve(x, y *big.Int) bool {
	/**
	yy = xxx + ax + b
	*/

	A, B, P := c.A, c.B, c.P

	var X, Y field.GFp
	X.SetModulus(P).SetBigInt(x)
	Y.SetModulus(P).SetBigInt(y)

	var lhs, rhs field.GFp
	lhs.SetModulus(P)
	rhs.SetModulus(P)

	// rhs = xxx + ax + b
	rhs.Mul(rhs.Sqr(&X), &X).ReduceSoft()      // rhs = xxx
	rhs.Add(&rhs, B)                           // rhs = xxx + b
	rhs.Add(&rhs, lhs.Mul(A, &X).ReduceSoft()) // rhs = xxx + b + ax
	rhs.Reduce()

	// lhs = yy
	lhs.Sqr(&Y).Reduce() // lhs = yy

	return lhs.Cmp(&rhs) == 0
}

func (c *Curve) ComputeY(x *big.Int, largeY bool) *big.Int {
	/**
	yy = xxx + ax + b
	y = sqrt(xxx + ax + b)
	*/

	A, B, P := c.A, c.B, c.P

	var X field.GFp
	X.SetModulus(P).SetBigInt(x)

	var y1, y2 field.GFp
	y1.SetModulus(P)
	y2.SetModulus(P)

	y1.Mul(y1.Sqr(&X), &X).ReduceSoft()     // Y = xxx
	y1.Add(&y1, y2.Mul(A, &X).ReduceSoft()) // Y = xxx + ax
	y1.Add(&y1, B)                          // Y = xxx + ax + b
	y1.Reduce()                             //
	if y1.Sqrt(&y1) == nil {                // Y = sqrt(xxx + ax + b)
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
