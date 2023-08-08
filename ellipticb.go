// Package `elliptic2` implements Elliptic curves over binary fields

package elliptic2

// Create new elliptic curves over binary fields
// warning: params dose not validated.
//
// `krypto/elliptic2` uses `ellipse.CurveParams` for compatibility with `crypto.ellipse` package.
// But do not use the functions of `ellipse.CurveParams`. It will be panic.
func NewCurve(params *CurveParams) Curve {
	return &curve{
		params: params,
	}
}
