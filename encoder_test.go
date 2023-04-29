package quack

import (
	"testing"
)

var header = map[string][]string{
	"Content-Type":     {"text/html; charset=utf-8"},
	"Content-Length":   {"1024"},
	"Content-Encoding": {"br"},
	"Vary":             {"Content-Encoding"},
}

func TestEncoder(t *testing.T) {
	var buf [1024]byte

	e := NewEncoder()

	_, _ = e.NewRequest(buf[:0], "GET", "https", "localhost", "/", header)
}

func BenchmarkEncoder(b *testing.B) {
	var buf [1024]byte

	b.ReportAllocs()
	b.ResetTimer()

	e := NewEncoder()
	for i := 0; i < b.N; i++ {
		_, _ = e.NewRequest(buf[:0], "GET", "https", "localhost", "/", header)
	}
}
