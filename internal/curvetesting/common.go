package curvetesting

import (
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/RyuaNerin/elliptic2"
	"github.com/RyuaNerin/elliptic2/internal/curve"
	"github.com/RyuaNerin/elliptic2/internal/curve/edwards"
	"github.com/RyuaNerin/elliptic2/internal/curve/montgomery"
	"github.com/RyuaNerin/elliptic2/internal/curve/twistededwards"
	"github.com/RyuaNerin/elliptic2/internal/curve/weierstrassbinary"
	"github.com/RyuaNerin/elliptic2/internal/curve/weierstrassprime"
	"github.com/stretchr/testify/require"
)

func W[
	TCurveArithmetic curve.CurveArithmeticBase,
	TOperator any,
](
	fnNewOp func(c TCurveArithmetic) TOperator,
) func(c curve.CurveArithmeticBase) TOperator {
	return func(c curve.CurveArithmeticBase) TOperator {
		return fnNewOp(c.(TCurveArithmetic))
	}
}

func GetCurveType(c elliptic2.Curve) string {
	if base := curve.GetBase(c); base != nil {
		switch base.(type) {
		case *edwards.Curve:
			return "Edwards"
		case *montgomery.Curve:
			return "Montgomery"
		case *twistededwards.Curve:
			return "Twisted Edwards"
		case *weierstrassprime.Curve:
			return "Weierstrass Prime"
		case *weierstrassbinary.Curve:
			return "Weierstrass Binary"
		}
	}

	if _, ok := c.(*elliptic.CurveParams); ok {
		return "Weierstrass Prime (Standard Library)"
	}

	panic(fmt.Sprintf("unsupported curve: %T", c))
}

func GetGenerator(c elliptic2.Curve) *Point {
	if base := curve.GetBase(c); base != nil {
		gx, gy, ok := base.Generator()
		if !ok {
			return nil
		}
		return &Point{X: gx, Y: gy}
	}

	if curveStd, ok := c.(elliptic.Curve); ok {
		params := curveStd.Params()
		if params.Gx == nil || params.Gy == nil || params.Gx.Sign() == 0 || params.Gy.Sign() == 0 {
			return nil
		}
		return &Point{X: params.Gx, Y: params.Gy}
	}

	panic(fmt.Sprintf("unsupported curve: %T", c))
}

func GetName(c elliptic2.Curve) string {
	if base := curve.GetBase(c); base != nil {
		return base.Params().Name
	}

	if curveStd, ok := c.(elliptic.Curve); ok {
		return curveStd.Params().Name
	}

	panic(fmt.Sprintf("unsupported curve: %T", c))
}

func GetRandomK(t testing.TB, c elliptic2.Curve) []byte {
	var n *big.Int

	if base := curve.GetBase(c); base != nil {
		n = base.Params().N
	}

	if curveStd, ok := c.(elliptic.Curve); ok {
		n = curveStd.Params().N
	}
	require.NotNil(t, n, "curve has no order")

	k, err := rand.Int(rand.Reader, n)
	require.NoError(t, err)
	return k.Bytes()
}
