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
	"accept":                           {"*/*", "application/dns-message"},
	"accept-encoding":                  {"gzip, deflate, br"},
	"accept-language":                  nil,
	"accept-ranges":                    {"bytes"},
	"access-control-allow-credentials": {"FALSE", "TRUE"},
	"access-control-allow-headers":     {"cache-control", "content-type", "*"},
	"access-control-allow-methods":     {"get", "get, post, options", "options"},
	"access-control-allow-origin":      {"*"},
	"access-control-expose-headers":    {"content-length"},
	"access-control-request-headers":   {"content-type"},
	"access-control-request-method":    {"get", "post"},
	"age":                              {"0"},
	"alt-svc":                          {"clear"},
	"authorization":                    nil,
	"cache-control":                    {"max-age=0", "max-age=2592000", "max-age=604800", "no-cache", "no-store", "public, max-age=31536000"},
	"content-disposition":              nil,
	"content-encoding":                 {"br", "gzip"},
	"content-length":                   {"0"},
	"content-security-policy":          {"script-src 'none'; object-src 'none'; base-uri 'none'"},
	"content-type":                     {"application/dns-message", "application/javascript", "application/json", "application/x-www-form-urlencoded", "image/gif", "image/jpeg", "image/png", "text/css", "text/html; charset=utf-8", "text/plain", "text/plain;charset=utf-8"},
	"cookie":                           nil,
	"date":                             nil,
	"early-data":                       {"1"},
	"etag":                             nil,
	"expect-ct":                        nil,
	"forwarded":                        nil,
	"if-modified-since":                nil,
	"if-none-match":                    nil,
	"if-range":                         nil,
	"last-modified":                    nil,
	"link":                             nil,
	"location":                         nil,
	"origin":                           nil,
	"purpose":                          {"prefetch"},
	"range":                            {"bytes=0-"},
	"referer":                          nil,
	"server":                           nil,
	"set-cookie":                       nil,
	"strict-transport-security":        {"max-age=31536000", "max-age=31536000; includesubdomains", "max-age=31536000; includesubdomains; preload"},
	"timing-allow-origin":              {"*"},
	"upgrade-insecure-requests":        {"1"},
	"user-agent":                       nil,
	"vary":                             {"accept-encoding", "origin"},
	"x-content-type-options":           {"nosniff"},
	"x-forwarded-for":                  nil,
	"x-frame-options":                  {"deny", "sameorigin"},
	"x-xss-protection":                 {"1; mode=block"},
}

func isPseudo(name string) bool { return strings.HasPrefix(name, ":") }

func testHeaderField(t *testing.T, name, value string) {
	var buf [1 << 10]byte

	e := new(Encoder)

	b := buf[:2]

	t.Helper()
	if isPseudo(name) {
		switch name {
		case ":status":
			status, err := strconv.Atoi(value)
			if err != nil {
				t.Fatalf("error parsing csv status code: %v", err)
			}
			b = appendStatus(b, status)
		case ":method":
			b = appendMethod(b, value)
		case ":path":
			b = appendPath(b, value)
		case ":authority":
			b = appendAuthority(b, value)
		case ":scheme":
			b = appendScheme(b, value)
		default:
			t.Fatalf("unknown pseudo header %s", name)
		}
	} else {
		b = e.encodeHeaderField(b, name, value)
	}
	d := new(Decoder)
	err := d.Decode(b, func(k, v string) {
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

func testStaticTable(t *testing.T) {
	for key, values := range http3Table {
		t.Run(key, func(t *testing.T) {
			if values == nil {
				testHeaderField(t, key, "")
				return
			}
			for _, value := range values {
				testHeaderField(t, key, value)
			}
		})
	}
}

func TestStaticTable(t *testing.T) {
	testStaticTable(t)
}
