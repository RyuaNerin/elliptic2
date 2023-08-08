package internal

import (
	"math/big"
	"strings"
)

func h(s string) string {
	var sb strings.Builder
	sb.Grow(len(s))
	s = strings.TrimPrefix(s, "0x")
	for _, c := range s {
		if '0' <= c && c <= '9' {
			sb.WriteRune(c)
		} else if 'a' <= c && c <= 'f' {
			sb.WriteRune(c)
		} else if 'A' <= c && c <= 'F' {
			sb.WriteRune(c)
		}
	}

	return sb.String()
}

// hex to *big.Int
func HI(s string) *big.Int {
	s = h(s)
	if s == "" {
		return new(big.Int)
	}
	result, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic(s)
	}
	return result
}
