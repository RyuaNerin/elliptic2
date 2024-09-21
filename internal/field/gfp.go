package field

import (
	"math/big"

	"github.com/RyuaNerin/elliptic2/internal"
	"github.com/RyuaNerin/elliptic2/internal/field/simd"
)

var (
	GFpZero  = NewGFpInt64(0)
	GFpOne   = NewGFpInt64(1)
	GFpTwo   = NewGFpInt64(2)
	GFpThree = NewGFpInt64(3)
)

func NewGFp() *GFp              { return new(GFp) }
func NewGFpInt64(x int64) *GFp  { return NewGFp().SetInt64(x) }
func NewGFpInt(x *big.Int) *GFp { return NewGFp().SetBigInt(x) }
func ParseGFpHex(s string) *GFp { return NewGFp().SetBigInt(internal.HI(s)) }

type GFp struct {
	value      big.Int
	modulus    *Modulus
	modulesInt big.Int
}

func (z *GFp) Grow() *GFp {
	var wlen int
	if z == nil {
		wlen = 4 + simd.Words
	} else {
		wlen = 4 + z.modulus.WordLen()
	}
	z.value.SetBits(make([]big.Word, wlen))
	return z
}

//////////////////////////////////////////////////

func (z *GFp) Modulus() *Modulus { return z.modulus }

func (z *GFp) SetModulus(p *Modulus) *GFp {
	z.modulus = p
	z.modulus.toBigInt(&z.modulesInt)
	if z.BitLen() == 0 && cap(z.value.Bits()) < p.WordLen()+4 {
		z.value.SetBits(make([]big.Word, p.WordLen()+4))
	}
	return z
}

func (z *GFp) String() string       { return z.value.String() }
func (z *GFp) Text(base int) string { return z.value.Text(base) }

func (z *GFp) ToBigInt(x *big.Int) *big.Int {
	if x == nil {
		x = new(big.Int)
	}
	x.Set(&z.value)
	return x
}

//////////////////////////////////////////////////

func (z *GFp) Set(x *GFp) *GFp {
	z.value.Set(&x.value)
	return z
}

func (z *GFp) SetBits(words []big.Word) *GFp {
	z.value.SetBits(words)
	return z
}

func (z *GFp) SetBytes(b []byte) *GFp {
	z.value.SetBytes(b)
	return z
}

func (z *GFp) SetUint64(x uint64) *GFp {
	z.value.SetUint64(x)
	return z
}

func (z *GFp) SetInt64(x int64) *GFp {
	z.value.SetInt64(x)
	return z
}

func (z *GFp) SetBigInt(x *big.Int) *GFp {
	z.value.Set(x)
	return z
}

//////////////////////////////////////////////////

func (z *GFp) Cmp(x *GFp) int { return z.value.Cmp(&x.value) }
func (z *GFp) Sign() int      { return z.value.Sign() }
func (z *GFp) IsZero() bool   { return z.value.Sign() == 0 }

//////////////////////////////////////////////////

func (z *GFp) Reduce() *GFp {
	z.value.Mod(&z.value, &z.modulesInt)
	return z
}

func (z *GFp) ReduceSoft() *GFp {
	n := z.value.BitLen()
	if n > z.modulus.BitLen()*2 {
		z.Reduce()
	}
	return z
}

//////////////////////////////////////////////////

func (z *GFp) BitLen() int    { return z.value.BitLen() }
func (z *GFp) Bit(i int) uint { return z.value.Bit(i) }

func (z *GFp) SetBit(x *GFp, i int, b uint) *GFp {
	z.value.SetBit(&x.value, i, b)
	return z
}

//////////////////////////////////////////////////

func (z *GFp) Lsh(x *GFp, n uint) *GFp {
	z.value.Lsh(&x.value, n)
	return z
}

func (z *GFp) Rsh(x *GFp, n uint) *GFp {
	z.value.Rsh(&x.value, n)
	return z
}

//////////////////////////////////////////////////

func (z *GFp) Add(x, y *GFp) *GFp {
	z.value.Add(&x.value, &y.value)
	return z
}

func (z *GFp) Sub(x, y *GFp) *GFp {
	z.value.Sub(&x.value, &y.value)
	return z
}

func (z *GFp) Mul(x, y *GFp) *GFp {
	z.value.Mul(&x.value, &y.value)
	return z
}

func (z *GFp) Sqr(x *GFp) *GFp {
	z.value.Mul(&x.value, &x.value)
	return z
}

func (z *GFp) Neg(x *GFp) *GFp {
	z.value.Neg(&x.value)
	return z
}

func (z *GFp) Inv(x *GFp) *GFp {
	r := z.value.ModInverse(&x.value, &z.modulesInt)
	if r == nil {
		return nil
	}

	return z
}

func (z *GFp) Sqrt(x *GFp) *GFp {
	if z.value.ModSqrt(&x.value, &z.modulesInt) == nil {
		return nil
	}
	return z
}
