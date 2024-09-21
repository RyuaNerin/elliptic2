package weierstrassbinary_test

import (
	"testing"

	. "github.com/RyuaNerin/elliptic2/internal/curve/weierstrassbinary"
	. "github.com/RyuaNerin/elliptic2/internal/curvetesting"
)

func TestLambdaConvertPoint(t *testing.T) { TestCoordinates(t, W(NewLambda), c...) }

func TestLambdaAdd(t *testing.T)            { TcOp(t, BuildOp, NewLambda, TAdd, tc...) }
func TestLambdaDouble(t *testing.T)         { TcOp(t, BuildOp, NewLambda, TDouble, tc...) }
func TestLambdaScalarMult(t *testing.T)     { TcOp(t, BuildOp, NewLambda, TMult, tc...) }
func TestLambdaScalarBaseMult(t *testing.T) { TcOp(t, BuildOp, NewLambda, TBaseMult, tc...) }

func BenchmarkLambdaAdd(b *testing.B)            { BOp(b, BuildOp, NewLambda, BAdd, c...) }
func BenchmarkLambdaDouble(b *testing.B)         { BOp(b, BuildOp, NewLambda, BDouble, c...) }
func BenchmarkLambdaScalarMult(b *testing.B)     { BOp(b, BuildOp, NewLambda, BMult, c...) }
func BenchmarkLambdaScalarBaseMult(b *testing.B) { BOp(b, BuildOp, NewLambda, BBaseMult, c...) }
