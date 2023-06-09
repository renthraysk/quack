package field

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

func TestFieldEncoderStaticTable(t *testing.T) {
	for name, values := range http3Table {
		t.Run(name, func(t *testing.T) {
			if values == nil {
				values = []string{""}
			}
			e := &Encoder{}
			buf := make([]byte, 2, 1<<10)

			for _, value := range values {
				var p []byte

				if isPseudo(name) {
					switch name {
					case ":status":
						status, err := strconv.Atoi(value)
						if err != nil {
							t.Fatalf("error parsing csv status code: %v", err)
						}
						p = appendStatus(buf, status)
					case ":method":
						p = appendMethod(buf, value)
					case ":path":
						p = appendPath(buf, value)
					case ":authority":
						p = appendAuthority(buf, value)
					case ":scheme":
						p = appendScheme(buf, value)
					default:
						t.Fatalf("unknown pseudo header %s", name)
					}
				} else {
					p = e.appendFieldLine(buf, name, value)
				}
				d := &Decoder{}
				err := d.Decode(p, func(k, v string) {
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
		})
	}
}
