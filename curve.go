package elliptic2

import (
	"crypto/elliptic"

	"github.com/RyuaNerin/elliptic2/internal"
)

type CurveParams = internal.CurveParams

type Curve interface {
	elliptic.Curve
	BinaryParams() *CurveParams
}

// Create new elliptic curves over binary fields
// warning: params dose not validated.
//
// `krypto/elliptic2` uses `ellipse.CurveParams` for compatibility with `crypto.ellipse` package.
// But do not use the functions of `ellipse.CurveParams`. It will be panic.
func NewCurve(params *CurveParams) Curve {
	return internal.NewCurve(params)
}
