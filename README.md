[![PkgGoDev](https://pkg.go.dev/badge/github.com/RyuaNerin/elliptic2)](https://pkg.go.dev/github.com/RyuaNerin/elliptic2)

# elliptic2

- `elliptic2` supports ***Elliptic curves over binary fields*** in go

- based on [kurtbrose/elliptic.py](https://gist.github.com/kurtbrose/4423605) ([fork](https://gist.github.com/RyuaNerin/777416c68a5419102d4833ee54cb4426))

## Warning

- **Not all curves have been validated.**

- `elliptic2` uses `elliptic.CurveParams` for compatibility with `crypto/elliptic` package.
    - But do not use the functions of `elliptic.CurveParams`. It will be panic.

- Multiple invocations of this function will return the same value like `crypto/elliptic`
    - so it can be used for equality checks and switch statements.

## Curves

### [ANSI x9.62](https://neuromancer.sk/std/x962/)

- ANSI x9.62 example curves.

| Name | OID | Field Type | Function | Tested | As known |
|:----:|:---:|:----------:|:--------:|:------:|:--------:|
| [`c2pnb176w1`](https://neuromancer.sk/std/x962/c2pnb176w1) | | `Binary` | `x962.C2pnb176w1()` | |  |
| [`c2pnb163v1`](https://neuromancer.sk/std/x962/c2pnb163v1) | `1.2.840.10045.3.0.1` | `Binary` | `x962.C2pnb163v1()` | | `wtls/wap-wsg-idm-ecid-wtls5` |
| [`c2pnb163v2`](https://neuromancer.sk/std/x962/c2pnb163v2) | | `Binary` | `x962.C2pnb163v2()` | |  |
| [`c2pnb163v3`](https://neuromancer.sk/std/x962/c2pnb163v3) | | `Binary` | `x962.C2pnb163v3()` | |  |
| [`c2pnb208w1`](https://neuromancer.sk/std/x962/c2pnb208w1) | | `Binary` | `x962.C2pnb208w1()` | |  |
| [`c2tnb191v3`](https://neuromancer.sk/std/x962/c2tnb191v3) | | `Binary` | `x962.C2tnb191v3()` | |  |
| [`c2tnb191v2`](https://neuromancer.sk/std/x962/c2tnb191v2) | | `Binary` | `x962.C2tnb191v2()` | |  |
| [`c2tnb191v1`](https://neuromancer.sk/std/x962/c2tnb191v1) | | `Binary` | `x962.C2tnb191v1()` | |  |
| [`c2tnb239v3`](https://neuromancer.sk/std/x962/c2tnb239v3) | | `Binary` | `x962.C2tnb239v3()` | |  |
| [`c2tnb239v2`](https://neuromancer.sk/std/x962/c2tnb239v2) | | `Binary` | `x962.C2tnb239v2()` | |  |
| [`c2tnb239v1`](https://neuromancer.sk/std/x962/c2tnb239v1) | | `Binary` | `x962.C2tnb239v1()` | |  |
| [`c2pnb272w1`](https://neuromancer.sk/std/x962/c2pnb272w1) | | `Binary` | `x962.C2pnb272w1()` | |  |
| [`c2pnb304w1`](https://neuromancer.sk/std/x962/c2pnb304w1) | | `Binary` | `x962.C2pnb304w1()` | |  |
| [`c2pnb368w1`](https://neuromancer.sk/std/x962/c2pnb368w1) | | `Binary` | `x962.C2pnb368w1()` | |  |
| [`c2tnb359v1`](https://neuromancer.sk/std/x962/c2tnb359v1) | | `Binary` | `x962.C2tnb359v1()` | |  |
| [`c2tnb431r1`](https://neuromancer.sk/std/x962/c2tnb431r1) | | `Binary` | `x962.C2tnb431r1()` | |  |
| [`c2onb191v4`](https://neuromancer.sk/std/x962/c2onb191v4) | | `Binary` | `x962.C2onb191v4()` | |  |
| [`c2onb191v5`](https://neuromancer.sk/std/x962/c2onb191v5) | | `Binary` | `x962.C2onb191v5()` | |  |
| [`c2onb239v4`](https://neuromancer.sk/std/x962/c2onb239v4) | | `Binary` | `x962.C2onb239v4()` | |  |
| [`c2onb239v5`](https://neuromancer.sk/std/x962/c2onb239v5) | | `Binary` | `x962.C2onb239v5()` | |  |

### [ANSI X9.63](https://neuromancer.sk/std/x963/)

- ANSI x9.63 example curves.

| Name | OID | Field Type | Function | Tested | As known |
|:----:|:---:|:----------:|:--------:|:------:|:--------:|
| [`ansit163k1`](https://neuromancer.sk/std/x963/ansit163k1) | | `Binary` | `x963.Ansit163k1()` | |  |
| [`ansit163r1`](https://neuromancer.sk/std/x963/ansit163r1) | `1.3.132.0.2` | `Binary` | `x963.Ansit163r1()` | | `secg/sect163r1` |
| [`ansit163r2`](https://neuromancer.sk/std/x963/ansit163r2) | `1.3.132.0.15` | `Binary` | `x963.Ansit163r2()` | | `secg/sect163r2`, `nist/B-163` |
| [`ansit193r1`](https://neuromancer.sk/std/x963/ansit193r1) | `1.3.132.0.24` | `Binary` | `x963.Ansit193r1()` | | `secg/sect193r1` |
| [`ansit193r2`](https://neuromancer.sk/std/x963/ansit193r2) | `1.3.132.0.25` | `Binary` | `x963.Ansit193r2()` | | `secg/sect193r2` |
| [`ansit233k1`](https://neuromancer.sk/std/x963/ansit233k1) | `1.3.132.0.26` | `Binary` | `x963.Ansit233k1()` | | `secg/sect233k1`, `nist/K-233`, `wtls/wap-wsg-idm-ecid-wtls10` |
| [`ansit233r1`](https://neuromancer.sk/std/x963/ansit233r1) | `1.3.132.0.27` | `Binary` | `x963.Ansit233r1()` | | `wtls/wap-wsg-idm-ecid-wtls11`, `nist/B-233`, `secg/sect233r1` |
| [`ansit239k1`](https://neuromancer.sk/std/x963/ansit239k1) | `1.3.132.0.3` | `Binary` | `x963.Ansit239k1()` | | `secg/sect239k1` |
| [`ansit283k1`](https://neuromancer.sk/std/x963/ansit283k1) | `1.3.132.0.16` | `Binary` | `x963.Ansit283k1()` | | `nist/K-283`, `secg/sect283k1` |
| [`ansit283r1`](https://neuromancer.sk/std/x963/ansit283r1) | `1.3.132.0.17` | `Binary` | `x963.Ansit283r1()` | | `nist/B-283`, `secg/sect283r1` |
| [`ansit409k1`](https://neuromancer.sk/std/x963/ansit409k1) | `1.3.132.0.36` | `Binary` | `x963.Ansit409k1()` | | `nist/K-409`, `secg/sect409k1` |
| [`ansit409r1`](https://neuromancer.sk/std/x963/ansit409r1) | `1.3.132.0.37` | `Binary` | `x963.Ansit409r1()` | | `nist/B-409`, `secg/sect409r1` |
| [`ansit571k1`](https://neuromancer.sk/std/x963/ansit571k1) | `1.3.132.0.38` | `Binary` | `x963.Ansit571k1()` | | `nist/K-571`, `secg/sect571k1` |
| [`ansit571r1`](https://neuromancer.sk/std/x963/ansit571r1) | `1.3.132.0.39` | `Binary` | `x963.Ansit571r1()` | | `nist/B-571`, `secg/sect571r1` |

### [NIST](https://neuromancer.sk/std/nist/)

- RECOMMENDED ELLIPTIC CURVES FOR FEDERAL GOVERNMENT USE  July 1999

| Name | OID | Field Type | Function | Tested | As known |
|:----:|:---:|:----------:|:--------:|:------:|:--------:|
| [`K-163`](https://neuromancer.sk/std/nist/K-163) | `1.3.132.0.1` | `Binary` | `nist.K163()` | | `secg/sect163k1` |
| [`B-163`](https://neuromancer.sk/std/nist/B-163) | `1.3.132.0.15` | `Binary` | `nist.B163()` | | `secg/sect163r2`, `x963/ansit163r2` |
| [`K-233`](https://neuromancer.sk/std/nist/K-233) | `1.3.132.0.26` | `Binary` | `nist.K233()` | | `secg/sect233k1`, `wtls/wap-wsg-idm-ecid-wtls10`, `x963/ansit233k1` |
| [`B-233`](https://neuromancer.sk/std/nist/B-233) | `1.3.132.0.27` | `Binary` | `nist.B233()` | | `secg/sect233r1`, `wtls/wap-wsg-idm-ecid-wtls11`, `x963/ansit233r1` |
| [`K-283`](https://neuromancer.sk/std/nist/K-283) | `1.3.132.0.16` | `Binary` | `nist.K283()` | | `secg/sect283k1`, `x963/ansit283k1` |
| [`B-283`](https://neuromancer.sk/std/nist/B-283) | `1.3.132.0.17` | `Binary` | `nist.B283()` | | `secg/sect283r1`, `x963/ansit283r1` |
| [`K-409`](https://neuromancer.sk/std/nist/K-409) | `1.3.132.0.36` | `Binary` | `nist.K409()` | | `secg/sect409k1`, `x963/ansit409k1` |
| [`B-409`](https://neuromancer.sk/std/nist/B-409) | `1.3.132.0.37` | `Binary` | `nist.B409()` | | `secg/sect409r1`, `x963/ansit409r1` |
| [`K-571`](https://neuromancer.sk/std/nist/K-571) | `1.3.132.0.38` | `Binary` | `nist.K571()` | | `secg/sect571k1`, `x963/ansit571k1` |
| [`B-571`](https://neuromancer.sk/std/nist/B-571) | `1.3.132.0.39` | `Binary` | `nist.B571()` | | `secg/sect571r1`, `x963/ansit571r1` |

### [SECG](https://neuromancer.sk/std/secg/)

- SEC 2: Recommended Elliptic Curve Domain Parameters version 2.0  January 27, 2010

| Name | OID | Field Type | Function | Tested | As known |
|:----:|:---:|:----------:|:--------:|:------:|:--------:|
| [`sect113r1`](https://neuromancer.sk/std/secg/sect113r1) | `1.3.132.0.4` | `Binary` | `secg.Sect113r1()` | | `wtls/wap-wsg-idm-ecid-wtls4` |
| [`sect113r2`](https://neuromancer.sk/std/secg/sect113r2) | | `Binary` | `secg.Sect113r2()` | |  |
| [`sect131r1`](https://neuromancer.sk/std/secg/sect131r1) | | `Binary` | `secg.Sect131r1()` | |  |
| [`sect131r2`](https://neuromancer.sk/std/secg/sect131r2) | | `Binary` | `secg.Sect131r2()` | |  |
| [`sect163k1`](https://neuromancer.sk/std/secg/sect163k1) | `1.3.132.0.1` | `Binary` | `secg.Sect163k1()` | | `nist/K-163`, `x963/ansit163k1`, `wtls/wap-wsg-idm-ecid-wtls3` |
| [`sect163r1`](https://neuromancer.sk/std/secg/sect163r1) | `1.3.132.0.2` | `Binary` | `secg.Sect163r1()` | | `x963/ansit163r1` |
| [`sect163r2`](https://neuromancer.sk/std/secg/sect163r2) | `1.3.132.0.15` | `Binary` | `secg.Sect163r2()` | | `nist/B-163`, `x963/ansit163r2` |
| [`sect193r1`](https://neuromancer.sk/std/secg/sect193r1) | `1.3.132.0.24` | `Binary` | `secg.Sect193r1()` | | `x963/ansit193r1` |
| [`sect193r2`](https://neuromancer.sk/std/secg/sect193r2) | `1.3.132.0.25` | `Binary` | `secg.Sect193r2()` | | `x963/ansit193r2` |
| [`sect233k1`](https://neuromancer.sk/std/secg/sect233k1) | `1.3.132.0.26` | `Binary` | `secg.Sect233k1()` | | `nist/K-233`, `wtls/wap-wsg-idm-ecid-wtls10`, `x963/ansit233k1` |
| [`sect233r1`](https://neuromancer.sk/std/secg/sect233r1) | `1.3.132.0.27` | `Binary` | `secg.Sect233r1()` | | `nist/B-233`, `wtls/wap-wsg-idm-ecid-wtls11`, `x963/ansit233r1` |
| [`sect239k1`](https://neuromancer.sk/std/secg/sect239k1) | `1.3.132.0.3` | `Binary` | `secg.Sect239k1()` | | `x963/ansit239k1` |
| [`sect283k1`](https://neuromancer.sk/std/secg/sect283k1) | `1.3.132.0.16` | `Binary` | `secg.Sect283k1()` | | `nist/K-283`, `x963/ansit283k1` |
| [`sect283r1`](https://neuromancer.sk/std/secg/sect283r1) | `1.3.132.0.17` | `Binary` | `secg.Sect283r1()` | | `nist/B-283`, `x963/ansit283r1` |
| [`sect409k1`](https://neuromancer.sk/std/secg/sect409k1) | `1.3.132.0.36` | `Binary` | `secg.Sect409k1()` | | `nist/K-409`, `x963/ansit409k1` |
| [`sect409r1`](https://neuromancer.sk/std/secg/sect409r1) | `1.3.132.0.37` | `Binary` | `secg.Sect409r1()` | | `nist/B-409`, `x963/ansit409r1` |
| [`sect571k1`](https://neuromancer.sk/std/secg/sect571k1) | `1.3.132.0.38` | `Binary` | `secg.Sect571k1()` | | `nist/K-571`, `x963/ansit571k1` |
| [`sect571r1`](https://neuromancer.sk/std/secg/sect571r1) | `1.3.132.0.39` | `Binary` | `secg.Sect571r1()` | | `nist/B-571`, `x963/ansit571r1` |

### [WTLS](https://neuromancer.sk/std/WTLS/)

- Wireless Application Protocol - Wireless Transport Layer Security (WAP-WTLS) curves: <https://www.wapforum.org/tech/documents/WAP-199-WTLS-20000218-a.pdf>

| Name | OID | Field Type | Function | Tested | As known |
|:----:|:---:|:----------:|:--------:|:------:|:--------:|
| [`wap-wsg-idm-ecid-wtls1`](https://neuromancer.sk/std/WTLS/wap-wsg-idm-ecid-wtls1) | | `Binary` | `WTLS.WapWsgIdmEcidWtls1()` | |  |
| [`wap-wsg-idm-ecid-wtls3`](https://neuromancer.sk/std/WTLS/wap-wsg-idm-ecid-wtls3) | `2.23.43.1.4.3` | `Binary` | `WTLS.WapWsgIdmEcidWtls3()` | | `nist/K-163`, `secg/sect163k1`, `x963/ansit163k1` |
| [`wap-wsg-idm-ecid-wtls4`](https://neuromancer.sk/std/WTLS/wap-wsg-idm-ecid-wtls4) | `2.23.43.1.4.4` | `Binary` | `WTLS.WapWsgIdmEcidWtls4()` | | `secg/sect113r1` |
| [`wap-wsg-idm-ecid-wtls5`](https://neuromancer.sk/std/WTLS/wap-wsg-idm-ecid-wtls5) | `2.23.43.1.4.5` | `Binary` | `WTLS.WapWsgIdmEcidWtls5()` | | `x962/c2pnb163v1` |
| [`wap-wsg-idm-ecid-wtls10`](https://neuromancer.sk/std/WTLS/wap-wsg-idm-ecid-wtls10) | `2.23.43.1.4.10` | `Binary` | `WTLS.WapWsgIdmEcidWtls10()` | | `secg/sect233k1`, `nist/K-233`, `x963/ansit233k1` |
| [`wap-wsg-idm-ecid-wtls11`](https://neuromancer.sk/std/WTLS/wap-wsg-idm-ecid-wtls11) | `2.23.43.1.4.11` | `Binary` | `WTLS.WapWsgIdmEcidWtls11()` | | `secg/sect233r1`, `nist/B-233`, `x963/ansit233r1` |

## [LICENSE](/LICENSE)

- MIT License
