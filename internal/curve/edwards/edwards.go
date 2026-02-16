package edwards

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
		C, D    *field.GFp     // c, d
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

var defaultOpFuncSelector = SelectProjectiveOp

func Build(params *CurveParams) *Curve {
	return BuildOp(params, defaultOpFuncSelector(params))
}

func BuildOp(params *CurveParams, fnNewOp func(c *Curve) curve.GFpOperator) *Curve {
	c := &Curve{
		CurveParams: *params,
		newOperator: fnNewOp,
	}

	c.C.SetModulus(c.P)
	c.D.SetModulus(c.P)

	if params.Gx != nil && params.Gx.Sign() != 0 && params.Gy != nil && params.Gy.Sign() != 0 {
		if !c.IsOnCurve(params.Gx, params.Gy) {
			panic("edwards: invalid generator point")
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
	return &elliptic2.CurveParams{
		Type:    elliptic2.CurveTypeEdwards,
		Name:    strings.Clone(c.Name),
		BitSize: c.BitSize,
		P:       c.P.ToBigInt(nil),
		N:       c.N.ToBigInt(nil),
		C:       c.C.ToBigInt(nil),
		D:       c.D.ToBigInt(nil),
		Gx:      new(big.Int).Set(c.Gx),
		Gy:      new(big.Int).Set(c.Gy),
	}
}

func (c *Curve) IsOnCurve(x, y *big.Int) bool {
	/**
	xx + yy = cc(1 + dxxyy)
	*/

	C, D, P := c.C, c.D, c.P

	var X, Y field.GFp
	X.SetModulus(P).SetBigInt(x)
	Y.SetModulus(P).SetBigInt(y)

	var xx, yy, lhs, rhs field.GFp
	xx.SetModulus(P)
	yy.SetModulus(P)
	lhs.SetModulus(P)
	rhs.SetModulus(P)

	xx.Sqr(&X).ReduceSoft() // xx = xx
	yy.Sqr(&Y).ReduceSoft() // yy = yy

	// lhs = xx + yy
	lhs.Add(&xx, &yy) // lsh = xx + yy
	lhs.Reduce()

	// rhs = cc(1 + dxxyy)
	rhs.Mul(&xx, &yy).ReduceSoft() // rhs =         xxyy
	rhs.Mul(&rhs, D).ReduceSoft()  // rhs =        dxxyy
	rhs.Add(&rhs, field.GFpOne)    // rhs =    1 + dxxyy
	yy.Sqr(C).ReduceSoft()         // yy  = cc
	rhs.Mul(&rhs, &yy)             // rhs = cc(1 + dxxyy)
	rhs.Reduce()

	return lhs.Cmp(&rhs) == 0
}

func (c *Curve) ComputeY(x *big.Int, largeY bool) *big.Int {
	/**
	xx + yy = cc(1 + dxxyy)
	xx + yy = cc + ccdxxyy
	yy - ccdxxyy = cc - xx
	yy(1 - ccdxx) = cc - xx
	yy = (cc - xx) / (1 - ccdxx)
	y = sqrt( (cc - xx) / (1 - ccdxx) )
	*/

	C, D, P := c.C, c.D, c.P

	var X field.GFp
	X.SetModulus(P).SetBigInt(x)

	var y1, y2, tm field.GFp
	y1.SetModulus(P)
	y2.SetModulus(P)
	tm.SetModulus(P)

	tm.Sqr(C).ReduceSoft()                 // tm =          cc
	y2.Sqr(&X).ReduceSoft()                // y2 =      xx
	y1.Sub(&tm, &y2).ReduceSoft()          // y1 = cc - xx
	tm.Mul(&tm, D).ReduceSoft()            // tm =          ccd
	tm.Mul(&tm, &y2).ReduceSoft()          // tm =          ccdxx
	tm.Sub(field.GFpOne, &tm).ReduceSoft() // tm =      1 - ccdxx
	if tm.Inv(&tm) == nil {                // tm = 1 / (1 - ccdxx)
		return nil
	}

	y1.Mul(&y1, &tm).ReduceSoft() // y1 = (cc - xx) / (1 - ccdxx)
	if y1.Sqrt(&y1) == nil {
		return nil
	}

	y2.Neg(&y1).Reduce() // y2 = -y1

	if y1.Cmp(&y2) < 0 { // y1 < y2
		if largeY {
			return y2.ToBigInt(nil)
		} else {
			return y1.ToBigInt(nil)
		}
	} else { // y1 > y2
		if largeY {
			return y1.ToBigInt(nil)
		} else {
			return y2.ToBigInt(nil)
		}
	}
}
