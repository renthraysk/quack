package field

import (
	"time"

	"github.com/renthraysk/quack/huffman"
	"github.com/renthraysk/quack/internal/inst"
)

// This is largely code generated, with a handful type specific methods added
// to enable conversion directly to QPACK encoding without need for
// an intermediary string (allocation).

// The pseudo headers

// AppendAuthority appends an :authority pseudo header field to p
func AppendAuthority(p []byte, authority string) []byte {
	if authority == "" {
		return inst.AppendIndexedLine(p, 0, true)
	}
	p = inst.AppendNamedReference(p, 0, false, true)
	return inst.AppendStringLiteral(p, authority, true)
}

// AppendPath appends a :path pseudo header field to p
func AppendPath(p []byte, path string) []byte {
	if path == "/" {
		return inst.AppendIndexedLine(p, 1, true)
	}
	p = inst.AppendNamedReference(p, 1, false, true)
	return inst.AppendStringLiteral(p, path, true)
}

// appendStatus appends a :status pseudo header field to p
func AppendStatus(p []byte, statusCode int) []byte {
	switch statusCode {
	case 100:
		return inst.AppendIndexedLine(p, 63, true)
	case 103:
		return inst.AppendIndexedLine(p, 24, true)
	case 200:
		return inst.AppendIndexedLine(p, 25, true)
	case 204:
		return inst.AppendIndexedLine(p, 64, true)
	case 206:
		return inst.AppendIndexedLine(p, 65, true)
	case 302:
		return inst.AppendIndexedLine(p, 66, true)
	case 304:
		return inst.AppendIndexedLine(p, 26, true)
	case 400:
		return inst.AppendIndexedLine(p, 67, true)
	case 403:
		return inst.AppendIndexedLine(p, 68, true)
	case 404:
		return inst.AppendIndexedLine(p, 27, true)
	case 421:
		return inst.AppendIndexedLine(p, 69, true)
	case 425:
		return inst.AppendIndexedLine(p, 70, true)
	case 500:
		return inst.AppendIndexedLine(p, 71, true)
	case 503:
		return inst.AppendIndexedLine(p, 28, true)
	}
	p = inst.AppendNamedReference(p, 24, false, true)
	return appendInt(p, int64(statusCode))
}

// appendInt appends the QPACK string literal representation of int64 i.
func appendInt(p []byte, i int64) []byte {
	// H HuffmanEncoded
	const H = 0b1000_0000

	if -9 <= i && i <= 99 {
		// No savings from huffman encoding 2 characters.
		if i < 0 {
			return append(p, 2, '-', byte('0'-i))
		}
		if i <= 9 {
			return append(p, 1, byte(i)+'0')
		}
		j := i / 10
		return append(p, 2, byte(j)+'0', byte(i-10*j)+'0')
	}

	j := len(p)
	p = append(p, 0)
	p = huffman.AppendInt(p, i)
	p[j] = H | uint8(len(p)-j-1)
	return p
}

// AppendMethod appends a :method pseudo header field to p
func AppendMethod(p []byte, method string) []byte {
	switch method {
	case "CONNECT":
		return inst.AppendIndexedLine(p, 15, true)
	case "DELETE":
		return inst.AppendIndexedLine(p, 16, true)
	case "GET":
		return inst.AppendIndexedLine(p, 17, true)
	case "HEAD":
		return inst.AppendIndexedLine(p, 18, true)
	case "OPTIONS":
		return inst.AppendIndexedLine(p, 19, true)
	case "POST":
		return inst.AppendIndexedLine(p, 20, true)
	case "PUT":
		return inst.AppendIndexedLine(p, 21, true)
	}
	p = inst.AppendNamedReference(p, 15, false, true)
	return inst.AppendStringLiteral(p, method, true)
}

// appendScheme appends a :scheme pseudo header field to p
func AppendScheme(p []byte, scheme string) []byte {
	switch scheme {
	case "http":
		return inst.AppendIndexedLine(p, 22, true)
	case "https":
		return inst.AppendIndexedLine(p, 23, true)
	}
	p = inst.AppendNamedReference(p, 22, false, true)
	return inst.AppendStringLiteral(p, scheme, true)
}

// Regular headers

// appendDate appends a Date header field with time t.
func AppendDate(p []byte, t time.Time) []byte {
	const StaticTableIndex = 6

	const H = 0b1000_0000

	p = inst.AppendNamedReference(p, StaticTableIndex, false, true)
	// RFC1123 time length is less 0x7F so only need a single byte for length
	i := len(p)
	p = append(p, 0)
	p = huffman.AppendRFC1123Time(p, t)
	p[i] = H | uint8(len(p)-i-1)
	return p
}

func staticLookup(name, value string) (uint64, Match) {
	switch name {
	case "Accept":
		switch value {
		case "*/*":
			return 29, MatchNameValue
		case "application/dns-message":
			return 30, MatchNameValue
		}
		return 29, MatchName
	case "Accept-Encoding":
		if value == "gzip, deflate, br" {
			return 31, MatchNameValue
		}
		return 31, MatchName
	case "Accept-Language":
		if value == "" {
			return 72, MatchNameValue
		}
		return 72, MatchName
	case "Accept-Ranges":
		if value == "bytes" {
			return 32, MatchNameValue
		}
		return 32, MatchName
	case "Access-Control-Allow-Credentials":
		switch value {
		case "FALSE":
			return 73, MatchNameValue
		case "TRUE":
			return 74, MatchNameValue
		}
		return 73, MatchName
	case "Access-Control-Allow-Headers":
		switch value {
		case "cache-control":
			return 33, MatchNameValue
		case "content-type":
			return 34, MatchNameValue
		case "*":
			return 75, MatchNameValue
		}
		return 33, MatchName
	case "Access-Control-Allow-Methods":
		switch value {
		case "get":
			return 76, MatchNameValue
		case "get, post, options":
			return 77, MatchNameValue
		case "options":
			return 78, MatchNameValue
		}
		return 76, MatchName
	case "Access-Control-Allow-Origin":
		if value == "*" {
			return 35, MatchNameValue
		}
		return 35, MatchName
	case "Access-Control-Expose-Headers":
		if value == "content-length" {
			return 79, MatchNameValue
		}
		return 79, MatchName
	case "Access-Control-Request-Headers":
		if value == "content-type" {
			return 80, MatchNameValue
		}
		return 80, MatchName
	case "Access-Control-Request-Method":
		switch value {
		case "get":
			return 81, MatchNameValue
		case "post":
			return 82, MatchNameValue
		}
		return 81, MatchName
	case "Age":
		if value == "0" {
			return 2, MatchNameValue
		}
		return 2, MatchName
	case "Alt-Svc":
		if value == "clear" {
			return 83, MatchNameValue
		}
		return 83, MatchName
	case "Authorization":
		if value == "" {
			return 84, MatchNameValue
		}
		return 84, MatchName
	case "Cache-Control":
		switch value {
		case "no-store":
			return 40, MatchNameValue
		case "public, max-age=31536000":
			return 41, MatchNameValue
		case "max-age=0":
			return 36, MatchNameValue
		case "max-age=2592000":
			return 37, MatchNameValue
		case "max-age=604800":
			return 38, MatchNameValue
		case "no-cache":
			return 39, MatchNameValue
		}
		return 40, MatchName
	case "Content-Disposition":
		if value == "" {
			return 3, MatchNameValue
		}
		return 3, MatchName
	case "Content-Encoding":
		switch value {
		case "gzip":
			return 43, MatchNameValue
		case "br":
			return 42, MatchNameValue
		}
		return 43, MatchName
	case "Content-Length":
		if value == "0" {
			return 4, MatchNameValue
		}
		return 4, MatchName
	case "Content-Security-Policy":
		if value == "script-src 'none'; object-src 'none'; base-uri 'none'" {
			return 85, MatchNameValue
		}
		return 85, MatchName
	case "Content-Type":
		switch value {
		case "text/plain;charset=utf-8":
			return 54, MatchNameValue
		case "application/dns-message":
			return 44, MatchNameValue
		case "application/javascript":
			return 45, MatchNameValue
		case "image/jpeg":
			return 49, MatchNameValue
		case "image/png":
			return 50, MatchNameValue
		case "text/css":
			return 51, MatchNameValue
		case "text/html; charset=utf-8":
			return 52, MatchNameValue
		case "text/plain":
			return 53, MatchNameValue
		case "application/json":
			return 46, MatchNameValue
		case "application/x-www-form-urlencoded":
			return 47, MatchNameValue
		case "image/gif":
			return 48, MatchNameValue
		}
		return 54, MatchName
	case "Cookie":
		if value == "" {
			return 5, MatchNameValue
		}
		return 5, MatchName
	case "Date":
		if value == "" {
			return 6, MatchNameValue
		}
		return 6, MatchName
	case "Early-Data":
		if value == "1" {
			return 86, MatchNameValue
		}
		return 86, MatchName
	case "Etag":
		if value == "" {
			return 7, MatchNameValue
		}
		return 7, MatchName
	case "Expect-Ct":
		if value == "" {
			return 87, MatchNameValue
		}
		return 87, MatchName
	case "Forwarded":
		if value == "" {
			return 88, MatchNameValue
		}
		return 88, MatchName
	case "If-Modified-Since":
		if value == "" {
			return 8, MatchNameValue
		}
		return 8, MatchName
	case "If-None-Match":
		if value == "" {
			return 9, MatchNameValue
		}
		return 9, MatchName
	case "If-Range":
		if value == "" {
			return 89, MatchNameValue
		}
		return 89, MatchName
	case "Last-Modified":
		if value == "" {
			return 10, MatchNameValue
		}
		return 10, MatchName
	case "Link":
		if value == "" {
			return 11, MatchNameValue
		}
		return 11, MatchName
	case "Location":
		if value == "" {
			return 12, MatchNameValue
		}
		return 12, MatchName
	case "Origin":
		if value == "" {
			return 90, MatchNameValue
		}
		return 90, MatchName
	case "Purpose":
		if value == "prefetch" {
			return 91, MatchNameValue
		}
		return 91, MatchName
	case "Range":
		if value == "bytes=0-" {
			return 55, MatchNameValue
		}
		return 55, MatchName
	case "Referer":
		if value == "" {
			return 13, MatchNameValue
		}
		return 13, MatchName
	case "Server":
		if value == "" {
			return 92, MatchNameValue
		}
		return 92, MatchName
	case "Set-Cookie":
		if value == "" {
			return 14, MatchNameValue
		}
		return 14, MatchName
	case "Strict-Transport-Security":
		switch value {
		case "max-age=31536000; includesubdomains":
			return 57, MatchNameValue
		case "max-age=31536000; includesubdomains; preload":
			return 58, MatchNameValue
		case "max-age=31536000":
			return 56, MatchNameValue
		}
		return 57, MatchName
	case "Timing-Allow-Origin":
		if value == "*" {
			return 93, MatchNameValue
		}
		return 93, MatchName
	case "Upgrade-Insecure-Requests":
		if value == "1" {
			return 94, MatchNameValue
		}
		return 94, MatchName
	case "User-Agent":
		if value == "" {
			return 95, MatchNameValue
		}
		return 95, MatchName
	case "Vary":
		switch value {
		case "accept-encoding":
			return 59, MatchNameValue
		case "origin":
			return 60, MatchNameValue
		}
		return 59, MatchName
	case "X-Content-Type-Options":
		if value == "nosniff" {
			return 61, MatchNameValue
		}
		return 61, MatchName
	case "X-Forwarded-For":
		if value == "" {
			return 96, MatchNameValue
		}
		return 96, MatchName
	case "X-Frame-Options":
		switch value {
		case "deny":
			return 97, MatchNameValue
		case "sameorigin":
			return 98, MatchNameValue
		}
		return 97, MatchName
	case "X-Xss-Protection":
		if value == "1; mode=block" {
			return 62, MatchNameValue
		}
		return 62, MatchName
	}
	return 0, MatchNone
}
