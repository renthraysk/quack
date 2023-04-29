package quack

import (
	"time"

	"github.com/renthraysk/quack/huffman"
)

// This is largely code generated, with a handful type specific methods added
// to enable conversion directly to QPACK encoding without need for
// an intermediary string (allocation).

// The pseudo headers

// appendAuthority appends an :authority pseudo header field to p
func appendAuthority(p []byte, authority string) []byte {
	if authority == "" {
		return append(p, 0xC0)
	}
	p = append(p, 0x50)
	return appendStringLiteral(p, authority, 0)
}

// appendPath appends a :path pseudo header field to p
func appendPath(p []byte, path string) []byte {
	if path == "/" {
		return append(p, 0xC1)
	}
	p = append(p, 0x51)
	return appendStringLiteral(p, path, 0)
}

// appendStatus appends a :status pseudo header field to p
func appendStatus(p []byte, statusCode int) []byte {
	var b byte

	switch statusCode {
	case 100:
		b = 0x00
	case 103:
		return append(p, 0xD8)
	case 200:
		return append(p, 0xD9)
	case 204:
		b = 0x01
	case 206:
		b = 0x02
	case 302:
		b = 0x03
	case 304:
		return append(p, 0xDA)
	case 400:
		b = 0x04
	case 403:
		b = 0x05
	case 404:
		return append(p, 0xDB)
	case 421:
		b = 0x06
	case 425:
		b = 0x07
	case 500:
		b = 0x08
	case 503:
		return append(p, 0xDC)
	default:
		p = append(p, 0x5F, 0x09)
		return appendInt(p, int64(statusCode))
	}
	return append(p, 0xFF, b)
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
	var b uint8
	switch method {
	case "DELETE":
		b = 0xD0
	case "GET":
		b = 0xD1
	case "HEAD":
		b = 0xD2
	case "OPTIONS":
		b = 0xD3
	case "POST":
		b = 0xD4
	case "PUT":
		b = 0xD5
	case "CONNECT":
		b = 0xCF
	default:
		p = append(p, 0x5F, 0x00)
		return appendStringLiteral(p, method, 0)
	}
	return append(p, b)
}

// appendScheme appends a :scheme pseudo header field to p
func appendScheme(p []byte, scheme string) []byte {
	var b uint8
	switch scheme {
	case "http":
		b = 0xD6
	case "https":
		b = 0xD7
	default:
		p = append(p, 0x5F, 0x07)
		return appendStringLiteral(p, scheme, 0)
	}
	return append(p, b)
}

// Regular headers

// appendDate appends a Date header field with time t.
func appendDate(p []byte, t time.Time) []byte {
	const (
		// '01' 2-bit Pattern
		P = 0b0100_0000
		// Never index bit
		N = 0b0010_0000
		// Static table bit
		T = 0b0001_0000
	)

	// String Literal
	const (
		// H HuffmanEncoded
		H = 0b1000_0000
	)

	// D Date static table index
	const D = 6

	// RFC1123 time length is less 0x7F so only need a single byte for length
	p = append(p, P|N|T|D, 0)
	i := len(p) - 1
	p = huffman.AppendRFC1123Time(p, t)
	p[i] = H | uint8(len(p)-i-1)
	return p
}

func staticLookup(name, value string) (uint64, match) {
	switch name {
	case "accept":
		switch value {
		case "*/*":
			return 29, matchNameValue
		case "application/dns-message":
			return 30, matchNameValue
		}
		return 29, matchName
	case "accept-encoding":
		if value == "gzip, deflate, br" {
			return 31, matchNameValue
		}
		return 31, matchName
	case "accept-language":
		if value == "" {
			return 72, matchNameValue
		}
		return 72, matchName
	case "accept-ranges":
		if value == "bytes" {
			return 32, matchNameValue
		}
		return 32, matchName
	case "access-control-allow-credentials":
		switch value {
		case "FALSE":
			return 73, matchNameValue
		case "TRUE":
			return 74, matchNameValue
		}
		return 73, matchName
	case "access-control-allow-headers":
		switch value {
		case "cache-control":
			return 33, matchNameValue
		case "content-type":
			return 34, matchNameValue
		case "*":
			return 75, matchNameValue
		}
		return 33, matchName
	case "access-control-allow-methods":
		switch value {
		case "get":
			return 76, matchNameValue
		case "get, post, options":
			return 77, matchNameValue
		case "options":
			return 78, matchNameValue
		}
		return 76, matchName
	case "access-control-allow-origin":
		if value == "*" {
			return 35, matchNameValue
		}
		return 35, matchName
	case "access-control-expose-headers":
		if value == "content-length" {
			return 79, matchNameValue
		}
		return 79, matchName
	case "access-control-request-headers":
		if value == "content-type" {
			return 80, matchNameValue
		}
		return 80, matchName
	case "access-control-request-method":
		switch value {
		case "get":
			return 81, matchNameValue
		case "post":
			return 82, matchNameValue
		}
		return 81, matchName
	case "age":
		if value == "0" {
			return 2, matchNameValue
		}
		return 2, matchName
	case "alt-svc":
		if value == "clear" {
			return 83, matchNameValue
		}
		return 83, matchName
	case "authorization":
		if value == "" {
			return 84, matchNameValue
		}
		return 84, matchName
	case "cache-control":
		switch value {
		case "max-age=0":
			return 36, matchNameValue
		case "max-age=2592000":
			return 37, matchNameValue
		case "max-age=604800":
			return 38, matchNameValue
		case "no-cache":
			return 39, matchNameValue
		case "no-store":
			return 40, matchNameValue
		case "public, max-age=31536000":
			return 41, matchNameValue
		}
		return 36, matchName
	case "content-disposition":
		if value == "" {
			return 3, matchNameValue
		}
		return 3, matchName
	case "content-encoding":
		switch value {
		case "br":
			return 42, matchNameValue
		case "gzip":
			return 43, matchNameValue
		}
		return 42, matchName
	case "content-length":
		if value == "0" {
			return 4, matchNameValue
		}
		return 4, matchName
	case "content-security-policy":
		if value == "script-src 'none'; object-src 'none'; base-uri 'none'" {
			return 85, matchNameValue
		}
		return 85, matchName
	case "content-type":
		switch value {
		case "text/css":
			return 51, matchNameValue
		case "text/plain;charset=utf-8":
			return 54, matchNameValue
		case "application/dns-message":
			return 44, matchNameValue
		case "application/x-www-form-urlencoded":
			return 47, matchNameValue
		case "image/gif":
			return 48, matchNameValue
		case "image/png":
			return 50, matchNameValue
		case "text/html; charset=utf-8":
			return 52, matchNameValue
		case "text/plain":
			return 53, matchNameValue
		case "application/javascript":
			return 45, matchNameValue
		case "application/json":
			return 46, matchNameValue
		case "image/jpeg":
			return 49, matchNameValue
		}
		return 51, matchName
	case "cookie":
		if value == "" {
			return 5, matchNameValue
		}
		return 5, matchName
	case "date":
		if value == "" {
			return 6, matchNameValue
		}
		return 6, matchName
	case "early-data":
		if value == "1" {
			return 86, matchNameValue
		}
		return 86, matchName
	case "etag":
		if value == "" {
			return 7, matchNameValue
		}
		return 7, matchName
	case "expect-ct":
		if value == "" {
			return 87, matchNameValue
		}
		return 87, matchName
	case "forwarded":
		if value == "" {
			return 88, matchNameValue
		}
		return 88, matchName
	case "if-modified-since":
		if value == "" {
			return 8, matchNameValue
		}
		return 8, matchName
	case "if-none-match":
		if value == "" {
			return 9, matchNameValue
		}
		return 9, matchName
	case "if-range":
		if value == "" {
			return 89, matchNameValue
		}
		return 89, matchName
	case "last-modified":
		if value == "" {
			return 10, matchNameValue
		}
		return 10, matchName
	case "link":
		if value == "" {
			return 11, matchNameValue
		}
		return 11, matchName
	case "location":
		if value == "" {
			return 12, matchNameValue
		}
		return 12, matchName
	case "origin":
		if value == "" {
			return 90, matchNameValue
		}
		return 90, matchName
	case "purpose":
		if value == "prefetch" {
			return 91, matchNameValue
		}
		return 91, matchName
	case "range":
		if value == "bytes=0-" {
			return 55, matchNameValue
		}
		return 55, matchName
	case "referer":
		if value == "" {
			return 13, matchNameValue
		}
		return 13, matchName
	case "server":
		if value == "" {
			return 92, matchNameValue
		}
		return 92, matchName
	case "set-cookie":
		if value == "" {
			return 14, matchNameValue
		}
		return 14, matchName
	case "strict-transport-security":
		switch value {
		case "max-age=31536000":
			return 56, matchNameValue
		case "max-age=31536000; includesubdomains":
			return 57, matchNameValue
		case "max-age=31536000; includesubdomains; preload":
			return 58, matchNameValue
		}
		return 56, matchName
	case "timing-allow-origin":
		if value == "*" {
			return 93, matchNameValue
		}
		return 93, matchName
	case "upgrade-insecure-requests":
		if value == "1" {
			return 94, matchNameValue
		}
		return 94, matchName
	case "user-agent":
		if value == "" {
			return 95, matchNameValue
		}
		return 95, matchName
	case "vary":
		switch value {
		case "accept-encoding":
			return 59, matchNameValue
		case "origin":
			return 60, matchNameValue
		}
		return 59, matchName
	case "x-content-type-options":
		if value == "nosniff" {
			return 61, matchNameValue
		}
		return 61, matchName
	case "x-forwarded-for":
		if value == "" {
			return 96, matchNameValue
		}
		return 96, matchName
	case "x-frame-options":
		switch value {
		case "deny":
			return 97, matchNameValue
		case "sameorigin":
			return 98, matchNameValue
		}
		return 97, matchName
	case "x-xss-protection":
		if value == "1; mode=block" {
			return 62, matchNameValue
		}
		return 62, matchName
	}
	return 0, matchNone
}
