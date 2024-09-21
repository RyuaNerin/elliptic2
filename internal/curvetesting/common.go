package curvetesting

import (
	"bufio"
	"crypto/elliptic"
	"fmt"
	"math/big"
	"math/rand"

	"github.com/RyuaNerin/elliptic2"
	"github.com/RyuaNerin/elliptic2/internal"
	"github.com/RyuaNerin/elliptic2/internal/curve"
	"github.com/RyuaNerin/elliptic2/internal/curve/edwards"
	"github.com/RyuaNerin/elliptic2/internal/curve/montgomery"
	"github.com/RyuaNerin/elliptic2/internal/curve/twistededwards"
	"github.com/RyuaNerin/elliptic2/internal/curve/weierstrassbinary"
	"github.com/RyuaNerin/elliptic2/internal/curve/weierstrassprime"
)

// var Random = bufio.NewReaderSize(rand.Reader, 1<<15)
var Random = bufio.NewReaderSize(rand.New(rand.NewSource(0)), 1<<15)

type TestingT interface {
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	FailNow()
	SkipNow()
}

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

func GetRandomK(c elliptic2.Curve) []byte {
	var n *big.Int

	if base := curve.GetBase(c); base != nil {
		n = base.Params().N
	}

	if curveStd, ok := c.(elliptic.Curve); ok {
		n = curveStd.Params().N
	}

	if n == nil {
		panic(fmt.Sprintf("unsupported curve: %T", c))
	}

	k, _ := internal.Int(Random, n)
	return k.Bytes()
}

func RequireIsOnCurve(t TestingT, curve elliptic2.Curve, x, y *big.Int) bool {
	if !curve.IsOnCurve(x, y) {
		t.Errorf("point not on curve")
		t.FailNow()
		return false
	}
	return true
}

func RequireGenerator(t TestingT, c elliptic2.Curve) bool {
	if base := curve.GetBase(c); base != nil {
		_, _, ok := base.Generator()
		if ok {
			return true
		}
	}

	curveStd, ok := c.(elliptic.Curve)
	if ok {
		p := curveStd.Params()
		if p.Gx != nil && p.Gy != nil && p.Gx.Sign() != 0 && p.Gy.Sign() != 0 {
			return true
		}
	}

	t.Logf("skipping test: curve %s has no generator", GetName(c))
	t.SkipNow()
	return false
}
