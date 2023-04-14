package quack

import (
	"strconv"
	"strings"
	"testing"
)

var http3Table = map[string][]string{
	":authority":                       nil,
	":method":                          {"CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "POST", "PUT"},
	":path":                            {"/"},
	":scheme":                          {"http", "https"},
	":status":                          {"103", "200", "304", "404", "503", "100", "204", "206", "302", "400", "403", "421", "425", "500"},
	"Accept":                           {"*/*", "application/dns-message"},
	"Accept-Encoding":                  {"gzip, deflate, br"},
	"Accept-Language":                  nil,
	"Accept-Ranges":                    {"bytes"},
	"Access-Control-Allow-Credentials": {"FALSE", "TRUE"},
	"Access-Control-Allow-Headers":     {"cache-control", "content-type", "*"},
	"Access-Control-Allow-Methods":     {"get", "get, post, options", "options"},
	"Access-Control-Allow-Origin":      {"*"},
	"Access-Control-Expose-Headers":    {"content-length"},
	"Access-Control-Request-Headers":   {"content-type"},
	"Access-Control-Request-Method":    {"get", "post"},
	"Age":                              {"0"},
	"Alt-Svc":                          {"clear"},
	"Authorization":                    nil,
	"Cache-Control":                    {"max-age=0", "max-age=2592000", "max-age=604800", "no-cache", "no-store", "public, max-age=31536000"},
	"Content-Disposition":              nil,
	"Content-Encoding":                 {"br", "gzip"},
	"Content-Length":                   {"0"},
	"Content-Security-Policy":          {"script-src 'none'; object-src 'none'; base-uri 'none'"},
	"Content-Type":                     {"application/dns-message", "application/javascript", "application/json", "application/x-www-form-urlencoded", "image/gif", "image/jpeg", "image/png", "text/css", "text/html; charset=utf-8", "text/plain", "text/plain;charset=utf-8"},
	"Cookie":                           nil,
	"Date":                             nil,
	"Early-Data":                       {"1"},
	"Etag":                             nil,
	"Expect-Ct":                        nil,
	"Forwarded":                        nil,
	"If-Modified-Since":                nil,
	"If-None-Match":                    nil,
	"If-Range":                         nil,
	"Last-Modified":                    nil,
	"Link":                             nil,
	"Location":                         nil,
	"Origin":                           nil,
	"Purpose":                          {"prefetch"},
	"Range":                            {"bytes=0-"},
	"Referer":                          nil,
	"Server":                           nil,
	"Set-Cookie":                       nil,
	"Strict-Transport-Security":        {"max-age=31536000", "max-age=31536000; includesubdomains", "max-age=31536000; includesubdomains; preload"},
	"Timing-Allow-Origin":              {"*"},
	"Upgrade-Insecure-Requests":        {"1"},
	"User-Agent":                       nil,
	"Vary":                             {"accept-encoding", "origin"},
	"X-Content-Type-Options":           {"nosniff"},
	"X-Forwarded-For":                  nil,
	"X-Frame-Options":                  {"deny", "sameorigin"},
	"X-Xss-Protection":                 {"1; mode=block"},
}

func isPseudo(name string) bool { return strings.HasPrefix(name, ":") }

func testHeaderField(t *testing.T, name, value string, neverIndex bool) {
	var buf [1 << 10]byte
	var e Encoder

	b := buf[:2]

	t.Helper()
	if isPseudo(name) {
		switch name {
		case ":status":
			status, err := strconv.Atoi(value)
			if err != nil {
				t.Fatalf("error parsing csv status code: %v", err)
			}
			b = e.appendStatus(b, status)
		case ":method":
			b = e.appendMethod(b, value)
		case ":path":
			b = e.appendPath(b, value)
		case ":authority":
			b = e.appendAuthority(b, value)
		case ":scheme":
			b = e.appendScheme(b, value)
		default:
			t.Fatalf("unknown pseudo header %s", name)
		}
	} else {
		b = e.AppendCanonicalHeaderField(b, name, value, neverIndex)
	}

	d := &decoder{}
	err := d.decode(b, func(k, v string) {
		if name != k {
			t.Errorf("expected name %q, got %q", name, k)
		}
		if value != v {
			t.Errorf("expected value %q, got %q", value, v)
		}
	})
	if err != nil {
		t.Errorf("decode error: %v", err)
	}
}

func testStaticTable(t *testing.T, neverIndex bool) {
	for key, values := range http3Table {
		t.Run(key, func(t *testing.T) {
			if values == nil {
				testHeaderField(t, key, "", neverIndex)
				return
			}
			for _, value := range values {
				testHeaderField(t, key, value, neverIndex)
			}
		})
	}
}

func TestStaticTable(t *testing.T) {
	testStaticTable(t, false)
	testStaticTable(t, true)
}
