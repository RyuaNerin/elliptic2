package curve

import (
	"crypto/elliptic"
	"encoding/asn1"
	"fmt"
	"math/big"
	"strings"

	"github.com/RyuaNerin/elliptic2"
	"github.com/RyuaNerin/elliptic2/internal/field"
)

const (
	/**
	benchmark on AMD Ryzen 7 9800X3D 8-Core Processor
	Window   P-192    P-224    P-256    P-384    P-521     avg
	Default  194.1%   251.3%   283.1%   724.0%   1222.0%   534.9%
	2        108.4%   140.7%   166.2%   346.5%    624.1%   277.2%
	3        100.6%   136.0%   150.7%   323.4%    575.7%   257.3%
	4        100.0%   130.4%   143.5%   315.2%    563.3%   250.5%
	5        106.1%   132.5%   145.8%   314.8%    560.8%   252.0%
	6        114.9%   144.1%   157.9%   339.7%    567.2%   264.7%
	*/
	maddWindowSize = 4 // must be >= 2
	maddTableSize  = 1 << (maddWindowSize - 1)
)

func init() {
	if maddWindowSize < 2 {
		panic("MaddWindowSize must be >= 2")
	}
}

type (
	FieldType int

	CurveArithmeticBase interface {
		FieldType() FieldType

		RawParams() any

		Modulus() *field.Modulus
		N() *big.Int

		Generator() (x, y *big.Int, ok bool)
		IsInfinity(x, y *big.Int) bool
		Identity() (x, y *big.Int)
		Params() *elliptic.CurveParams   // for elliptic.Curve compatibility
		Params2() *elliptic2.CurveParams // for CurveParamsSupplier compatibility

		IsOnCurve(x, y *big.Int) bool
		ComputeY(x *big.Int, largeY bool) *big.Int
	}
	CurveArithmetic[C any, CP Coordinate[C], OP Operator[C, CP]] interface {
		CurveArithmeticBase
		NewOperator() OP
	}
	Operator[C any, CP Coordinate[C]] interface {
		IsInfinity(p *C) bool
		SetInfinity(dst *C)

		ToCoordinate(dst *C, x, y *big.Int)
		ToAffinePoint(x, y *big.Int, p *C)

		Add(dst, p1, p2 *C)
		Double(dst, p1 *C)
	}
	MaddOperator[C any, CP Coordinate[C]] interface {
		Operator[C, CP]
		Madd(dst, p1, p2 *C)
		ScaleZ(p *C)
		Neg(dst, coords *C)
	}

	option    func(*curveBase)
	curveBase struct {
		oid asn1.ObjectIdentifier
	}
)

const (
	_ FieldType = iota
	FieldTypeGFp
	FieldTypeGF2m
)

func WithOID(oidStr string) option {
	var oid asn1.ObjectIdentifier
	for _, s := range strings.Split(oidStr, ".") {
		var v int
		_, err := fmt.Sscanf(s, "%d", &v)
		if err != nil {
			panic(fmt.Sprintf("invalid OID string: %s", oidStr))
		}
		oid = append(oid, v)
	}
	return func(o *curveBase) {
		o.oid = oid
	}
}

func (c *curveBase) OID() asn1.ObjectIdentifier {
	if c.oid == nil {
		return nil
	}
	return append(asn1.ObjectIdentifier(nil), c.oid...)
}

func NewCurve(base CurveArithmeticBase, opts ...option) elliptic.Curve {
	cb := curveBase{}
	for _, fn := range opts {
		fn(&cb)
	}

	switch base.FieldType() {
	case FieldTypeGFp:
		b, ok := base.(GFpCurveArithmetic)
		if !ok {
			panic("base is not GFpCurveBase")
		}
		c := curveGFp{curveBase: cb, base: b}

		op := b.NewOperator()
		if _, ok := op.(GFpMaddOperator); !ok {
			return &c
		}

		return &curveGFpMadd{curveGFp: c}

	case FieldTypeGF2m:
		b, ok := base.(GF2mCurveArithmetic)
		if !ok {
			panic("base is not GF2mCurveBase")
		}
		c := curveGF2m{curveBase: cb, base: b}

		op := b.NewOperator()
		if _, ok := op.(GF2mMaddOperator); !ok {
			return &c
		}

		return &curveGF2mMadd{curveGF2m: c}

	default:
		panic("unsupported field type")
	}
}

// for testing
func NewCurveSimple(base CurveArithmeticBase) elliptic.Curve {
	switch base.FieldType() {
	case FieldTypeGFp:
		b, ok := base.(GFpCurveArithmetic)
		if !ok {
			panic("base is not GFpCurveBase")
		}
		return &curveGFp{base: b}

	case FieldTypeGF2m:
		b, ok := base.(GF2mCurveArithmetic)
		if !ok {
			panic("base is not GF2mCurveBase")
		}
		return &curveGF2m{base: b}

	default:
		panic("unsupported field type")
	}
}

// for testing
func NewCurveMadd(base CurveArithmeticBase) elliptic.Curve {
	switch base.FieldType() {
	case FieldTypeGFp:
		b, ok := base.(GFpCurveArithmetic)
		if !ok {
			panic("base is not GFpCurveBase")
		}

		if _, ok := b.NewOperator().(GFpMaddOperator); !ok {
			return nil
		}

		return &curveGFpMadd{curveGFp: curveGFp{base: b}}

	case FieldTypeGF2m:
		b, ok := base.(GF2mCurveArithmetic)
		if !ok {
			panic("base is not GF2mCurveBase")
		}

		if _, ok := b.NewOperator().(GF2mMaddOperator); !ok {
			return nil
		}

		return &curveGF2mMadd{curveGF2m: curveGF2m{base: b}}

	default:
		panic("unsupported field type")
	}
}

func GetBase(c elliptic2.Curve) CurveArithmeticBase {
	switch p := c.(type) {
	case *curveGFp:
		return p.base
	case *curveGFpMadd:
		return p.base
	case *curveGF2m:
		return p.base
	case *curveGF2mMadd:
		return p.base
	}

	return nil
}

// computeWNAF computes the width-w Non-Adjacent Form of scalar k.
// k is big-endian byte slice. Result is little-endian (LSB at index 0).
func computeWNAF(naf []int8, k []byte) []int8 {
	const w = maddWindowSize
	const twoW int16 = 1 << w           // 2^w (예: w=4면 16)
	const halfTwoW int16 = 1 << (w - 1) // 2^(w-1) (예: w=4면 8)
	const mask int16 = twoW - 1         // 2^w - 1 (예: w=4면 15)
	const extraBits = maddTableSize     // process fixed extra digits (upper bound)

	// big-endian → little-endian 복사
	// int16 사용: 빼기 연산 시 음수/오버플로우 처리 용이
	buf := make([]int16, len(k)+1)
	for i, b := range k {
		buf[len(k)-1-i] = int16(b)
	}

	totalBits := len(k) * 8

	for i := range totalBits {
		bytePos := i >> 3
		bitPos := uint(i & 7)

		// w비트 윈도우 추출 (현재 위치부터 w비트)
		var window int16
		for j := range w {
			bp := (i + j) >> 3
			bt := uint((i + j) & 7)
			window |= ((buf[bp] >> bt) & 1) << j
		}

		// signed 표현: window >= 2^(w-1)이면 음수로 변환
		digit := window & mask
		subtractSigned := -((digit - halfTwoW) >> 15)
		digit -= twoW & subtractSigned

		isBit := window & 1
		isBitMask := -isBit
		digit &= isBitMask

		naf = append(naf, int8(digit))

		// buf에서 digit 빼기 (현재 비트 위치에서)
		// digit은 현재 비트 위치 기준이므로 그대로 빼면 됨
		buf[bytePos] -= digit << bitPos

		// carry/borrow 전파(상수시간)
		carry := int32(0)
		for j := range len(buf) - 1 {
			v := int32(buf[j]) + carry
			isNeg := v >> 31
			borrow := ((-v + 255) >> 8) & isNeg
			v += borrow << 8

			vCarry := v >> 8
			vCarry &= ^isNeg
			carry = vCarry - borrow

			buf[j] = int16(v - (vCarry << 8))
		}
		buf[len(buf)-1] = int16(int32(buf[len(buf)-1]) + carry)
	}

	for range extraBits {
		window := buf[len(buf)-1] & mask

		digit := window
		subtractSigned := -((digit - halfTwoW) >> 15)
		digit -= twoW & subtractSigned

		isBit := window & 1
		isBitMask := -isBit
		digit &= isBitMask

		naf = append(naf, int8(digit))
		buf[len(buf)-1] = (buf[len(buf)-1] - digit) >> 1
	}

	return naf
}

func normalizeScalar(scalar []byte, n *big.Int) []byte {
	byteSize := (n.BitLen() + 7) / 8
	if len(scalar) == byteSize {
		return scalar
	}
	s := new(big.Int).SetBytes(scalar)
	if len(scalar) > byteSize {
		s.Mod(s, n)
	}
	out := make([]byte, byteSize)
	return s.FillBytes(out)
}
