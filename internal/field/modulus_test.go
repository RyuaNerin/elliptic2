package field_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/RyuaNerin/elliptic2/internal/field"
)

func TestFx(t *testing.T) {
	want := []int{283, 12, 7, 5, 0}

	b := field.GF2mPolynomials(want...)
	got := field.ParseGF2mPolynomials(b)

	require.Len(t, got, len(want), "length mismatch")
	require.Equal(t, want, got, "parsed polynomials do not match expected values")
}
