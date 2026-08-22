package curve

import (
	"crypto/elliptic"
	"math/big"
	"reflect"
	"slices"

	"github.com/RyuaNerin/elliptic2"
)

func isEquivalentParams(p1, p2 *elliptic.CurveParams) bool {
	eq := func(a, b *big.Int) bool {
		if a == nil && b == nil {
			return true
		}
		if a == nil || b == nil {
			return false
		}
		return a.Cmp(b) == 0
	}

	switch {
	case !eq(p1.P, p2.P):
	case !eq(p1.N, p2.N):
	case !eq(p1.B, p2.B):
	case !eq(p1.Gx, p2.Gx):
	case !eq(p1.Gy, p2.Gy):
	case p1.BitSize != p2.BitSize:
	default:
		return true
	}

	return false
}

func isEquivalentParams2(p1, p2 *elliptic2.CurveParams) bool {
	eq := func(a, b *big.Int) bool {
		if a == nil && b == nil {
			return true
		}
		if (a == nil) != (b == nil) {
			return false
		}
		return a.Cmp(b) == 0
	}

	switch {
	case p1.Type != p2.Type:
	case p1.BitSize != p2.BitSize:
	case !eq(p1.A, p2.A):
	case !eq(p1.B, p2.B):
	case !eq(p1.C, p2.C):
	case !eq(p1.D, p2.D):
	case !eq(p1.A2, p2.A2):
	case !eq(p1.A6, p2.A6):
	case !eq(p1.P, p2.P):
	case slices.Compare(p1.Poly, p2.Poly) != 0:
	case !eq(p1.N, p2.N):
	case !eq(p1.Gx, p2.Gx):
	case !eq(p1.Gy, p2.Gy):
	case !eq(p1.InfX, p2.InfX):
	case !eq(p1.InfY, p2.InfY):
	default:
		return true
	}

	return false
}

func isNilCurve(c elliptic2.Curve) bool {
	if c == nil {
		return true
	}
	v := reflect.ValueOf(c)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

func isSameCurveInstance(c1, c2 elliptic2.Curve) bool {
	v1, v2 := reflect.ValueOf(c1), reflect.ValueOf(c2)
	return v1.IsValid() && v2.IsValid() &&
		v1.Type() == v2.Type() && v1.Comparable() && v1.Equal(v2)
}
