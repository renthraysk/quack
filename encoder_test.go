package quack

import (
	"encoding/base64"
	"math"
	"net/http"
	"strconv"
	"testing"
	"time"
)

var header = http.Header{
	"Content-Type":     []string{"text/html; charset=utf-8"},
	"Content-Length":   []string{"1024"},
	"Content-Encoding": []string{"br"},
	"Vary":             []string{"Content-Encoding"},
}

func canName(s string) string {
	b := make([]byte, 0, 32)
	nextA := 'a'
	for i, c := range []byte(s) {
		if c-byte(nextA) < 26 {
			if len(b) == 0 {
				b = append(b, s...)
			}
			b[i] = c ^ 0x20
		}
		nextA = 'A'
		if c == '-' {
			nextA = 'a'
		}
	}
	if len(b) > 0 {
		return string(b)
	}
	return s
}

func FuzzEncodeDecode(f *testing.F) {
	f.Fuzz(func(t *testing.T, nameB, valueB []byte) {
		if len(nameB) == 0 {
			return
		}
		name := base64.RawURLEncoding.EncodeToString(nameB)

		const m = 1<<'\x00' | 1<<'\n' | 1<<'\r'

		for i, c := range valueB {
			if c < 64 && (1<<c)&m != 0 {
				valueB[i] = ' '
			}
		}
		value := string(valueB)

		e := NewEncoder()
		d := NewDecoder()
		hf := make([]headerField, 0, 1)
		encoded := e.appendHeaderField([]byte{0, 0}, name, value, false)
		err := d.Decode(encoded, func(name, value string) {
			hf = append(hf, headerField{name, value})
		})
		if err != nil {
			t.Fatalf("decode error: %v", err)
		}
		if len(hf) != 1 {
			t.Fatalf("expected 1 headerField, got %d", len(hf))
		}
		if n := canName(name); hf[0].Name != n {
			t.Errorf("expected name %q, got %q", n, hf[0].Name)
		}
		if hf[0].Value != value {
			t.Errorf("expected value %q, got %q", value, hf[0].Value)
		}
	})
}

func TestIntAppend(t *testing.T) {
	for _, v := range []int64{0, 1000, math.MaxInt64} {

		e := NewEncoder()
		r := e.NewRequest(nil, "GET", "https", "localhost", "/")
		r = e.AppendContentLength(r, v, true)

		d := NewDecoder()
		got := ""
		err := d.Decode(r, func(name, value string) {
			if name == "Content-Length" {
				got = value
			}
		})
		if err != nil {
			t.Errorf("decode failed: %v", err)
		}
		if expected := strconv.FormatInt(v, 10); got != expected {
			t.Errorf("expected: %v got %v", expected, got)
		}
	}
}

func TestTimeAppend(t *testing.T) {
	for _, expected := range []string{
		"Thu, 02 Feb 2006 15:04:05 GMT",
	} {
		tt, err := time.Parse(time.RFC1123, expected)
		if err != nil {
			t.Fatalf("time.Parse failed: %v", err)
		}
		e := NewEncoder()
		r := e.NewRequest(nil, "GET", "https", "localhost", "/")
		r = e.AppendTimeHeaderField(r, "Date", tt, true)
		r = e.AppendTimeHeaderField(r, "Last-Modified", tt, true)

		d := NewDecoder()
		gotD := ""
		gotLM := ""
		err = d.Decode(r, func(name, value string) {
			if name == "Date" {
				gotD = value
			}
			if name == "Last-Modified" {
				gotLM = value
			}
		})
		if err != nil {
			t.Errorf("decode failed: %v", err)
		}
		if gotD != expected {
			t.Errorf("Date expected: %v got %v", expected, gotD)
		}
		if gotLM != expected {
			t.Errorf("Last-Modified expected: %v got %v", expected, gotLM)
		}
	}
}

func BenchmarkTimeAppend(b *testing.B) {
	now := time.Now()
	var buf [32]byte

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = appendTime(buf[:0], now)
	}
}

func BenchmarkTimeStdlibFormat(b *testing.B) {
	now := time.Now()
	var buf [32]byte

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = appendStringLiteral(buf[:0], now.UTC().Format(time.RFC1123))
	}
}

func BenchmarkAppendInt(b *testing.B) {
	var buf [32]byte

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = appendInt(buf[:0], int64(i))
	}
}

func BenchmarkStdlibInt(b *testing.B) {
	var buf [32]byte

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = appendStringLiteral(buf[:0], strconv.FormatInt(int64(i), 10))
	}
}

func BenchmarkEncoder(b *testing.B) {
	var buf [1024]byte

	b.ReportAllocs()
	b.ResetTimer()

	e := NewEncoder()
	for i := 0; i < b.N; i++ {
		_ = e.NewRequest(buf[:0], "GET", "https", "localhost", "/")
	}
}
