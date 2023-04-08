package quack

import (
	"time"
)

// This is largely code generated, with a handful type specific methods added
// to enable conversion directly to QPACK encoding without need for
// an intermediary string (allocation).

// AppendCanonicalHeaderField will append the QPACK encoding of (name, value)
// pair to p. name is expected to be in canonical (Content-Type) form, and if
// name is not encoded as an index, will be converted to all lower case.
// Does not handle pseudo headers.
func (e *Encoder) AppendCanonicalHeaderField(p []byte, name, value string, neverIndex bool) []byte {
	const NeverIndex = 0b0010_0000 // 01N1_XXXX: Literal Field Line with Name Reference in static table

	var prefix byte

	if neverIndex {
		prefix = NeverIndex
	}
	switch name {
	case "Accept-Encoding":
		return e.appendAcceptEncoding(p, value, prefix)
	case "Accept-Language":
		return e.appendAcceptLanguage(p, value, prefix)
	case "Accept-Ranges":
		return e.appendAcceptRanges(p, value, prefix)
	case "Accept":
		return e.appendAccept(p, value, prefix)
	case "Access-Control-Allow-Credentials":
		return e.appendAccessControlAllowCredentials(p, value, prefix)
	case "Access-Control-Allow-Headers":
		return e.appendAccessControlAllowHeaders(p, value, prefix)
	case "Access-Control-Allow-Methods":
		return e.appendAccessControlAllowMethods(p, value, prefix)
	case "Access-Control-Allow-Origin":
		return e.appendAccessControlAllowOrigin(p, value, prefix)
	case "Access-Control-Expose-Headers":
		return e.appendAccessControlExposeHeaders(p, value, prefix)
	case "Access-Control-Request-Headers":
		return e.appendAccessControlRequestHeaders(p, value, prefix)
	case "Access-Control-Request-Method":
		return e.appendAccessControlRequestMethod(p, value, prefix)
	case "Age":
		return e.appendAge(p, value, prefix)
	case "Alt-Svc":
		return e.appendAltSvc(p, value, prefix)
	case "Authorization":
		return e.appendAuthorization(p, value, prefix)
	case "Cache-Control":
		return e.appendCacheControl(p, value, prefix)
	case "Content-Disposition":
		return e.appendContentDisposition(p, value, prefix)
	case "Content-Encoding":
		return e.appendContentEncoding(p, value, prefix)
	case "Content-Length":
		return e.appendContentLengthString(p, value, prefix)
	case "Content-Security-Policy":
		return e.appendContentSecurityPolicy(p, value, prefix)
	case "Content-Type":
		return e.AppendContentType(p, value, prefix)
	case "Cookie":
		return e.appendCookie(p, value, prefix)
	case "Date":
		return e.appendDateString(p, value, prefix)
	case "Early-Data":
		return e.appendEarlyData(p, value, prefix)
	case "Etag":
		return e.appendEtag(p, value, prefix)
	case "Expect-Ct":
		return e.appendExpectCt(p, value, prefix)
	case "Forwarded":
		return e.appendForwarded(p, value, prefix)
	case "If-Modified-Since":
		return e.appendIfModifiedSinceString(p, value, prefix)
	case "If-None-Match":
		return e.appendIfNoneMatch(p, value, prefix)
	case "If-Range":
		return e.appendIfRange(p, value, prefix)
	case "Last-Modified":
		return e.appendLastModifiedString(p, value, prefix)
	case "Link":
		return e.appendLink(p, value, prefix)
	case "Location":
		return e.appendLocation(p, value, prefix)
	case "Origin":
		return e.appendOrigin(p, value, prefix)
	case "Purpose":
		return e.appendPurpose(p, value, prefix)
	case "Range":
		return e.appendRange(p, value, prefix)
	case "Referer":
		return e.appendReferer(p, value, prefix)
	case "Server":
		return e.appendServer(p, value, prefix)
	case "Set-Cookie":
		return e.appendSetCookie(p, value, prefix)
	case "Strict-Transport-Security":
		return e.appendStrictTransportSecurity(p, value, prefix)
	case "Timing-Allow-Origin":
		return e.appendTimingAllowOrigin(p, value, prefix)
	case "Upgrade-Insecure-Requests":
		return e.appendUpgradeInsecureRequests(p, value, prefix)
	case "User-Agent":
		return e.appendUserAgent(p, value, prefix)
	case "Vary":
		return e.appendVary(p, value, prefix)
	case "X-Content-Type-Options":
		return e.appendXContentTypeOptions(p, value, prefix)
	case "X-Forwarded-For":
		return e.appendXForwardedFor(p, value, prefix)
	case "X-Frame-Options":
		return e.appendXFrameOptions(p, value, prefix)
	case "X-Xss-Protection":
		return e.appendXXssProtection(p, value, prefix)
	}
	return e.appendHeaderField(p, name, value, neverIndex)
}

// The pseudo headers

// appendAuthority appends an :authority pseudo header field to p
func (e *Encoder) appendAuthority(p []byte, authority string) []byte {
	if authority == "" {
		return append(p, 0xC0)
	}
	p = append(p, 0x50)
	return appendStringLiteral(p, authority)
}

// appendPath appends a :path pseudo header field to p
func (e *Encoder) appendPath(p []byte, path string) []byte {
	if path == "/" {
		return append(p, 0xC1)
	}
	p = append(p, 0x51)
	return appendStringLiteral(p, path)
}

// appendStatus appends a :status pseudo header field to p
func (e *Encoder) appendStatus(p []byte, statusCode int) []byte {
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

// appendMethod appends a :method pseudo header field to p
func (e *Encoder) appendMethod(p []byte, method string) []byte {
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
		return appendStringLiteral(p, method)
	}
	return append(p, b)
}

// appendScheme appends a :scheme pseudo header field to p
func (e *Encoder) appendScheme(p []byte, scheme string) []byte {
	var b uint8
	switch scheme {
	case "http":
		b = 0xD6
	case "https":
		b = 0xD7
	default:
		p = append(p, 0x5F, 0x07)
		return appendStringLiteral(p, scheme)
	}
	return append(p, b)
}

// Regular headers

func (e *Encoder) appendAge(p []byte, age string, neverIndex byte) []byte {
	if age == "0" {
		return append(p, 0xC2)
	}
	p = append(p, neverIndex|0x52)
	return appendStringLiteral(p, age)
}

func (e *Encoder) appendContentDisposition(p []byte, value string, neverIndex byte) []byte {
	if value == "" {
		return append(p, 0xC3)
	}
	p = append(p, neverIndex|0x53)
	return appendStringLiteral(p, value)
}

func (e *Encoder) appendContentLengthString(p []byte, contentLength string, neverIndex byte) []byte {
	if contentLength == "0" {
		return append(p, 0xC4)
	}
	p = append(p, neverIndex|0x54)
	return appendStringLiteral(p, contentLength)
}

func (e *Encoder) appendContentLength(p []byte, contentLength int64, neverIndex byte) []byte {
	if contentLength == 0 {
		return append(p, 0xC4)
	}
	p = append(p, neverIndex|0x54)
	return appendInt(p, contentLength)
}

func (e *Encoder) appendCookie(p []byte, value string, neverIndex byte) []byte {
	if value == "" {
		return append(p, 0xC5)
	}
	p = append(p, neverIndex|0x55)
	return appendStringLiteral(p, value)
}

func (e *Encoder) appendDateString(p []byte, s string, neverIndex byte) []byte {
	if s == "" {
		return append(p, 0xC6)
	}
	p = append(p, neverIndex|0x56)
	return appendStringLiteral(p, s)
}

func (e *Encoder) AppendDate(p []byte, t time.Time, neverIndex byte) []byte {
	p = append(p, neverIndex|0x56)
	return appendTime(p, t)
}

func (e *Encoder) appendEtag(p []byte, eTag string, neverIndex byte) []byte {
	if eTag == "" {
		return append(p, 0xC7)
	}
	p = append(p, neverIndex|0x57)
	return appendStringLiteral(p, eTag)
}

func (e *Encoder) appendIfModifiedSinceString(p []byte, t string, neverIndex byte) []byte {
	if t == "" {
		return append(p, 0xC8)
	}
	p = append(p, neverIndex|0x58)
	return appendStringLiteral(p, t)
}

func (e *Encoder) appendIfModifiedSince(p []byte, t time.Time, neverIndex byte) []byte {
	p = append(p, neverIndex|0x58)
	return appendTime(p, t)
}

func (e *Encoder) appendIfNoneMatch(p []byte, eTag string, neverIndex byte) []byte {
	if eTag == "" {
		return append(p, 0xC9)
	}
	p = append(p, neverIndex|0x59)
	return appendStringLiteral(p, eTag)
}

func (e *Encoder) appendLastModifiedString(p []byte, t string, neverIndex byte) []byte {
	if t == "" {
		return append(p, 0xCA)
	}
	p = append(p, neverIndex|0x5A)
	return appendStringLiteral(p, t)
}

func (e *Encoder) appendLastModified(p []byte, t time.Time, neverIndex byte) []byte {
	p = append(p, neverIndex|0x5A)
	return appendTime(p, t)
}

func (e *Encoder) appendLink(p []byte, link string, neverIndex byte) []byte {
	if link == "" {
		return append(p, 0xCB)
	}
	p = append(p, neverIndex|0x5B)
	return appendStringLiteral(p, link)
}

func (e *Encoder) appendLocation(p []byte, location string, neverIndex byte) []byte {
	if location == "" {
		return append(p, 0xCC)
	}
	p = append(p, neverIndex|0x5C)
	return appendStringLiteral(p, location)
}

func (e *Encoder) appendReferer(p []byte, referer string, neverIndex byte) []byte {
	if referer == "" {
		return append(p, 0xCD)
	}
	p = append(p, neverIndex|0x5D)
	return appendStringLiteral(p, referer)
}

func (e *Encoder) appendSetCookie(p []byte, cookie string, neverIndex byte) []byte {
	if cookie == "" {
		return append(p, 0xCE)
	}
	p = append(p, neverIndex|0x5E)
	return appendStringLiteral(p, cookie)
}

func (e *Encoder) appendAccept(p []byte, accept string, neverIndex byte) []byte {
	var b byte
	switch accept {
	case "*/*":
		b = 0xDD
	case "application/dns-message":
		b = 0xDE
	default:
		p = append(p, neverIndex|0x5F, 0x0E)
		return appendStringLiteral(p, accept)
	}
	return append(p, b)
}

func (e *Encoder) appendAcceptEncoding(p []byte, encoding string, neverIndex byte) []byte {
	if encoding == "gzip, deflate, br" {
		return append(p, 0xDF)
	}
	p = append(p, neverIndex|0x5F, 0x10)
	return appendStringLiteral(p, encoding)
}

func (e *Encoder) appendAcceptRanges(p []byte, ranges string, neverIndex byte) []byte {
	if ranges == "bytes" {
		return append(p, 0xE0)
	}
	p = append(p, neverIndex|0x5F, 0x11)
	return appendStringLiteral(p, ranges)
}

func (e *Encoder) appendAccessControlAllowHeaders(p []byte, headers string, neverIndex byte) []byte {
	var b byte
	switch headers {
	case "*":
		return append(p, 0xFF, 0x0C)
	case "cache-control":
		b = 0xE1
	case "content-type":
		b = 0xE2
	default:
		p = append(p, neverIndex|0x5F, 0x12)
		return appendStringLiteral(p, headers)
	}
	return append(p, b)
}

func (e *Encoder) appendAccessControlAllowOrigin(p []byte, origin string, neverIndex byte) []byte {
	if origin == "*" {
		return append(p, 0xE3)
	}
	p = append(p, neverIndex|0x5F, 0x14)
	return appendStringLiteral(p, origin)
}

func (e *Encoder) appendCacheControl(p []byte, cacheControl string, neverIndex byte) []byte {
	var b byte
	switch cacheControl {
	case "max-age=0":
		b = 0xE4
	case "max-age=2592000":
		b = 0xE5
	case "max-age=604800":
		b = 0xE6
	case "no-cache":
		b = 0xE7
	case "no-store":
		b = 0xE8
	case "public, max-age=31536000":
		b = 0xE9
	default:
		p = append(p, neverIndex|0x5F, 0x15)
		return appendStringLiteral(p, cacheControl)
	}
	return append(p, b)
}

func (e *Encoder) appendContentEncoding(p []byte, contentEncoding string, neverIndex byte) []byte {
	var b byte
	switch contentEncoding {
	case "br":
		b = 0xEA
	case "gzip":
		b = 0xEB
	default:
		p = append(p, neverIndex|0x5F, 0x1B)
		return appendStringLiteral(p, contentEncoding)
	}
	return append(p, b)
}

func (e *Encoder) AppendContentType(p []byte, contentType string, neverIndex byte) []byte {
	var b byte
	switch contentType {
	case "image/jpeg":
		b = 0xF1
	case "image/png":
		b = 0xF2
	case "text/css":
		b = 0xF3
	case "application/dns-message":
		b = 0xEC
	case "application/javascript":
		b = 0xED
	case "application/json":
		b = 0xEE
	case "application/x-www-form-urlencoded":
		b = 0xEF
	case "image/gif":
		b = 0xF0
	case "text/html; charset=utf-8":
		b = 0xF4
	case "text/plain;charset=utf-8":
		b = 0xF6
	case "text/plain":
		b = 0xF5
	default:
		p = append(p, neverIndex|0x5F, 0x1D)
		return appendStringLiteral(p, contentType)
	}
	return append(p, b)
}

func (e *Encoder) appendRange(p []byte, rng string, neverIndex byte) []byte {
	if rng == "bytes=0-" {
		return append(p, 0xF7)
	}
	p = append(p, neverIndex|0x5F, 0x28)
	return appendStringLiteral(p, rng)
}

func (e *Encoder) appendStrictTransportSecurity(p []byte, tps string, neverIndex byte) []byte {
	var b byte

	switch tps {
	case "max-age=31536000":
		b = 0xF8
	case "max-age=31536000; includesubdomains":
		b = 0xF9
	case "max-age=31536000; includesubdomains; preload":
		b = 0xFA
	default:
		p = append(p, neverIndex|0x5F, 0x29)
		return appendStringLiteral(p, tps)
	}
	return append(p, b)
}

func (e *Encoder) appendVary(p []byte, vary string, neverIndex byte) []byte {
	var b byte

	switch vary {
	case "accept-encoding":
		b = 0xFB
	case "origin":
		b = 0xFC
	default:
		p = append(p, neverIndex|0x5F, 0x2C)
		return appendStringLiteral(p, vary)
	}
	return append(p, b)
}

func (e *Encoder) appendXContentTypeOptions(p []byte, value string, neverIndex byte) []byte {
	if value == "nosniff" {
		return append(p, 0xFD)
	}
	p = append(p, neverIndex|0x5F, 0x2E)
	return appendStringLiteral(p, value)
}

func (e *Encoder) appendXXssProtection(p []byte, value string, neverIndex byte) []byte {
	if value == "1; mode=block" {
		return append(p, 0xFE)
	}
	p = append(p, neverIndex|0x5F, 0x2F)
	return appendStringLiteral(p, value)
}

func (e *Encoder) appendAcceptLanguage(p []byte, lang string, neverIndex byte) []byte {
	if lang == "" {
		return append(p, 0xFF, 0x09)
	}
	p = append(p, neverIndex|0x5F, 0x39)
	return appendStringLiteral(p, lang)
}

func (e *Encoder) appendAccessControlAllowCredentials(p []byte, allow string, neverIndex byte) []byte {
	var b uint8
	switch allow {
	case "TRUE":
		b = 0x0B
	case "FALSE":
		b = 0x0A
	default:
		p = append(p, neverIndex|0x5F, 0x3A)
		return appendStringLiteral(p, allow)
	}
	return append(p, 0xFF, b)
}

func (e *Encoder) appendAccessControlAllowMethods(p []byte, methods string, neverIndex byte) []byte {
	var b uint8
	switch methods {
	case "get":
		b = 0x0D
	case "get, post, options":
		b = 0x0E
	case "options":
		b = 0x0F
	default:
		p = append(p, neverIndex|0x5F, 0x3D)
		return appendStringLiteral(p, methods)
	}
	return append(p, 0xFF, b)
}

func (e *Encoder) appendAccessControlExposeHeaders(p []byte, headers string, neverIndex byte) []byte {
	if headers == "content-length" {
		return append(p, 0xFF, 0x10)
	}
	p = append(p, neverIndex|0x5F, 0x40)
	return appendStringLiteral(p, headers)
}

func (e *Encoder) appendAccessControlRequestHeaders(p []byte, headers string, neverIndex byte) []byte {
	if headers == "content-type" {
		return append(p, 0xFF, 0x11)
	}
	p = append(p, neverIndex|0x5F, 0x41)
	return appendStringLiteral(p, headers)
}

func (e *Encoder) appendAccessControlRequestMethod(p []byte, method string, neverIndex byte) []byte {
	var b byte
	switch method {
	case "get":
		b = 0x12
	case "post":
		b = 0x13
	default:
		p = append(p, neverIndex|0x5F, 0x42)
		return appendStringLiteral(p, method)
	}
	return append(p, 0xFF, b)
}

func (e *Encoder) appendAltSvc(p []byte, svc string, neverIndex byte) []byte {
	if svc == "clear" {
		return append(p, 0xFF, 0x14)
	}
	p = append(p, neverIndex|0x5F, 0x44)
	return appendStringLiteral(p, svc)
}

func (e *Encoder) appendAuthorization(p []byte, authorization string, neverIndex byte) []byte {
	if authorization == "" {
		return append(p, 0xFF, 0x15)
	}
	p = append(p, neverIndex|0x5F, 0x45)
	return appendStringLiteral(p, authorization)
}

func (e *Encoder) appendContentSecurityPolicy(p []byte, policy string, neverIndex byte) []byte {
	if policy == "script-src 'none'; object-src 'none'; base-uri 'none'" {
		return append(p, 0xFF, 0x16)
	}
	p = append(p, neverIndex|0x5F, 0x46)
	return appendStringLiteral(p, policy)
}

func (e *Encoder) appendEarlyData(p []byte, value string, neverIndex byte) []byte {
	if value == "1" {
		return append(p, 0xFF, 0x17)
	}
	p = append(p, neverIndex|0x5F, 0x47)
	return appendStringLiteral(p, value)
}

func (e *Encoder) appendExpectCt(p []byte, value string, neverIndex byte) []byte {
	if value == "" {
		return append(p, 0xFF, 0x18)
	}
	p = append(p, neverIndex|0x5F, 0x48)
	return appendStringLiteral(p, value)
}

func (e *Encoder) appendForwarded(p []byte, forward string, neverIndex byte) []byte {
	if forward == "" {
		return append(p, 0xFF, 0x19)
	}
	p = append(p, neverIndex|0x5F, 0x49)
	return appendStringLiteral(p, forward)
}

func (e *Encoder) appendIfRange(p []byte, rng string, neverIndex byte) []byte {
	if rng == "" {
		return append(p, 0xFF, 0x1A)
	}
	p = append(p, neverIndex|0x5F, 0x4A)
	return appendStringLiteral(p, rng)
}

func (e *Encoder) appendOrigin(p []byte, origin string, neverIndex byte) []byte {
	if origin == "" {
		return append(p, 0xFF, 0x1B)
	}
	p = append(p, neverIndex|0x5F, 0x4B)
	return appendStringLiteral(p, origin)
}

func (e *Encoder) appendPurpose(p []byte, purpose string, neverIndex byte) []byte {
	if purpose == "prefetch" {
		return append(p, 0xFF, 0x1C)
	}
	p = append(p, neverIndex|0x5F, 0x4C)
	return appendStringLiteral(p, purpose)
}

func (e *Encoder) appendServer(p []byte, server string, neverIndex byte) []byte {
	if server == "" {
		return append(p, 0xFF, 0x1D)
	}
	p = append(p, neverIndex|0x5F, 0x4D)
	return appendStringLiteral(p, server)
}

func (e *Encoder) appendTimingAllowOrigin(p []byte, value string, neverIndex byte) []byte {
	if value == "*" {
		return append(p, 0xFF, 0x1E)
	}
	p = append(p, neverIndex|0x5F, 0x4E)
	return appendStringLiteral(p, value)
}

func (e *Encoder) appendUpgradeInsecureRequests(p []byte, value string, neverIndex byte) []byte {
	if value == "1" {
		return append(p, 0xFF, 0x1F)
	}
	p = append(p, neverIndex|0x5F, 0x4F)
	return appendStringLiteral(p, value)
}

func (e *Encoder) appendUserAgent(p []byte, userAgent string, neverIndex byte) []byte {
	if userAgent == "" {
		return append(p, 0xFF, 0x20)
	}
	p = append(p, neverIndex|0x5F, 0x50)
	return appendStringLiteral(p, userAgent)
}

func (e *Encoder) appendXForwardedFor(p []byte, xForwardedFor string, neverIndex byte) []byte {
	if xForwardedFor == "" {
		return append(p, 0xFF, 0x21)
	}
	p = append(p, neverIndex|0x5F, 0x51)
	return appendStringLiteral(p, xForwardedFor)
}

func (e *Encoder) appendXFrameOptions(p []byte, xFrameOptions string, neverIndex byte) []byte {
	var b byte
	switch xFrameOptions {
	case "deny":
		b = 0x22
	case "sameorigin":
		b = 0x23
	default:
		p = append(p, neverIndex|0x5F, 0x52)
		return appendStringLiteral(p, xFrameOptions)
	}
	return append(p, 0xFF, b)
}
