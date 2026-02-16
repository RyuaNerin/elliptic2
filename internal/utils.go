package internal

import (
	"encoding/hex"
	"fmt"
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
		switch {
		case '0' <= c && c <= '9':
		case 'a' <= c && c <= 'f':
		case 'A' <= c && c <= 'F':
		default:
			panic(fmt.Sprintf("invalid hex string: %s", s))
		}
		sb.WriteRune(c)
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
