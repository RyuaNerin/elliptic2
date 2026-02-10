[![PkgGoDev](https://pkg.go.dev/badge/github.com/RyuaNerin/elliptic2/v2)](https://pkg.go.dev/github.com/RyuaNerin/elliptic2/v2)

# elliptic2

- `elliptic2` is a Go package providing elliptic curve implementations and utilities written in pure Go.

- Supports a wide range of curves including NIST, Brainpool, ANSSI, GOST, BLS, BN, MNT, NUMS, Oakley, Oscaa, Secg, and more (see the list below).

- Supports go 1.22+.

- **⚠️ Warning: This library is not constant-time and is vulnerable to timing side-channel attacks when handling secret data.**

## Installation

```bash
$ go get github.com/RyuaNerin/elliptic2/v2
```

## Usage

- For more detailed examples, see the [example/main.go](example/main.go) file.

- The `nist.B233()`, `nist.P256()`, and similar functions return unique interface instances.

  - These can be used to identify elliptic curves in if or switch statements.

  ```go
  curve := nist.P256()

  switch curve {
  case nist.P256():
      fmt.Println("Using P-256 curve")
  case nist.B233():
      fmt.Println("Using B-233 curve")
  default:
      fmt.Println("Unknown curve")
  }
  ```

- Compatibility with Go's ECDSA Package

  ```go
  import (
      "crypto/ecdsa"
      "crypto/rand"

      "github.com/RyuaNerin/elliptic2/v2/nist"
  )

  func main() {
      // Generate ECDSA key pair using elliptic2 curve
      privateKey, err := ecdsa.GenerateKey(nist.P256(), rand.Reader)
      if err != nil {
          panic(err)
      }

      // Sign and verify as usual
      message := []byte("hello")
      signature, err := ecdsa.SignASN1(rand.Reader, privateKey, message)
      if err != nil {
          panic(err)
      }

      valid := ecdsa.VerifyASN1(&privateKey.PublicKey, message, signature)
      fmt.Println("Signature valid:", valid)
  }
  ```

## [LICENSE](/LICENSE)

- Apache License 2.0

## Supported Curve List

| Curve                                                                    | Field Type | Package                                                                                                                                     | Fomular                   |
|:-------------------------------------------------------------------------|:----------:|:--------------------------------------------------------------------------------------------------------------------------------------------|:--------------------------|
| [Edwards](https://hyperelliptic.org/EFD/g1p/auto-edwards.html)           | Prime      | [`github.com/RyuaNerin/elliptic2/v2/curve/edwards`](https://pkg.go.dev/github.com/RyuaNerin/elliptic2/v2/curve/edwards)                     | `xx + yy = cc(1 + dxxyy)` |
| [Twisted Edwards](https://hyperelliptic.org/EFD/g1p/auto-twisted.html)   | Prime      | [`github.com/RyuaNerin/elliptic2/v2/curve/twistededwards`](https://pkg.go.dev/github.com/RyuaNerin/elliptic2/v2/curve/twistededwards)       | `axx + yy = 1 + dxxyy`    |
| [Montgomery](https://hyperelliptic.org/EFD/g1p/auto-montgom.html)        | Prime      | [`github.com/RyuaNerin/elliptic2/v2/curve/montgomery`](https://pkg.go.dev/github.com/RyuaNerin/elliptic2/v2/curve/montgomery)               | `byy = xxx + axx + x`     |
| [Short Weierstrass](https://hyperelliptic.org/EFD/g1p/auto-shortw.html)  | Prime      | [`github.com/RyuaNerin/elliptic2/v2/curve/weierstrassprime`](https://pkg.go.dev/github.com/RyuaNerin/elliptic2/v2/curve/weierstrassprime)   | `yy = xxx + ax + b`       |
| [Short Weierstrass](https://hyperelliptic.org/EFD/g12o/auto-shortw.html) | Binary     | [`github.com/RyuaNerin/elliptic2/v2/curve/weierstrassbinary`](https://pkg.go.dev/github.com/RyuaNerin/elliptic2/v2/curve/weierstrassbinary) | `yy + xy = xxx + axx + b` |

## Standard Curves

<!-- curves start -->
### [ANSI X9.63](https://neuromancer.sk/std/x963/)

- ANSI x9.63 example curves.

| Name | OID | Field Type | Generator | Function | As known |
|:-----|:----|:----------:|:---------:|:---------|:---------|
| [`ansip160k1`](https://neuromancer.sk/std/x963/ansip160k1) | `1.3.132.0.9` | `prime field` | O | `x963.Ansip160k1()` | `secg/secp160k1` |
| [`ansip160r1`](https://neuromancer.sk/std/x963/ansip160r1) | `1.3.132.0.8` | `prime field` | O | `x963.Ansip160r1()` | `secg/secp160r1`, `wtls/wap-wsg-idm-ecid-wtls7` |
| [`ansip160r2`](https://neuromancer.sk/std/x963/ansip160r2) | `1.3.132.0.30` | `prime field` | O | `x963.Ansip160r2()` | `secg/secp160r2` |
| [`ansip192k1`](https://neuromancer.sk/std/x963/ansip192k1) | `1.3.132.0.31` | `prime field` | O | `x963.Ansip192k1()` | `secg/secp192k1` |
| [`ansip224k1`](https://neuromancer.sk/std/x963/ansip224k1) | `1.3.132.0.32` | `prime field` | O | `x963.Ansip224k1()` | `secg/secp224k1` |
| [`ansip224r1`](https://neuromancer.sk/std/x963/ansip224r1) | `1.3.132.0.33` | `prime field` | O | `x963.Ansip224r1()` | `wtls/wap-wsg-idm-ecid-wtls12`, `nist/P-224`, `secg/secp224r1` |
| [`ansip256k1`](https://neuromancer.sk/std/x963/ansip256k1) | `1.3.132.0.10` | `prime field` | O | `x963.Ansip256k1()` | `secg/secp256k1` |
| [`ansip384r1`](https://neuromancer.sk/std/x963/ansip384r1) | `1.3.132.0.34` | `prime field` | O | `x963.Ansip384r1()` | `secg/secp384r1`, `nist/P-384` |
| [`ansip521r1`](https://neuromancer.sk/std/x963/ansip521r1) | `1.3.132.0.35` | `prime field` | O | `x963.Ansip521r1()` | `secg/secp521r1`, `nist/P-521` |
| [`ansit163k1`](https://neuromancer.sk/std/x963/ansit163k1) | `1.3.132.0.1` | `binary field` | O | `x963.Ansit163k1()` | `secg/sect163k1`, `nist/k-163`, `wtls/wap-wsg-idm-ecid-wtls3` |
| [`ansit163r1`](https://neuromancer.sk/std/x963/ansit163r1) | `1.3.132.0.2` | `binary field` | O | `x963.Ansit163r1()` | `secg/sect163r1` |
| [`ansit163r2`](https://neuromancer.sk/std/x963/ansit163r2) | `1.3.132.0.15` | `binary field` | O | `x963.Ansit163r2()` | `secg/sect163r2`, `nist/B-163` |
| [`ansit193r1`](https://neuromancer.sk/std/x963/ansit193r1) | `1.3.132.0.24` | `binary field` | O | `x963.Ansit193r1()` | `secg/sect193r1` |
| [`ansit193r2`](https://neuromancer.sk/std/x963/ansit193r2) | `1.3.132.0.25` | `binary field` | O | `x963.Ansit193r2()` | `secg/sect193r2` |
| [`ansit233k1`](https://neuromancer.sk/std/x963/ansit233k1) | `1.3.132.0.26` | `binary field` | O | `x963.Ansit233k1()` | `secg/sect233k1`, `nist/K-233`, `wtls/wap-wsg-idm-ecid-wtls10` |
| [`ansit233r1`](https://neuromancer.sk/std/x963/ansit233r1) | `1.3.132.0.27` | `binary field` | O | `x963.Ansit233r1()` | `wtls/wap-wsg-idm-ecid-wtls11`, `nist/B-233`, `secg/sect233r1` |
| [`ansit239k1`](https://neuromancer.sk/std/x963/ansit239k1) | `1.3.132.0.3` | `binary field` | O | `x963.Ansit239k1()` | `secg/sect239k1` |
| [`ansit283k1`](https://neuromancer.sk/std/x963/ansit283k1) | `1.3.132.0.16` | `binary field` | O | `x963.Ansit283k1()` | `nist/K-283`, `secg/sect283k1` |
| [`ansit283r1`](https://neuromancer.sk/std/x963/ansit283r1) | `1.3.132.0.17` | `binary field` | O | `x963.Ansit283r1()` | `nist/B-283`, `secg/sect283r1` |
| [`ansit409k1`](https://neuromancer.sk/std/x963/ansit409k1) | `1.3.132.0.36` | `binary field` | O | `x963.Ansit409k1()` | `nist/K-409`, `secg/sect409k1` |
| [`ansit409r1`](https://neuromancer.sk/std/x963/ansit409r1) | `1.3.132.0.37` | `binary field` | O | `x963.Ansit409r1()` | `nist/B-409`, `secg/sect409r1` |
| [`ansit571k1`](https://neuromancer.sk/std/x963/ansit571k1) | `1.3.132.0.38` | `binary field` | O | `x963.Ansit571k1()` | `nist/K-571`, `secg/sect571k1` |
| [`ansit571r1`](https://neuromancer.sk/std/x963/ansit571r1) | `1.3.132.0.39` | `binary field` | O | `x963.Ansit571r1()` | `nist/B-571`, `secg/sect571r1` |

### [ANSI x9.62](https://neuromancer.sk/std/x962/)

- ANSI x9.62 example curves.

| Name | OID | Field Type | Generator | Function | As known |
|:-----|:----|:----------:|:---------:|:---------|:---------|
| [`c2pnb163v1`](https://neuromancer.sk/std/x962/c2pnb163v1) | `1.2.840.10045.3.0.1` | `binary field` | O | `x962.C2pnb163v1()` | `wtls/wap-wsg-idm-ecid-wtls5` |
| [`c2pnb163v2`](https://neuromancer.sk/std/x962/c2pnb163v2) | `1.2.840.10045.3.0.2` | `binary field` | O | `x962.C2pnb163v2()` | |
| [`c2pnb163v3`](https://neuromancer.sk/std/x962/c2pnb163v3) | `1.2.840.10045.3.0.3` | `binary field` | O | `x962.C2pnb163v3()` | |
| [`c2pnb176w1`](https://neuromancer.sk/std/x962/c2pnb176w1) | `1.2.840.10045.3.0.4` | `binary field` | O | `x962.C2pnb176w1()` | |
| [`c2pnb208w1`](https://neuromancer.sk/std/x962/c2pnb208w1) | `1.2.840.10045.3.0.10` | `binary field` | O | `x962.C2pnb208w1()` | |
| [`c2pnb272w1`](https://neuromancer.sk/std/x962/c2pnb272w1) | `1.2.840.10045.3.0.16` | `binary field` | O | `x962.C2pnb272w1()` | |
| [`c2pnb304w1`](https://neuromancer.sk/std/x962/c2pnb304w1) | `1.2.840.10045.3.0.17` | `binary field` | O | `x962.C2pnb304w1()` | |
| [`c2pnb368w1`](https://neuromancer.sk/std/x962/c2pnb368w1) | `1.2.840.10045.3.0.19` | `binary field` | O | `x962.C2pnb368w1()` | |
| [`c2tnb191v1`](https://neuromancer.sk/std/x962/c2tnb191v1) | `1.2.840.10045.3.0.5` | `binary field` | O | `x962.C2tnb191v1()` | |
| [`c2tnb191v2`](https://neuromancer.sk/std/x962/c2tnb191v2) | `1.2.840.10045.3.0.6` | `binary field` | O | `x962.C2tnb191v2()` | |
| [`c2tnb191v3`](https://neuromancer.sk/std/x962/c2tnb191v3) | `1.2.840.10045.3.0.7` | `binary field` | O | `x962.C2tnb191v3()` | |
| [`c2tnb239v1`](https://neuromancer.sk/std/x962/c2tnb239v1) | `1.2.840.10045.3.0.11` | `binary field` | O | `x962.C2tnb239v1()` | |
| [`c2tnb239v2`](https://neuromancer.sk/std/x962/c2tnb239v2) | `1.2.840.10045.3.0.12` | `binary field` | O | `x962.C2tnb239v2()` | |
| [`c2tnb239v3`](https://neuromancer.sk/std/x962/c2tnb239v3) | `1.2.840.10045.3.0.13` | `binary field` | O | `x962.C2tnb239v3()` | |
| [`c2tnb359v1`](https://neuromancer.sk/std/x962/c2tnb359v1) | `1.2.840.10045.3.0.18` | `binary field` | O | `x962.C2tnb359v1()` | |
| [`c2tnb431r1`](https://neuromancer.sk/std/x962/c2tnb431r1) | `1.2.840.10045.3.0.20` | `binary field` | O | `x962.C2tnb431r1()` | |
| [`prime192v1`](https://neuromancer.sk/std/x962/prime192v1) | `1.2.840.10045.3.1.1` | `prime field` | O | `x962.Prime192v1()` | `secg/secp192r1`, `nist/P-192` |
| [`prime192v2`](https://neuromancer.sk/std/x962/prime192v2) | `1.2.840.10045.3.1.2` | `prime field` | O | `x962.Prime192v2()` | |
| [`prime192v3`](https://neuromancer.sk/std/x962/prime192v3) | `1.2.840.10045.3.1.3` | `prime field` | O | `x962.Prime192v3()` | |
| [`prime239v1`](https://neuromancer.sk/std/x962/prime239v1) | `1.2.840.10045.3.1.4` | `prime field` | O | `x962.Prime239v1()` | |
| [`prime239v2`](https://neuromancer.sk/std/x962/prime239v2) | `1.2.840.10045.3.1.5` | `prime field` | O | `x962.Prime239v2()` | |
| [`prime239v3`](https://neuromancer.sk/std/x962/prime239v3) | `1.2.840.10045.3.1.6` | `prime field` | O | `x962.Prime239v3()` | |
| [`prime256v1`](https://neuromancer.sk/std/x962/prime256v1) | `1.2.840.10045.3.1.7` | `prime field` | O | `x962.Prime256v1()` | `secg/secp256r1`, `nist/P-256` |

- Those curves are not avilable in this namespace:
  | Name | Reason |
  |:-----|:-------|
  | [`c2onb191v4`](https://neuromancer.sk/std/x962/c2onb191v4) | unsupported curve type |
  | [`c2onb191v5`](https://neuromancer.sk/std/x962/c2onb191v5) | unsupported curve type |
  | [`c2onb239v4`](https://neuromancer.sk/std/x962/c2onb239v4) | unsupported curve type |
  | [`c2onb239v5`](https://neuromancer.sk/std/x962/c2onb239v5) | unsupported curve type |

### [ANSSI](https://neuromancer.sk/std/anssi/)

- Agence nationale de la s챕curit챕 des syst챔mes d'information: Publication d'un param챕trage de courbe elliptique visant des applications de passeport 챕lectronique et de l'administration 챕lectronique fran챌aise. 21 November 2011

| Name | OID | Field Type | Generator | Function | As known |
|:-----|:----|:----------:|:---------:|:---------|:---------|
| [`FRP256v1`](https://neuromancer.sk/std/anssi/FRP256v1) | `1.2.250.1.223.101.256.1` | `prime field` | O | `anssi.FRP256v1()` | |

### [Barreto-Lynn-Scott](https://neuromancer.sk/std/bls/)

- BLS curves. A family of pairing friendly curves, with embedding degree = 12 or 24.

| Name | OID | Field Type | Generator | Function | As known |
|:-----|:----|:----------:|:---------:|:---------|:---------|
| [`BLS12-377`](https://neuromancer.sk/std/bls/BLS12-377) || `prime field` | O | `bls.BLS12_377()` | |
| [`BLS12-381`](https://neuromancer.sk/std/bls/BLS12-381) || `prime field` | O | `bls.BLS12_381()` | |
| [`BLS12-446`](https://neuromancer.sk/std/bls/BLS12-446) || `prime field` | O | `bls.BLS12_446()` | |
| [`BLS12-455`](https://neuromancer.sk/std/bls/BLS12-455) || `prime field` | O | `bls.BLS12_455()` | |
| [`BLS12-638`](https://neuromancer.sk/std/bls/BLS12-638) || `prime field` | O | `bls.BLS12_638()` | |
| [`BLS24-477`](https://neuromancer.sk/std/bls/BLS24-477) || `prime field` | O | `bls.BLS24_477()` | |
| [`Bandersnatch`](https://neuromancer.sk/std/bls/Bandersnatch) || `prime field` || `bls.Bandersnatch()` | |

### [Barreto-Naehrig](https://neuromancer.sk/std/bn/)

- BN (Barreto, Naehrig curves) from: A Family of Implementation-Friendly BN Elliptic Curves - <https://eprint.iacr.org/2010/429.pdf>.

| Name | OID | Field Type | Generator | Function | As known |
|:-----|:----|:----------:|:---------:|:---------|:---------|
| [`bn158`](https://neuromancer.sk/std/bn/bn158) || `prime field` | O | `bn.Bn158()` | |
| [`bn190`](https://neuromancer.sk/std/bn/bn190) || `prime field` | O | `bn.Bn190()` | |
| [`bn222`](https://neuromancer.sk/std/bn/bn222) || `prime field` | O | `bn.Bn222()` | |
| [`bn254`](https://neuromancer.sk/std/bn/bn254) || `prime field` | O | `bn.Bn254()` | `other/Fp254BNb` |
| [`bn286`](https://neuromancer.sk/std/bn/bn286) || `prime field` | O | `bn.Bn286()` | |
| [`bn318`](https://neuromancer.sk/std/bn/bn318) || `prime field` | O | `bn.Bn318()` | |
| [`bn350`](https://neuromancer.sk/std/bn/bn350) || `prime field` | O | `bn.Bn350()` | |
| [`bn382`](https://neuromancer.sk/std/bn/bn382) || `prime field` | O | `bn.Bn382()` | |
| [`bn414`](https://neuromancer.sk/std/bn/bn414) || `prime field` | O | `bn.Bn414()` | |
| [`bn446`](https://neuromancer.sk/std/bn/bn446) || `prime field` | O | `bn.Bn446()` | |
| [`bn478`](https://neuromancer.sk/std/bn/bn478) || `prime field` | O | `bn.Bn478()` | |
| [`bn510`](https://neuromancer.sk/std/bn/bn510) || `prime field` | O | `bn.Bn510()` | |
| [`bn542`](https://neuromancer.sk/std/bn/bn542) || `prime field` | O | `bn.Bn542()` | |
| [`bn574`](https://neuromancer.sk/std/bn/bn574) || `prime field` | O | `bn.Bn574()` | |
| [`bn606`](https://neuromancer.sk/std/bn/bn606) || `prime field` | O | `bn.Bn606()` | |
| [`bn638`](https://neuromancer.sk/std/bn/bn638) || `prime field` | O | `bn.Bn638()` | |

### [Brainpool](https://neuromancer.sk/std/brainpool/)

- ECC Brainpool Standard Curves and Curve Generation v. 1.0  19.10.2005

| Name | OID | Field Type | Generator | Function | As known |
|:-----|:----|:----------:|:---------:|:---------|:---------|
| [`brainpoolP160r1`](https://neuromancer.sk/std/brainpool/brainpoolP160r1) | `1.3.36.3.3.2.8.1.1.1` | `prime field` | O | `brainpool.BrainpoolP160r1()` | |
| [`brainpoolP160t1`](https://neuromancer.sk/std/brainpool/brainpoolP160t1) | `1.3.36.3.3.2.8.1.1.2` | `prime field` | O | `brainpool.BrainpoolP160t1()` | |
| [`brainpoolP192r1`](https://neuromancer.sk/std/brainpool/brainpoolP192r1) | `1.3.36.3.3.2.8.1.1.3` | `prime field` | O | `brainpool.BrainpoolP192r1()` | |
| [`brainpoolP192t1`](https://neuromancer.sk/std/brainpool/brainpoolP192t1) | `1.3.36.3.3.2.8.1.1.4` | `prime field` | O | `brainpool.BrainpoolP192t1()` | |
| [`brainpoolP224r1`](https://neuromancer.sk/std/brainpool/brainpoolP224r1) | `1.3.36.3.3.2.8.1.1.5` | `prime field` | O | `brainpool.BrainpoolP224r1()` | |
| [`brainpoolP224t1`](https://neuromancer.sk/std/brainpool/brainpoolP224t1) | `1.3.36.3.3.2.8.1.1.6` | `prime field` | O | `brainpool.BrainpoolP224t1()` | |
| [`brainpoolP256r1`](https://neuromancer.sk/std/brainpool/brainpoolP256r1) | `1.3.36.3.3.2.8.1.1.7` | `prime field` | O | `brainpool.BrainpoolP256r1()` | |
| [`brainpoolP256t1`](https://neuromancer.sk/std/brainpool/brainpoolP256t1) | `1.3.36.3.3.2.8.1.1.8` | `prime field` | O | `brainpool.BrainpoolP256t1()` | |
| [`brainpoolP320r1`](https://neuromancer.sk/std/brainpool/brainpoolP320r1) | `1.3.36.3.3.2.8.1.1.9` | `prime field` | O | `brainpool.BrainpoolP320r1()` | |
| [`brainpoolP320t1`](https://neuromancer.sk/std/brainpool/brainpoolP320t1) | `1.3.36.3.3.2.8.1.1.10` | `prime field` | O | `brainpool.BrainpoolP320t1()` | |
| [`brainpoolP384r1`](https://neuromancer.sk/std/brainpool/brainpoolP384r1) | `1.3.36.3.3.2.8.1.1.11` | `prime field` | O | `brainpool.BrainpoolP384r1()` | |
| [`brainpoolP384t1`](https://neuromancer.sk/std/brainpool/brainpoolP384t1) | `1.3.36.3.3.2.8.1.1.12` | `prime field` | O | `brainpool.BrainpoolP384t1()` | |
| [`brainpoolP512r1`](https://neuromancer.sk/std/brainpool/brainpoolP512r1) | `1.3.36.3.3.2.8.1.1.13` | `prime field` | O | `brainpool.BrainpoolP512r1()` | |
| [`brainpoolP512t1`](https://neuromancer.sk/std/brainpool/brainpoolP512t1) | `1.3.36.3.3.2.8.1.1.14` | `prime field` | O | `brainpool.BrainpoolP512t1()` | |

### [GOST](https://neuromancer.sk/std/gost/)

- GOST R 34.10-2001: RFC5832, GOST R 34.10-2012: RFC7836

| Name | OID | Field Type | Generator | Function | As known |
|:-----|:----|:----------:|:---------:|:---------|:---------|
| [`gost256`](https://neuromancer.sk/std/gost/gost256) || `prime field` | O | `gost.Gost256()` | |
| [`gost512`](https://neuromancer.sk/std/gost/gost512) || `prime field` | O | `gost.Gost512()` | |
| [`id-GostR3410-2001-CryptoPro-A-ParamSet`](https://neuromancer.sk/std/gost/id-GostR3410-2001-CryptoPro-A-ParamSet) || `prime field` | O | `gost.IdGostR3410_2001CryptoProAParamSet()` | |
| [`id-GostR3410-2001-CryptoPro-B-ParamSet`](https://neuromancer.sk/std/gost/id-GostR3410-2001-CryptoPro-B-ParamSet) || `prime field` | O | `gost.IdGostR3410_2001CryptoProBParamSet()` | |
| [`id-GostR3410-2001-CryptoPro-C-ParamSet`](https://neuromancer.sk/std/gost/id-GostR3410-2001-CryptoPro-C-ParamSet) || `prime field` || `gost.IdGostR3410_2001CryptoProCParamSet()` | |
| [`id-tc26-gost-3410-12-512-paramSetA`](https://neuromancer.sk/std/gost/id-tc26-gost-3410-12-512-paramSetA) || `prime field` | O | `gost.IdTc26Gost3410_12_512ParamSetA()` | |
| [`id-tc26-gost-3410-12-512-paramSetB`](https://neuromancer.sk/std/gost/id-tc26-gost-3410-12-512-paramSetB) || `prime field` | O | `gost.IdTc26Gost3410_12_512ParamSetB()` | |
| [`id-tc26-gost-3410-2012-256-paramSetA`](https://neuromancer.sk/std/gost/id-tc26-gost-3410-2012-256-paramSetA) || `prime field` | O | `gost.IdTc26Gost3410_2012_256ParamSetA()` | |
| [`id-tc26-gost-3410-2012-512-paramSetC`](https://neuromancer.sk/std/gost/id-tc26-gost-3410-2012-512-paramSetC) || `prime field` | O | `gost.IdTc26Gost3410_2012_512ParamSetC()` | |

### [Miyaji-Nakabayashi-Takano](https://neuromancer.sk/std/mnt/)

- MNT (Miyaji, Nakabayashi, and Takano curves) example curves from: New explicit conditions of elliptic curve traces for FR-reduction - https://dspace.jaist.ac.jp/dspace/bitstream/10119/4432/1/73-48.pdf.

| Name | OID | Field Type | Generator | Function | As known |
|:-----|:----|:----------:|:---------:|:---------|:---------|
| [`mnt1`](https://neuromancer.sk/std/mnt/mnt1) || `prime field` | O | `mnt.Mnt1()` | |
| [`mnt2/1`](https://neuromancer.sk/std/mnt/mnt2/1) || `prime field` | O | `mnt.Mnt2_1()` | |
| [`mnt2/2`](https://neuromancer.sk/std/mnt/mnt2/2) || `prime field` | O | `mnt.Mnt2_2()` | |
| [`mnt3/1`](https://neuromancer.sk/std/mnt/mnt3/1) || `prime field` | O | `mnt.Mnt3_1()` | |
| [`mnt3/2`](https://neuromancer.sk/std/mnt/mnt3/2) || `prime field` | O | `mnt.Mnt3_2()` | |
| [`mnt3/3`](https://neuromancer.sk/std/mnt/mnt3/3) || `prime field` | O | `mnt.Mnt3_3()` | |
| [`mnt4`](https://neuromancer.sk/std/mnt/mnt4) || `prime field` | O | `mnt.Mnt4()` | |
| [`mnt5/1`](https://neuromancer.sk/std/mnt/mnt5/1) || `prime field` | O | `mnt.Mnt5_1()` | |
| [`mnt5/2`](https://neuromancer.sk/std/mnt/mnt5/2) || `prime field` | O | `mnt.Mnt5_2()` | |
| [`mnt5/3`](https://neuromancer.sk/std/mnt/mnt5/3) || `prime field` | O | `mnt.Mnt5_3()` | |

### [NIST](https://neuromancer.sk/std/nist/)

- RECOMMENDED ELLIPTIC CURVES FOR FEDERAL GOVERNMENT USE  July 1999

| Name | OID | Field Type | Generator | Function | As known |
|:-----|:----|:----------:|:---------:|:---------|:---------|
| [`B-163`](https://neuromancer.sk/std/nist/B-163) | `1.3.132.0.15` | `binary field` | O | `nist.B163()` | `secg/sect163r2`, `x963/ansit163r2` |
| [`B-233`](https://neuromancer.sk/std/nist/B-233) | `1.3.132.0.27` | `binary field` | O | `nist.B233()` | `secg/sect233r1`, `wtls/wap-wsg-idm-ecid-wtls11`, `x963/ansit233r1` |
| [`B-283`](https://neuromancer.sk/std/nist/B-283) | `1.3.132.0.17` | `binary field` | O | `nist.B283()` | `secg/sect283r1`, `x963/ansit283r1` |
| [`B-409`](https://neuromancer.sk/std/nist/B-409) | `1.3.132.0.37` | `binary field` | O | `nist.B409()` | `secg/sect409r1`, `x963/ansit409r1` |
| [`B-571`](https://neuromancer.sk/std/nist/B-571) | `1.3.132.0.39` | `binary field` | O | `nist.B571()` | `secg/sect571r1`, `x963/ansit571r1` |
| [`K-163`](https://neuromancer.sk/std/nist/K-163) | `1.3.132.0.1` | `binary field` | O | `nist.K163()` | `secg/sect163k1`, `x963/ansit163k1`, `wtls/wap-wsg-idm-ecid-wtls3` |
| [`K-233`](https://neuromancer.sk/std/nist/K-233) | `1.3.132.0.26` | `binary field` | O | `nist.K233()` | `secg/sect233k1`, `wtls/wap-wsg-idm-ecid-wtls10`, `x963/ansit233k1` |
| [`K-283`](https://neuromancer.sk/std/nist/K-283) | `1.3.132.0.16` | `binary field` | O | `nist.K283()` | `secg/sect283k1`, `x963/ansit283k1` |
| [`K-409`](https://neuromancer.sk/std/nist/K-409) | `1.3.132.0.36` | `binary field` | O | `nist.K409()` | `secg/sect409k1`, `x963/ansit409k1` |
| [`K-571`](https://neuromancer.sk/std/nist/K-571) | `1.3.132.0.38` | `binary field` | O | `nist.K571()` | `secg/sect571k1`, `x963/ansit571k1` |
| [`P-192`](https://neuromancer.sk/std/nist/P-192) | `1.2.840.10045.3.1.1` | `prime field` | O | `nist.P192()` | `secg/secp192r1`, `x962/prime192v1` |
| [`P-224`](https://neuromancer.sk/std/nist/P-224) | `1.3.132.0.33` | `prime field` | O | `nist.P224()` | `secg/secp224r1`, `wtls/wap-wsg-idm-ecid-wtls12`, `x963/ansip224r1` |
| [`P-256`](https://neuromancer.sk/std/nist/P-256) | `1.2.840.10045.3.1.7` | `prime field` | O | `nist.P256()` | `secg/secp256r1`, `x962/prime256v1` |
| [`P-384`](https://neuromancer.sk/std/nist/P-384) | `1.3.132.0.34` | `prime field` | O | `nist.P384()` | `secg/secp384r1`, `x963/ansip384r1` |
| [`P-521`](https://neuromancer.sk/std/nist/P-521) | `1.3.132.0.35` | `prime field` | O | `nist.P521()` | `secg/secp521r1`, `x963/ansip521r1` |

### [NUMS](https://neuromancer.sk/std/nums/)

- Microsoft Nothing Up My Sleeve (NUMS) curves from: <https://eprint.iacr.org/2014/130> and <https://tools.ietf.org/html/draft-black-numscurves-02>

| Name | OID | Field Type | Generator | Function | As known |
|:-----|:----|:----------:|:---------:|:---------|:---------|
| [`ed-254-mont`](https://neuromancer.sk/std/nums/ed-254-mont) || `prime field` || `nums.Ed254Mont()` | |
| [`ed-255-mers`](https://neuromancer.sk/std/nums/ed-255-mers) || `prime field` || `nums.Ed255Mers()` | |
| [`ed-256-mont`](https://neuromancer.sk/std/nums/ed-256-mont) || `prime field` || `nums.Ed256Mont()` | |
| [`ed-382-mont`](https://neuromancer.sk/std/nums/ed-382-mont) || `prime field` || `nums.Ed382Mont()` | |
| [`ed-383-mers`](https://neuromancer.sk/std/nums/ed-383-mers) || `prime field` || `nums.Ed383Mers()` | |
| [`ed-384-mont`](https://neuromancer.sk/std/nums/ed-384-mont) || `prime field` || `nums.Ed384Mont()` | |
| [`ed-510-mont`](https://neuromancer.sk/std/nums/ed-510-mont) || `prime field` || `nums.Ed510Mont()` | |
| [`ed-511-mers`](https://neuromancer.sk/std/nums/ed-511-mers) || `prime field` || `nums.Ed511Mers()` | |
| [`ed-512-mont`](https://neuromancer.sk/std/nums/ed-512-mont) || `prime field` || `nums.Ed512Mont()` | |
| [`numsp256d1`](https://neuromancer.sk/std/nums/numsp256d1) || `prime field` | O | `nums.Numsp256d1()` | |
| [`numsp256t1`](https://neuromancer.sk/std/nums/numsp256t1) || `prime field` | O | `nums.Numsp256t1()` | |
| [`numsp384d1`](https://neuromancer.sk/std/nums/numsp384d1) || `prime field` | O | `nums.Numsp384d1()` | |
| [`numsp384t1`](https://neuromancer.sk/std/nums/numsp384t1) || `prime field` | O | `nums.Numsp384t1()` | |
| [`numsp512d1`](https://neuromancer.sk/std/nums/numsp512d1) || `prime field` | O | `nums.Numsp512d1()` | |
| [`numsp512t1`](https://neuromancer.sk/std/nums/numsp512t1) || `prime field` | O | `nums.Numsp512t1()` | |
| [`w-254-mont`](https://neuromancer.sk/std/nums/w-254-mont) || `prime field` || `nums.W254Mont()` | |
| [`w-255-mers`](https://neuromancer.sk/std/nums/w-255-mers) || `prime field` || `nums.W255Mers()` | |
| [`w-256-mont`](https://neuromancer.sk/std/nums/w-256-mont) || `prime field` || `nums.W256Mont()` | |
| [`w-382-mont`](https://neuromancer.sk/std/nums/w-382-mont) || `prime field` || `nums.W382Mont()` | |
| [`w-383-mers`](https://neuromancer.sk/std/nums/w-383-mers) || `prime field` || `nums.W383Mers()` | |
| [`w-384-mont`](https://neuromancer.sk/std/nums/w-384-mont) || `prime field` || `nums.W384Mont()` | |
| [`w-510-mont`](https://neuromancer.sk/std/nums/w-510-mont) || `prime field` || `nums.W510Mont()` | |
| [`w-511-mers`](https://neuromancer.sk/std/nums/w-511-mers) || `prime field` || `nums.W511Mers()` | |
| [`w-512-mont`](https://neuromancer.sk/std/nums/w-512-mont) || `prime field` || `nums.W512Mont()` | |

### [OSCAA](https://neuromancer.sk/std/oscaa/)

- http://gmssl.org/english.html

| Name | OID | Field Type | Generator | Function | As known |
|:-----|:----|:----------:|:---------:|:---------|:---------|
| [`SM2`](https://neuromancer.sk/std/oscaa/SM2) | `1.2.156.10197.1.301` | `prime field` | O | `oscaa.SM2()` | |

### [Oakley](https://neuromancer.sk/std/oakley/)

- Oakley groups from <https://tools.ietf.org/html/rfc2409> and <https://tools.ietf.org/html/rfc5114>

| Name | OID | Field Type | Generator | Function | As known |
|:-----|:----|:----------:|:---------:|:---------|:---------|
| [`192-bit Random ECP Group`](https://neuromancer.sk/std/oakley/192-bit%20Random%20ECP%20Group) || `prime field` | O | `oakley.Oakley192BitRandomECPGroup()` | |
| [`224-bit Random ECP Group`](https://neuromancer.sk/std/oakley/224-bit%20Random%20ECP%20Group) || `prime field` | O | `oakley.Oakley224BitRandomECPGroup()` | |
| [`256-bit Random ECP Group`](https://neuromancer.sk/std/oakley/256-bit%20Random%20ECP%20Group) || `prime field` | O | `oakley.Oakley256BitRandomECPGroup()` | |
| [`384-bit Random ECP Group`](https://neuromancer.sk/std/oakley/384-bit%20Random%20ECP%20Group) || `prime field` | O | `oakley.Oakley384BitRandomECPGroup()` | |
| [`521-bit Random ECP Group`](https://neuromancer.sk/std/oakley/521-bit%20Random%20ECP%20Group) || `prime field` | O | `oakley.Oakley521BitRandomECPGroup()` | |
| [`Oakley Group 3`](https://neuromancer.sk/std/oakley/Oakley%20Group%203) || `binary field` || `oakley.OakleyGroup3()` | |
| [`Oakley Group 4`](https://neuromancer.sk/std/oakley/Oakley%20Group%204) || `binary field` || `oakley.OakleyGroup4()` | |

### [SECG](https://neuromancer.sk/std/secg/)

- SEC 2: Recommended Elliptic Curve Domain Parameters version 2.0  January 27, 2010

| Name | OID | Field Type | Generator | Function | As known |
|:-----|:----|:----------:|:---------:|:---------|:---------|
| [`secp112r1`](https://neuromancer.sk/std/secg/secp112r1) | `1.3.132.0.6` | `prime field` | O | `secg.Secp112r1()` | `wtls/wap-wsg-idm-ecid-wtls6` |
| [`secp112r2`](https://neuromancer.sk/std/secg/secp112r2) | `1.3.132.0.7` | `prime field` | O | `secg.Secp112r2()` | |
| [`secp128r1`](https://neuromancer.sk/std/secg/secp128r1) | `1.3.132.0.28` | `prime field` | O | `secg.Secp128r1()` | |
| [`secp128r2`](https://neuromancer.sk/std/secg/secp128r2) | `1.3.132.0.29` | `prime field` | O | `secg.Secp128r2()` | |
| [`secp160k1`](https://neuromancer.sk/std/secg/secp160k1) | `1.3.132.0.9` | `prime field` | O | `secg.Secp160k1()` | `x963/ansip160k1` |
| [`secp160r1`](https://neuromancer.sk/std/secg/secp160r1) | `1.3.132.0.8` | `prime field` | O | `secg.Secp160r1()` | `wtls/wap-wsg-idm-ecid-wtls7`, `x963/ansip160r1` |
| [`secp160r2`](https://neuromancer.sk/std/secg/secp160r2) | `1.3.132.0.30` | `prime field` | O | `secg.Secp160r2()` | `x963/ansip160r2` |
| [`secp192k1`](https://neuromancer.sk/std/secg/secp192k1) | `1.3.132.0.31` | `prime field` | O | `secg.Secp192k1()` | `x963/ansip192k1` |
| [`secp192r1`](https://neuromancer.sk/std/secg/secp192r1) | `1.2.840.10045.3.1.1` | `prime field` | O | `secg.Secp192r1()` | `nist/P-192`, `x962/prime192v1` |
| [`secp224k1`](https://neuromancer.sk/std/secg/secp224k1) | `1.3.132.0.32` | `prime field` | O | `secg.Secp224k1()` | `x963/ansip224k1` |
| [`secp224r1`](https://neuromancer.sk/std/secg/secp224r1) | `1.3.132.0.33` | `prime field` | O | `secg.Secp224r1()` | `nist/P-224`, `wtls/wap-wsg-idm-ecid-wtls12`, `x963/ansip224r1` |
| [`secp256k1`](https://neuromancer.sk/std/secg/secp256k1) | `1.3.132.0.10` | `prime field` | O | `secg.Secp256k1()` | `x963/ansip256k1` |
| [`secp256r1`](https://neuromancer.sk/std/secg/secp256r1) | `1.2.840.10045.3.1.7` | `prime field` | O | `secg.Secp256r1()` | `nist/P-256`, `x962/prime256v1` |
| [`secp384r1`](https://neuromancer.sk/std/secg/secp384r1) | `1.3.132.0.34` | `prime field` | O | `secg.Secp384r1()` | `nist/P-384`, `x963/ansip384r1` |
| [`secp521r1`](https://neuromancer.sk/std/secg/secp521r1) | `1.3.132.0.35` | `prime field` | O | `secg.Secp521r1()` | `nist/P-521`, `x963/ansip521r1` |
| [`sect113r1`](https://neuromancer.sk/std/secg/sect113r1) | `1.3.132.0.4` | `binary field` | O | `secg.Sect113r1()` | `wtls/wap-wsg-idm-ecid-wtls4` |
| [`sect113r2`](https://neuromancer.sk/std/secg/sect113r2) | `1.3.132.0.5` | `binary field` | O | `secg.Sect113r2()` | |
| [`sect131r1`](https://neuromancer.sk/std/secg/sect131r1) | `1.3.132.0.22` | `binary field` | O | `secg.Sect131r1()` | |
| [`sect131r2`](https://neuromancer.sk/std/secg/sect131r2) | `1.3.132.0.23` | `binary field` | O | `secg.Sect131r2()` | |
| [`sect163k1`](https://neuromancer.sk/std/secg/sect163k1) | `1.3.132.0.1` | `binary field` | O | `secg.Sect163k1()` | `nist/K-163`, `x963/ansit163k1`, `wtls/wap-wsg-idm-ecid-wtls3` |
| [`sect163r1`](https://neuromancer.sk/std/secg/sect163r1) | `1.3.132.0.2` | `binary field` | O | `secg.Sect163r1()` | `x963/ansit163r1` |
| [`sect163r2`](https://neuromancer.sk/std/secg/sect163r2) | `1.3.132.0.15` | `binary field` | O | `secg.Sect163r2()` | `nist/B-163`, `x963/ansit163r2` |
| [`sect193r1`](https://neuromancer.sk/std/secg/sect193r1) | `1.3.132.0.24` | `binary field` | O | `secg.Sect193r1()` | `x963/ansit193r1` |
| [`sect193r2`](https://neuromancer.sk/std/secg/sect193r2) | `1.3.132.0.25` | `binary field` | O | `secg.Sect193r2()` | `x963/ansit193r2` |
| [`sect233k1`](https://neuromancer.sk/std/secg/sect233k1) | `1.3.132.0.26` | `binary field` | O | `secg.Sect233k1()` | `nist/K-233`, `wtls/wap-wsg-idm-ecid-wtls10`, `x963/ansit233k1` |
| [`sect233r1`](https://neuromancer.sk/std/secg/sect233r1) | `1.3.132.0.27` | `binary field` | O | `secg.Sect233r1()` | `nist/B-233`, `wtls/wap-wsg-idm-ecid-wtls11`, `x963/ansit233r1` |
| [`sect239k1`](https://neuromancer.sk/std/secg/sect239k1) | `1.3.132.0.3` | `binary field` | O | `secg.Sect239k1()` | `x963/ansit239k1` |
| [`sect283k1`](https://neuromancer.sk/std/secg/sect283k1) | `1.3.132.0.16` | `binary field` | O | `secg.Sect283k1()` | `nist/K-283`, `x963/ansit283k1` |
| [`sect283r1`](https://neuromancer.sk/std/secg/sect283r1) | `1.3.132.0.17` | `binary field` | O | `secg.Sect283r1()` | `nist/B-283`, `x963/ansit283r1` |
| [`sect409k1`](https://neuromancer.sk/std/secg/sect409k1) | `1.3.132.0.36` | `binary field` | O | `secg.Sect409k1()` | `nist/K-409`, `x963/ansit409k1` |
| [`sect409r1`](https://neuromancer.sk/std/secg/sect409r1) | `1.3.132.0.37` | `binary field` | O | `secg.Sect409r1()` | `nist/B-409`, `x963/ansit409r1` |
| [`sect571k1`](https://neuromancer.sk/std/secg/sect571k1) | `1.3.132.0.38` | `binary field` | O | `secg.Sect571k1()` | `nist/K-571`, `x963/ansit571k1` |
| [`sect571r1`](https://neuromancer.sk/std/secg/sect571r1) | `1.3.132.0.39` | `binary field` | O | `secg.Sect571r1()` | `nist/B-571`, `x963/ansit571r1` |

### [WTLS](https://neuromancer.sk/std/wtls/)

- Wireless Application Protocol - Wireless Transport Layer Security (WAP-WTLS) curves: <https://www.wapforum.org/tech/documents/WAP-199-WTLS-20000218-a.pdf>

| Name | OID | Field Type | Generator | Function | As known |
|:-----|:----|:----------:|:---------:|:---------|:---------|
| [`wap-wsg-idm-ecid-wtls1`](https://neuromancer.sk/std/wtls/wap-wsg-idm-ecid-wtls1) | `2.23.43.1.4.1` | `binary field` | O | `wtls.WapWsgIdmEcidWtls1()` | |
| [`wap-wsg-idm-ecid-wtls10`](https://neuromancer.sk/std/wtls/wap-wsg-idm-ecid-wtls10) | `2.23.43.1.4.10` | `binary field` | O | `wtls.WapWsgIdmEcidWtls10()` | `secg/sect233k1`, `nist/K-233`, `x963/ansit233k1` |
| [`wap-wsg-idm-ecid-wtls11`](https://neuromancer.sk/std/wtls/wap-wsg-idm-ecid-wtls11) | `2.23.43.1.4.11` | `binary field` | O | `wtls.WapWsgIdmEcidWtls11()` | `secg/sect233r1`, `nist/B-233`, `x963/ansit233r1` |
| [`wap-wsg-idm-ecid-wtls12`](https://neuromancer.sk/std/wtls/wap-wsg-idm-ecid-wtls12) | `2.23.43.1.4.12` | `prime field` | O | `wtls.WapWsgIdmEcidWtls12()` | `secg/secp224r1`, `nist/P-224`, `x963/ansip224r1` |
| [`wap-wsg-idm-ecid-wtls3`](https://neuromancer.sk/std/wtls/wap-wsg-idm-ecid-wtls3) | `2.23.43.1.4.3` | `binary field` | O | `wtls.WapWsgIdmEcidWtls3()` | `nist/K-163`, `secg/sect163k1`, `x963/ansit163k1` |
| [`wap-wsg-idm-ecid-wtls4`](https://neuromancer.sk/std/wtls/wap-wsg-idm-ecid-wtls4) | `2.23.43.1.4.4` | `binary field` | O | `wtls.WapWsgIdmEcidWtls4()` | `secg/sect113r1` |
| [`wap-wsg-idm-ecid-wtls5`](https://neuromancer.sk/std/wtls/wap-wsg-idm-ecid-wtls5) | `2.23.43.1.4.5` | `binary field` | O | `wtls.WapWsgIdmEcidWtls5()` | `x962/c2pnb163v1` |
| [`wap-wsg-idm-ecid-wtls6`](https://neuromancer.sk/std/wtls/wap-wsg-idm-ecid-wtls6) | `2.23.43.1.4.6` | `prime field` | O | `wtls.WapWsgIdmEcidWtls6()` | `secg/secp112r1` |
| [`wap-wsg-idm-ecid-wtls7`](https://neuromancer.sk/std/wtls/wap-wsg-idm-ecid-wtls7) | `2.23.43.1.4.7` | `prime field` | O | `wtls.WapWsgIdmEcidWtls7()` | `secg/secp160r1`, `x963/ansip160r1` |
| [`wap-wsg-idm-ecid-wtls8`](https://neuromancer.sk/std/wtls/wap-wsg-idm-ecid-wtls8) | `2.23.43.1.4.8` | `prime field` | O | `wtls.WapWsgIdmEcidWtls8()` | |
| [`wap-wsg-idm-ecid-wtls9`](https://neuromancer.sk/std/wtls/wap-wsg-idm-ecid-wtls9) | `2.23.43.1.4.9` | `prime field` | O | `wtls.WapWsgIdmEcidWtls9()` | |

### [other](https://neuromancer.sk/std/other/)

- An assortment of some other curves.

| Name | OID | Field Type | Generator | Function | As known |
|:-----|:----|:----------:|:---------:|:---------|:---------|
| [`BADA55-R-256`](https://neuromancer.sk/std/other/BADA55-R-256) || `prime field` || `other.BADA55R256()` | |
| [`BADA55-VPR-224`](https://neuromancer.sk/std/other/BADA55-VPR-224) || `prime field` || `other.BADA55VPR224()` | |
| [`BADA55-VPR2-224`](https://neuromancer.sk/std/other/BADA55-VPR2-224) || `prime field` || `other.BADA55VPR2_224()` | |
| [`BADA55-VR-224`](https://neuromancer.sk/std/other/BADA55-VR-224) || `prime field` || `other.BADA55VR224()` | |
| [`BADA55-VR-256`](https://neuromancer.sk/std/other/BADA55-VR-256) || `prime field` || `other.BADA55VR256()` | |
| [`BADA55-VR-384`](https://neuromancer.sk/std/other/BADA55-VR-384) || `prime field` || `other.BADA55VR384()` | |
| [`Curve1174`](https://neuromancer.sk/std/other/Curve1174) || `prime field` | O | `other.Curve1174()` | |
| [`Curve22103`](https://neuromancer.sk/std/other/Curve22103) || `prime field` | O | `other.Curve22103()` | |
| [`Curve25519`](https://neuromancer.sk/std/other/Curve25519) || `prime field` | O | `other.Curve25519()` | |
| [`Curve383187`](https://neuromancer.sk/std/other/Curve383187) || `prime field` | O | `other.Curve383187()` | |
| [`Curve41417`](https://neuromancer.sk/std/other/Curve41417) || `prime field` | O | `other.Curve41417()` | |
| [`Curve4417`](https://neuromancer.sk/std/other/Curve4417) || `prime field` | O | `other.Curve4417()` | |
| [`Curve448`](https://neuromancer.sk/std/other/Curve448) || `prime field` | O | `other.Curve448()` | |
| [`Curve67254`](https://neuromancer.sk/std/other/Curve67254) || `prime field` | O | `other.Curve67254()` | |
| [`E-222`](https://neuromancer.sk/std/other/E-222) || `prime field` | O | `other.E222()` | |
| [`E-382`](https://neuromancer.sk/std/other/E-382) || `prime field` | O | `other.E382()` | |
| [`E-521`](https://neuromancer.sk/std/other/E-521) || `prime field` | O | `other.E521()` | |
| [`Ed25519`](https://neuromancer.sk/std/other/Ed25519) || `prime field` | O | `other.Ed25519()` | |
| [`Ed448`](https://neuromancer.sk/std/other/Ed448) || `prime field` | O | `other.Ed448()` | |
| [`Ed448-Goldilocks`](https://neuromancer.sk/std/other/Ed448-Goldilocks) || `prime field` | O | `other.Ed448Goldilocks()` | |
| [`Fp224BN`](https://neuromancer.sk/std/other/Fp224BN) || `prime field` | O | `other.Fp224BN()` | |
| [`Fp254BNa`](https://neuromancer.sk/std/other/Fp254BNa) || `prime field` | O | `other.Fp254BNa()` | |
| [`Fp254BNb`](https://neuromancer.sk/std/other/Fp254BNb) || `prime field` | O | `other.Fp254BNb()` | `bn/bn254` |
| [`Fp256BN`](https://neuromancer.sk/std/other/Fp256BN) || `prime field` | O | `other.Fp256BN()` | |
| [`Fp384BN`](https://neuromancer.sk/std/other/Fp384BN) || `prime field` | O | `other.Fp384BN()` | |
| [`Fp512BN`](https://neuromancer.sk/std/other/Fp512BN) || `prime field` | O | `other.Fp512BN()` | |
| [`JubJub`](https://neuromancer.sk/std/other/JubJub) || `prime field` | O | `other.JubJub()` | |
| [`M-221`](https://neuromancer.sk/std/other/M-221) || `prime field` | O | `other.M221()` | |
| [`M-383`](https://neuromancer.sk/std/other/M-383) || `prime field` | O | `other.M383()` | |
| [`M-511`](https://neuromancer.sk/std/other/M-511) || `prime field` | O | `other.M511()` | |
| [`MDC201601`](https://neuromancer.sk/std/other/MDC201601) || `prime field` | O | `other.MDC201601()` | |
| [`Pallas`](https://neuromancer.sk/std/other/Pallas) || `prime field` | O | `other.Pallas()` | |
| [`Tom-256`](https://neuromancer.sk/std/other/Tom-256) || `prime field` | O | `other.Tom256()` | |
| [`Tom-384`](https://neuromancer.sk/std/other/Tom-384) || `prime field` || `other.Tom384()` | |
| [`Tom-521`](https://neuromancer.sk/std/other/Tom-521) || `prime field` || `other.Tom521()` | |
| [`Tweedledee`](https://neuromancer.sk/std/other/Tweedledee) || `prime field` | O | `other.Tweedledee()` | |
| [`Tweedledum`](https://neuromancer.sk/std/other/Tweedledum) || `prime field` | O | `other.Tweedledum()` | |
| [`Vesta`](https://neuromancer.sk/std/other/Vesta) || `prime field` | O | `other.Vesta()` | |
| [`ssc-160`](https://neuromancer.sk/std/other/ssc-160) || `prime field` || `other.Ssc160()` | |
| [`ssc-192`](https://neuromancer.sk/std/other/ssc-192) || `prime field` || `other.Ssc192()` | |
| [`ssc-224`](https://neuromancer.sk/std/other/ssc-224) || `prime field` || `other.Ssc224()` | |
| [`ssc-256`](https://neuromancer.sk/std/other/ssc-256) || `prime field` || `other.Ssc256()` | |
| [`ssc-288`](https://neuromancer.sk/std/other/ssc-288) || `prime field` || `other.Ssc288()` | |
| [`ssc-320`](https://neuromancer.sk/std/other/ssc-320) || `prime field` || `other.Ssc320()` | |
| [`ssc-384`](https://neuromancer.sk/std/other/ssc-384) || `prime field` || `other.Ssc384()` | |
| [`ssc-512`](https://neuromancer.sk/std/other/ssc-512) || `prime field` || `other.Ssc512()` | |

- Those curves are not avilable in this namespace:
  | Name | Reason |
  |:-----|:-------|
  | [`FourQ`](https://neuromancer.sk/std/other/FourQ) | unsupported curve type |
  | [`Fp254n2BNa`](https://neuromancer.sk/std/other/Fp254n2BNa) | unsupported curve type |

<!-- curves end -->
