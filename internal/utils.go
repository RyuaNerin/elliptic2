package internal

import (
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"strings"
)

func H(s string) (ss string, isNegative bool) {
	if strings.HasPrefix(s, "-") {
		isNegative = true
		s = s[1:]
	}

	var sb strings.Builder
	sb.Grow(len(s))
	s = strings.TrimPrefix(s, "0x")

	if len(s)%2 != 0 {
		sb.WriteByte('0')
	}
	for _, c := range s {
		if '0' <= c && c <= '9' {
			sb.WriteRune(c)
		} else if 'a' <= c && c <= 'f' {
			sb.WriteRune(c)
		} else if 'A' <= c && c <= 'F' {
			sb.WriteRune(c)
		} else {
			panic(fmt.Sprintf("invalid hex string: %s", s))
		}
	}

	ss = sb.String()
	if ss == "" {
		ss = "0"
	}

	return ss, isNegative
}

// hex to *big.Int
func HI(s string) *big.Int {
	s, isNegative := H(s)
	if s == "" {
		return new(big.Int).SetInt64(0)
	}
	result, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic(s)
	}

	if isNegative {
		result.Neg(result)
	}

	return result
}

func HB(s string) []byte {
	s, isNegative := H(s)
	if isNegative {
		panic("negative hex string")
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return b
}

// Int returns a uniform random value in [0, max). It panics if max <= 0, and
// returns an error if rand.Read returns one.
func Int(rand io.Reader, max *big.Int) (n *big.Int, err error) {
	if max.Sign() <= 0 {
		panic("crypto/rand: argument to Int is <= 0")
	}
	n = new(big.Int)
	n.Sub(max, n.SetUint64(1))
	// bitLen is the maximum bit length needed to encode a value < max.
	bitLen := n.BitLen()
	if bitLen == 0 {
		// the only valid result is 0
		return
	}
	// k is the maximum byte length needed to encode a value < max.
	k := (bitLen + 7) / 8
	// b is the number of bits in the most significant byte of max-1.
	b := uint(bitLen % 8)
	if b == 0 {
		b = 8
	}

	bytes := make([]byte, k)

	for {
		_, err = io.ReadFull(rand, bytes)
		if err != nil {
			return nil, err
		}

		// Clear bits in the first byte to increase the probability
		// that the candidate is < max.
		bytes[0] &= uint8(int(1<<b) - 1)

		n.SetBytes(bytes)
		if n.Cmp(max) < 0 {
			return
		}
	}
}
