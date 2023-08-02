// Code generated by "static_codegen" DO NOT EDIT.
package field

func staticLookup(name, value string) (index uint64, m match) {
	switch name {
	case "Accept":
		switch value {
		case "*/*":
			return 29, matchNameValue
		case "application/dns-message":
			return 30, matchNameValue
		}
		return 29, matchName
	case "Accept-Encoding":
		if value == "gzip, deflate, br" {
			return 31, matchNameValue
		}
		return 31, matchName
	case "Accept-Language":
		if value == "" {
			return 72, matchNameValue
		}
		return 72, matchName
	case "Accept-Ranges":
		if value == "bytes" {
			return 32, matchNameValue
		}
		return 32, matchName
	case "Access-Control-Allow-Credentials":
		switch value {
		case "FALSE":
			return 73, matchNameValue
		case "TRUE":
			return 74, matchNameValue
		}
		return 73, matchName
	case "Access-Control-Allow-Headers":
		switch value {
		case "cache-control":
			return 33, matchNameValue
		case "content-type":
			return 34, matchNameValue
		case "*":
			return 75, matchNameValue
		}
		return 33, matchName
	case "Access-Control-Allow-Methods":
		switch value {
		case "get":
			return 76, matchNameValue
		case "get, post, options":
			return 77, matchNameValue
		case "options":
			return 78, matchNameValue
		}
		return 76, matchName
	case "Access-Control-Allow-Origin":
		if value == "*" {
			return 35, matchNameValue
		}
		return 35, matchName
	case "Access-Control-Expose-Headers":
		if value == "content-length" {
			return 79, matchNameValue
		}
		return 79, matchName
	case "Access-Control-Request-Headers":
		if value == "content-type" {
			return 80, matchNameValue
		}
		return 80, matchName
	case "Access-Control-Request-Method":
		switch value {
		case "get":
			return 81, matchNameValue
		case "post":
			return 82, matchNameValue
		}
		return 81, matchName
	case "Age":
		if value == "0" {
			return 2, matchNameValue
		}
		return 2, matchName
	case "Alt-Svc":
		if value == "clear" {
			return 83, matchNameValue
		}
		return 83, matchName
	case "Authorization":
		if value == "" {
			return 84, matchNameValue
		}
		return 84, matchName
	case "Cache-Control":
		switch value {
		case "max-age=604800":
			return 38, matchNameValue
		case "no-cache":
			return 39, matchNameValue
		case "no-store":
			return 40, matchNameValue
		case "public, max-age=31536000":
			return 41, matchNameValue
		case "max-age=0":
			return 36, matchNameValue
		case "max-age=2592000":
			return 37, matchNameValue
		}
		return 38, matchName
	case "Content-Disposition":
		if value == "" {
			return 3, matchNameValue
		}
		return 3, matchName
	case "Content-Encoding":
		switch value {
		case "br":
			return 42, matchNameValue
		case "gzip":
			return 43, matchNameValue
		}
		return 42, matchName
	case "Content-Length":
		if value == "0" {
			return 4, matchNameValue
		}
		return 4, matchName
	case "Content-Security-Policy":
		if value == "script-src 'none'; object-src 'none'; base-uri 'none'" {
			return 85, matchNameValue
		}
		return 85, matchName
	case "Content-Type":
		switch value {
		case "application/javascript":
			return 45, matchNameValue
		case "application/json":
			return 46, matchNameValue
		case "application/x-www-form-urlencoded":
			return 47, matchNameValue
		case "image/png":
			return 50, matchNameValue
		case "text/css":
			return 51, matchNameValue
		case "text/plain":
			return 53, matchNameValue
		case "application/dns-message":
			return 44, matchNameValue
		case "image/gif":
			return 48, matchNameValue
		case "image/jpeg":
			return 49, matchNameValue
		case "text/html; charset=utf-8":
			return 52, matchNameValue
		case "text/plain;charset=utf-8":
			return 54, matchNameValue
		}
		return 45, matchName
	case "Cookie":
		if value == "" {
			return 5, matchNameValue
		}
		return 5, matchName
	case "Date":
		if value == "" {
			return 6, matchNameValue
		}
		return 6, matchName
	case "Early-Data":
		if value == "1" {
			return 86, matchNameValue
		}
		return 86, matchName
	case "Etag":
		if value == "" {
			return 7, matchNameValue
		}
		return 7, matchName
	case "Expect-Ct":
		if value == "" {
			return 87, matchNameValue
		}
		return 87, matchName
	case "Forwarded":
		if value == "" {
			return 88, matchNameValue
		}
		return 88, matchName
	case "If-Modified-Since":
		if value == "" {
			return 8, matchNameValue
		}
		return 8, matchName
	case "If-None-Match":
		if value == "" {
			return 9, matchNameValue
		}
		return 9, matchName
	case "If-Range":
		if value == "" {
			return 89, matchNameValue
		}
		return 89, matchName
	case "Last-Modified":
		if value == "" {
			return 10, matchNameValue
		}
		return 10, matchName
	case "Link":
		if value == "" {
			return 11, matchNameValue
		}
		return 11, matchName
	case "Location":
		if value == "" {
			return 12, matchNameValue
		}
		return 12, matchName
	case "Origin":
		if value == "" {
			return 90, matchNameValue
		}
		return 90, matchName
	case "Purpose":
		if value == "prefetch" {
			return 91, matchNameValue
		}
		return 91, matchName
	case "Range":
		if value == "bytes=0-" {
			return 55, matchNameValue
		}
		return 55, matchName
	case "Referer":
		if value == "" {
			return 13, matchNameValue
		}
		return 13, matchName
	case "Server":
		if value == "" {
			return 92, matchNameValue
		}
		return 92, matchName
	case "Set-Cookie":
		if value == "" {
			return 14, matchNameValue
		}
		return 14, matchName
	case "Strict-Transport-Security":
		switch value {
		case "max-age=31536000":
			return 56, matchNameValue
		case "max-age=31536000; includesubdomains":
			return 57, matchNameValue
		case "max-age=31536000; includesubdomains; preload":
			return 58, matchNameValue
		}
		return 56, matchName
	case "Timing-Allow-Origin":
		if value == "*" {
			return 93, matchNameValue
		}
		return 93, matchName
	case "Upgrade-Insecure-Requests":
		if value == "1" {
			return 94, matchNameValue
		}
		return 94, matchName
	case "User-Agent":
		if value == "" {
			return 95, matchNameValue
		}
		return 95, matchName
	case "Vary":
		switch value {
		case "accept-encoding":
			return 59, matchNameValue
		case "origin":
			return 60, matchNameValue
		}
		return 59, matchName
	case "X-Content-Type-Options":
		if value == "nosniff" {
			return 61, matchNameValue
		}
		return 61, matchName
	case "X-Forwarded-For":
		if value == "" {
			return 96, matchNameValue
		}
		return 96, matchName
	case "X-Frame-Options":
		switch value {
		case "deny":
			return 97, matchNameValue
		case "sameorigin":
			return 98, matchNameValue
		}
		return 97, matchName
	case "X-Xss-Protection":
		if value == "1; mode=block" {
			return 62, matchNameValue
		}
		return 62, matchName
	}
	return 0, matchNone
}

const intern string = "" +
	"script-src 'none'; object-src 'none'; base-uri 'none'max-age=315" +
	"36000; includesubdomains; preloadapplication/x-www-form-urlencod" +
	"edAccess-Control-Allow-CredentialsAccess-Control-Request-Headers" +
	"Access-Control-Expose-HeadersAccess-Control-Request-MethodAccess" +
	"-Control-Allow-HeadersAccess-Control-Allow-MethodsAccess-Control" +
	"-Allow-OriginStrict-Transport-SecurityUpgrade-Insecure-Requestsp" +
	"ublic, max-age=31536000text/html; charset=utf-8text/plain;charse" +
	"t=utf-8Content-Security-Policyapplication/dns-messageX-Content-T" +
	"ype-Optionsapplication/javascriptContent-DispositionTiming-Allow" +
	"-Originget, post, optionsIf-Modified-Sincegzip, deflate, brConte" +
	"nt-EncodingX-Xss-Protectionapplication/jsonAccept-EncodingAccept" +
	"-LanguageX-Forwarded-ForX-Frame-Optionsaccept-encodingmax-age=25" +
	"92000Content-Lengthcontent-lengthmax-age=6048001; mode=blockAcce" +
	"pt-RangesAuthorizationCache-ControlIf-None-MatchLast-Modifiedcac" +
	"he-controlcontent-type:authorityEarly-DataSet-CookieUser-Agentim" +
	"age/jpegsameoriginExpect-Ctimage/gifimage/pngmax-age=0If-RangeLo" +
	"cationbytes=0-no-cacheno-storeprefetchtext/css:method:scheme:sta" +
	"tusAlt-SvcCONNECTOPTIONSPurposeReferernosniffDELETEServer:pathFA" +
	"LSEclearhttpsDateEtagHEADLinkPOSTTRUEVarydeny*/*1001032042063023" +
	"04400403404421425500503GETPUT"

var staticTable = [...]header{
	{name: intern[918:928]},                             // 0 :authority
	{name: intern[1145:1150], value: intern[108:109]},   // 1 :path: /
	{name: intern[953:956], value: intern[66:67]},       // 2 Age: 0
	{name: intern[545:564]},                             // 3 Content-Disposition
	{name: intern[773:787], value: intern[66:67]},       // 4 Content-Length: 0
	{name: intern[942:948]},                             // 5 Cookie
	{name: intern[1165:1169]},                           // 6 Date
	{name: intern[1169:1173]},                           // 7 Etag
	{name: intern[601:618]},                             // 8 If-Modified-Since
	{name: intern[867:880]},                             // 9 If-None-Match
	{name: intern[880:893]},                             // 10 Last-Modified
	{name: intern[1177:1181]},                           // 11 Link
	{name: intern[1022:1030]},                           // 12 Location
	{name: intern[1119:1126]},                           // 13 Referer
	{name: intern[938:948]},                             // 14 Set-Cookie
	{name: intern[1070:1077], value: intern[1098:1105]}, // 15 :method: CONNECT
	{name: intern[1070:1077], value: intern[1133:1139]}, // 16 :method: DELETE
	{name: intern[1070:1077], value: intern[1239:1242]}, // 17 :method: GET
	{name: intern[1070:1077], value: intern[1173:1177]}, // 18 :method: HEAD
	{name: intern[1070:1077], value: intern[1105:1112]}, // 19 :method: OPTIONS
	{name: intern[1070:1077], value: intern[1181:1185]}, // 20 :method: POST
	{name: intern[1070:1077], value: intern[1242:1245]}, // 21 :method: PUT
	{name: intern[1077:1084], value: intern[1160:1164]}, // 22 :scheme: http
	{name: intern[1077:1084], value: intern[1160:1165]}, // 23 :scheme: https
	{name: intern[1084:1091], value: intern[1203:1206]}, // 24 :status: 103
	{name: intern[1084:1091], value: intern[769:772]},   // 25 :status: 200
	{name: intern[1084:1091], value: intern[1215:1218]}, // 26 :status: 304
	{name: intern[1084:1091], value: intern[1224:1227]}, // 27 :status: 404
	{name: intern[1084:1091], value: intern[1236:1239]}, // 28 :status: 503
	{name: intern[683:689], value: intern[1197:1200]},   // 29 Accept: */*
	{name: intern[683:689], value: intern[478:501]},     // 30 Accept: application/dns-message
	{name: intern[683:698], value: intern[618:635]},     // 31 Accept-Encoding: gzip, deflate, br
	{name: intern[828:841], value: intern[1030:1035]},   // 32 Accept-Ranges: bytes
	{name: intern[250:278], value: intern[893:906]},     // 33 Access-Control-Allow-Headers: cache-control
	{name: intern[250:278], value: intern[906:918]},     // 34 Access-Control-Allow-Headers: content-type
	{name: intern[306:333], value: intern[1197:1198]},   // 35 Access-Control-Allow-Origin: *
	{name: intern[854:867], value: intern[1005:1014]},   // 36 Cache-Control: max-age=0
	{name: intern[854:867], value: intern[758:773]},     // 37 Cache-Control: max-age=2592000
	{name: intern[854:867], value: intern[801:815]},     // 38 Cache-Control: max-age=604800
	{name: intern[854:867], value: intern[1038:1046]},   // 39 Cache-Control: no-cache
	{name: intern[854:867], value: intern[1046:1054]},   // 40 Cache-Control: no-store
	{name: intern[854:867], value: intern[383:407]},     // 41 Cache-Control: public, max-age=31536000
	{name: intern[635:651], value: intern[633:635]},     // 42 Content-Encoding: br
	{name: intern[635:651], value: intern[618:622]},     // 43 Content-Encoding: gzip
	{name: intern[503:515], value: intern[478:501]},     // 44 Content-Type: application/dns-message
	{name: intern[503:515], value: intern[523:545]},     // 45 Content-Type: application/javascript
	{name: intern[503:515], value: intern[667:683]},     // 46 Content-Type: application/json
	{name: intern[503:515], value: intern[97:130]},      // 47 Content-Type: application/x-www-form-urlencoded
	{name: intern[503:515], value: intern[987:996]},     // 48 Content-Type: image/gif
	{name: intern[503:515], value: intern[958:968]},     // 49 Content-Type: image/jpeg
	{name: intern[503:515], value: intern[996:1005]},    // 50 Content-Type: image/png
	{name: intern[503:515], value: intern[1062:1070]},   // 51 Content-Type: text/css
	{name: intern[503:515], value: intern[407:431]},     // 52 Content-Type: text/html; charset=utf-8
	{name: intern[503:515], value: intern[431:441]},     // 53 Content-Type: text/plain
	{name: intern[503:515], value: intern[431:455]},     // 54 Content-Type: text/plain;charset=utf-8
	{name: intern[835:840], value: intern[1030:1038]},   // 55 Range: bytes=0-
	{name: intern[333:358], value: intern[53:69]},       // 56 Strict-Transport-Security: max-age=31536000
	{name: intern[333:358], value: intern[53:88]},       // 57 Strict-Transport-Security: max-age=31536000; includesubdomains
	{name: intern[333:358], value: intern[53:97]},       // 58 Strict-Transport-Security: max-age=31536000; includesubdomains; preload
	{name: intern[1189:1193], value: intern[743:758]},   // 59 Vary: accept-encoding
	{name: intern[1189:1193], value: intern[972:978]},   // 60 Vary: origin
	{name: intern[501:523], value: intern[1126:1133]},   // 61 X-Content-Type-Options: nosniff
	{name: intern[651:667], value: intern[815:828]},     // 62 X-Xss-Protection: 1; mode=block
	{name: intern[1084:1091], value: intern[1200:1203]}, // 63 :status: 100
	{name: intern[1084:1091], value: intern[1206:1209]}, // 64 :status: 204
	{name: intern[1084:1091], value: intern[1209:1212]}, // 65 :status: 206
	{name: intern[1084:1091], value: intern[1212:1215]}, // 66 :status: 302
	{name: intern[1084:1091], value: intern[1218:1221]}, // 67 :status: 400
	{name: intern[1084:1091], value: intern[1221:1224]}, // 68 :status: 403
	{name: intern[1084:1091], value: intern[1227:1230]}, // 69 :status: 421
	{name: intern[1084:1091], value: intern[1230:1233]}, // 70 :status: 425
	{name: intern[1084:1091], value: intern[1233:1236]}, // 71 :status: 500
	{name: intern[698:713]},                             // 72 Accept-Language
	{name: intern[130:162], value: intern[1150:1155]},   // 73 Access-Control-Allow-Credentials: FALSE
	{name: intern[130:162], value: intern[1185:1189]},   // 74 Access-Control-Allow-Credentials: TRUE
	{name: intern[250:278], value: intern[1197:1198]},   // 75 Access-Control-Allow-Headers: *
	{name: intern[278:306], value: intern[583:586]},     // 76 Access-Control-Allow-Methods: get
	{name: intern[278:306], value: intern[583:601]},     // 77 Access-Control-Allow-Methods: get, post, options
	{name: intern[278:306], value: intern[594:601]},     // 78 Access-Control-Allow-Methods: options
	{name: intern[192:221], value: intern[787:801]},     // 79 Access-Control-Expose-Headers: content-length
	{name: intern[162:192], value: intern[906:918]},     // 80 Access-Control-Request-Headers: content-type
	{name: intern[221:250], value: intern[583:586]},     // 81 Access-Control-Request-Method: get
	{name: intern[221:250], value: intern[588:592]},     // 82 Access-Control-Request-Method: post
	{name: intern[1091:1098], value: intern[1155:1160]}, // 83 Alt-Svc: clear
	{name: intern[841:854]},                             // 84 Authorization
	{name: intern[455:478], value: intern[0:53]},        // 85 Content-Security-Policy: script-src 'none'; object-src 'none'; base-uri 'none'
	{name: intern[928:938], value: intern[62:63]},       // 86 Early-Data: 1
	{name: intern[978:987]},                             // 87 Expect-Ct
	{name: intern[715:724]},                             // 88 Forwarded
	{name: intern[1014:1022]},                           // 89 If-Range
	{name: intern[327:333]},                             // 90 Origin
	{name: intern[1112:1119], value: intern[1054:1062]}, // 91 Purpose: prefetch
	{name: intern[1139:1145]},                           // 92 Server
	{name: intern[564:583], value: intern[1197:1198]},   // 93 Timing-Allow-Origin: *
	{name: intern[358:383], value: intern[62:63]},       // 94 Upgrade-Insecure-Requests: 1
	{name: intern[948:958]},                             // 95 User-Agent
	{name: intern[713:728]},                             // 96 X-Forwarded-For
	{name: intern[728:743], value: intern[1193:1197]},   // 97 X-Frame-Options: deny
	{name: intern[728:743], value: intern[968:978]},     // 98 X-Frame-Options: sameorigin
}
