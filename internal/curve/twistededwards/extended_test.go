package twistededwards_test

import (
	"testing"

	. "github.com/RyuaNerin/elliptic2/internal/curve"
	. "github.com/RyuaNerin/elliptic2/internal/curve/twistededwards"
	. "github.com/RyuaNerin/elliptic2/internal/curvetesting"
)

func TestExtendedConvertPoint(t *testing.T) { TestCoordinates(t, W(NewExtended), c...) }
func TestProjectiveMadd(t *testing.T) {
	t.Run("A", func(t *testing.T) { TestMadd[GFpMaddOperator](t, W(NewExtended), c...) })
	t.Run("A=-1", func(t *testing.T) { TestMadd[GFpMaddOperator](t, W(NewExtendedAm1), cAm1) })
}

func TestExtendedAdd(t *testing.T) {
	t.Run("A", func(t *testing.T) { TcOp(t, BuildOp, NewExtended, TAdd, tc...) })
	t.Run("A=-1", func(t *testing.T) { TcOp(t, BuildOp, NewExtendedAm1, TAdd, tcAm1) })
}

func TestExtendedDouble(t *testing.T) {
	t.Run("A", func(t *testing.T) { TcOp(t, BuildOp, NewExtended, TDouble, tc...) })
	t.Run("A=-1", func(t *testing.T) { TcOp(t, BuildOp, NewExtendedAm1, TDouble, tcAm1) })
}

func TestExtendedScalarMult(t *testing.T) {
	t.Run("A", func(t *testing.T) { TcOp(t, BuildOp, NewExtended, TMult, tc...) })
	t.Run("A=-1", func(t *testing.T) { TcOp(t, BuildOp, NewExtendedAm1, TMult, tcAm1) })
}

func TestExtendedScalarBaseMult(t *testing.T) {
	t.Run("A", func(t *testing.T) { TcOp(t, BuildOp, NewExtended, TBaseMult, tc...) })
	t.Run("A=-1", func(t *testing.T) { TcOp(t, BuildOp, NewExtendedAm1, TBaseMult, tcAm1) })
}

func BenchmarkExtendedAdd(b *testing.B) {
	b.Run("A", func(b *testing.B) { BOp(b, BuildOp, NewExtended, BAdd, c...) })
	b.Run("A=-1", func(b *testing.B) { BOp(b, BuildOp, NewExtendedAm1, BAdd, cAm1) })
}

func BenchmarkExtendedDouble(b *testing.B) {
	b.Run("A", func(b *testing.B) { BOp(b, BuildOp, NewExtended, BDouble, c...) })
	b.Run("A=-1", func(b *testing.B) { BOp(b, BuildOp, NewExtendedAm1, BDouble, cAm1) })
}

func BenchmarkExtendedScalarMult(b *testing.B) {
	b.Run("A", func(b *testing.B) { BOp(b, BuildOp, NewExtended, BMult, c...) })
	b.Run("A=-1", func(b *testing.B) { BOp(b, BuildOp, NewExtendedAm1, BMult, cAm1) })
}

func BenchmarkExtendedScalarBaseMult(b *testing.B) {
	b.Run("A", func(b *testing.B) { BOp(b, BuildOp, NewExtended, BBaseMult, c...) })
	b.Run("A=-1", func(b *testing.B) { BOp(b, BuildOp, NewExtendedAm1, BBaseMult, cAm1) })
}
