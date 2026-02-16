package curvetesting

import (
	"crypto/elliptic"
	"fmt"
	"math/big"
	"testing"

	"github.com/RyuaNerin/elliptic2"
	"github.com/RyuaNerin/elliptic2/internal/curve"
	"github.com/stretchr/testify/require"
)

func RequireGenerator(t testing.TB, c elliptic2.Curve) bool {
	t.Helper()

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

func RequireIsOnCurve(t testing.TB, curve elliptic2.Curve, x, y *big.Int, msg string, args ...any) {
	t.Helper()
	prefix := fmt.Sprintf(msg, args...)
	require.True(t,
		curve.IsOnCurve(x, y),
		"%s is not on the curve:\n  X: %s\n  Y: %s",
		prefix, x.String(), y.String(),
	)
}

func RequireNotIsOnCurve(t testing.TB, curve elliptic2.Curve, x, y *big.Int, msg string, args ...any) {
	t.Helper()
	prefix := fmt.Sprintf(msg, args...)
	require.False(t,
		curve.IsOnCurve(x, y),
		"%s is on the curve:\n  X: %s\n  Y: %s",
		prefix, x.String(), y.String(),
	)
}

func RequireEqual(t testing.TB, want, got *big.Int, msg string, args ...any) {
	t.Helper()
	prefix := fmt.Sprintf(msg, args...)
	require.Zero(t,
		want.Cmp(got),
		"%s mismatch:\n  want: %s\n  got:  %s",
		prefix, want.String(), got.String(),
	)
}

func RequireXYEquals(t testing.TB, want, got *Point, msg string, args ...any) {
	t.Helper()
	prefix := fmt.Sprintf(msg, args...)
	require.Zero(t,
		want.X.Cmp(got.X),
		"%s.X mismatch:\n  want: %s\n  got:  %s",
		prefix, want.X.String(), got.X.String(),
	)
	require.Zero(t,
		want.Y.Cmp(got.Y),
		"%s.Y mismatch:\n  want: %s\n  got:  %s",
		prefix, want.Y.String(), got.Y.String(),
	)
}

func RequireUnmodified(t testing.TB, saved, current *big.Int, msg string, args ...any) {
	t.Helper()
	prefix := fmt.Sprintf(msg, args...)
	require.Zero(t, saved.Cmp(current),
		"%s modified:\n  want: %s\n  got:  %s",
		prefix, saved.String(), current.String(),
	)
}

func RequireXYUnmodified(t testing.TB, saved, current *Point, msg string, args ...any) {
	t.Helper()
	prefix := fmt.Sprintf(msg, args...)
	require.Zero(t, saved.X.Cmp(current.X),
		"%s.X modified:\n  want: %s\n  got:  %s",
		prefix, saved.X.String(), current.X.String(),
	)
	require.Zero(t, saved.Y.Cmp(current.Y),
		"%s.Y modified:\n  want: %s\n  got:  %s",
		prefix, saved.Y.String(), current.Y.String(),
	)
}
