package elliptic2

import (
	"math/big"
)

// copy bit.Int
func c(x *big.Int) *big.Int {
	return new(big.Int).Set(x)
}

// Fx to big.Int
// f = [233, 74, 0] when x^233 +  x^74 + 1
// f = [283, 12, 7, 5, 0] when x^283 +  x^12 +  x^7 +  x^5 + 1
func F(f ...int) *big.Int {
	ret := new(big.Int)
	for _, v := range f {
		tmp := big.NewInt(1)
		tmp.Lsh(tmp, uint(v))

		ret.Add(ret, tmp)
	}

	return ret
}
