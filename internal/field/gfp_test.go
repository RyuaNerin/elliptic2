package field_test

import (
	"testing"

	. "github.com/RyuaNerin/elliptic2/internal"
	. "github.com/RyuaNerin/elliptic2/internal/field"
)

func gfp(m *GFp) {}

func TestGFpInvValidation(t *testing.T) { testInvValidation(t, gfp, gfpModulus) }
func TestGFpSqrValidation(t *testing.T) { testSqrValidation(t, gfp, gfpModulus) }

func BenchmarkGFpAdd(b *testing.B) { bArg(b, gfp, gfpModulus, add) }
func BenchmarkGFpMul(b *testing.B) { bArg(b, gfp, gfpModulus, mul) }
func BenchmarkGFpSqr(b *testing.B) { bArg(b, gfp, gfpModulus, sqr) }
func BenchmarkGFpInv(b *testing.B) { bArg(b, gfp, gfpModulus, inv) }

var gfpModulus = []*Modulus{
	NewGFpModulus(HI(`fffffffffffffffffffffffffffffffeffffffffffffffff`)),
	NewGFpModulus(HI(`ffffffffffffffffffffffffffffffff000000000000000000000001`)),
	NewGFpModulus(HI(`ffffffff00000001000000000000000000000000ffffffffffffffffffffffff`)),
	NewGFpModulus(HI(`fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffeffffffff0000000000000000ffffffff`)),
	NewGFpModulus(HI(`1ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff`)),
}
