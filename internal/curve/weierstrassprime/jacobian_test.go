package weierstrassprime_test

import (
	"testing"

	. "github.com/RyuaNerin/elliptic2/internal/curve"
	. "github.com/RyuaNerin/elliptic2/internal/curve/weierstrassprime"
	. "github.com/RyuaNerin/elliptic2/internal/curvetesting"
)

func TestJacobianConvertPoint(t *testing.T) { TestCoordinates(t, W(NewJacobian), c...) }
func TestJacobianMadd(t *testing.T) {
	t.Run("A", func(t *testing.T) { TestMadd[GFpMaddOperator](t, W(NewJacobian), c...) })
	t.Run("A=0", func(t *testing.T) { TestMadd[GFpMaddOperator](t, W(NewJacobianA0), cA0) })
	t.Run("A=-3", func(t *testing.T) { TestMadd[GFpMaddOperator](t, W(NewJacobianAm3), cAm3) })
}

func TestJacobianAdd(t *testing.T) {
	t.Run("A", func(t *testing.T) { TcOp(t, BuildOp, NewJacobian, TAdd, tc...) })
	t.Run("A=0", func(t *testing.T) { TcOp(t, BuildOp, NewJacobianA0, TAdd, tcA0) })
	t.Run("A=-3", func(t *testing.T) { TcOp(t, BuildOp, NewJacobianAm3, TAdd, tcAm3) })
}

func TestJacobianDouble(t *testing.T) {
	t.Run("A", func(t *testing.T) { TcOp(t, BuildOp, NewJacobian, TDouble, tc...) })
	t.Run("A=0", func(t *testing.T) { TcOp(t, BuildOp, NewJacobianA0, TDouble, tcA0) })
	t.Run("A=-3", func(t *testing.T) { TcOp(t, BuildOp, NewJacobianAm3, TDouble, tcAm3) })
}

func TestJacobianScalarMult(t *testing.T) {
	t.Run("A", func(t *testing.T) { TcOp(t, BuildOp, NewJacobian, TMult, tc...) })
	t.Run("A=0", func(t *testing.T) { TcOp(t, BuildOp, NewJacobianA0, TMult, tcA0) })
	t.Run("A=-3", func(t *testing.T) { TcOp(t, BuildOp, NewJacobianAm3, TMult, tcAm3) })
}

func TestJacobianScalarBaseMult(t *testing.T) {
	t.Run("A", func(t *testing.T) { TcOp(t, BuildOp, NewJacobian, TBaseMult, tc...) })
	t.Run("A=0", func(t *testing.T) { TcOp(t, BuildOp, NewJacobianA0, TBaseMult, tcA0) })
	t.Run("A=-3", func(t *testing.T) { TcOp(t, BuildOp, NewJacobianAm3, TBaseMult, tcAm3) })
}

func BenchmarkJacobianAdd(b *testing.B) {
	b.Run("A", func(b *testing.B) { BOp(b, BuildOp, NewJacobian, BAdd, c...) })
	b.Run("A=0", func(b *testing.B) { BOp(b, BuildOp, NewJacobianA0, BAdd, cA0) })
	b.Run("A=-3", func(b *testing.B) { BOp(b, BuildOp, NewJacobianAm3, BAdd, cAm3) })
}

func BenchmarkJacobianDouble(b *testing.B) {
	b.Run("A", func(b *testing.B) { BOp(b, BuildOp, NewJacobian, BDouble, c...) })
	b.Run("A=0", func(b *testing.B) { BOp(b, BuildOp, NewJacobianA0, BDouble, cA0) })
	b.Run("A=-3", func(b *testing.B) { BOp(b, BuildOp, NewJacobianAm3, BDouble, cAm3) })
}

func BenchmarkJacobianScalarMult(b *testing.B) {
	b.Run("A", func(b *testing.B) { BOp(b, BuildOp, NewJacobian, BMult, c...) })
	b.Run("A=0", func(b *testing.B) { BOp(b, BuildOp, NewJacobianA0, BMult, cA0) })
	b.Run("A=-3", func(b *testing.B) { BOp(b, BuildOp, NewJacobianAm3, BMult, cAm3) })
}

func BenchmarkJacobianScalarBaseMult(b *testing.B) {
	b.Run("A", func(b *testing.B) { BOp(b, BuildOp, NewJacobian, BBaseMult, c...) })
	b.Run("A=0", func(b *testing.B) { BOp(b, BuildOp, NewJacobianA0, BBaseMult, cA0) })
	b.Run("A=-3", func(b *testing.B) { BOp(b, BuildOp, NewJacobianAm3, BBaseMult, cAm3) })
}
