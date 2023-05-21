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
	_, err := e.AppendRequest(buf[:0], "GET", "https", "localhost", "/", header)
	if err != nil {
		t.Errorf("NewRequest failed: %v", err)
	}
}

func BenchmarkEncoder(b *testing.B) {
	var buf [1024]byte

	b.ReportAllocs()
	b.ResetTimer()

	e := NewEncoder()
	for i := 0; i < b.N; i++ {
		_, _ = e.AppendRequest(buf[:0], "GET", "https", "localhost", "/", header)
	}
}
