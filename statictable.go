package quack

const intern = "" +
	"script-src 'none'; object-src 'none'; base-uri 'none'max-age=315" +
	"36000; includesubdomains; preloadapplication/x-www-form-urlencod" +
	"edaccess-control-allow-credentialsaccess-control-request-headers" +
	"access-control-expose-headersaccess-control-request-methodaccess" +
	"-control-allow-headersaccess-control-allow-methodsaccess-control" +
	"-allow-originstrict-transport-securityupgrade-insecure-requestsp" +
	"ublic, max-age=31536000text/html; charset=utf-8text/plain;charse" +
	"t=utf-8application/dns-messagecontent-security-policyapplication" +
	"/javascriptx-content-type-optionscontent-dispositiontiming-allow" +
	"-originget, post, optionsif-modified-sincegzip, deflate, brx-xss" +
	"-protectioncontent-encodingapplication/jsonx-frame-optionsmax-ag" +
	"e=2592000x-forwarded-foraccept-languageaccept-encodingcontent-le" +
	"ngthmax-age=604800cache-controllast-modifiedaccept-rangesif-none" +
	"-matchauthorization1; mode=blockimage/jpegearly-datasameoriginse" +
	"t-cookie:authorityuser-agentmax-age=0expect-ctimage/gifimage/png" +
	"if-rangeprefetchno-cachelocationno-storetext/cssbytes=0-:status:" +
	"methodOPTIONSpurpose:schemealt-svcreferernosniffCONNECTDELETEser" +
	"verFALSEclear:pathhttpsTRUEHEADdenydateetagvaryPOSTlink302404100" +
	"204206500400503421GET*/*103PUT403304425"

var staticTable = [...]headerField{
	{Name: intern[904:914]},                             // 0 :authority
	{Name: intern[1101:1106], Value: intern[108:109]},   // 1 :path: /
	{Name: intern[57:60], Value: intern[66:67]},         // 2 age: 0
	{Name: intern[545:564]},                             // 3 content-disposition
	{Name: intern[758:772], Value: intern[66:67]},       // 4 content-length: 0
	{Name: intern[898:904]},                             // 5 cookie
	{Name: intern[1123:1127]},                           // 6 date
	{Name: intern[1127:1131]},                           // 7 etag
	{Name: intern[601:618]},                             // 8 if-modified-since
	{Name: intern[825:838]},                             // 9 if-none-match
	{Name: intern[799:812]},                             // 10 last-modified
	{Name: intern[1139:1143]},                           // 11 link
	{Name: intern[984:992]},                             // 12 location
	{Name: intern[1058:1065]},                           // 13 referer
	{Name: intern[894:904]},                             // 14 set-cookie
	{Name: intern[1023:1030], Value: intern[1072:1079]}, // 15 :method: CONNECT
	{Name: intern[1023:1030], Value: intern[1079:1085]}, // 16 :method: DELETE
	{Name: intern[1023:1030], Value: intern[1170:1173]}, // 17 :method: GET
	{Name: intern[1023:1030], Value: intern[1115:1119]}, // 18 :method: HEAD
	{Name: intern[1023:1030], Value: intern[1030:1037]}, // 19 :method: OPTIONS
	{Name: intern[1023:1030], Value: intern[1135:1139]}, // 20 :method: POST
	{Name: intern[1023:1030], Value: intern[1179:1182]}, // 21 :method: PUT
	{Name: intern[1044:1051], Value: intern[1106:1110]}, // 22 :scheme: http
	{Name: intern[1044:1051], Value: intern[1106:1111]}, // 23 :scheme: https
	{Name: intern[1016:1023], Value: intern[1176:1179]}, // 24 :status: 103
	{Name: intern[1016:1023], Value: intern[709:712]},   // 25 :status: 200
	{Name: intern[1016:1023], Value: intern[1185:1188]}, // 26 :status: 304
	{Name: intern[1016:1023], Value: intern[1146:1149]}, // 27 :status: 404
	{Name: intern[1016:1023], Value: intern[1164:1167]}, // 28 :status: 503
	{Name: intern[728:734], Value: intern[1173:1176]},   // 29 accept: */*
	{Name: intern[728:734], Value: intern[455:478]},     // 30 accept: application/dns-message
	{Name: intern[743:758], Value: intern[618:635]},     // 31 accept-encoding: gzip, deflate, br
	{Name: intern[812:825], Value: intern[1008:1013]},   // 32 accept-ranges: bytes
	{Name: intern[250:278], Value: intern[786:799]},     // 33 access-control-allow-headers: cache-control
	{Name: intern[250:278], Value: intern[525:537]},     // 34 access-control-allow-headers: content-type
	{Name: intern[306:333], Value: intern[1173:1174]},   // 35 access-control-allow-origin: *
	{Name: intern[786:799], Value: intern[924:933]},     // 36 cache-control: max-age=0
	{Name: intern[786:799], Value: intern[698:713]},     // 37 cache-control: max-age=2592000
	{Name: intern[786:799], Value: intern[772:786]},     // 38 cache-control: max-age=604800
	{Name: intern[786:799], Value: intern[976:984]},     // 39 cache-control: no-cache
	{Name: intern[786:799], Value: intern[992:1000]},    // 40 cache-control: no-store
	{Name: intern[786:799], Value: intern[383:407]},     // 41 cache-control: public, max-age=31536000
	{Name: intern[651:667], Value: intern[633:635]},     // 42 content-encoding: br
	{Name: intern[651:667], Value: intern[618:622]},     // 43 content-encoding: gzip
	{Name: intern[525:537], Value: intern[455:478]},     // 44 content-type: application/dns-message
	{Name: intern[525:537], Value: intern[501:523]},     // 45 content-type: application/javascript
	{Name: intern[525:537], Value: intern[667:683]},     // 46 content-type: application/json
	{Name: intern[525:537], Value: intern[97:130]},      // 47 content-type: application/x-www-form-urlencoded
	{Name: intern[525:537], Value: intern[942:951]},     // 48 content-type: image/gif
	{Name: intern[525:537], Value: intern[864:874]},     // 49 content-type: image/jpeg
	{Name: intern[525:537], Value: intern[951:960]},     // 50 content-type: image/png
	{Name: intern[525:537], Value: intern[1000:1008]},   // 51 content-type: text/css
	{Name: intern[525:537], Value: intern[407:431]},     // 52 content-type: text/html; charset=utf-8
	{Name: intern[525:537], Value: intern[431:441]},     // 53 content-type: text/plain
	{Name: intern[525:537], Value: intern[431:455]},     // 54 content-type: text/plain;charset=utf-8
	{Name: intern[819:824], Value: intern[1008:1016]},   // 55 range: bytes=0-
	{Name: intern[333:358], Value: intern[53:69]},       // 56 strict-transport-security: max-age=31536000
	{Name: intern[333:358], Value: intern[53:88]},       // 57 strict-transport-security: max-age=31536000; includesubdomains
	{Name: intern[333:358], Value: intern[53:97]},       // 58 strict-transport-security: max-age=31536000; includesubdomains; preload
	{Name: intern[1131:1135], Value: intern[743:758]},   // 59 vary: accept-encoding
	{Name: intern[1131:1135], Value: intern[327:333]},   // 60 vary: origin
	{Name: intern[523:545], Value: intern[1065:1072]},   // 61 x-content-type-options: nosniff
	{Name: intern[635:651], Value: intern[851:864]},     // 62 x-xss-protection: 1; mode=block
	{Name: intern[1016:1023], Value: intern[1149:1152]}, // 63 :status: 100
	{Name: intern[1016:1023], Value: intern[1152:1155]}, // 64 :status: 204
	{Name: intern[1016:1023], Value: intern[1155:1158]}, // 65 :status: 206
	{Name: intern[1016:1023], Value: intern[1143:1146]}, // 66 :status: 302
	{Name: intern[1016:1023], Value: intern[1161:1164]}, // 67 :status: 400
	{Name: intern[1016:1023], Value: intern[1182:1185]}, // 68 :status: 403
	{Name: intern[1016:1023], Value: intern[1167:1170]}, // 69 :status: 421
	{Name: intern[1016:1023], Value: intern[1188:1191]}, // 70 :status: 425
	{Name: intern[1016:1023], Value: intern[1158:1161]}, // 71 :status: 500
	{Name: intern[728:743]},                             // 72 accept-language
	{Name: intern[130:162], Value: intern[1091:1096]},   // 73 access-control-allow-credentials: FALSE
	{Name: intern[130:162], Value: intern[1111:1115]},   // 74 access-control-allow-credentials: TRUE
	{Name: intern[250:278], Value: intern[1173:1174]},   // 75 access-control-allow-headers: *
	{Name: intern[278:306], Value: intern[583:586]},     // 76 access-control-allow-methods: get
	{Name: intern[278:306], Value: intern[583:601]},     // 77 access-control-allow-methods: get, post, options
	{Name: intern[278:306], Value: intern[538:545]},     // 78 access-control-allow-methods: options
	{Name: intern[192:221], Value: intern[758:772]},     // 79 access-control-expose-headers: content-length
	{Name: intern[162:192], Value: intern[525:537]},     // 80 access-control-request-headers: content-type
	{Name: intern[221:250], Value: intern[583:586]},     // 81 access-control-request-method: get
	{Name: intern[221:250], Value: intern[588:592]},     // 82 access-control-request-method: post
	{Name: intern[1051:1058], Value: intern[1096:1101]}, // 83 alt-svc: clear
	{Name: intern[838:851]},                             // 84 authorization
	{Name: intern[478:501], Value: intern[0:53]},        // 85 content-security-policy: script-src 'none'; object-src 'none'; base-uri 'none'
	{Name: intern[874:884], Value: intern[62:63]},       // 86 early-data: 1
	{Name: intern[933:942]},                             // 87 expect-ct
	{Name: intern[715:724]},                             // 88 forwarded
	{Name: intern[960:968]},                             // 89 if-range
	{Name: intern[327:333]},                             // 90 origin
	{Name: intern[1037:1044], Value: intern[968:976]},   // 91 purpose: prefetch
	{Name: intern[1085:1091]},                           // 92 server
	{Name: intern[564:583], Value: intern[1173:1174]},   // 93 timing-allow-origin: *
	{Name: intern[358:383], Value: intern[62:63]},       // 94 upgrade-insecure-requests: 1
	{Name: intern[914:924]},                             // 95 user-agent
	{Name: intern[713:728]},                             // 96 x-forwarded-for
	{Name: intern[683:698], Value: intern[1119:1123]},   // 97 x-frame-options: deny
	{Name: intern[683:698], Value: intern[884:894]},     // 98 x-frame-options: sameorigin
}
