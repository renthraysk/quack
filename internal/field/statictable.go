package field

const intern = "" +
	"script-src 'none'; object-src 'none'; base-uri 'none'max-age=315" +
	"36000; includesubdomains; preloadapplication/x-www-form-urlencod" +
	"edAccess-Control-Allow-CredentialsAccess-Control-Request-Headers" +
	"Access-Control-Expose-HeadersAccess-Control-Request-MethodAccess" +
	"-Control-Allow-HeadersAccess-Control-Allow-MethodsAccess-Control" +
	"-Allow-OriginStrict-Transport-SecurityUpgrade-Insecure-Requestsp" +
	"ublic, max-age=31536000text/html; charset=utf-8text/plain;charse" +
	"t=utf-8application/dns-messageContent-Security-Policyapplication" +
	"/javascriptX-Content-Type-OptionsContent-DispositionTiming-Allow" +
	"-Originget, post, optionsIf-Modified-Sincegzip, deflate, brX-Xss" +
	"-ProtectionContent-Encodingapplication/jsonX-Frame-Optionsmax-ag" +
	"e=2592000X-Forwarded-ForAccept-Languageaccept-encodingAccept-Enc" +
	"odingcontent-lengthmax-age=604800Content-LengthCache-ControlLast" +
	"-ModifiedAccept-RangesIf-None-Matchcache-controlAuthorization1; " +
	"mode=blockcontent-typeimage/jpegEarly-DatasameoriginSet-Cookie:a" +
	"uthorityUser-Agentmax-age=0Expect-Ctimage/gifimage/pngIf-Rangepr" +
	"efetchno-cacheLocationno-storetext/cssbytes=0-:status:methodOPTI" +
	"ONSPurpose:schemeAlt-SvcReferernosniffCONNECTDELETEServerFALSEcl" +
	"ear:pathhttpsTRUEHEADdenyDateEtagVaryPOSTLink3024041002042065004" +
	"00503421GET*/*103PUT403304425"

var staticTable = [...]header{
	{name: intern[958:968]},                             // 0 :authority
	{name: intern[1155:1160], value: intern[108:109]},   // 1 :path: /
	{name: intern[973:976], value: intern[66:67]},       // 2 Age: 0
	{name: intern[545:564]},                             // 3 Content-Disposition
	{name: intern[801:815], value: intern[66:67]},       // 4 Content-Length: 0
	{name: intern[952:958]},                             // 5 Cookie
	{name: intern[1177:1181]},                           // 6 Date
	{name: intern[1181:1185]},                           // 7 Etag
	{name: intern[601:618]},                             // 8 If-Modified-Since
	{name: intern[854:867]},                             // 9 If-None-Match
	{name: intern[828:841]},                             // 10 Last-Modified
	{name: intern[1193:1197]},                           // 11 Link
	{name: intern[1038:1046]},                           // 12 Location
	{name: intern[1112:1119]},                           // 13 Referer
	{name: intern[948:958]},                             // 14 Set-Cookie
	{name: intern[1077:1084], value: intern[1126:1133]}, // 15 :method: CONNECT
	{name: intern[1077:1084], value: intern[1133:1139]}, // 16 :method: DELETE
	{name: intern[1077:1084], value: intern[1224:1227]}, // 17 :method: GET
	{name: intern[1077:1084], value: intern[1169:1173]}, // 18 :method: HEAD
	{name: intern[1077:1084], value: intern[1084:1091]}, // 19 :method: OPTIONS
	{name: intern[1077:1084], value: intern[1189:1193]}, // 20 :method: POST
	{name: intern[1077:1084], value: intern[1233:1236]}, // 21 :method: PUT
	{name: intern[1098:1105], value: intern[1160:1164]}, // 22 :scheme: http
	{name: intern[1098:1105], value: intern[1160:1165]}, // 23 :scheme: https
	{name: intern[1070:1077], value: intern[1230:1233]}, // 24 :status: 103
	{name: intern[1070:1077], value: intern[709:712]},   // 25 :status: 200
	{name: intern[1070:1077], value: intern[1239:1242]}, // 26 :status: 304
	{name: intern[1070:1077], value: intern[1200:1203]}, // 27 :status: 404
	{name: intern[1070:1077], value: intern[1218:1221]}, // 28 :status: 503
	{name: intern[728:734], value: intern[1227:1230]},   // 29 Accept: */*
	{name: intern[728:734], value: intern[455:478]},     // 30 Accept: application/dns-message
	{name: intern[758:773], value: intern[618:635]},     // 31 Accept-Encoding: gzip, deflate, br
	{name: intern[841:854], value: intern[1062:1067]},   // 32 Accept-Ranges: bytes
	{name: intern[250:278], value: intern[867:880]},     // 33 Access-Control-Allow-Headers: cache-control
	{name: intern[250:278], value: intern[906:918]},     // 34 Access-Control-Allow-Headers: content-type
	{name: intern[306:333], value: intern[1227:1228]},   // 35 Access-Control-Allow-Origin: *
	{name: intern[815:828], value: intern[978:987]},     // 36 Cache-Control: max-age=0
	{name: intern[815:828], value: intern[698:713]},     // 37 Cache-Control: max-age=2592000
	{name: intern[815:828], value: intern[787:801]},     // 38 Cache-Control: max-age=604800
	{name: intern[815:828], value: intern[1030:1038]},   // 39 Cache-Control: no-cache
	{name: intern[815:828], value: intern[1046:1054]},   // 40 Cache-Control: no-store
	{name: intern[815:828], value: intern[383:407]},     // 41 Cache-Control: public, max-age=31536000
	{name: intern[651:667], value: intern[633:635]},     // 42 Content-Encoding: br
	{name: intern[651:667], value: intern[618:622]},     // 43 Content-Encoding: gzip
	{name: intern[525:537], value: intern[455:478]},     // 44 Content-Type: application/dns-message
	{name: intern[525:537], value: intern[501:523]},     // 45 Content-Type: application/javascript
	{name: intern[525:537], value: intern[667:683]},     // 46 Content-Type: application/json
	{name: intern[525:537], value: intern[97:130]},      // 47 Content-Type: application/x-www-form-urlencoded
	{name: intern[525:537], value: intern[996:1005]},    // 48 Content-Type: image/gif
	{name: intern[525:537], value: intern[918:928]},     // 49 Content-Type: image/jpeg
	{name: intern[525:537], value: intern[1005:1014]},   // 50 Content-Type: image/png
	{name: intern[525:537], value: intern[1054:1062]},   // 51 Content-Type: text/css
	{name: intern[525:537], value: intern[407:431]},     // 52 Content-Type: text/html; charset=utf-8
	{name: intern[525:537], value: intern[431:441]},     // 53 Content-Type: text/plain
	{name: intern[525:537], value: intern[431:455]},     // 54 Content-Type: text/plain;charset=utf-8
	{name: intern[848:853], value: intern[1062:1070]},   // 55 Range: bytes=0-
	{name: intern[333:358], value: intern[53:69]},       // 56 Strict-Transport-Security: max-age=31536000
	{name: intern[333:358], value: intern[53:88]},       // 57 Strict-Transport-Security: max-age=31536000; includesubdomains
	{name: intern[333:358], value: intern[53:97]},       // 58 Strict-Transport-Security: max-age=31536000; includesubdomains; preload
	{name: intern[1185:1189], value: intern[743:758]},   // 59 Vary: accept-encoding
	{name: intern[1185:1189], value: intern[942:948]},   // 60 Vary: origin
	{name: intern[523:545], value: intern[1119:1126]},   // 61 X-Content-Type-Options: nosniff
	{name: intern[635:651], value: intern[893:906]},     // 62 X-Xss-Protection: 1; mode=block
	{name: intern[1070:1077], value: intern[1203:1206]}, // 63 :status: 100
	{name: intern[1070:1077], value: intern[1206:1209]}, // 64 :status: 204
	{name: intern[1070:1077], value: intern[1209:1212]}, // 65 :status: 206
	{name: intern[1070:1077], value: intern[1197:1200]}, // 66 :status: 302
	{name: intern[1070:1077], value: intern[1215:1218]}, // 67 :status: 400
	{name: intern[1070:1077], value: intern[1236:1239]}, // 68 :status: 403
	{name: intern[1070:1077], value: intern[1221:1224]}, // 69 :status: 421
	{name: intern[1070:1077], value: intern[1242:1245]}, // 70 :status: 425
	{name: intern[1070:1077], value: intern[1212:1215]}, // 71 :status: 500
	{name: intern[728:743]},                             // 72 Accept-Language
	{name: intern[130:162], value: intern[1145:1150]},   // 73 Access-Control-Allow-Credentials: FALSE
	{name: intern[130:162], value: intern[1165:1169]},   // 74 Access-Control-Allow-Credentials: TRUE
	{name: intern[250:278], value: intern[1227:1228]},   // 75 Access-Control-Allow-Headers: *
	{name: intern[278:306], value: intern[583:586]},     // 76 Access-Control-Allow-Methods: get
	{name: intern[278:306], value: intern[583:601]},     // 77 Access-Control-Allow-Methods: get, post, options
	{name: intern[278:306], value: intern[594:601]},     // 78 Access-Control-Allow-Methods: options
	{name: intern[192:221], value: intern[773:787]},     // 79 Access-Control-Expose-Headers: content-length
	{name: intern[162:192], value: intern[906:918]},     // 80 Access-Control-Request-Headers: content-type
	{name: intern[221:250], value: intern[583:586]},     // 81 Access-Control-Request-Method: get
	{name: intern[221:250], value: intern[588:592]},     // 82 Access-Control-Request-Method: post
	{name: intern[1105:1112], value: intern[1150:1155]}, // 83 Alt-Svc: clear
	{name: intern[880:893]},                             // 84 Authorization
	{name: intern[478:501], value: intern[0:53]},        // 85 Content-Security-Policy: script-src 'none'; object-src 'none'; base-uri 'none'
	{name: intern[928:938], value: intern[62:63]},       // 86 Early-Data: 1
	{name: intern[987:996]},                             // 87 Expect-Ct
	{name: intern[715:724]},                             // 88 Forwarded
	{name: intern[1014:1022]},                           // 89 If-Range
	{name: intern[327:333]},                             // 90 Origin
	{name: intern[1091:1098], value: intern[1022:1030]}, // 91 Purpose: prefetch
	{name: intern[1139:1145]},                           // 92 Server
	{name: intern[564:583], value: intern[1227:1228]},   // 93 Timing-Allow-Origin: *
	{name: intern[358:383], value: intern[62:63]},       // 94 Upgrade-Insecure-Requests: 1
	{name: intern[968:978]},                             // 95 User-Agent
	{name: intern[713:728]},                             // 96 X-Forwarded-For
	{name: intern[683:698], value: intern[1173:1177]},   // 97 X-Frame-Options: deny
	{name: intern[683:698], value: intern[938:948]},     // 98 X-Frame-Options: sameorigin
}
