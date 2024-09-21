package edwards_test

import (
	"testing"

	. "github.com/RyuaNerin/elliptic2/internal/curve/edwards"
	. "github.com/RyuaNerin/elliptic2/internal/curvetesting"
)

func TestAffineAdd(t *testing.T)            { TcOp(t, BuildOp, NewAffine, TAdd, tc...) }
func TestAffineDouble(t *testing.T)         { TcOp(t, BuildOp, NewAffine, TDouble, tc...) }
func TestAffineScalarMult(t *testing.T)     { TcOp(t, BuildOp, NewAffine, TMult, tc...) }
func TestAffineScalarBaseMult(t *testing.T) { TcOp(t, BuildOp, NewAffine, TBaseMult, tc...) }

func BenchmarkAffineAdd(b *testing.B)            { BOp(b, BuildOp, NewAffine, BAdd, c...) }
func BenchmarkAffineDouble(b *testing.B)         { BOp(b, BuildOp, NewAffine, BDouble, c...) }
func BenchmarkAffineScalarMult(b *testing.B)     { BOp(b, BuildOp, NewAffine, BMult, c...) }
func BenchmarkAffineScalarBaseMult(b *testing.B) { BOp(b, BuildOp, NewAffine, BBaseMult, c...) }
