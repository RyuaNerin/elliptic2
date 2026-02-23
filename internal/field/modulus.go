package field

import (
	"fmt"
	"math/big"
	"sync/atomic"
	"unsafe"

	"github.com/RyuaNerin/elliptic2/internal"
	"github.com/RyuaNerin/elliptic2/internal/field/gf2mreduce"
	"github.com/RyuaNerin/elliptic2/internal/field/simd"
)

// GF2mPolynomials to big.Int
// f = [233, 74, 0] when x^233 +  x^74 + 1
// f = [283, 12, 7, 5, 0] when x^283 +  x^12 +  x^7 +  x^5 + 1
func GF2mPolynomials(f ...int) *big.Int {
	mavP := 0
	for _, v := range f {
		if mavP < v {
			mavP = v
		}
	}

	ret := new(big.Int)
	ret.SetBits(make([]big.Word, (mavP+simd.WordBitSize-1)/simd.WordBitSize))
	for _, v := range f {
		ret.SetBit(ret, v, 1)
	}
	return ret
}

func ParseGF2mPolynomials(f *big.Int) []int {
	var ret []int
	for i := f.BitLen() - 1; i >= 0; i-- {
		if f.Bit(i) != 0 {
			ret = append(ret, i)
		}
	}
	return ret
}

type Modulus struct {
	nat
	bits   int
	reduce gf2mreduce.ReduceFunc

	gf2mPoly  []int
	gf2mGamma unsafe.Pointer // *GF2m
	gf2mSqrtX *GF2m          // *GF2m, sqrt(x)
}

func NewGFpModulus(p *big.Int) *Modulus {
	bits := p.BitLen()
	if bits > simd.MaxBitSize {
		panic(fmt.Sprintf("modulus bit size too large. got=%d max=%d", bits, simd.MaxBitSize))
	}

	m := &Modulus{
		bits: bits,
	}
	copy(m.words[:], p.Bits())
	return m
}

func NewGFpModulusFromHex(hexStr string) *Modulus {
	return NewGFpModulus(internal.HI(hexStr))
}

func NewGF2mModulus(modPoly *big.Int) *Modulus {
	poly := ParseGF2mPolynomials(modPoly)
	if poly[0] >= simd.MaxBitSize {
		panic(fmt.Sprintf("modulus bit size too large. got=%d max=%d", poly[0], simd.MaxBitSize))
	}

	p := &Modulus{
		bits:     poly[0],
		reduce:   gf2mreduce.GetReduceFunction(poly),
		gf2mPoly: make([]int, len(poly)),
	}
	copy(p.gf2mPoly, poly)
	p.precomputeSqrt()

	copy(p.words[:], modPoly.Bits())

	return p
}

func NewGF2mModulusFromPolynomials(f ...int) *Modulus {
	return NewGF2mModulus(GF2mPolynomials(f...))
}

func (p *Modulus) String() string {
	return fmt.Sprintf("Modulus{bits:%d,poly:%v}", p.bits, p.gf2mPoly)
}

func (p *Modulus) Bit(n uint) uint {
	wordIdx := n / simd.WordBitSize
	bitIdx := n % simd.WordBitSize
	return uint((p.words[wordIdx] >> bitIdx) & 1)
}

func (p *Modulus) Poly() []int  { return p.gf2mPoly }
func (p *Modulus) BitLen() int  { return p.bits }
func (p *Modulus) Degree() int  { return p.bits }
func (p *Modulus) WordLen() int { return (p.bits + simd.WordBitSize - 1) / simd.WordBitSize }

func (p *Modulus) NewGF2m() *GF2m                     { return new(GF2m).SetModulus(p) }
func (p *Modulus) NewGF2mFromBigInt(x *big.Int) *GF2m { return p.NewGF2m().SetBigInt(x) }
func (p *Modulus) NewGFp() *GFp                       { return new(GFp).SetModulus(p) }
func (p *Modulus) NewGFpFromBigInt(x *big.Int) *GFp   { return p.NewGFp().SetBigInt(x) }

// ToBigInt converts to big.Int
func (z *Modulus) ToBigInt(x *big.Int) *big.Int { return z.toBigInt(x) }

func (p *Modulus) precomputeSqrt() {
	// sqrtX = x^{2^{m-1}} mod P
	p.gf2mSqrtX = p.NewGF2m()
	p.gf2mSqrtX.SetUint64(2)

	// Sqrt(x) = x^{2^(m-1)}
	for range p.bits - 1 {
		p.gf2mSqrtX.Sqr(p.gf2mSqrtX)
	}
}

func (poly *Modulus) FindGF2mGamma() *GF2m {
	gamma := atomic.LoadPointer(&poly.gf2mGamma)
	if gamma != nil {
		return (*GF2m)(gamma)
	}

	gammaVal := poly.findGamma()
	if atomic.CompareAndSwapPointer(&poly.gf2mGamma, nil, unsafe.Pointer(gammaVal)) {
		return gammaVal
	}

	return (*GF2m)(atomic.LoadPointer(&poly.gf2mGamma))
}

func (poly *Modulus) findGamma() *GF2m {
	var gamma GF2m
	gamma.SetModulus(poly)

	shifted := big.NewInt(1)
	for range poly.bits {
		gamma.SetBigInt(shifted)
		if gamma.Trace() == 1 {
			result := new(GF2m)
			result.SetModulus(poly)
			result.Set(&gamma)
			return result
		}
		shifted.Lsh(shifted, 1) // shifted <<= 1
	}

	panic("no gamma found")
}
