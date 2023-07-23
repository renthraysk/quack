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

// appendAuthority appends an :authority pseudo header field to p
func appendAuthority(p []byte, authority string) []byte {
	if authority == "" {
		return inst.AppendIndexedLine(p, 0, true)
	}
	p = inst.AppendNamedReference(p, 0, false, true)
	return inst.AppendStringLiteral(p, authority, true)
}

// appendPath appends a :path pseudo header field to p
func appendPath(p []byte, path string) []byte {
	if path == "/" {
		return inst.AppendIndexedLine(p, 1, true)
	}
	p = inst.AppendNamedReference(p, 1, false, true)
	return inst.AppendStringLiteral(p, path, true)
}

// appendStatus appends a :status pseudo header field to p
func appendStatus(p []byte, statusCode int) []byte {
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

// appendMethod appends a :method pseudo header field to p
func appendMethod(p []byte, method string) []byte {
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
func appendScheme(p []byte, scheme string) []byte {
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
func appendDate(p []byte, t time.Time) []byte {
	const StaticTableIndex = 6

	const H = 0b1000_0000

	p = inst.AppendNamedReference(p, StaticTableIndex, false, true)
	// RFC1123 time length is less 0x7F so only need a single byte for length
	i := len(p)
	p = append(p, 0)
	p = huffman.AppendHttpTime(p, t)
	p[i] = H | uint8(len(p)-i-1)
	return p
}

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
		case "no-store":
			return 40, matchNameValue
		case "public, max-age=31536000":
			return 41, matchNameValue
		case "max-age=0":
			return 36, matchNameValue
		case "max-age=2592000":
			return 37, matchNameValue
		case "max-age=604800":
			return 38, matchNameValue
		case "no-cache":
			return 39, matchNameValue
		}
		return 40, matchName
	case "Content-Disposition":
		if value == "" {
			return 3, matchNameValue
		}
		return 3, matchName
	case "Content-Encoding":
		switch value {
		case "gzip":
			return 43, matchNameValue
		case "br":
			return 42, matchNameValue
		}
		return 43, matchName
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
		case "text/plain;charset=utf-8":
			return 54, matchNameValue
		case "application/dns-message":
			return 44, matchNameValue
		case "application/javascript":
			return 45, matchNameValue
		case "image/jpeg":
			return 49, matchNameValue
		case "image/png":
			return 50, matchNameValue
		case "text/css":
			return 51, matchNameValue
		case "text/html; charset=utf-8":
			return 52, matchNameValue
		case "text/plain":
			return 53, matchNameValue
		case "application/json":
			return 46, matchNameValue
		case "application/x-www-form-urlencoded":
			return 47, matchNameValue
		case "image/gif":
			return 48, matchNameValue
		}
		return 54, matchName
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
		case "max-age=31536000; includesubdomains":
			return 57, matchNameValue
		case "max-age=31536000; includesubdomains; preload":
			return 58, matchNameValue
		case "max-age=31536000":
			return 56, matchNameValue
		}
		return 57, matchName
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
