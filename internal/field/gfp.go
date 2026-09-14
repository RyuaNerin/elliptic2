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
	nat
	modulus *Modulus
	neg     bool
}

//////////////////////////////////////////////////

func (z *GFp) Modulus() *Modulus { return z.modulus }

func (z *GFp) SetModulus(p *Modulus) *GFp {
	z.modulus = p
	return z.Reduce()
}

func (z *GFp) String() string       { return z.ToBigInt(nil).String() }
func (z *GFp) Text(base int) string { return z.ToBigInt(nil).Text(base) }

func (z *GFp) ToBigInt(x *big.Int) *big.Int {
	r := z.toBigInt(x)
	if z.neg {
		r.Neg(r)
	}
	return r
}

func (z *GFp) wordLen() int {
	if z.modulus == nil {
		return len(z.words)
	}
	return z.modulus.gfpWordLen
}

//////////////////////////////////////////////////

func (z *GFp) Set(x *GFp) *GFp {
	z.set(&x.nat)
	z.neg = x.neg
	if z.neg && z.modulus != nil {
		z.Reduce()
	}
	return z
}

func (z *GFp) SetBits(w []big.Word) *GFp {
	z.setBits(w)
	z.neg = false
	return z.Reduce()
}

func (z *GFp) SetBytes(b []byte) *GFp {
	if len(b) > len(z.words)*simd.WordByteSize {
		return z.SetBigInt(new(big.Int).SetBytes(b))
	}
	z.setBytes(b)
	z.neg = false
	return z.Reduce()
}

func (z *GFp) SetUint64(x uint64) *GFp {
	z.setUint64(x)
	z.neg = false
	return z.Reduce()
}

func (z *GFp) SetInt64(x int64) *GFp {
	z.setUint64(uint64(max(x, -x)))
	z.neg = x < 0
	return z.Reduce()
}

func (z *GFp) SetBigInt(x *big.Int) *GFp {
	if z.modulus == nil {
		z.setBigInt(x)
		z.neg = x.Sign() < 0
		return z
	}

	z.neg = false

	if len(x.Bits()) > len(z.words) {
		var t big.Int
		t.Mod(x, z.modulus.gfpBig)
		z.setBigInt(&t)
		return z
	}

	z.setBigInt(x) // 절댓값
	z.neg = x.Sign() < 0
	return z.Reduce()
}

//////////////////////////////////////////////////

func (z *GFp) Cmp(x *GFp) int {
	if z.neg != x.neg {
		if z.neg {
			return -1
		}
		return 1
	}
	if c := z.cmp(&x.nat); c != 0 {
		if z.neg {
			return -c
		}
		return c
	}
	return 0
}

func (z *GFp) Sign() int {
	if z.isZero() {
		return 0
	}
	if z.neg {
		return -1
	}
	return 1
}

func (z *GFp) IsZero() bool { return z.isZero() }

func (z *GFp) isOne() bool {
	if z.words[0] != 1 {
		return false
	}
	for _, w := range z.words[1:] {
		if w != 0 {
			return false
		}
	}
	return true
}

//////////////////////////////////////////////////

func (z *GFp) Reduce() *GFp {
	if z.modulus == nil {
		return z
	}

	n := z.modulus.gfpWordLen
	z.modulus.mod(z.words[:n], z.words[:])
	clear(z.words[n:])

	if z.neg {
		z.neg = false
		if !z.isZero() {
			simd.SubVV(z.words[:n], z.modulus.words[:n], z.words[:n])
		}
	}
	return z
}

func (z *GFp) ReduceSoft() *GFp { return z }

//////////////////////////////////////////////////

func (z *GFp) BitLen() int { return bitLen(z.words[:]) }

func (z *GFp) Bit(i int) uint {
	if i < 0 || i >= len(z.words)*simd.WordBitSize {
		return 0
	}
	return uint((z.words[i/simd.WordBitSize] >> (i % simd.WordBitSize)) & 1)
}

func (z *GFp) SetBit(x *GFp, i int, b uint) *GFp {
	if i < 0 || i >= len(z.words)*simd.WordBitSize {
		panic("elliptic2: bit index out of range")
	}

	z.Set(x)
	z.neg = false

	wordIdx := i / simd.WordBitSize
	bitIdx := i % simd.WordBitSize
	if b == 0 {
		z.words[wordIdx] &^= 1 << bitIdx
	} else {
		z.words[wordIdx] |= 1 << bitIdx
	}
	return z.Reduce()
}

//////////////////////////////////////////////////

func (z *GFp) Lsh(x *GFp, s uint) *GFp {
	n := z.wordLen()
	z.set(&x.nat)

	if z.modulus == nil {
		panic("gfp: modulus is not set")
	}

	var t [1 + simd.Words]big.Word
	for s > 0 {
		k := s
		if k >= simd.WordBitSize {
			k = simd.WordBitSize - 1
		}

		t[n] = simd.ShlVU(t[:n], z.words[:n], k)
		z.modulus.mod(z.words[:n], t[:n+1])
		clear(z.words[n:])

		s -= k
	}
	return z
}

func (z *GFp) Rsh(x *GFp, s uint) *GFp {
	z.set(&x.nat)

	for s > 0 {
		k := s
		if k >= simd.WordBitSize {
			k = simd.WordBitSize - 1
		}

		simd.ShrVU(z.words[:], z.words[:], k)
		s -= k
	}
	return z
}

//////////////////////////////////////////////////

func (z *GFp) Add(x, y *GFp) *GFp {
	p := z.modulus
	n := p.gfpWordLen

	c := simd.AddVV(z.words[:n], x.words[:n], y.words[:n])
	if c != 0 || cmpWords(z.words[:n], p.words[:n]) >= 0 {
		simd.SubVV(z.words[:n], z.words[:n], p.words[:n])
	}
	clear(z.words[n:])
	return z
}

func (z *GFp) Sub(x, y *GFp) *GFp {
	p := z.modulus
	n := p.gfpWordLen

	if simd.SubVV(z.words[:n], x.words[:n], y.words[:n]) != 0 {
		simd.AddVV(z.words[:n], z.words[:n], p.words[:n])
	}
	clear(z.words[n:])
	return z
}

func (z *GFp) Mul(x, y *GFp) *GFp {
	p := z.modulus
	n := p.gfpWordLen

	var prod [2 * simd.Words]big.Word
	basicMul(prod[:2*n], x.words[:n], y.words[:n])

	p.mod(z.words[:n], prod[:2*n])
	clear(z.words[n:])
	return z
}

func (z *GFp) Sqr(x *GFp) *GFp {
	p := z.modulus
	n := p.gfpWordLen

	var prod [2 * simd.Words]big.Word
	basicMul(prod[:2*n], x.words[:n], x.words[:n])

	p.mod(z.words[:n], prod[:2*n])
	clear(z.words[n:])
	return z
}

func (z *GFp) Neg(x *GFp) *GFp {
	p := z.modulus
	n := p.gfpWordLen

	if x.isZero() {
		z.setZero()
		return z
	}

	simd.SubVV(z.words[:n], p.words[:n], x.words[:n])
	clear(z.words[n:])
	return z
}

func (z *GFp) Inv(x *GFp) *GFp {
	p := z.modulus
	n := p.gfpWordLen

	if x.isZero() {
		return nil
	}
	if p.words[0]&1 == 0 {
		panic("gfp: even modulus is not supported")
	}

	var uArr, vArr, aArr, cArr [1 + simd.Words]big.Word
	u, v, a, c := uArr[:n], vArr[:n], aArr[:n], cArr[:n]

	copy(u, x.words[:n]) // u = x
	copy(v, p.words[:n]) // v = p
	a[0] = 1             // a = 1 : a * x = u (mod p)
	//                      c = 0 : c * x = v (mod p)

	for !isZeroWords(u) {
		for u[0]&1 == 0 {
			simd.ShrVU(u, u, 1)
			halveMod(a, p.words[:n])
		}
		for v[0]&1 == 0 {
			simd.ShrVU(v, v, 1)
			halveMod(c, p.words[:n])
		}

		if cmpWords(u, v) >= 0 {
			simd.SubVV(u, u, v)
			subMod(a, a, c, p.words[:n])
		} else {
			simd.SubVV(v, v, u)
			subMod(c, c, a, p.words[:n])
		}
	}

	if !isOneWords(v) { // gcd(x, p) != 1
		return nil
	}

	copy(z.words[:n], c)
	clear(z.words[n:])
	return z
}

func (z *GFp) Sqrt(x *GFp) *GFp {
	p := z.modulus
	n := p.gfpWordLen

	if x.isZero() {
		z.setZero()
		return z
	}
	if p.bits <= 2 { // p == 2
		return z.Set(x)
	}

	var xv GFp
	xv.modulus = p
	xv.set(&x.nat)
	x = &xv

	// p % 4 == 3  =>  sqrt(x) = x^((p+1)/4)
	if p.words[0]&3 == 3 {
		var e [1 + simd.Words]big.Word
		copy(e[:n], p.words[:n])
		e[n] = simd.AddVW(e[:n], e[:n], 1) // e = p + 1
		simd.ShrVU(e[:n+1], e[:n+1], 2)    // e = (p + 1) / 4

		z.expWords(x, e[:n+1])
		return z.checkSqrt(x)
	}

	// pm1 = p - 1, half = (p - 1) / 2
	var pm1, half [1 + simd.Words]big.Word
	copy(pm1[:n], p.words[:n])
	simd.SubVW(pm1[:n], pm1[:n], 1)
	simd.ShrVU(half[:n], pm1[:n], 1)

	// Euler's criterion : x^((p-1)/2) == 1
	tmp := p.NewGFp()
	if !tmp.expWords(x, half[:n]).isOne() {
		return nil
	}

	// negOne = p - 1
	negOne := p.NewGFp()
	copy(negOne.words[:n], pm1[:n])

	// p - 1 = q * 2^s  (q 는 홀수)
	var q [1 + simd.Words]big.Word
	copy(q[:n], pm1[:n])
	s := 0
	for q[0]&1 == 0 {
		simd.ShrVU(q[:n], q[:n], 1)
		s++
	}

	// qh = (q + 1) / 2
	var qh [1 + simd.Words]big.Word
	copy(qh[:n], q[:n])
	qh[n] = simd.AddVW(qh[:n], qh[:n], 1)
	simd.ShrVU(qh[:n+1], qh[:n+1], 1)

	// 이차 비잉여 (quadratic non-residue) 탐색
	nonResidue := p.NewGFp()
	for w := uint64(2); ; w++ {
		nonResidue.SetUint64(w)
		if tmp.expWords(nonResidue, half[:n]).Cmp(negOne) == 0 {
			break
		}
	}

	m := s
	c := p.NewGFp().expWords(nonResidue, q[:n]) // c = nonResidue^q
	t := p.NewGFp().expWords(x, q[:n])          // t = x^q
	r := p.NewGFp().expWords(x, qh[:n+1])       // r = x^((q+1)/2)

	b := p.NewGFp()
	for !t.isOne() {
		// t^(2^i) == 1 인 가장 작은 i (0 < i < m) 를 찾는다.
		i := 0
		tmp.Set(t)
		for !tmp.isOne() {
			tmp.Sqr(tmp)
			i++
			if i >= m {
				return nil
			}
		}

		b.Set(c)
		for range m - i - 1 {
			b.Sqr(b)
		}

		m = i
		c.Sqr(b)
		t.Mul(t, c)
		r.Mul(r, b)
	}

	z.Set(r)
	return z.checkSqrt(x)
}

func (z *GFp) checkSqrt(x *GFp) *GFp {
	t := z.modulus.NewGFp()
	if t.Sqr(z).Cmp(x) != 0 {
		return nil
	}
	return z
}

func (z *GFp) expWords(x *GFp, e []big.Word) *GFp {
	n := z.modulus.gfpWordLen

	r := z.modulus.NewGFp()
	r.words[0] = 1

	base := z.modulus.NewGFp()
	base.set(&x.nat)

	for i := bitLen(e) - 1; i >= 0; i-- {
		r.Sqr(r)
		if (e[i/simd.WordBitSize]>>(i%simd.WordBitSize))&1 != 0 {
			r.Mul(r, base)
		}
	}

	copy(z.words[:n], r.words[:n])
	clear(z.words[n:])
	return z
}

//////////////////////////////////////////////////

func cmpWords(x, y []big.Word) int {
	for i := len(x) - 1; i >= 0; i-- {
		switch {
		case x[i] < y[i]:
			return -1
		case x[i] > y[i]:
			return 1
		}
	}
	return 0
}

func isOneWords(x []big.Word) bool {
	if len(x) == 0 || x[0] != 1 {
		return false
	}
	return isZeroWords(x[1:])
}

func halveMod(z, p []big.Word) {
	var c big.Word
	if z[0]&1 != 0 {
		c = simd.AddVV(z, z, p)
	}
	simd.ShrVU(z, z, 1)
	if c != 0 {
		z[len(z)-1] |= 1 << (simd.WordBitSize - 1)
	}
}

func subMod(z, x, y, p []big.Word) {
	if simd.SubVV(z, x, y) != 0 {
		simd.AddVV(z, z, p)
	}
}

// z = x * y
// math/big > nat.go > basicMul
func basicMul(z, x, y []big.Word) {
	clear(z)
	for i, d := range y {
		if d != 0 {
			z[len(x)+i] = simd.AddMulVVW(z[i:i+len(x)], x, d)
		}
	}
}

// z = x mod p
//
// math/big > nat.go > divBasic
func (p *Modulus) mod(z, x []big.Word) {
	n := p.gfpWordLen

	xn := len(x)
	for xn > 0 && x[xn-1] == 0 {
		xn--
	}

	if xn < n || (xn == n && cmpWords(x[:n], p.words[:n]) < 0) {
		if &z[0] != &x[0] {
			copy(z[:xn], x[:xn])
		}
		clear(z[xn:n])
		return
	}

	// 단일 워드 modulus
	if n == 1 {
		var r big.Word
		d := p.words[0]
		for i := xn - 1; i >= 0; i-- {
			_, r = simd.DivWW(r, x[i], d)
		}
		z[0] = r
		return
	}

	// D1. 정규화
	s := p.gfpShift
	v := p.gfpNorm[:n]

	var uArr [2*simd.Words + 1]big.Word
	u := uArr[:xn+1]
	u[xn] = simd.ShlVU(u[:xn], x[:xn], s)

	var qhatvArr [1 + simd.Words]big.Word
	qhatv := qhatvArr[:n+1]

	vn1 := v[n-1]
	vn2 := v[n-2]

	for j := xn - n; j >= 0; j-- {
		// D3. qhat 추정
		qhat := simd.WordMaxValue
		if ujn := u[j+n]; ujn != vn1 {
			var rhat big.Word
			qhat, rhat = simd.DivWWRecip(ujn, u[j+n-1], vn1, p.gfpRecip)

			// qhat * v[n-2] > rhat:u[j+n-2] 인 동안 qhat 을 줄인다.
			x1, x2 := simd.MulWW(qhat, vn2)
			for greaterThan(x1, x2, rhat, u[j+n-2]) {
				qhat--
				prevRhat := rhat
				rhat += vn1
				if rhat < prevRhat { // overflow
					break
				}
				x1, x2 = simd.MulWW(qhat, vn2)
			}
		}

		// D4. u[j:j+n+1] -= qhat * v
		qhatv[n] = simd.MulAddVWW(qhatv[:n], v, qhat, 0)
		if simd.SubVV(u[j:j+n+1], u[j:j+n+1], qhatv) != 0 {
			// D6. qhat 이 1 크므로 되돌린다.
			c := simd.AddVV(u[j:j+n], u[j:j+n], v)
			u[j+n] += c
		}
	}

	// D8. 나머지 역정규화
	simd.ShrVU(z[:n], u[:n], s)
}

func greaterThan(x1, x2, y1, y2 big.Word) bool {
	return x1 > y1 || x1 == y1 && x2 > y2
}
