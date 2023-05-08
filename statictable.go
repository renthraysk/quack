package quack

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

var staticTable = [...]headerField{
	{Name: intern[958:968]},                             // 0 :authority
	{Name: intern[1155:1160], Value: intern[108:109]},   // 1 :path: /
	{Name: intern[973:976], Value: intern[66:67]},       // 2 Age: 0
	{Name: intern[545:564]},                             // 3 Content-Disposition
	{Name: intern[801:815], Value: intern[66:67]},       // 4 Content-Length: 0
	{Name: intern[952:958]},                             // 5 Cookie
	{Name: intern[1177:1181]},                           // 6 Date
	{Name: intern[1181:1185]},                           // 7 Etag
	{Name: intern[601:618]},                             // 8 If-Modified-Since
	{Name: intern[854:867]},                             // 9 If-None-Match
	{Name: intern[828:841]},                             // 10 Last-Modified
	{Name: intern[1193:1197]},                           // 11 Link
	{Name: intern[1038:1046]},                           // 12 Location
	{Name: intern[1112:1119]},                           // 13 Referer
	{Name: intern[948:958]},                             // 14 Set-Cookie
	{Name: intern[1077:1084], Value: intern[1126:1133]}, // 15 :method: CONNECT
	{Name: intern[1077:1084], Value: intern[1133:1139]}, // 16 :method: DELETE
	{Name: intern[1077:1084], Value: intern[1224:1227]}, // 17 :method: GET
	{Name: intern[1077:1084], Value: intern[1169:1173]}, // 18 :method: HEAD
	{Name: intern[1077:1084], Value: intern[1084:1091]}, // 19 :method: OPTIONS
	{Name: intern[1077:1084], Value: intern[1189:1193]}, // 20 :method: POST
	{Name: intern[1077:1084], Value: intern[1233:1236]}, // 21 :method: PUT
	{Name: intern[1098:1105], Value: intern[1160:1164]}, // 22 :scheme: http
	{Name: intern[1098:1105], Value: intern[1160:1165]}, // 23 :scheme: https
	{Name: intern[1070:1077], Value: intern[1230:1233]}, // 24 :status: 103
	{Name: intern[1070:1077], Value: intern[709:712]},   // 25 :status: 200
	{Name: intern[1070:1077], Value: intern[1239:1242]}, // 26 :status: 304
	{Name: intern[1070:1077], Value: intern[1200:1203]}, // 27 :status: 404
	{Name: intern[1070:1077], Value: intern[1218:1221]}, // 28 :status: 503
	{Name: intern[728:734], Value: intern[1227:1230]},   // 29 Accept: */*
	{Name: intern[728:734], Value: intern[455:478]},     // 30 Accept: application/dns-message
	{Name: intern[758:773], Value: intern[618:635]},     // 31 Accept-Encoding: gzip, deflate, br
	{Name: intern[841:854], Value: intern[1062:1067]},   // 32 Accept-Ranges: bytes
	{Name: intern[250:278], Value: intern[867:880]},     // 33 Access-Control-Allow-Headers: cache-control
	{Name: intern[250:278], Value: intern[906:918]},     // 34 Access-Control-Allow-Headers: content-type
	{Name: intern[306:333], Value: intern[1227:1228]},   // 35 Access-Control-Allow-Origin: *
	{Name: intern[815:828], Value: intern[978:987]},     // 36 Cache-Control: max-age=0
	{Name: intern[815:828], Value: intern[698:713]},     // 37 Cache-Control: max-age=2592000
	{Name: intern[815:828], Value: intern[787:801]},     // 38 Cache-Control: max-age=604800
	{Name: intern[815:828], Value: intern[1030:1038]},   // 39 Cache-Control: no-cache
	{Name: intern[815:828], Value: intern[1046:1054]},   // 40 Cache-Control: no-store
	{Name: intern[815:828], Value: intern[383:407]},     // 41 Cache-Control: public, max-age=31536000
	{Name: intern[651:667], Value: intern[633:635]},     // 42 Content-Encoding: br
	{Name: intern[651:667], Value: intern[618:622]},     // 43 Content-Encoding: gzip
	{Name: intern[525:537], Value: intern[455:478]},     // 44 Content-Type: application/dns-message
	{Name: intern[525:537], Value: intern[501:523]},     // 45 Content-Type: application/javascript
	{Name: intern[525:537], Value: intern[667:683]},     // 46 Content-Type: application/json
	{Name: intern[525:537], Value: intern[97:130]},      // 47 Content-Type: application/x-www-form-urlencoded
	{Name: intern[525:537], Value: intern[996:1005]},    // 48 Content-Type: image/gif
	{Name: intern[525:537], Value: intern[918:928]},     // 49 Content-Type: image/jpeg
	{Name: intern[525:537], Value: intern[1005:1014]},   // 50 Content-Type: image/png
	{Name: intern[525:537], Value: intern[1054:1062]},   // 51 Content-Type: text/css
	{Name: intern[525:537], Value: intern[407:431]},     // 52 Content-Type: text/html; charset=utf-8
	{Name: intern[525:537], Value: intern[431:441]},     // 53 Content-Type: text/plain
	{Name: intern[525:537], Value: intern[431:455]},     // 54 Content-Type: text/plain;charset=utf-8
	{Name: intern[848:853], Value: intern[1062:1070]},   // 55 Range: bytes=0-
	{Name: intern[333:358], Value: intern[53:69]},       // 56 Strict-Transport-Security: max-age=31536000
	{Name: intern[333:358], Value: intern[53:88]},       // 57 Strict-Transport-Security: max-age=31536000; includesubdomains
	{Name: intern[333:358], Value: intern[53:97]},       // 58 Strict-Transport-Security: max-age=31536000; includesubdomains; preload
	{Name: intern[1185:1189], Value: intern[743:758]},   // 59 Vary: accept-encoding
	{Name: intern[1185:1189], Value: intern[942:948]},   // 60 Vary: origin
	{Name: intern[523:545], Value: intern[1119:1126]},   // 61 X-Content-Type-Options: nosniff
	{Name: intern[635:651], Value: intern[893:906]},     // 62 X-Xss-Protection: 1; mode=block
	{Name: intern[1070:1077], Value: intern[1203:1206]}, // 63 :status: 100
	{Name: intern[1070:1077], Value: intern[1206:1209]}, // 64 :status: 204
	{Name: intern[1070:1077], Value: intern[1209:1212]}, // 65 :status: 206
	{Name: intern[1070:1077], Value: intern[1197:1200]}, // 66 :status: 302
	{Name: intern[1070:1077], Value: intern[1215:1218]}, // 67 :status: 400
	{Name: intern[1070:1077], Value: intern[1236:1239]}, // 68 :status: 403
	{Name: intern[1070:1077], Value: intern[1221:1224]}, // 69 :status: 421
	{Name: intern[1070:1077], Value: intern[1242:1245]}, // 70 :status: 425
	{Name: intern[1070:1077], Value: intern[1212:1215]}, // 71 :status: 500
	{Name: intern[728:743]},                             // 72 Accept-Language
	{Name: intern[130:162], Value: intern[1145:1150]},   // 73 Access-Control-Allow-Credentials: FALSE
	{Name: intern[130:162], Value: intern[1165:1169]},   // 74 Access-Control-Allow-Credentials: TRUE
	{Name: intern[250:278], Value: intern[1227:1228]},   // 75 Access-Control-Allow-Headers: *
	{Name: intern[278:306], Value: intern[583:586]},     // 76 Access-Control-Allow-Methods: get
	{Name: intern[278:306], Value: intern[583:601]},     // 77 Access-Control-Allow-Methods: get, post, options
	{Name: intern[278:306], Value: intern[594:601]},     // 78 Access-Control-Allow-Methods: options
	{Name: intern[192:221], Value: intern[773:787]},     // 79 Access-Control-Expose-Headers: content-length
	{Name: intern[162:192], Value: intern[906:918]},     // 80 Access-Control-Request-Headers: content-type
	{Name: intern[221:250], Value: intern[583:586]},     // 81 Access-Control-Request-Method: get
	{Name: intern[221:250], Value: intern[588:592]},     // 82 Access-Control-Request-Method: post
	{Name: intern[1105:1112], Value: intern[1150:1155]}, // 83 Alt-Svc: clear
	{Name: intern[880:893]},                             // 84 Authorization
	{Name: intern[478:501], Value: intern[0:53]},        // 85 Content-Security-Policy: script-src 'none'; object-src 'none'; base-uri 'none'
	{Name: intern[928:938], Value: intern[62:63]},       // 86 Early-Data: 1
	{Name: intern[987:996]},                             // 87 Expect-Ct
	{Name: intern[715:724]},                             // 88 Forwarded
	{Name: intern[1014:1022]},                           // 89 If-Range
	{Name: intern[327:333]},                             // 90 Origin
	{Name: intern[1091:1098], Value: intern[1022:1030]}, // 91 Purpose: prefetch
	{Name: intern[1139:1145]},                           // 92 Server
	{Name: intern[564:583], Value: intern[1227:1228]},   // 93 Timing-Allow-Origin: *
	{Name: intern[358:383], Value: intern[62:63]},       // 94 Upgrade-Insecure-Requests: 1
	{Name: intern[968:978]},                             // 95 User-Agent
	{Name: intern[713:728]},                             // 96 X-Forwarded-For
	{Name: intern[683:698], Value: intern[1173:1177]},   // 97 X-Frame-Options: deny
	{Name: intern[683:698], Value: intern[938:948]},     // 98 X-Frame-Options: sameorigin
}
