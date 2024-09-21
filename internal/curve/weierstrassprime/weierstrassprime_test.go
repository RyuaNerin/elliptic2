package weierstrassprime_test

import (
	"testing"

	"github.com/RyuaNerin/elliptic2"
	. "github.com/RyuaNerin/elliptic2/internal"
	. "github.com/RyuaNerin/elliptic2/internal/curve"
	. "github.com/RyuaNerin/elliptic2/internal/curve/weierstrassprime"
	. "github.com/RyuaNerin/elliptic2/internal/curvetesting"
	. "github.com/RyuaNerin/elliptic2/internal/field"
)

func TestInOnCurve(t *testing.T) { TcIsOnCurve(t, tc) }
func TestComputeY(t *testing.T)  { TcComputeY(t, tc) }

var (
	c = []elliptic2.Curve{
		cA,
		cA0,
		cAm3,
	}

	cA = NewCurve(Build(&CurveParams{
		Name:    "secp112r2",
		BitSize: 112,
		P:       NewGFpModulus(HI("db7c2abf62e35e668076bead208b")),
		A:       ParseGFpHex("6127c24c05f38a0aaaf65c0ef02c"),
		B:       ParseGFpHex("51def1815db5ed74fcc34c85d709"),
		N:       ParseGFpHex("36df0aafd8b8d7597ca10520d04b"),
		Gx:      HI("4ba30ab5e892b4e1649dd0928643"),
		Gy:      HI("adcd46f5882e3747def36e956e97"),
	}))
	cA0 = NewCurve(Build(&CurveParams{
		Name:    `wap-wsg-idm-ecid-wtls8`,
		BitSize: 112,
		P:       NewGFpModulus(HI(`fffffffffffffffffffffffffde7`)),
		A:       ParseGFpHex(``),
		B:       ParseGFpHex(`3`),
		N:       ParseGFpHex(`100000000000001ecea551ad837e9`),
		Gx:      HI(`1`),
		Gy:      HI(`2`),
	}))
	cAm3 = NewCurve(Build(&CurveParams{
		Name:    "secp112r1",
		BitSize: 112,
		P:       NewGFpModulus(HI("db7c2abf62e35e668076bead208b")),
		A:       ParseGFpHex("db7c2abf62e35e668076bead2088"),
		B:       ParseGFpHex("659ef8ba043916eede8911702b22"),
		N:       ParseGFpHex("db7c2abf62e35e7628dfac6561c5"),
		Gx:      HI("9487239995a5ee76b55f9c2f098"),
		Gy:      HI("a89ce5af8724c0a23e0e0ff77500"),
	}))
)

var (
	tc = []CurveTestCases{
		tcA,
		tcA0,
		tcAm3,
	}

	tcA = CurveTestCases{
		Curve: cA,
		K: [][]byte{
			HB(`4127929b42996e48b0d0ae842ba1`),
			HB(`3d0a0f0ee3d1e54e75d14cd271f4`),
			HB(`7920531a0cd0efd128079ee496a4`),
			HB(`8def260e99a5231ab6a7529c0742`),
			HB(`a82c8bc40c4b5b070a1748484cf8`),
			HB(`2d266fd03f25c475f8aa61a44059`),
			HB(`6dc6d4a392b12b0082084e0fee04`),
			HB(`479a2dffca5378e25ef0e80cb5cf`),
			HB(`5353d597f9a6dc4b21f9b83d3c18`),
			HB(`6b7c4a09ade26d2b732141b70ec8`),
		},
		P: []Point{
			P(`b9d5f3919601417931260df09489`, `afb70f78d855509b4616a2228a31`),
			P(`5eddb9ffaeaa161396514a8bfb79`, `1672eae7d432be5b9a9b7db72630`),
			P(`2404cd255b1f1387b82a65489f38`, `488510ed61eca9d31cad5689775a`),
			P(`141df27c47db085c1fd8ad47b3e7`, `37df46542c95fb6ec454e5bb7669`),
			P(`3fdb0d82cc104a4bea3abd01b093`, `8a9b1b6e3d883a29b1371fa56ac6`),
			P(`95243ee573c79660959b4b67e071`, `9fd00838fbd858cc4eeffbad2d93`),
			P(`49eea9cdca25258979dd2def01b9`, `8bb778c72d90232b0cb36c1cf21f`),
			P(`54a1a5354bb305cc9c81891fda5a`, `c50ed880bdd3e7969865a239fa4f`),
			P(`603320dd087f96ae21742bb4fe68`, `dac3331946766439687098906e1d`),
			P(`3b35545c2b725fc9a013878dd899`, `397aefb8eb0a642d659bdd9b4126`),
		},
		Add: []Point{
			P(`102f841a07eb37ed216e0c9d131d`, `ca3c8e23925ee3e306fb0f57c6da`),
			P(`628d7db85f0bc157702dc1424c27`, `4d1a1685c39124cbf937b1d8fa2`),
			P(`c3a733171a498a14f4ca090f8599`, `4d91cb3eae8a5d592c6926e6d050`),
			P(`86a4c8b19d882ba41c15c9f7a684`, `797ac1bac880227a67f9ac05f71a`),
			P(`8428072e95c2825f39fce1129fce`, `b9066f6b8f00a22fd3176763a8cc`),
			P(`d7862a0ffacc426060c3546b6956`, `cf1bd1f685e54289d2610f06807f`),
			P(`6ef8b9ae8034d24e2b9dc336df28`, `5e5535a2db623e88219e65b5faa5`),
			P(`2972d19089e22208b1cd9bd1d48a`, `a906b8e24731adfa11feacdbea24`),
			P(`72a095217b8958a111e5081a9bcb`, `a6022f79f845b902179d4db4a631`),
			P(`74f62323381bc1d6c891eae79a18`, `60251aadf5772e8ecc1b00a8c3eb`),
		},
		Double: []Point{
			P(`102f841a07eb37ed216e0c9d131d`, `ca3c8e23925ee3e306fb0f57c6da`),
			P(`a04192cb8796535d16e1366bdc1c`, `43a02dfa3f2c95c8adc8748643`),
			P(`243c0d14d7845409c1f930e6d9e9`, `11decd636829ce15dfc88b78acbc`),
			P(`5d6930105b623c10732bda3cf8d`, `d2f4c151ef13d2a85beac53716e6`),
			P(`cd1963659f020e9c0411511d547`, `3ddc031892c3c1c184f986537282`),
			P(`1ffcb7831e89aaca30f796b3b1c0`, `8836a504fdaba19f5752c7d8eb6a`),
			P(`b5e58b961c95d497adf99329a42c`, `9be0d1d2e0db035bdfd0d78c9719`),
			P(`ad5d71dd744a7e7d43e43f63a7b8`, `c5a53f2c81018251c0762c802caa`),
			P(`42ffceaf3b85f865f774c7702a2e`, `d4f3dc736567afbc418ea85335fe`),
			P(`c3926a85cf23fe78af7ef8450d48`, `c8dc9bac17173994c2058be1c888`),
		},
		ScalarMult: []Point{
			P(`b911e9f235a91e045bb8e39f1573`, `9a2b2e2e484f29808d3d94b5b3f0`),
			P(`33b707542a609848a3a930a292e1`, `e3045ad5f7dac58e9fb7fdcbc60`),
			P(`142d37d55d16cdbeda542481ab30`, `8a0d57e7b60381d5fbb524369d4d`),
			P(`32c81c47852f2f95c4668b05f93`, `6785d51c40e7c272fe492d334e9f`),
			P(`5685a1b615c7d320cf1a88062825`, `a2b4d98632eae1c78af881e6e31f`),
			P(`a4c482aa32184519cd45ee7387da`, `2cb55816115c8cb77fe923d0752b`),
			P(`6953266e4e9eb501e7b33a189831`, `3fd27687d2ffede69cd6e1a0f100`),
			P(`6d2e6e4ddf7f55c2ebf26b79226c`, `722645865e9276022570f53612c9`),
			P(`a60559b4e73181cd85c16c00793e`, `79d03e2070c7baf570d6f97ffe70`),
			P(`b83daefcc40d92161eb46a2ba364`, `11638340996b18b75ed9525f9665`),
		},
		InvalidP: []Point{
			P(`49ba38e7298db1c0030a590a1eb5`, `9dd7f6249ce00d397bdcd2cef002`),
			P(`6ceae71007fdd51f130c072aacc3`, `7fad6bbd03e4aeeac02ffaccaac5`),
			P(`4d9c7f11e3e5cceffdd36789d7be`, `58061bade9e17f96c9ac6e2c0651`),
			P(`899fcda07eaaddbc27c6e2643121`, `6f428d158408cd0890d4f12b6eba`),
			P(`62def074764efe7bfae85f7c13de`, `1e5500bb7bbac005d5db7679d98f`),
			P(`48ec72c129cc694ba4147308a34d`, `490ddb44762446b2cc652e30433c`),
			P(`709f6e8d4ebeb0ada5f4e5e6506c`, `99db159f09b42d3a746b21b05c89`),
			P(`b20c1d51cef5e2a4e8c9615fe15a`, `22f349463bad5a297d8d0b65f748`),
			P(`557cfb3af2dd0f3426f75a598be8`, `5f178115e0ae711db37ad4bb2b01`),
			P(`3599294c17c545fedd96f009071b`, `39e99fb054527f16bd4d31402c76`),
		},
	}
	tcA0 = CurveTestCases{
		Curve: cA0,
		K: [][]byte{
			HB(`824f929b42996e48b0d0ae842ba1`),
			HB(`7a150f0ee3d1e54e75d14cd271f4`),
			HB(`f240531a0cd0efd128079ee496a4`),
			HB(`5a4d6fd03f25c475f8aa61a44059`),
			HB(`db8dd4a392b12b0082084e0fee04`),
			HB(`8f352dffca5378e25ef0e80cb5cf`),
			HB(`a6a7d597f9a6dc4b21f9b83d3c18`),
			HB(`d6f84a09ade26d2b732141b70ec8`),
			HB(`937438e7298db1c0030a590a1eb5`),
			HB(`d9d5e71007fdd51f130c072aacc3`),
		},
		P: []Point{
			P(`c3d74fb8c15ee8cd95789218b178`, `77ee1b2d8671f215a329c83e21c9`),
			P(`52d149080830a0f7c81064b0bc35`, `31c29f1abba1d62bc3a3385a70c5`),
			P(`602c0b2ea0908aa8eae4173d9f7`, `62b621f35cf09ad97a04fea6017`),
			P(`40bd7e70075a041e2c6ae8bff9a4`, `ed655f42bd7a6d4197fcbce17b9b`),
			P(`1f306f2851cec705650361037799`, `3736c517574d6f02d0c03fbab006`),
			P(`ee676f63df6d96066020599f2f4a`, `bdbbdcc708dcdb3c085d7e374359`),
			P(`a388e112d68a4976aa18db5b0c1c`, `4f3fbe39a7743c1e88ae9905370c`),
			P(`591c1bf8646c6b18ca1c6e19264c`, `7fb11191f13a3125aa4c9a96e314`),
			P(`ef95a545c9aaf24887cc849ac0bb`, `705be07ae1d5d70620c18c49d71`),
			P(`bb9e92391e71d3d787e8610ef492`, `49e6f07e699c67c57d71746a8eb9`),
		},
		Add: []Point{
			P(`dd5acb95811270acfb82c8e9b304`, `b8ff7496254bf1c67b822137324a`),
			P(`ab7174b3ba712e3ac19512d57944`, `117f966898c9f471d1651316064c`),
			P(`f61e1bffab1e0a376a4b77756152`, `538ea52f23506190198ed10351f4`),
			P(`1708751364bd90c4d75ddccaec7f`, `80f61d6f8894e3f70fcb24798f04`),
			P(`406fe6fafd4b96d0a693cbdec022`, `c2be132a16e80c50fb3a3b99dae`),
			P(`613f7ee463a4700a664ee79948c3`, `90a7c5c7a3bc5f6aa582e2976c65`),
			P(`dfae3aedbd834d923549c0227ccd`, `26504a652d4b5f674ec9486c311a`),
			P(`e191e7d05aea85fa22d4080dd8ff`, `43deaa9cb2ef18f493e51de0ad6b`),
			P(`bcfb23e3363440d4835ea212a612`, `63bb08fbb580bc0977a0d94aed77`),
			P(`9454562595210654b0aea072c394`, `9c38ec7ac24e92c153e5b3c69473`),
		},
		Double: []Point{
			P(`dd5acb95811270acfb82c8e9b304`, `b8ff7496254bf1c67b822137324a`),
			P(`59826723505780971e53da40e6fb`, `3048059e616011c076891a4241c8`),
			P(`29975f3393745d008e18e5c8332c`, `82ea6037e6710f458c5f5f8ce2d7`),
			P(`99d2afa681482ebcf1eb3e75ef6e`, `fe56ff20162b410b1df4b02a3426`),
			P(`65020eed301abb4a45ee646bc457`, `a17bb99c9f9891b15082ef49d404`),
			P(`ab2629b788b23b5a24730c82748`, `9dc8b3c2a8dfb3fed39d1aec14bd`),
			P(`ca288df5a4162a33e5ea2ce6adc6`, `1e688f6d4c754b1f5e1863e9c8a7`),
			P(`1cb7cc9c95ab7d88f923728af95d`, `518564beb3944cfc243e1e58ad1c`),
			P(`ffe79186cb2a0badda27b3de587e`, `130848bea52f207dcb1fa4d8d906`),
			P(`22b903ca85f0dad6f87d36b80616`, `97c59fc427465caf2dd89cf66100`),
		},
		ScalarMult: []Point{
			P(`1e981cd20b3ff6d55ba2f28937ae`, `cc4a2f0f90ebe2993d7b67a25c95`),
			P(`91745baf5c5581129cccdec7068`, `504b5c7755f7ece19319ac695db5`),
			P(`19cca947973bf3e66e346a381b0a`, `aa0dd70f8be22000800ae200b59b`),
			P(`cc2fb42f72da09bb471641516250`, `919e896bf6dcbfe9f1ee4a4ae3a6`),
			P(`8a7ca52a82bfcc9b516bba1d37d6`, `1928f2b57e65918c1dc802056511`),
			P(`e0b6c7f348382edbeb1eb86cf139`, `ab146a9d538edd6407dbbcc29fb2`),
			P(`21d665a7624821c00b72405a09`, `ee7a9c91309f540d96bdb0c55620`),
			P(`5cafe0dbec7376243ffbd7a1dce`, `46e97af23a94a4bdc597e9a3b769`),
			P(`eb82583f19969d64cafb9a6d84c2`, `3eea3ef5ad481b0598239fb5c066`),
			P(`c9ae3eaac8c59b5390fc740cc925`, `7ff36bfd6a8baed870134f70314a`),
		},
		InvalidP: []Point{
			P(`4d9c7f11e3e5cceffdd36789d7be`, `7cdb3c4904a230a31c2319c13f64`),
			P(`8a17e01196aaa130b624ea26fe9b`, `977b481cd03b650a5ea442ea881c`),
			P(`65909c7df110b19e6b3118fc7b32`, `7cc6bd2600ee444f362d9a0f14f0`),
			P(`2c216628c1f4c36b66b8ed783035`, `c94e8f05de93e4df20925b6f120a`),
			P(`110e2b62094dda3c22fecefa8b9e`, `cbda102bf3692227c7f783dd6376`),
			P(`48ec72c129cc694ba4147308a34d`, `ae66d6b4dd74e07ddd4facf1b870`),
			P(`dde7cdf09513f102081eab2f004c`, `46e4e7b7fe0cdd2a09f948a18588`),
			P(`a459553e812e7cae7a101062dfe5`, `f4359e5807ea26660749eea8421b`),
			P(`709f6e8d4ebeb0ada5f4e5e6506c`, `a0a6e0344688b76a7682b6cf80d3`),
			P(`b20c1d51cef5e2a4e8c9615fe15a`, `41823baaec04bf818adeac68f090`),
		},
	}
	tcAm3 = CurveTestCases{
		Curve: cAm3,
		K: [][]byte{
			HB(`2d266fd03f25c475f8aa61a44059`),
			HB(`6dc6d4a392b12b0082084e0fee04`),
			HB(`479a2dffca5378e25ef0e80cb5cf`),
			HB(`5353d597f9a6dc4b21f9b83d3c18`),
			HB(`6b7c4a09ade26d2b732141b70ec8`),
			HB(`a8e13a767ef92c4cc9061dad7237`),
			HB(`88f2e6c9596fd7ea401fe61c42e9`),
			HB(`49ba38e7298db1c0030a590a1eb5`),
			HB(`6ceae71007fdd51f130c072aacc3`),
			HB(`4d9c7f11e3e5cceffdd36789d7be`),
		},
		P: []Point{
			P(`6b933ba17550d2ad2f74a5236673`, `593c92748c07c1277d35a9dff406`),
			P(`e6e0d99ce3cd5c4a4fb6aa04ee8`, `d16912c9d2e19bde11a38434bb8b`),
			P(`a1caef028aa8382669753c72b591`, `424cae00c340fc386bb857cbe60f`),
			P(`6908f927543148a1bfe50e2add3c`, `8a729e8edc178a782665e4b8c124`),
			P(`7ec0f2a4c0a93d712fe83e7fb603`, `baec09b37e405f59e24454fcaa76`),
			P(`536008d12269893de31fd8941a5b`, `cb34ed924395ec0420a680295c9e`),
			P(`9438537f5a08f16d99c0a1955a0e`, `cf8dd34f091380c179713537c797`),
			P(`244fafbb03a945743efad89952f6`, `41be12a363c3b3970da85c579d34`),
			P(`118bcfe817e8f3700fd25d651c84`, `40133f28caec7f674110cf25f4e7`),
			P(`2a78dc8e610a4555318b80b3284a`, `27736fc05d247c057da1f24f18e4`),
		},
		Add: []Point{
			P(`c3e2ef95f62963766cbc8d6d5032`, `50472341d638cd2388356f67f549`),
			P(`8a9c461150787c26045e836e5db1`, `bd887381a23591da20dafc5b0aac`),
			P(`17714820b71f8b0862694b936fd4`, `2f31288faf36d25b89a605909602`),
			P(`2c8de70a5cc058f4391611af38cc`, `64798f925a4415ccbc09caa9d947`),
			P(`c189fdb94fcff56b35d1b22b8c06`, `1581e708436d44b78717e20b3c8c`),
			P(`543a1d14478a769bcd464e47a1ac`, `2001f54f9fbb30e5b364b15ccb20`),
			P(`98dcd094a31bb61fce46ed3a9627`, `c9fd6c3aefdbdb36186a0446d5c1`),
			P(`2105688a0214a97e89216cf5ea1a`, `4e4813c3b541f51a5f1118868ea5`),
			P(`a3a978533aef96491bfd2d171114`, `4f52f2c5447d8e14d6041dd607db`),
			P(`bc64a2c911c8ea0db01821942375`, `11ecb7127b557cdfb3b9b8d2b53c`),
		},
		Double: []Point{
			P(`c3e2ef95f62963766cbc8d6d5032`, `50472341d638cd2388356f67f549`),
			P(`57003020c6272039086b264c7d5a`, `8cb67ba895bdf1787b68dbc1193e`),
			P(`a45483a6bd259ef80d66e161529e`, `116a1c812973d3ddb4d9fb3ed05a`),
			P(`46cbc412d1ea3f7210d1215dbdd3`, `66bd1a5d53d90b096a5d6b28e823`),
			P(`7570d2223f21cabac4128c7c062a`, `9644a98fbd37db7026e9cdcdad5b`),
			P(`fb0c5e63fe4d49eb264fd91ef26`, `4904ae558c77b0a65723d41967b`),
			P(`3bdd8b38e4f98f2c25e140365e1`, `47e2bb52be9e92a5e9ea196bf0bd`),
			P(`7a94eeced6a7d0c52531e8e0910a`, `17a057e5660f746bde18059d382f`),
			P(`60c34bb7fed1edab8bac8269c9e0`, `335a7da9b5c08ad9ec9c33711f6f`),
			P(`59728a92760c96e8a1bcfc51233a`, `982943d7d71c165a000d46310833`),
		},
		ScalarMult: []Point{
			P(`d52224467b8fc75128043f1bdd11`, `7753a6d07c95f6d3cbe61cb0945`),
			P(`c0ad1bbda0b640adb7d69f278fee`, `a6ccdb3aa8c920fd7a378ac78148`),
			P(`8ddc94de54b5a4e14d30f9473d86`, `89be0d642554466d50edccfc2dd6`),
			P(`9de9c65f1361d16fc9632d152232`, `a3107bd116bd47b7fa60d92ac68b`),
			P(`4bccdaba36c1e3196cfb2ebc38fc`, `759fa4a525f7de24551b716c1aaf`),
			P(`d624d55472ce26d78172aac17f95`, `2fe1763f9209234c48cdd40abac4`),
			P(`cb6438a1ca3584980fb07b0b04d2`, `9ff71a1d51c136c3bfb80bc4b4cc`),
			P(`68aac3d23fb0e1d4eb476c0de419`, `8752ae660567678051932f1ab23d`),
			P(`1a8e73de945f84a282ed22b9d048`, `7a2225525cf79365c7b205a2d76f`),
			P(`49f5b5559375f56421365b29d80e`, `947046f9f9788271d0acfa287c45`),
		},
		InvalidP: []Point{
			P(`899fcda07eaaddbc27c6e2643121`, `cda0c8beecc648f18b35cd043d62`),
			P(`8a17e01196aaa130b624ea26fe9b`, `c4c3554809ede75b3259e873b1df`),
			P(`7f5a0e06f8b29b24361e06be9ee2`, `4bd48048a9a5d1fcfb98861ba4dd`),
			P(`62def074764efe7bfae85f7c13de`, `346076d6953f5c5b28599bbd4294`),
			P(`92966f8d4afd405687354937a867`, `458fef24acf9d84c685a71bcdede`),
			P(`a459553e812e7cae7a101062dfe5`, `1007f5945eafbde4fca5aa27f254`),
			P(`709f6e8d4ebeb0ada5f4e5e6506c`, `d6aa3af7285eaf7311ea649375b3`),
			P(`b20c1d51cef5e2a4e8c9615fe15a`, `6995ee765551cec50b8a6a27e4df`),
			P(`97d02865845f57114574c4a54dae`, `3c89759187e577888e3d1e10fb09`),
			P(`d4fcc0b36ea2560cc0eb070221d2`, `3dff8c1a5919fb559c06b7753a0a`),
		},
	}
)
