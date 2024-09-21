package edwards_test

import (
	"testing"

	. "github.com/RyuaNerin/elliptic2/internal/curve"
	. "github.com/RyuaNerin/elliptic2/internal/curve/edwards"
	. "github.com/RyuaNerin/elliptic2/internal/curvetesting"
)

func TestProjectiveConvertPoint(t *testing.T) { TestCoordinates(t, W(NewProjective), c...) }
func TestProjectiveMadd(t *testing.T)         { TestMadd[GFpMaddOperator](t, W(NewProjective), c...) }

func TestProjectiveAdd(t *testing.T)      { TcOp(t, BuildOp, NewProjective, TAdd, tc...) }
func TestProjectiveDouble(t *testing.T)   { TcOp(t, BuildOp, NewProjective, TDouble, tc...) }
func TestProjectiveMult(t *testing.T)     { TcOp(t, BuildOp, NewProjective, TMult, tc...) }
func TestProjectiveBaseMult(t *testing.T) { TcOp(t, BuildOp, NewProjective, TBaseMult, tc...) }

func BenchmarkProjectiveAdd(b *testing.B)      { BOp(b, BuildOp, NewProjective, BAdd, c...) }
func BenchmarkProjectiveDouble(b *testing.B)   { BOp(b, BuildOp, NewProjective, BDouble, c...) }
func BenchmarkProjectiveMult(b *testing.B)     { BOp(b, BuildOp, NewProjective, BMult, c...) }
func BenchmarkProjectiveBaseMult(b *testing.B) { BOp(b, BuildOp, NewProjective, BBaseMult, c...) }
