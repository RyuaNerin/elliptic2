package elliptic2

import (
	"reflect"
)

// IsEquivalentCurve reports whether c1 and c2 have equivalent mathematical
// curve parameters. Curve names and OIDs are intentionally ignored.
func IsEquivalentCurve(c1, c2 Curve) bool {
	if isNilCurve(c1) || isNilCurve(c2) {
		return false
	}
	if isSameCurveInstance(c1, c2) {
		return true
	}
	if ce1, ok := c1.(CurveExtended); ok {
		return ce1.IsEquivalentCurve(c2)
	}
	if ce2, ok := c2.(CurveExtended); ok {
		return ce2.IsEquivalentCurve(c1)
	}
	return c1 == c2
}

// IsEquivalentCurveWithOID reports whether c1 and c2 have equivalent
// mathematical curve parameters and the same object identifier.
func IsEquivalentCurveWithOID(c1, c2 Curve) bool {
	if isNilCurve(c1) || isNilCurve(c2) {
		return false
	}
	if isSameCurveInstance(c1, c2) {
		return true
	}
	if ce1, ok := c1.(CurveExtended); ok {
		return ce1.IsEquivalentCurveWithOID(c2)
	}
	if ce2, ok := c2.(CurveExtended); ok {
		return ce2.IsEquivalentCurveWithOID(c1)
	}
	return c1 == c2
}

func isNilCurve(c Curve) bool {
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

func isSameCurveInstance(c1, c2 Curve) bool {
	v1, v2 := reflect.ValueOf(c1), reflect.ValueOf(c2)
	return v1.IsValid() && v2.IsValid() &&
		v1.Type() == v2.Type() && v1.Comparable() && v1.Equal(v2)
}
