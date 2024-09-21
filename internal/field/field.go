package field

import (
	"math/big"
)

type GFs interface {
	GFp | GF2m
}

type GF[Self any] interface {
	*Self

	//////////////////////////////////////////////////

	Modulus() *Modulus
	SetModulus(p *Modulus) *Self

	String() string
	Text(base int) string

	ToBigInt(x *big.Int) *big.Int

	//////////////////////////////////////////////////

	Set(x *Self) *Self
	SetBits(words []big.Word) *Self
	SetBytes(b []byte) *Self
	SetUint64(x uint64) *Self
	SetInt64(x int64) *Self
	SetBigInt(x *big.Int) *Self

	//////////////////////////////////////////////////

	Cmp(x *Self) int
	Sign() int
	IsZero() bool

	//////////////////////////////////////////////////

	Reduce() *Self
	ReduceSoft() *Self

	//////////////////////////////////////////////////

	BitLen() int
	Bit(i int) uint
	SetBit(x *Self, i int, b uint) *Self

	//////////////////////////////////////////////////

	Lsh(x *Self, n uint) *Self
	Rsh(x *Self, n uint) *Self

	//////////////////////////////////////////////////

	Add(x, y *Self) *Self
	Sub(x, y *Self) *Self
	Mul(x, y *Self) *Self
	Sqr(x *Self) *Self
	Neg(x *Self) *Self
	Inv(x *Self) *Self
	Sqrt(x *Self) *Self
}

func typecheck[Self any, P GF[Self]]() P { return nil }

var (
	_ = typecheck[GFp]()
	_ = typecheck[GF2m]()
)
