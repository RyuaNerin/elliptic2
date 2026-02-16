package field

import (
	"math/big"
	"math/bits"

	"github.com/RyuaNerin/elliptic2/internal"
	"github.com/RyuaNerin/elliptic2/internal/field/simd"
)

func NewGF2m() *GF2m              { return new(GF2m) }
func NewGF2mInt64(x int64) *GF2m  { return NewGF2m().SetInt64(x) }
func NewGF2mInt(x *big.Int) *GF2m { return NewGF2m().SetBigInt(x) }
func ParseGF2mHex(s string) *GF2m { return NewGF2m().SetBigInt(internal.HI(s)) }

type GF2m struct {
	nat
	modulus *Modulus
}

//////////////////////////////////////////////////

func (z *GF2m) Modulus() *Modulus { return z.modulus }

func (z *GF2m) SetModulus(p *Modulus) *GF2m {
	z.modulus = p
	return z
}

func (z *GF2m) String() string       { return z.toBigInt(nil).String() }
func (z *GF2m) Text(base int) string { return z.toBigInt(nil).Text(base) }

func (z *GF2m) ToBigInt(x *big.Int) *big.Int { return z.toBigInt(x) }

//////////////////////////////////////////////////

func (z *GF2m) Set(x *GF2m) *GF2m {
	z.set(&x.nat)
	return z
}

func (z *GF2m) SetBits(w []big.Word) *GF2m {
	z.setBits(w)
	return z
}

func (z *GF2m) SetBytes(b []byte) *GF2m {
	z.setBytes(b)
	return z
}

func (z *GF2m) SetUint64(x uint64) *GF2m {
	z.setUint64(x)
	return z
}

func (z *GF2m) SetInt64(x int64) *GF2m {
	z.setUint64(uint64(x))
	return z
}

func (z *GF2m) SetBigInt(x *big.Int) *GF2m {
	z.setBigInt(x)
	return z
}

//////////////////////////////////////////////////

func (z *GF2m) Cmp(x *GF2m) int { return z.cmp(&x.nat) }

func (z *GF2m) Sign() int {
	if z.IsZero() {
		return 0
	}
	return 1
}
func (z *GF2m) IsZero() bool { return z.isZero() }

//////////////////////////////////////////////////

func (z *GF2m) Reduce() *GF2m {
	var words [2 * simd.Words]big.Word
	copy(words[:], z.words[:])
	z.modulus.reduce(z.words[:], words[:])
	return z
}

func (z *GF2m) ReduceSoft() *GF2m {
	return z.Reduce()
}

//////////////////////////////////////////////////

func (z *GF2m) BitLen() int { return z.modulus.bits }

func (z *GF2m) Bit(i int) uint {
	wordIdx := i / simd.WordBitSize
	bitIdx := i % simd.WordBitSize
	return uint((z.words[wordIdx] >> bitIdx) & 1)
}

func (z *GF2m) SetBit(x *GF2m, i int, b uint) *GF2m {
	copy(z.words[:], x.words[:])
	z.modulus = x.modulus

	wordIdx := i / simd.WordBitSize
	bitIdx := i % simd.WordBitSize
	if b == 0 {
		z.words[wordIdx] &^= 1 << bitIdx
	} else {
		z.words[wordIdx] |= 1 << bitIdx
	}
	return z
}

//////////////////////////////////////////////////

func (GF2m) Lsh(*GF2m, uint) *GF2m {
	panic("gf2m: Lsh is not supported")
}

func (GF2m) Rsh(*GF2m, uint) *GF2m {
	panic("gf2m: Rsh is not supported")
}

//////////////////////////////////////////////////

func (z *GF2m) Add(x, y *GF2m) *GF2m {
	for idx := range z.modulus.WordLen() {
		z.words[idx] = x.words[idx] ^ y.words[idx]
	}
	return z
}

func (z *GF2m) Sub(x, y *GF2m) *GF2m {
	return z.Add(x, y)
}

func (z *GF2m) Mul(x, y *GF2m) *GF2m {
	n := z.modulus.WordLen()

	var prod [2 * simd.Words]big.Word

	for a := range n {
		if x.words[a] == 0 {
			continue
		}
		for b := range n {
			if y.words[b] == 0 {
				continue
			}
			lo, hi := simd.CLMUL(x.words[a], y.words[b])

			prod[a+b] ^= lo
			prod[a+b+1] ^= hi
		}
	}

	z.modulus.reduce(z.words[:], prod[:2*n])
	return z
}

func (z *GF2m) Sqr(x *GF2m) *GF2m {
	n := z.modulus.WordLen()

	var prod [2 * simd.Words]big.Word

	for i := range n {
		lo, hi := simd.ExpandBits(x.words[i])
		prod[2*i+0] = lo
		prod[2*i+1] = hi
	}

	z.modulus.reduce(z.words[:], prod[:])
	return z
}

func (z *GF2m) Neg(x *GF2m) *GF2m {
	return z.Set(x)
}

func (z *GF2m) Sqrt(x *GF2m) *GF2m {
	if x.IsZero() {
		return z.SetUint64(1)
	}

	z.Set(x)
	for range z.modulus.bits - 1 {
		z.Sqr(z)
	}

	return z
}

func (z *GF2m) HalfTrace(a *GF2m) *GF2m {
	m := z.modulus.bits

	var tmp GF2m
	tmp.SetModulus(z.modulus)
	tmp.Set(a)

	z.Set(a)
	for range m / 2 {
		tmp.Sqr(&tmp)
		tmp.Sqr(&tmp)
		z.Add(z, &tmp)
	}

	return z
}

func (z *GF2m) SolveQuadraticEven(w, gamma *GF2m) *GF2m {
	m := z.modulus.bits

	sum := z.modulus.NewGF2m()
	inner := z.modulus.NewGF2m()
	gammaPow := z.modulus.NewGF2m()

	z.setZero()
	inner.setZero()
	gammaPow.Set(gamma)

	for i := range m {
		if i > 0 {
			var wPow GF2m
			wPow.SetModulus(z.modulus)
			wPow.Set(w)
			for range i - 1 {
				wPow.Sqr(&wPow)
			}
			inner.Add(inner, &wPow)
		}

		sum.Mul(inner, gammaPow)
		z.Add(z, sum)

		gammaPow.Sqr(gammaPow)
	}

	return z
}

func (z *GF2m) Trace() uint {
	tmp := z.modulus.NewGF2m()
	tmp.Set(z)

	for i := 1; i < z.modulus.bits; i++ {
		tmp.Sqr(tmp)
		tmp.Add(tmp, z)
	}

	// tmp는 0 또는 1
	if tmp.IsZero() {
		return 0
	}
	return 1
}

func (z *GF2m) Inv(x *GF2m) *GF2m {
	if x.IsZero() {
		panic("gf2m: inverse of zero")
	}

	n := z.modulus.WordLen()

	var uV, vV, g1V, g2V [1 + simd.Words]big.Word
	u, v, g1, g2 := uV[:], vV[:], g1V[:], g2V[:]

	copy(u, x.words[:n])         // u = x
	copy(v, z.modulus.words[:n]) // v = mod_poly
	g1V[0] = 1                   // g1 = 1
	// g2 = 0

	workLen := n + 1

	// Extended GCD
	for !isZeroWords(u[:workLen]) {
		degU := bitLen(u[:workLen]) - 1
		degV := bitLen(v[:workLen]) - 1

		if degU < degV {
			u, v = v, u
			g1, g2 = g2, g1
			degU, degV = degV, degU
		}

		if degV < 0 {
			break
		}

		shift := degU - degV

		// u ^= v << shift
		lshXor(u[:workLen], v[:workLen], uint(shift))
		// g1 ^= g2 << shift
		lshXor(g1[:workLen], g2[:workLen], uint(shift))
	}

	// Result is in g2 (since v should be 1)
	copy(z.words[:], g2[:n])
	return z
}

func bitLen(words []big.Word) int {
	for i := len(words) - 1; i >= 0; i-- {
		if words[i] != 0 {
			return i*simd.WordBitSize + bits.Len(uint(words[i]))
		}
	}
	return 0
}

func isZeroWords(words []big.Word) bool {
	for _, w := range words {
		if w != 0 {
			return false
		}
	}
	return true
}

// lshXor computes dst ^= src << shift
func lshXor(dst, src []big.Word, shift uint) {
	if shift == 0 {
		for i := range dst {
			dst[i] ^= src[i]
		}
		return
	}

	wordShift := int(shift / simd.WordBitSize)
	bitShift := shift % simd.WordBitSize

	// Clear tmp
	var tmp [1 + simd.Words]big.Word

	// Shift src into tmp
	for i := range src {
		if src[i] == 0 {
			continue
		}
		targetWord := i + wordShift
		if targetWord < len(tmp) {
			tmp[targetWord] ^= src[i] << bitShift
		}
		if bitShift > 0 && targetWord+1 < len(tmp) {
			tmp[targetWord+1] ^= src[i] >> (simd.WordBitSize - bitShift)
		}
	}

	// XOR into dst
	for i := range dst {
		dst[i] ^= tmp[i]
	}
}
