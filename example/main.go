package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/RyuaNerin/elliptic2"
	"github.com/RyuaNerin/elliptic2/curve/weierstrassbinary"
	"github.com/RyuaNerin/elliptic2/nist"
)

func main() {
	testB233()

	testP256()

	testCustomCurve()
}

func testP256() {
	msg := []byte("message")

	fmt.Println("==============================")
	fmt.Println("Testing interoperability between elliptic2 and standard library for P-256 curve")

	//////////////////////////////////////////////////

	fmt.Println()
	fmt.Println("1. Creating NIST P-256 curve using elliptic2/nist package")
	privLib, err := ecdsa.GenerateKey(nist.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	fmt.Println("Private Key D:", fmt.Sprintf("%x", privLib.D))
	fmt.Println("Public Key X:", fmt.Sprintf("%x", privLib.X))
	fmt.Println("Public Key Y:", fmt.Sprintf("%x", privLib.Y))
	fmt.Println("Curve Name :", privLib.Curve.Params().Name)
	fmt.Println("Curve == standard library P-256 ?", privLib.Curve == elliptic.P256())
	fmt.Println("Curve == elliptic2/nist P-256 ?", privLib.Curve == nist.P256())

	//////////////////////////////////////////////////

	fmt.Println()
	fmt.Println("2. Converting elliptic2 PublicKey to standard library PublicKey")
	pubStd := &privLib.PublicKey
	pubStd.Curve = elliptic.P256()
	fmt.Println("Converted Curve Name :", pubStd.Curve.Params().Name)
	fmt.Println("Curve == standard library P-256 ?", pubStd.Curve == elliptic.P256())
	fmt.Println("Curve == elliptic2/nist P-256 ?", pubStd.Curve == nist.P256())

	//////////////////////////////////////////////////

	fmt.Println()
	fmt.Println("3. Signing with elliptic2 P-256 curve")
	sig, err := ecdsa.SignASN1(rand.Reader, privLib, msg)
	if err != nil {
		panic(err)
	}
	fmt.Println("Signature (ASN.1):", fmt.Sprintf("%x", sig))

	//////////////////////////////////////////////////

	fmt.Println()
	fmt.Println("4. Verifying sig with standard library P-256 curve")
	valid := ecdsa.VerifyASN1(pubStd, msg, sig)
	fmt.Println("Signature valid:", valid)
}

func testB233() {
	msg := []byte("message")

	fmt.Println("==============================")
	fmt.Println("Testing elliptic2 B-233 curve")

	//////////////////////////////////////////////////

	fmt.Println("1. Creating NIST B-233 curve using elliptic2/nist package")
	curve := nist.B233()

	// if the curve created by elliptic2, it should implement CurveExtended
	curveExt, ok := curve.(elliptic2.CurveExtended)
	if !ok {
		panic("not CurveExtended")
	}

	//////////////////////////////////////////////////

	fmt.Println()
	fmt.Println("2. Scalar multiplication on B-233 curve")

	k := []byte{0x1, 0x2, 0x3, 0x4, 0x5}
	x, y := curveExt.ScalarBaseMult(k)
	fmt.Println("k =", fmt.Sprintf("%x", k))
	fmt.Printf("k*G = (%s, %s)\n", x.Text(16), y.Text(16))

	priv, err := ecdsa.GenerateKey(curveExt, rand.Reader)
	if err != nil {
		panic(err)
	}

	//////////////////////////////////////////////////

	fmt.Println()
	fmt.Println("3. Signing and verifying on B-233 curve")

	sig, err := ecdsa.SignASN1(rand.Reader, priv, msg)
	if err != nil {
		panic(err)
	}

	fmt.Println("Signature (ASN.1):", fmt.Sprintf("%x", sig))

	//////////////////////////////////////////////////

	fmt.Println()
	fmt.Println("4. Verifying signature on B-233 curve")

	valid := ecdsa.VerifyASN1(&priv.PublicKey, msg, sig)
	fmt.Println("Signature valid:", valid)
}

func testCustomCurve() {
	// B-163 based
	const CurveName = "My Curve"
	const BitSize = 163
	p := weierstrassbinary.Polynomials(163, 7, 6, 3, 0) // x^163 + x^7 + x^6 + x^3 + 1
	a, _ := new(big.Int).SetString("00000000000000000000000000000222222222222", 16)
	b, _ := new(big.Int).SetString("21111111111111111111111111111111111111111", 16)
	n, _ := new(big.Int).SetString("40000000000333333333333333333333333333333", 16)
	const myOID = "1.2.3.4.5.6.7"

	msg := []byte("message")

	var gx *big.Int
	var gy *big.Int

	// First, create a curve without G point
	fmt.Println("==============================")
	fmt.Println("1. Generating random G point on the curve")
	{
		curveWithoutG := weierstrassbinary.NewCurve(
			weierstrassbinary.CurveParams{
				Name:    CurveName,
				BitSize: BitSize,
				Poly:    p,
				A:       a,
				B:       b,
				N:       n,
				OID:     myOID,
			},
		)

		curveWithoutGExt := curveWithoutG.(elliptic2.CurveExtended)

		// Generate Random G point
		for range 1000 {
			fmt.Println("Generating random point G on the curve...")
			var err error
			gx, err = rand.Int(rand.Reader, p)
			if err != nil {
				panic(err)
			}

			gy = curveWithoutGExt.ComputeY(gx, false)
			if gy == nil {
				fmt.Println("No valid Y for the given X, retrying...")
				continue
			}
			if !curveWithoutGExt.IsOnCurve(gx, gy) {
				gy = nil
				fmt.Println("Computed point is not on the curve, retrying...")
				continue
			}

			fmt.Println("Found valid point G on the curve.")
			fmt.Println("Gx:", gx.Text(16))
			fmt.Println("Gy:", gy.Text(16))
			break
		}
		if gy == nil {
			panic("failed to find valid G point")
		}
	}

	//////////////////////////////////////////////////

	// Now, create a curve with the generated G point
	fmt.Println()
	fmt.Println("2. Creating custom curve with generated G point")
	c := weierstrassbinary.NewCurve(
		weierstrassbinary.CurveParams{
			Name:    CurveName,
			BitSize: BitSize,
			Poly:    p,
			A:       a,
			B:       b,
			N:       n,
			Gx:      gx,
			Gy:      gy,
			OID:     myOID,
		},
	)
	cExt := c.(elliptic2.CurveExtended)

	params := cExt.Params2()

	fmt.Println("Curve Name:", params.Name)
	fmt.Println("Field Type:", params.Type)
	fmt.Println("Bit Size :", params.BitSize)
	fmt.Println("Order    :", params.N.Text(16))
	fmt.Println("P        :", params.P.Text(16))
	fmt.Println("A        :", params.A.Text(16))
	fmt.Println("B        :", params.B.Text(16))
	fmt.Println("C        :", params.C.Text(16))
	fmt.Println("D        :", params.D.Text(16))
	fmt.Println("Gx       :", params.Gx.Text(16))
	fmt.Println("Gy       :", params.Gy.Text(16))
	fmt.Println("OID      :", cExt.OID().String())

	//////////////////////////////////////////////////

	fmt.Println()
	fmt.Println("3. Creating ECDSA key using elliptic2 custom curve")
	key, err := ecdsa.GenerateKey(nist.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	fmt.Println("Private Key D:", fmt.Sprintf("%x", key.D))
	fmt.Println("Public Key X:", fmt.Sprintf("%x", key.X))
	fmt.Println("Public Key Y:", fmt.Sprintf("%x", key.Y))
	fmt.Println("Curve Name :", key.Curve.Params().Name)
	fmt.Println("Curve == standard library P-256 ?", key.Curve == elliptic.P256())
	fmt.Println("Curve == elliptic2/nist P-256 ?", key.Curve == nist.P256())

	//////////////////////////////////////////////////

	fmt.Println()
	fmt.Println("4. Signing with custom curve")
	sig, err := ecdsa.SignASN1(rand.Reader, key, msg)
	if err != nil {
		panic(err)
	}
	fmt.Println("Signature (ASN.1):", fmt.Sprintf("%x", sig))

	//////////////////////////////////////////////////

	fmt.Println()
	fmt.Println("5. Verifying sig with standard library P-256 curve")
	valid := ecdsa.VerifyASN1(&key.PublicKey, msg, sig)
	fmt.Println("Signature valid:", valid)
}
