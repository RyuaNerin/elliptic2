package field_test

import (
	"math/big"
	"testing"

	. "github.com/RyuaNerin/elliptic2/internal"
	. "github.com/RyuaNerin/elliptic2/internal/field"
)

func gf2m(m *GF2m) {}

func TestGF2mInvValidation(t *testing.T) { testInvValidation(t, gf2m, gf2mModulus) }
func TestGF2mSqrValidation(t *testing.T) { testSqrValidation(t, gf2m, gf2mModulus) }

func TestGF2mAdd(t *testing.T) { tArg(t, gf2m, gf2mTC, tcArgIndexes.add, add, addWant) }
func TestGF2mMul(t *testing.T) { tArg(t, gf2m, gf2mTC, tcArgIndexes.mul, mul, mulWant) }
func TestGF2mSqr(t *testing.T) { tArg(t, gf2m, gf2mTC, tcArgIndexes.sqr, sqr, sqrWant) }
func TestGF2mInv(t *testing.T) { tArg(t, gf2m, gf2mTC, tcArgIndexes.inv, inv, invWant) }

func BenchmarkGF2mAdd(b *testing.B) { bArg(b, gf2m, gf2mModulus, add) }
func BenchmarkGF2mMul(b *testing.B) { bArg(b, gf2m, gf2mModulus, mul) }
func BenchmarkGF2mSqr(b *testing.B) { bArg(b, gf2m, gf2mModulus, sqr) }
func BenchmarkGF2mInv(b *testing.B) { bArg(b, gf2m, gf2mModulus, inv) }

var gf2mModulus = []*Modulus{
	NewGF2mModulus(GF2mPolynomials(163, 7, 6, 3, 0)),
	NewGF2mModulus(GF2mPolynomials(233, 74, 0)),
	NewGF2mModulus(GF2mPolynomials(283, 12, 7, 5, 0)),
	NewGF2mModulus(GF2mPolynomials(409, 87, 0)),
	NewGF2mModulus(GF2mPolynomials(571, 10, 5, 2, 0)),
}

var gf2mTC = []testCase{
	{
		modulus: gf2mModulus[0],
		arg: [3]*big.Int{
			HI(`172fa6b1260`),
			HI(`155a0228596f9a161a211`),
			HI(`2535fe684404a6e5a1dd54e6fb16c0ad6dfb300c6`),
		},
		add: [7]*big.Int{
			HI(`0`),
			HI(`0`),
			HI(`0`),
			HI(`155a02285978b5b0ab071`),
			HI(`2535fe684404a6e5a1dd54e6fb16c0ba425d812a6`),
			HI(`2535fe684404a6e5a1dd41bcf93e99c2f7ed2a2d7`),
			HI(`2535fe684404a6e5a1dd41bcf93e99d5d84b9b0b7`),
		},
		mul: [7]*big.Int{
			HI(`115045544144501041400`),
			HI(`11111440004044011411455414401140144040101`),
			HI(`41b0838d55a71ddc39245c539206a48099ee4b7ca`),
			HI(`13d454adfbaa468a7b27f12c7e2f460`),
			HI(`190a45ab0477205ea33f923205016996618d2876c`),
			HI(`adabb1ef6f2aa100e082b3ace84973b44b30cd3f`),
			HI(`6d3d6bdbd7c578dd5ecbf1eed038fad7cefefad06`),
		},
		sqr: [3]*big.Int{
			HI(`115045544144501041400`),
			HI(`11111440004044011411455414401140144040101`),
			HI(`41b0838d55a71ddc39245c539206a48099ee4b7ca`),
		},
		inv: [3]*big.Int{
			HI(`3a5ffde7e979d7bd92481fd39a13b698082a0ab6e`),
			HI(`88b0398ca3b7396244e4c7323e1a94f021a607c9`),
			HI(`7a8851f549ad453c90814f14d4136b2fe7d9d5d63`),
		},
	},
	{
		modulus: gf2mModulus[0],
		arg: [3]*big.Int{
			HI(`3ce97773bf`),
			HI(`240ab6398b29749f319f`),
			HI(`172384e0c602d9fae9bb5c4ff3ddf1d3d5bfe516f`),
		},
		add: [7]*big.Int{
			HI(`0`),
			HI(`0`),
			HI(`0`),
			HI(`240ab6398b159de84220`),
			HI(`172384e0c602d9fae9bb5c4ff3ddf1d01b28922d0`),
			HI(`172384e0c602d9fae9bb5e0f58be696142f6160f0`),
			HI(`172384e0c602d9fae9bb5e0f58be69628c616134f`),
		},
		mul: [7]*big.Int{
			HI(`5505441151515054555`),
			HI(`410004445140541404504411510415505014155`),
			HI(`19fd4aebc90efd68ef7d1661467c4de6cc612813c`),
			HI(`76f0f020ebeef210d324018bb8cb5`),
			HI(`4640f0102ed02a2773d43c8f2fd200e269048ecd1`),
			HI(`93046dd449c9da9db93dcc1e1043a1ef798a4af3`),
			HI(`479f20885cc62f896d2c618001b0c44769e5bebc5`),
		},
		sqr: [3]*big.Int{
			HI(`5505441151515054555`),
			HI(`410004445140541404504411510415505014155`),
			HI(`19fd4aebc90efd68ef7d1661467c4de6cc612813c`),
		},
		inv: [3]*big.Int{
			HI(`2227628f137d687499bf4f70af99a2cbe3e08aea6`),
			HI(`2ee26b87a9a4222993e4074bedb551254510154e4`),
			HI(`51ff2c1c9a796be7a8099cfb4e8e80e395f23fbcd`),
		},
	},
	{
		modulus: gf2mModulus[0],
		arg: [3]*big.Int{
			HI(`1a1914747b5`),
			HI(`128ce84d8de8649ba5b45`),
			HI(`5dfd0f2a09417309b77d3be8c039c6abafe5a1414`),
		},
		add: [7]*big.Int{
			HI(`0`),
			HI(`0`),
			HI(`0`),
			HI(`128ce84d8df27d8fd1cf0`),
			HI(`5dfd0f2a09417309b77d3be8c039c6b1b6f1d53a1`),
			HI(`5dfd0f2a09417309b77d296428744b43cb7e04f51`),
			HI(`5dfd0f2a09417309b77d296428744b59d26a708e4`),
		},
		mul: [7]*big.Int{
			HI(`144014101101510154511`),
			HI(`10440505440105140515440141041454411451011`),
			HI(`1a1d42377a7c7edd644f65d8b88f246504f587749`),
			HI(`1981f566b50e84c7773aa1d6b5ad321`),
			HI(`2dd4fb07e30c899ece68ac727d583ada9e546513f`),
			HI(`526f4fec987dbf9ffac632807a067d59f7b64e62d`),
			HI(`3156976f13deee04e093c4b0881ca124c7ce0fd09`),
		},
		sqr: [3]*big.Int{
			HI(`144014101101510154511`),
			HI(`10440505440105140515440141041454411451011`),
			HI(`1a1d42377a7c7edd644f65d8b88f246504f587749`),
		},
		inv: [3]*big.Int{
			HI(`6f8abe4f0a8b0eecb4200acf8101156f71cbb0a03`),
			HI(`59b3d28935af674840b30c69423a7a25f8b29b3bb`),
			HI(`41c213f5c588d67d51ddee328835417bcf8651174`),
		},
	},
	{
		modulus: gf2mModulus[0],
		arg: [3]*big.Int{
			HI(`129e05aff8c`),
			HI(`25df64a20e12dc6b582bf`),
			HI(`f802eb3d334ea9a9c1c934ed620ec2d444fd10b5`),
		},
		add: [7]*big.Int{
			HI(`0`),
			HI(`0`),
			HI(`0`),
			HI(`25df64a20e00426ef7d33`),
			HI(`f802eb3d334ea9a9c1c934ed620ec3fda4a7ef39`),
			HI(`f802eb3d334ea9a9c1cb691b282e23f98248920a`),
			HI(`f802eb3d334ea9a9c1cb691b282e22d062126d86`),
		},
		mul: [7]*big.Int{
			HI(`104415400114455554050`),
			HI(`41151551410440400540104515014451140044555`),
			HI(`6a4f854322cb1da370dd864a8bac75f87fb1eb959`),
			HI(`207085a3c025322fe2252ebac19d584`),
			HI(`8872d992d4d7e96c7ad85cbc95f8881b22481b78`),
			HI(`19ba5b85eaf998d3fe5a1ea5c6d80292483e1c1db`),
			HI(`4da037760c7f042b488d4b1e5b28952a2118fc9bf`),
		},
		sqr: [3]*big.Int{
			HI(`104415400114455554050`),
			HI(`41151551410440400540104515014451140044555`),
			HI(`6a4f854322cb1da370dd864a8bac75f87fb1eb959`),
		},
		inv: [3]*big.Int{
			HI(`5d89b2900a5e87742f5bf65a20b34207270567df1`),
			HI(`2afd6ea8fef35782dd8a94b87285af89a57b8e819`),
			HI(`7e7fc6b90c2643cef32dca9dd96f672c50ed59101`),
		},
	},
	{
		modulus: gf2mModulus[0],
		arg: [3]*big.Int{
			HI(`266c3e3291`),
			HI(`180dbfdf90df11b4e4915`),
			HI(`14f8c35b689ee87c7f94df07826c2c79a185d1371`),
		},
		add: [7]*big.Int{
			HI(`0`),
			HI(`0`),
			HI(`0`),
			HI(`180dbfdf90dd777707b84`),
			HI(`14f8c35b689ee87c7f94df07826c2c7bc746321e0`),
			HI(`14f8c35b689ee87c7f94c70a3db3bca6b03135a64`),
			HI(`14f8c35b689ee87c7f94c70a3db3bca4d6f2d68f5`),
		},
		mul: [7]*big.Int{
			HI(`4141450055405044101`),
			HI(`14000514555515541005155010145105410410111`),
			HI(`1e0a1fa9a783944549e43d87cc8a86ab900f5a180`),
			HI(`3543aebae7c8f1a46822a9e88988c5`),
			HI(`6c5f0ee46b0a6c3eeb9c67fd4da35c4066c01ca94`),
			HI(`577b0497fa0d9b793e1449a7d8ba6340057ad3383`),
			HI(`63b03f1d762368210e6e8e102cbce5872e7f7acb2`),
		},
		sqr: [3]*big.Int{
			HI(`4141450055405044101`),
			HI(`14000514555515541005155010145105410410111`),
			HI(`1e0a1fa9a783944549e43d87cc8a86ab900f5a180`),
		},
		inv: [3]*big.Int{
			HI(`30830500d125333c3561fc66614d2a90a2df76c99`),
			HI(`7d3f2ca4ca4f361f6530cd05dc4dff9ef42f8c6c9`),
			HI(`7440acf4099c78f0db91cb209fba22bfa44461a74`),
		},
	},
}
