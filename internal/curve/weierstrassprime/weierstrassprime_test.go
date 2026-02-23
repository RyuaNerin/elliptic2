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
		N:       HI("36df0aafd8b8d7597ca10520d04b"),
		Gx:      HI("4ba30ab5e892b4e1649dd0928643"),
		Gy:      HI("adcd46f5882e3747def36e956e97"),
	}))
	cA0 = NewCurve(Build(&CurveParams{
		Name:    `wap-wsg-idm-ecid-wtls8`,
		BitSize: 112,
		P:       NewGFpModulus(HI(`fffffffffffffffffffffffffde7`)),
		A:       ParseGFpHex(``),
		B:       ParseGFpHex(`3`),
		N:       HI(`100000000000001ecea551ad837e9`),
		Gx:      HI(`1`),
		Gy:      HI(`2`),
	}))
	cAm3 = NewCurve(Build(&CurveParams{
		Name:    "secp112r1",
		BitSize: 112,
		P:       NewGFpModulus(HI("db7c2abf62e35e668076bead208b")),
		A:       ParseGFpHex("db7c2abf62e35e668076bead2088"),
		B:       ParseGFpHex("659ef8ba043916eede8911702b22"),
		N:       HI("db7c2abf62e35e7628dfac6561c5"),
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
			HB(`1049929b42996e48b0d0ae842ba1`),
			HB(`f420f0ee3d1e54e75d14cd271f4`),
			HB(`1e48531a0cd0efd128079ee496a4`),
			HB(`237b260e99a5231ab6a7529c0742`),
			HB(`2a0b8bc40c4b5b070a1748484cf8`),
			HB(`b496fd03f25c475f8aa61a44059`),
			HB(`1b71d4a392b12b0082084e0fee04`),
			HB(`11e62dffca5378e25ef0e80cb5cf`),
			HB(`14d4d597f9a6dc4b21f9b83d3c18`),
			HB(`1adf4a09ade26d2b732141b70ec8`),
		},
		P: []Point{
			P(`17d90cc7f51164ef374f9952cacf`, `19f5277fe3f93b45213668424c91`),
			P(`d7d6ab12e5282e0c4045209e479f`, `75b14928a2dc764e363e83034a47`),
			P(`d4c61688c88ba4091927e2391b24`, `5e49328b0f7cf78ffeb9bda19cd7`),
			P(`b49d6d08bc7d48d21717392d312a`, `469b8df018c9084d74d6c0589c05`),
			P(`187c3ea1c089598b83c55514e715`, `65f382b2282f6e13b79029675a11`),
			P(`44217879dea8a52578420568e0d6`, `ca6bd86ae9d90e755b2e229a710d`),
			P(`9423697f1e668508068f860dfbff`, `af6c66e1d6dd1a4c02cb5e69efcc`),
			P(`d04b258267bc80e5fd4d193de5a8`, `7a3a294824b4064cdfe81af92cb6`),
			P(`77766d7d897b1502ecdb46d9d676`, `59f759d1789794e00bdcd8e9b60c`),
			P(`5d488a38f09d33fc98ccaf950f9c`, `d7b67eb9b32e73f0faef0e12bc9c`),
		},
		Add: []Point{
			P(`c61654eb8f4ec6d637df3b9f916f`, `9482cd3a0b4fdd5af4d75652e017`),
			P(`592e737c8ead3dabcffb322ae72e`, `149496a92f836cd1e15920bafd8b`),
			P(`39b6b7c349d99ee4d2659fc68ac5`, `9954e9faa78d43d7b431b2fb8763`),
			P(`9e3ba5c405476f68d270928cbeea`, `166a88334a1d5d3e33b90f3e60ec`),
			P(`b832192f2b15f143ab53baf425d5`, `6cd4c260a314f7afe4fd8f585d7b`),
			P(`cc5067f7b542162b445a08ddf1d1`, `3ab07e052bceefbaea44f5f170e2`),
			P(`1d51bf9b9ff97e3f400c542e0eaf`, `9028c900ec76aa0b4099c9bfa7a5`),
			P(`b93b61a5fa56e1651f0448dcc090`, `ff599565c959d94253bfcfa9684`),
			P(`51e5b027a99252e2ae058b192ba`, `5c03073a1c6f98a54ded3893af82`),
			P(`6720f7105af1cbe6af1ec66eb579`, `c9095f26e2b283cf985569075754`),
		},
		Double: []Point{
			P(`c61654eb8f4ec6d637df3b9f916f`, `9482cd3a0b4fdd5af4d75652e017`),
			P(`da0fd3185d5fbddc7c39bba95973`, `7ac574423d71fd2566e55cf94697`),
			P(`5bd0d9c89f1b06b024f9eb2455d5`, `69bca7789db76c72bec5213e4731`),
			P(`2929f8965e483526f5743f1304fe`, `9a3672cc9c35a35775c96c66b9e8`),
			P(`73d0d3fbb4eff012d31c054bc1f2`, `aa9b0bc277df124f7c5284ac372e`),
			P(`aa7a9000fe623ae9f4e6ee2117e7`, `3d041d0066dcf6e392d595f97b7d`),
			P(`6fc35431e62ca662f114b4788e79`, `a5d2261dbb57027be50fe76f188b`),
			P(`39fa5678f8bcc63333708156a60b`, `62c44133576e71a9f1ac3d6fd37e`),
			P(`55e0dca16fd5f16f66c02e4ffc6a`, `bd91e8a3e426c09ff6e35cd23f23`),
			P(`58a74d60e69d14eb4de96bda04f9`, `4d36ae0f3e1c5ae8a135099a2959`),
		},
		ScalarMult: []Point{
			P(`aa5fd5eec5d65663d95b9c49c894`, `6b685f331cf3668c03c9dfaa8dc1`),
			P(`a93b109761d3bcc628670f17b2b3`, `2165f5924f9cf9dbc6c38af724bd`),
			P(`31a2fde637839cb8140664d6c490`, `15cb6c6b18f62e66d9f81c5c922b`),
			P(`3c946d580ef752febbf760943512`, `5a40cff2ca7862887176f010cd60`),
			P(`ccd911353ef8a6f7377c8136fb5`, `10858cfdbe138c668b4ee6c4c597`),
			P(`91ed67ce3dd2f7cf7a2d4bd4ccf9`, `4e2517665f1cdbef15655a555837`),
			P(`d94cb0146d9e0b87acca980d0222`, `e424a074163703dc406aaf8a2ab`),
			P(`259d281882e7c8ee95810a8b48c8`, `d04d0e62c1f6cac1f7d5f2a77729`),
			P(`3c235708f274789accde175a246f`, `34019c0809d661ea395b191efadc`),
			P(`2a2fee3dd3f4dbb2dfdfc545667c`, `a17b6b5e114a82e86be61adc823a`),
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
