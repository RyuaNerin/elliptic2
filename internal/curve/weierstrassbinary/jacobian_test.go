package weierstrassbinary_test

import (
	"testing"

	. "github.com/RyuaNerin/elliptic2/internal/curve"
	. "github.com/RyuaNerin/elliptic2/internal/curve/weierstrassbinary"
	. "github.com/RyuaNerin/elliptic2/internal/curvetesting"
)

func TestJacobianConvertPoint(t *testing.T) { TestCoordinates(t, W(NewJacobian), c...) }
func TestProjectiveMadd(t *testing.T)       { TestMadd[GF2mMaddOperator](t, W(NewJacobian), c...) }

func TestJacobianAdd(t *testing.T)            { TcOp(t, BuildOp, NewJacobian, TAdd, tc...) }
func TestJacobianDouble(t *testing.T)         { TcOp(t, BuildOp, NewJacobian, TDouble, tc...) }
func TestJacobianScalarMult(t *testing.T)     { TcOp(t, BuildOp, NewJacobian, TMult, tc...) }
func TestJacobianScalarBaseMult(t *testing.T) { TcOp(t, BuildOp, NewJacobian, TBaseMult, tc...) }

func BenchmarkJacobianAdd(b *testing.B)            { BOp(b, BuildOp, NewJacobian, BAdd, c...) }
func BenchmarkJacobianDouble(b *testing.B)         { BOp(b, BuildOp, NewJacobian, BDouble, c...) }
func BenchmarkJacobianScalarMult(b *testing.B)     { BOp(b, BuildOp, NewJacobian, BMult, c...) }
func BenchmarkJacobianScalarBaseMult(b *testing.B) { BOp(b, BuildOp, NewJacobian, BBaseMult, c...) }
