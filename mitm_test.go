package mitm

import (
	"crypto/tls"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestParseClientHello(t *testing.T) {
	for i, test := range []struct {
		inputHex string
		expected rawHelloInfo
	}{
		{
			// curl 7.51.0 (x86_64-apple-darwin16.0) libcurl/7.51.0 SecureTransport zlib/1.2.8
			inputHex: `010000a6030358a28c73a71bdfc1f09dee13fecdc58805dcce42ac44254df548f14645f7dc2c00004400ffc02cc02bc024c023c00ac009c008c030c02fc028c027c014c013c012009f009e006b0067003900330016009d009c003d003c0035002f000a00af00ae008d008c008b01000039000a00080006001700180019000b00020100000d00120010040102010501060104030203050306030005000501000000000012000000170000`,
			expected: rawHelloInfo{
				cipherSuites:       []uint16{255, 49196, 49195, 49188, 49187, 49162, 49161, 49160, 49200, 49199, 49192, 49191, 49172, 49171, 49170, 159, 158, 107, 103, 57, 51, 22, 157, 156, 61, 60, 53, 47, 10, 175, 174, 141, 140, 139},
				extensions:         []uint16{10, 11, 13, 5, 18, 23},
				compressionMethods: []byte{0},
				curves:             []tls.CurveID{23, 24, 25},
				points:             []uint8{0},
			},
		},
		{
			// Chrome 56
			inputHex: `010000c003031dae75222dae1433a5a283ddcde8ddabaefbf16d84f250eee6fdff48cdfff8a00000201a1ac02bc02fc02cc030cca9cca8cc14cc13c013c014009c009d002f0035000a010000777a7a0000ff010001000000000e000c0000096c6f63616c686f73740017000000230000000d00140012040308040401050308050501080606010201000500050100000000001200000010000e000c02683208687474702f312e3175500000000b00020100000a000a0008aaaa001d001700182a2a000100`,
			expected: rawHelloInfo{
				cipherSuites:       []uint16{6682, 49195, 49199, 49196, 49200, 52393, 52392, 52244, 52243, 49171, 49172, 156, 157, 47, 53, 10},
				extensions:         []uint16{31354, 65281, 0, 23, 35, 13, 5, 18, 16, 30032, 11, 10, 10794},
				compressionMethods: []byte{0},
				curves:             []tls.CurveID{43690, 29, 23, 24},
				points:             []uint8{0},
			},
		},
		{
			// Firefox 51
			inputHex: `010000bd030375f9022fc3a6562467f3540d68013b2d0b961979de6129e944efe0b35531323500001ec02bc02fcca9cca8c02cc030c00ac009c013c01400330039002f0035000a010000760000000e000c0000096c6f63616c686f737400170000ff01000100000a000a0008001d001700180019000b00020100002300000010000e000c02683208687474702f312e31000500050100000000ff030000000d0020001e040305030603020308040805080604010501060102010402050206020202`,
			expected: rawHelloInfo{
				cipherSuites:       []uint16{49195, 49199, 52393, 52392, 49196, 49200, 49162, 49161, 49171, 49172, 51, 57, 47, 53, 10},
				extensions:         []uint16{0, 23, 65281, 10, 11, 35, 16, 5, 65283, 13},
				compressionMethods: []byte{0},
				curves:             []tls.CurveID{29, 23, 24, 25},
				points:             []uint8{0},
			},
		},
		{
			// openssl s_client (OpenSSL 0.9.8zh 14 Jan 2016)
			inputHex: `0100012b03035d385236b8ca7b7946fa0336f164e76bf821ed90e8de26d97cc677671b6f36380000acc030c02cc028c024c014c00a00a500a300a1009f006b006a0069006800390038003700360088008700860085c032c02ec02ac026c00fc005009d003d00350084c02fc02bc027c023c013c00900a400a200a0009e00670040003f003e0033003200310030009a0099009800970045004400430042c031c02dc029c025c00ec004009c003c002f009600410007c011c007c00cc00200050004c012c008001600130010000dc00dc003000a00ff0201000055000b000403000102000a001c001a00170019001c001b0018001a0016000e000d000b000c0009000a00230000000d0020001e060106020603050105020503040104020403030103020303020102020203000f000101`,
			expected: rawHelloInfo{
				cipherSuites:       []uint16{49200, 49196, 49192, 49188, 49172, 49162, 165, 163, 161, 159, 107, 106, 105, 104, 57, 56, 55, 54, 136, 135, 134, 133, 49202, 49198, 49194, 49190, 49167, 49157, 157, 61, 53, 132, 49199, 49195, 49191, 49187, 49171, 49161, 164, 162, 160, 158, 103, 64, 63, 62, 51, 50, 49, 48, 154, 153, 152, 151, 69, 68, 67, 66, 49201, 49197, 49193, 49189, 49166, 49156, 156, 60, 47, 150, 65, 7, 49169, 49159, 49164, 49154, 5, 4, 49170, 49160, 22, 19, 16, 13, 49165, 49155, 10, 255},
				extensions:         []uint16{11, 10, 35, 13, 15},
				compressionMethods: []byte{1, 0},
				curves:             []tls.CurveID{23, 25, 28, 27, 24, 26, 22, 14, 13, 11, 12, 9, 10},
				points:             []uint8{0, 1, 2},
			},
		},
	} {
		data, err := hex.DecodeString(test.inputHex)
		if err != nil {
			t.Fatalf("Test %d: Could not decode hex data: %v", i, err)
		}
		actual := parseRawClientHello(data)
		if !reflect.DeepEqual(test.expected, actual) {
			t.Errorf("Test %d: Expected %+v; got %+v", i, test.expected, actual)
		}
	}
}

func TestHeuristicFunctionsAndHandler(t *testing.T) {
	// To test the heuristics, we assemble a collection of real
	// ClientHello messages from various TLS clients, both genuine
	// and intercepted. Please be sure to hex-encode them and
	// document the User-Agent associated with the connection
	// as well as any intercepting proxy as thoroughly as possible.
	//
	// If the TLS client used is not an HTTP client (e.g. s_client),
	// you can leave the userAgent blank, but please use a comment
	// to document crucial missing information such as client name,
	// version, and platform, maybe even the date you collected
	// the sample! Please group similar clients together, ordered
	// by version for convenience.

	// clientHello pairs a User-Agent string to its ClientHello message.
	type clientHello struct {
		userAgent    string
		helloHex     string      // do NOT include the header, just the ClientHello message
		interception bool        // if test case shows an interception, set to true
		reqHeaders   http.Header // if the request should set any headers to imitate a browser or proxy
	}

	// clientHellos groups samples of true (real) ClientHellos by the
	// name of the browser that produced them. We limit the set of
	// browsers to those we are programmed to protect, as well as a
	// category for "Other" which contains real ClientHello messages
	// from clients that we do not recognize, which may be used to
	// test or imitate interception scenarios.
	//
	// Please group similar clients and order by version for convenience
	// when adding to the test cases.
	clientHellos := map[string][]clientHello{
		"Chrome": {
			{
				userAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
				helloHex:     `010000c003031dae75222dae1433a5a283ddcde8ddabaefbf16d84f250eee6fdff48cdfff8a00000201a1ac02bc02fc02cc030cca9cca8cc14cc13c013c014009c009d002f0035000a010000777a7a0000ff010001000000000e000c0000096c6f63616c686f73740017000000230000000d00140012040308040401050308050501080606010201000500050100000000001200000010000e000c02683208687474702f312e3175500000000b00020100000a000a0008aaaa001d001700182a2a000100`,
				interception: false,
			},
			// TODO: Chrome on iOS will use iOS' TLS stack for requests that load
			// the web page (apparently required by the dev ToS) but will use its
			// own TLS stack for everything else, it seems. Figure out a decent way
			// to test this with a nice, unified corpus that allows for this variance.
			// {
			// 	// Chrome on iOS
			// 	userAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 10_0_2 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) CriOS/56.0.2924.79 Mobile/14A456 Safari/602.1",
			// 	helloHex:  `010000de030358b062c509b21410a6496b5a82bfec74436cdecebe8ea1da29799939bbd3c17200002c00ffc02cc02bc024c023c00ac009c008c030c02fc028c027c014c013c012009d009c003d003c0035002f000a0100008900000014001200000f66696e6572706978656c732e636f6d000a00080006001700180019000b00020100000d00120010040102010501060104030203050306033374000000100030002e0268320568322d31360568322d31350568322d313408737064792f332e3106737064792f3308687474702f312e310005000501000000000012000000170000`,
			// },
			// {
			// 	// Chrome on iOS (requesting favicon)
			// 	userAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 10_0_2 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) CriOS/56.0.2924.79 Mobile/14A456 Safari/602.1",
			// 	helloHex:  `010000c20303863eb64788e3b9638c261300318411cbdd8f09576d58eec1e744b6ce944f574f0000208a8acca9cca8cc14cc13c02bc02fc02cc030c013c014009c009d002f0035000a01000079baba0000ff0100010000000014001200000f66696e6572706978656c732e636f6d0017000000230000000d00140012040308040401050308050501080606010201000500050100000000001200000010000e000c02683208687474702f312e31000b00020100000a000a00083a3a001d001700184a4a000100`,
			// },
			{
				userAgent:    "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
				helloHex:     `010000c603036f717a88212c3e9e41940f82c42acb3473e0e4a64e8f52d9af33d34e972e08a30000206a6ac02bc02fc02cc030cca9cca8cc14cc13c013c014009c009d002f0035000a0100007d7a7a0000ff0100010000000014001200000f66696e6572706978656c732e636f6d0017000000230000000d00140012040308040401050308050501080606010201000500050100000000001200000010000e000c02683208687474702f312e3175500000000b00020100000a000a00087a7a001d001700188a8a000100`,
				interception: false,
			},
			{
				userAgent:    "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
				helloHex:     `010001fc030383141d213d1bf069171843489faf808028d282c9828e1ba87637c863833c730720a67e76e152f4b704523b72317ef4587e231f02e2395e0ecac6be9f28c35e6ce600208a8ac02bc02fc02cc030cca9cca8cc14cc13c013c014009c009d002f0035000a010001931a1a0000ff0100010000000014001200000f66696e6572706978656c732e636f6d00170000002300785e85429bf1764f33111cd3ad5d1c56d765976fd962b49dbecbb6f7865e2a8d8536ad854f1fa99a8bbbf998814fee54a63a0bf162869d2bba37e9778304e7c4140825718e191b574c6246a0611de6447bdd80417f83ff9d9b7124069a9f74b90394ecb89bec5f6a1a67c1b89e50b8674782f53dd51807651a000d00140012040308040401050308050501080606010201000500050100000000001200000010000e000c02683208687474702f312e3175500000000b00020100000a000a00081a1a001d001700182a2a0001000015009a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000`,
				interception: false,
			},
		},
		"Firefox": {
			{
				userAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:51.0) Gecko/20100101 Firefox/51.0",
				helloHex:     `010000bd030375f9022fc3a6562467f3540d68013b2d0b961979de6129e944efe0b35531323500001ec02bc02fcca9cca8c02cc030c00ac009c013c01400330039002f0035000a010000760000000e000c0000096c6f63616c686f737400170000ff01000100000a000a0008001d001700180019000b00020100002300000010000e000c02683208687474702f312e31000500050100000000ff030000000d0020001e040305030603020308040805080604010501060102010402050206020202`,
				interception: false,
			},
			{
				userAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:53.0) Gecko/20100101 Firefox/53.0",
				helloHex:     `010001fc0303c99d54ae0628bbb9fea3833a4244c6a712cac9d7738f4930b8b9d8e2f6bd578220f7936cedb48907981c9292fb08ceee6f59bd6fddb3d4271ccd7c12380c5038ab001ec02bc02fcca9cca8c02cc030c00ac009c013c01400330039002f0035000a01000195001500af000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000c0000096c6f63616c686f737400170000ff01000100000a000a0008001d001700180019000b000201000023007886da2d41843ff42131b856982c19a545837b70e604325423a817d925e9d95bd084737682cea6b804dfb7cbe336a3b27b8d520d57520c29cfe5f4f3d3236183b84b05c18f0ca30bf598111e390086fea00d9631f1f78527277eb7838b86e73c4e5d15b55d086b1a4a8aa29f12a55126c6274bcd499bbeb23a0010000e000c02683208687474702f312e31000500050100000000000d0018001604030503060308040805080604010501060102030201`,
				interception: false,
			},
			{
				userAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:53.0) Gecko/20100101 Firefox/53.0",
				helloHex:     `010000b1030365d899820b999245d571c2f7d6b850f63ad931d3c68ceb9cf5a508421a871dc500001ec02bc02fcca9cca8c02cc030c00ac009c013c01400330039002f0035000a0100006a0000000e000c0000096c6f63616c686f737400170000ff01000100000a000a0008001d001700180019000b00020100002300000010000e000c02683208687474702f312e31000500050100000000000d0018001604030503060308040805080604010501060102030201`,
				interception: false,
			},
			{
				// this was a Nightly release at the time
				userAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:55.0) Gecko/20100101 Firefox/55.0",
				helloHex:     `010001fc030331e380b7d12018e1202ef3327607203df5c5732b4fa5ab5abaf0b60034c2fb662070c836b9b89123e37f4f1074d152df438fa8ee8a0f89b036fd952f4fcc0b994f001c130113031302c02bc02fcca9cca8c02cc030c013c014002f0035000a0100019700000014001200000f63616464797365727665722e636f6d00170000ff01000100000a000e000c001d00170018001901000101000b0002010000230078c97e7716a041e2ea824571bef26a3dff2bf50a883cd15d904ab2d17deb514f6e0a079ee7c212c000178387ffafc2e530b6df6662f570aae134330f13c458a0eaad5a96a9696f572110918740b15db1143d19aaaa706942030b433a7e6150f62b443c0564e5b8f7ee9577bf3bf7faec8c67425b648ab54d880010000e000c02683208687474702f312e310005000501000000000028006b0069001d0020aee6e596155ee6f79f943e81ceabe0979d27fbbb8b9189ccb2ebc75226351f32001700410421875a44e510decac11ef1d7cfddd4dfe105d5cd3a2d42fba03ebde23e51e8ce65bda1b48be82d4848d1db2bfce68e94092e925a9ce0dbf5df35479558108489002b0009087f12030303020301000d0018001604030503060308040805080604010501060102030201002d000201010015002500000000000000000000000000000000000000000000000000000000000000000000000000`,
				interception: false,
			},
		},
		"Edge": {
			{
				userAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.79 Safari/537.36 Edge/14.14393",
				helloHex:     `010000bd030358a3c9bf05f734842e189fb6ce653b67b846e990bc1fc5fb8c397874d06020f1000038c02cc02bc030c02f009f009ec024c023c028c027c00ac009c014c01300390033009d009c003d003c0035002f000a006a00400038003200130100005c000500050100000000000a00080006001d00170018000b00020100000d00140012040105010201040305030203020206010603002300000010000e000c02683208687474702f312e310017000055000006000100020002ff01000100`,
				interception: false,
			},
		},
		"Safari": {
			{
				userAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8",
				helloHex:     `010000d2030358a295b513c8140c6ff880f4a8a73cc830ed2dab2c4f2068eb365228d828732e00002600ffc02cc02bc024c023c00ac009c030c02fc028c027c014c013009d009c003d003c0035002f010000830000000e000c0000096c6f63616c686f7374000a00080006001700180019000b00020100000d00120010040102010501060104030203050306033374000000100030002e0268320568322d31360568322d31350568322d313408737064792f332e3106737064792f3308687474702f312e310005000501000000000012000000170000`,
				interception: false,
			},
		},
		"Other": { // these are either non-browser clients or intercepted client hellos
			{
				// openssl s_client (OpenSSL 0.9.8zh 14 Jan 2016) - NOT an interception, but not a browser either
				helloHex: `0100012b03035d385236b8ca7b7946fa0336f164e76bf821ed90e8de26d97cc677671b6f36380000acc030c02cc028c024c014c00a00a500a300a1009f006b006a0069006800390038003700360088008700860085c032c02ec02ac026c00fc005009d003d00350084c02fc02bc027c023c013c00900a400a200a0009e00670040003f003e0033003200310030009a0099009800970045004400430042c031c02dc029c025c00ec004009c003c002f009600410007c011c007c00cc00200050004c012c008001600130010000dc00dc003000a00ff0201000055000b000403000102000a001c001a00170019001c001b0018001a0016000e000d000b000c0009000a00230000000d0020001e060106020603050105020503040104020403030103020303020102020203000f000101`,
				// NOTE: This test case is not actually an interception, but s_client is not a browser
				// or any client we support MITM checking for, either. Since it advertises heartbeat,
				// our heuristics still flag it as a MITM.
				interception: true,
			},
			{
				// curl 7.51.0 (x86_64-apple-darwin16.0) libcurl/7.51.0 SecureTransport zlib/1.2.8
				userAgent:    "curl/7.51.0",
				helloHex:     `010000a6030358a28c73a71bdfc1f09dee13fecdc58805dcce42ac44254df548f14645f7dc2c00004400ffc02cc02bc024c023c00ac009c008c030c02fc028c027c014c013c012009f009e006b0067003900330016009d009c003d003c0035002f000a00af00ae008d008c008b01000039000a00080006001700180019000b00020100000d00120010040102010501060104030203050306030005000501000000000012000000170000`,
				interception: false,
			},
			{
				// Avast 17.1.2286 (Feb. 2017) on Windows 10 x64 build 14393, intercepting Edge
				userAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.79 Safari/537.36 Edge/14.14393",
				helloHex:     `010000ce0303b418fdc4b6cf6436a5e2bfb06b96ed5faa7285c20c7b49341a78be962a9dc40000003ac02cc02bc030c02f009f009ec024c023c028c027c00ac009c014c01300390033009d009c003d003c0035002f000a006a004000380032001300ff0100006b00000014001200000f66696e6572706978656c732e636f6d000b000403000102000a00080006001d0017001800230000000d001400120401050102010403050302030202060106030005000501000000000010000e000c02683208687474702f312e310016000000170000`,
				interception: true,
			},
			{
				// Kaspersky Internet Security 17.0.0.611 on Windows 10 x64 build 14393, intercepting Edge
				userAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.79 Safari/537.36 Edge/14.14393",
				helloHex:     `010000eb030361ce302bf4b0d5adf1ff30b2cf433c4a4b68f33e07b2651695e7ae6ec3cf126400003ac02cc02bc030c02f009f009ec024c023c028c027c00ac009c014c01300390033009d009c003d003c0035002f000a006a004000380032001300ff0100008800000014001200000f66696e6572706978656c732e636f6d000b000403000102000a001c001a00170019001c001b0018001a0016000e000d000b000c0009000a00230000000d0020001e060106020603050105020503040104020403030103020303020102020203000500050100000000000f0001010010000e000c02683208687474702f312e31`,
				interception: true,
			},
			{
				// Kaspersky Internet Security 17.0.0.611 on Windows 10 x64 build 14393, intercepting Firefox 51
				userAgent:    "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:51.0) Gecko/20100101 Firefox/51.0",
				helloHex:     `010001fc0303768e3f9ea75194c7cb03d23e8e6371b95fb696d339b797be57a634309ec98a42200f2a7554098364b7f05d21a8c7f43f31a893a4fc5670051020408c8e4dc234dd001cc02bc02fc02cc030c00ac009c013c01400330039002f0035000a00ff0100019700000014001200000f66696e6572706978656c732e636f6d000b000403000102000a001c001a00170019001c001b0018001a0016000e000d000b000c0009000a00230078bf4e244d4de3d53c6331edda9672dfc4a17aae92b671e86da1368b1b5ae5324372817d8f3b7ffe1a7a1537a5049b86cd7c44863978c1e615b005942755da20fc3a4e34a16f78034aa3b1cffcef95f81a0995c522a53b0e95a4f98db84c43359d93d8647b2de2a69f3ebdcfc6bca452730cbd00179226dedf000d0020001e060106020603050105020503040104020403030103020303020102020203000500050100000000000f0001010010000e000c02683208687474702f312e3100150093000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000`,
				interception: true,
			},
			{
				// Kaspersky Internet Security 17.0.0.611 on Windows 10 x64 build 14393, intercepting Chrome 56
				userAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
				helloHex:     `010000c903033481e7af24e647ba5a79ec97e9264c1a1f990cf842f50effe22be52130d5af82000018c02bc02fc02cc030c013c014009c009d002f0035000a00ff0100008800000014001200000f66696e6572706978656c732e636f6d000b000403000102000a001c001a00170019001c001b0018001a0016000e000d000b000c0009000a00230000000d0020001e060106020603050105020503040104020403030103020303020102020203000500050100000000000f0001010010000e000c02683208687474702f312e31`,
				interception: true,
			},
			{
				// AVG 17.1.3006 (build 17.1.3354.20) on Windows 10 x64 build 14393, intercepting Edge
				userAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.79 Safari/537.36 Edge/14.14393",
				helloHex:     `010000ca0303fd83091207161eca6b4887db50587109c50e463beb190362736b1fcf9e05f807000036c02cc02bc030c02f009f009ec024c023c028c027c00ac009c014c01300390033009d009c003d003c0035002f006a00400038003200ff0100006b00000014001200000f66696e6572706978656c732e636f6d000b000403000102000a00080006001d0017001800230000000d001400120401050102010403050302030202060106030005000501000000000010000e000c02683208687474702f312e310016000000170000`,
				interception: true,
			},
			{
				// IE 11 on Windows 7, this connection was intercepted by Blue Coat
				// no sensible User-Agent value, since Blue Coat changes it to something super generic
				// By the way, here's another reason we hate Blue Coat: they break TLS 1.3:
				// https://twitter.com/FiloSottile/status/835269932929667072
				helloHex:     `010000b1030358a3f3bae627f464da8cb35976b88e9119640032d41e62a107d608ed8d3e62b9000034c028c027c014c013009f009e009d009cc02cc02bc024c023c00ac009003d003c0035002f006a004000380032000a0013000500040100005400000014001200000f66696e6572706978656c732e636f6d000500050100000000000a00080006001700180019000b00020100000d0014001206010603040105010201040305030203020200170000ff01000100`,
				interception: true,
				reqHeaders:   http.Header{"X-Bluecoat-Via": {"66808702E9A2CF4"}}, // actual field name would be "X-BlueCoat-Via" but Go canonicalizes field names
			},
			{
				// Firefox 51.0.1 being intercepted by burp 1.7.17
				userAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:51.0) Gecko/20100101 Firefox/51.0",
				helloHex:     `010000d8030358a92f4daca95acc2f6a10a9c50d736135eae39406d3090238464540d482677600003ac023c027003cc025c02900670040c009c013002fc004c00e00330032c02bc02f009cc02dc031009e00a2c008c012000ac003c00d0016001300ff01000075000a0034003200170001000300130015000600070009000a0018000b000c0019000d000e000f001000110002001200040005001400080016000b00020100000d00180016060306010503050104030401040202030201020201010000001700150000126a61677561722e6b796877616e612e6f7267`,
				interception: true,
			},
			{
				// Chrome 56 on Windows 10 being intercepted by Fortigate (on some public school network); note: I had to enable TLS 1.0 for this test (proxy was issuing a SHA-1 cert to client)
				userAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
				helloHex:     `010000e5030158ac612125c83bae95282113b2a4c572cf613c160d234350fb6d0ddce879ffec000064003300320039003800160013c013c009c014c00ac012c008002f0035000a00150012003d003c00670040006b006ac011c0070096009a009900410084004500440088008700ba00be00bd00c000c400c3c03cc044c042c03dc045c04300090005000400ff01000058000a003600340000000100020003000400050006000700080009000a000b000c000d000e000f0010001100120013001400150016001700180019000b0002010000000014001200000f66696e6572706978656c732e636f6d`,
				interception: true,
			},
			{
				// IE 11 on Windows 10, intercepted by Fortigate (same firewallas above)
				userAgent:    "Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko",
				helloHex:     `010000e5030158ac634c5278d7b17421f23a64cc91d68c470c6b247322fe867ba035b373d05c000064003300320039003800160013c013c009c014c00ac012c008002f0035000a00150012003d003c00670040006b006ac011c0070096009a009900410084004500440088008700ba00be00bd00c000c400c3c03cc044c042c03dc045c04300090005000400ff01000058000a003600340000000100020003000400050006000700080009000a000b000c000d000e000f0010001100120013001400150016001700180019000b0002010000000014001200000f66696e6572706978656c732e636f6d`,
				interception: true,
			},
			{
				// Edge 38.14393.0.0 on Windows 10, intercepted by Fortigate (same as above)
				userAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.79 Safari/537.36 Edge/14.14393",
				helloHex:     `010000e5030158ac6421a45794b8ade6a0ac6c910cde0f99c49bb1ba737b88638ec8dcf0d077000064003300320039003800160013c013c009c014c00ac012c008002f0035000a00150012003d003c00670040006b006ac011c0070096009a009900410084004500440088008700ba00be00bd00c000c400c3c03cc044c042c03dc045c04300090005000400ff01000058000a003600340000000100020003000400050006000700080009000a000b000c000d000e000f0010001100120013001400150016001700180019000b0002010000000014001200000f66696e6572706978656c732e636f6d`,
				interception: true,
			},
			{
				// Firefox 50.0.1 on Windows 10, intercepted by Fortigate (same as above)
				userAgent:    "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:50.0) Gecko/20100101 Firefox/50.0",
				helloHex:     `010000e5030158ac64e40495e77b7baf2031281451620bfe354b0c37521ebc0a40f5dc0c0cb6000064003300320039003800160013c013c009c014c00ac012c008002f0035000a00150012003d003c00670040006b006ac011c0070096009a009900410084004500440088008700ba00be00bd00c000c400c3c03cc044c042c03dc045c04300090005000400ff01000058000a003600340000000100020003000400050006000700080009000a000b000c000d000e000f0010001100120013001400150016001700180019000b0002010000000014001200000f66696e6572706978656c732e636f6d`,
				interception: true,
			},
		},
	}

	for client, chs := range clientHellos {
		for i, ch := range chs {
			hello, err := hex.DecodeString(ch.helloHex)
			if err != nil {
				t.Errorf("[%s] Test %d: Error decoding ClientHello: %v", client, i, err)
				continue
			}
			parsed := parseRawClientHello(hello)

			isChrome := parsed.looksLikeChrome()
			isFirefox := parsed.looksLikeFirefox()
			isSafari := parsed.looksLikeSafari()
			isEdge := parsed.looksLikeEdge()

			// we want each of the heuristic functions to be as
			// exclusive but as low-maintenance as possible;
			// in other words, if one returns true, the others
			// should return false, with as little logic as possible,
			// but with enough logic to force TLS proxies to do a
			// good job preserving characterstics of the handshake.
			var correct bool
			switch client {
			case "Chrome":
				correct = isChrome && !isFirefox && !isSafari && !isEdge
			case "Firefox":
				correct = !isChrome && isFirefox && !isSafari && !isEdge
			case "Safari":
				correct = !isChrome && !isFirefox && isSafari && !isEdge
			case "Edge":
				correct = !isChrome && !isFirefox && !isSafari && isEdge
			case "Other":
				correct = !isChrome && !isFirefox && !isSafari && !isEdge
			}

			if !correct {
				t.Errorf("[%s] Test %d: Chrome=%v, Firefox=%v, Safari=%v, Edge=%v; parsed hello: %+v",
					client, i, isChrome, isFirefox, isSafari, isEdge, parsed)
			}

			// test the handler too
			var got, checked bool
			want := ch.interception
			handler := NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got, checked = Check(r), Check(r)
			}), nil, nil, false)
			handler.listener.helloInfos[""] = parsed
			w := httptest.NewRecorder()
			r, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			r.Header.Set("User-Agent", ch.userAgent)
			if ch.reqHeaders != nil {
				for field, values := range ch.reqHeaders {
					r.Header[field] = values // NOTE: field names not standardized when setting directly like this!
				}
			}
			handler.ServeHTTP(w, r)
			if got != want {
				t.Errorf("[%s] Test %d: Expected MITM=%v but got %v (type assertion OK (checked)=%v)",
					client, i, want, got, checked)
			}
		}
	}
}
