package quack

const intern = "" +
	":authority:path/Age0Content-DispositionContent-LengthCookieDateE" +
	"tagIf-Modified-SinceIf-None-MatchLast-ModifiedLinkLocationRefere" +
	"rSet-Cookie:methodCONNECTDELETEGETHEADOPTIONSPOSTPUT:schemehttph" +
	"ttps:status103200304404503Accept*/*application/dns-messageAccept" +
	"-Encodinggzip, deflate, brAccept-RangesbytesAccess-Control-Allow" +
	"-Headerscache-controlcontent-typeAccess-Control-Allow-Origin*Cac" +
	"he-Controlmax-age=0max-age=2592000max-age=604800no-cacheno-store" +
	"public, max-age=31536000Content-EncodingbrgzipContent-Typeapplic" +
	"ation/javascriptapplication/jsonapplication/x-www-form-urlencode" +
	"dimage/gifimage/jpegimage/pngtext/csstext/html; charset=utf-8tex" +
	"t/plaintext/plain;charset=utf-8Rangebytes=0-Strict-Transport-Sec" +
	"uritymax-age=31536000max-age=31536000; includesubdomainsmax-age=" +
	"31536000; includesubdomains; preloadVaryaccept-encodingoriginX-C" +
	"ontent-Type-OptionsnosniffX-Xss-Protection1; mode=block100204206" +
	"302400403421425500Accept-LanguageAccess-Control-Allow-Credential" +
	"sFALSETRUEAccess-Control-Allow-Methodsgetget, post, optionsoptio" +
	"nsAccess-Control-Expose-Headerscontent-lengthAccess-Control-Requ" +
	"est-HeadersAccess-Control-Request-MethodpostAlt-SvcclearAuthoriz" +
	"ationContent-Security-Policyscript-src 'none'; object-src 'none'" +
	"; base-uri 'none'Early-Data1Expect-CtForwardedIf-RangeOriginPurp" +
	"oseprefetchServerTiming-Allow-OriginUpgrade-Insecure-RequestsUse" +
	"r-AgentX-Forwarded-ForX-Frame-Optionsdenysameorigin"

var staticTable = [...]struct {
	Name  string
	Value string
}{
	{Name: intern[0:10]},                                // 0 :authority
	{Name: intern[10:15], Value: intern[15:16]},         // 1 :path: /
	{Name: intern[16:19], Value: intern[19:20]},         // 2 Age: 0
	{Name: intern[20:39]},                               // 3 Content-Disposition
	{Name: intern[39:53], Value: intern[19:20]},         // 4 Content-Length: 0
	{Name: intern[53:59]},                               // 5 Cookie
	{Name: intern[59:63]},                               // 6 Date
	{Name: intern[63:67]},                               // 7 Etag
	{Name: intern[67:84]},                               // 8 If-Modified-Since
	{Name: intern[84:97]},                               // 9 If-None-Match
	{Name: intern[97:110]},                              // 10 Last-Modified
	{Name: intern[110:114]},                             // 11 Link
	{Name: intern[114:122]},                             // 12 Location
	{Name: intern[122:129]},                             // 13 Referer
	{Name: intern[129:139]},                             // 14 Set-Cookie
	{Name: intern[139:146], Value: intern[146:153]},     // 15 :method: CONNECT
	{Name: intern[139:146], Value: intern[153:159]},     // 16 :method: DELETE
	{Name: intern[139:146], Value: intern[159:162]},     // 17 :method: GET
	{Name: intern[139:146], Value: intern[162:166]},     // 18 :method: HEAD
	{Name: intern[139:146], Value: intern[166:173]},     // 19 :method: OPTIONS
	{Name: intern[139:146], Value: intern[173:177]},     // 20 :method: POST
	{Name: intern[139:146], Value: intern[177:180]},     // 21 :method: PUT
	{Name: intern[180:187], Value: intern[187:191]},     // 22 :scheme: http
	{Name: intern[180:187], Value: intern[191:196]},     // 23 :scheme: https
	{Name: intern[196:203], Value: intern[203:206]},     // 24 :status: 103
	{Name: intern[196:203], Value: intern[206:209]},     // 25 :status: 200
	{Name: intern[196:203], Value: intern[209:212]},     // 26 :status: 304
	{Name: intern[196:203], Value: intern[212:215]},     // 27 :status: 404
	{Name: intern[196:203], Value: intern[215:218]},     // 28 :status: 503
	{Name: intern[218:224], Value: intern[224:227]},     // 29 Accept: */*
	{Name: intern[218:224], Value: intern[227:250]},     // 30 Accept: application/dns-message
	{Name: intern[250:265], Value: intern[265:282]},     // 31 Accept-Encoding: gzip, deflate, br
	{Name: intern[282:295], Value: intern[295:300]},     // 32 Accept-Ranges: bytes
	{Name: intern[300:328], Value: intern[328:341]},     // 33 Access-Control-Allow-Headers: cache-control
	{Name: intern[300:328], Value: intern[341:353]},     // 34 Access-Control-Allow-Headers: content-type
	{Name: intern[353:380], Value: intern[380:381]},     // 35 Access-Control-Allow-Origin: *
	{Name: intern[381:394], Value: intern[394:403]},     // 36 Cache-Control: max-age=0
	{Name: intern[381:394], Value: intern[403:418]},     // 37 Cache-Control: max-age=2592000
	{Name: intern[381:394], Value: intern[418:432]},     // 38 Cache-Control: max-age=604800
	{Name: intern[381:394], Value: intern[432:440]},     // 39 Cache-Control: no-cache
	{Name: intern[381:394], Value: intern[440:448]},     // 40 Cache-Control: no-store
	{Name: intern[381:394], Value: intern[448:472]},     // 41 Cache-Control: public, max-age=31536000
	{Name: intern[472:488], Value: intern[488:490]},     // 42 Content-Encoding: br
	{Name: intern[472:488], Value: intern[490:494]},     // 43 Content-Encoding: gzip
	{Name: intern[494:506], Value: intern[227:250]},     // 44 Content-Type: application/dns-message
	{Name: intern[494:506], Value: intern[506:528]},     // 45 Content-Type: application/javascript
	{Name: intern[494:506], Value: intern[528:544]},     // 46 Content-Type: application/json
	{Name: intern[494:506], Value: intern[544:577]},     // 47 Content-Type: application/x-www-form-urlencoded
	{Name: intern[494:506], Value: intern[577:586]},     // 48 Content-Type: image/gif
	{Name: intern[494:506], Value: intern[586:596]},     // 49 Content-Type: image/jpeg
	{Name: intern[494:506], Value: intern[596:605]},     // 50 Content-Type: image/png
	{Name: intern[494:506], Value: intern[605:613]},     // 51 Content-Type: text/css
	{Name: intern[494:506], Value: intern[613:637]},     // 52 Content-Type: text/html; charset=utf-8
	{Name: intern[494:506], Value: intern[637:647]},     // 53 Content-Type: text/plain
	{Name: intern[494:506], Value: intern[647:671]},     // 54 Content-Type: text/plain;charset=utf-8
	{Name: intern[671:676], Value: intern[676:684]},     // 55 Range: bytes=0-
	{Name: intern[684:709], Value: intern[709:725]},     // 56 Strict-Transport-Security: max-age=31536000
	{Name: intern[684:709], Value: intern[725:760]},     // 57 Strict-Transport-Security: max-age=31536000; includesubdomains
	{Name: intern[684:709], Value: intern[760:804]},     // 58 Strict-Transport-Security: max-age=31536000; includesubdomains; preload
	{Name: intern[804:808], Value: intern[808:823]},     // 59 Vary: accept-encoding
	{Name: intern[804:808], Value: intern[823:829]},     // 60 Vary: origin
	{Name: intern[829:851], Value: intern[851:858]},     // 61 X-Content-Type-Options: nosniff
	{Name: intern[858:874], Value: intern[874:887]},     // 62 X-Xss-Protection: 1; mode=block
	{Name: intern[196:203], Value: intern[887:890]},     // 63 :status: 100
	{Name: intern[196:203], Value: intern[890:893]},     // 64 :status: 204
	{Name: intern[196:203], Value: intern[893:896]},     // 65 :status: 206
	{Name: intern[196:203], Value: intern[896:899]},     // 66 :status: 302
	{Name: intern[196:203], Value: intern[899:902]},     // 67 :status: 400
	{Name: intern[196:203], Value: intern[902:905]},     // 68 :status: 403
	{Name: intern[196:203], Value: intern[905:908]},     // 69 :status: 421
	{Name: intern[196:203], Value: intern[908:911]},     // 70 :status: 425
	{Name: intern[196:203], Value: intern[911:914]},     // 71 :status: 500
	{Name: intern[914:929]},                             // 72 Accept-Language
	{Name: intern[929:961], Value: intern[961:966]},     // 73 Access-Control-Allow-Credentials: FALSE
	{Name: intern[929:961], Value: intern[966:970]},     // 74 Access-Control-Allow-Credentials: TRUE
	{Name: intern[300:328], Value: intern[380:381]},     // 75 Access-Control-Allow-Headers: *
	{Name: intern[970:998], Value: intern[998:1001]},    // 76 Access-Control-Allow-Methods: get
	{Name: intern[970:998], Value: intern[1001:1019]},   // 77 Access-Control-Allow-Methods: get, post, options
	{Name: intern[970:998], Value: intern[1019:1026]},   // 78 Access-Control-Allow-Methods: options
	{Name: intern[1026:1055], Value: intern[1055:1069]}, // 79 Access-Control-Expose-Headers: content-length
	{Name: intern[1069:1099], Value: intern[341:353]},   // 80 Access-Control-Request-Headers: content-type
	{Name: intern[1099:1128], Value: intern[998:1001]},  // 81 Access-Control-Request-Method: get
	{Name: intern[1099:1128], Value: intern[1128:1132]}, // 82 Access-Control-Request-Method: post
	{Name: intern[1132:1139], Value: intern[1139:1144]}, // 83 Alt-Svc: clear
	{Name: intern[1144:1157]},                           // 84 Authorization
	{Name: intern[1157:1180], Value: intern[1180:1233]}, // 85 Content-Security-Policy: script-src 'none'; object-src 'none'; base-uri 'none'
	{Name: intern[1233:1243], Value: intern[1243:1244]}, // 86 Early-Data: 1
	{Name: intern[1244:1253]},                           // 87 Expect-Ct
	{Name: intern[1253:1262]},                           // 88 Forwarded
	{Name: intern[1262:1270]},                           // 89 If-Range
	{Name: intern[1270:1276]},                           // 90 Origin
	{Name: intern[1276:1283], Value: intern[1283:1291]}, // 91 Purpose: prefetch
	{Name: intern[1291:1297]},                           // 92 Server
	{Name: intern[1297:1316], Value: intern[380:381]},   // 93 Timing-Allow-Origin: *
	{Name: intern[1316:1341], Value: intern[1243:1244]}, // 94 Upgrade-Insecure-Requests: 1
	{Name: intern[1341:1351]},                           // 95 User-Agent
	{Name: intern[1351:1366]},                           // 96 X-Forwarded-For
	{Name: intern[1366:1381], Value: intern[1381:1385]}, // 97 X-Frame-Options: deny
	{Name: intern[1366:1381], Value: intern[1385:1395]}, // 98 X-Frame-Options: sameorigin
}
