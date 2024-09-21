package field_test

import (
	"testing"

	"github.com/RyuaNerin/elliptic2/internal/field"
)

func TestFx(t *testing.T) {
	want := []int{283, 12, 7, 5, 0}

	b := field.GF2mPolynomials(want...)
	got := field.ParseGF2mPolynomials(b)

	if len(want) != len(got) {
		t.Fatalf("length mismatch:\ngot:  %d\n want: %d", len(got), len(want))
		return
	}
	for i := range want {
		if want[i] != got[i] {
			t.Fatalf("mismatch at index %d:\ngot:  %d\nwant:  %d", i, got[i], want[i])
			return
		}
	}
}
