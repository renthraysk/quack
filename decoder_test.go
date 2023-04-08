package quack

import (
	"encoding/hex"
	"testing"
)

func headersEqual(a, b []headerField) bool {
	if len(a) != len(b) {
		return false
	}
	for i, h := range a {
		if h.Name != b[i].Name {
			return false
		}
		if h.Value != b[i].Value {
			return false
		}
	}
	return true
}

const requestBin = "" +
	"0000d1d75086a0e41d139d09c15f0ec0497ca589d34d1f43aeba0c41a4c7" +
	"a98f33a69a3fdf9a68fa1d75d0620d263d4c79a68fbed00177fe8d48e62b" +
	"03ee697e8d48e62b1e0b1d7f5f2c7cfdf6800bbddf5f398b2d4b62bbf45a" +
	"befb4005db"

var requestQuack = []headerField{
	{":method", "GET"},
	{":scheme", "https"},
	{":authority", "localhost"},
	{":path", "/"},
	{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"},
	{"Accept-Encoding", "gzip, deflate, br"},
	{"Accept-Language", "en-GB,en;q=0.5"}}

func newRequest(e *Encoder, method, scheme, authority, path string) ([]headerField, []byte) {
	expected := []headerField{
		{":method", method},
		{":scheme", scheme},
		{":authority", authority},
		{":path", path},
	}
	return expected, e.NewRequest(nil, method, scheme, authority, path)
}

func dehex(tb testing.TB, s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		tb.Fatalf("dehex failed: %v", err)
	}
	return b
}
func TestDecodeRequest(t *testing.T) {
	d := NewDecoder()

	got := make([]headerField, 0, 8)
	f := func(name, value string) {
		got = append(got, headerField{name, value})
	}
	err := d.Decode(dehex(t, requestBin), f)
	if err != nil {
		t.Errorf("decode error: %v", err)
	}
	if !headersEqual(requestQuack, got) {
		t.Errorf("expected %v, got %v", requestQuack, got)
	}
}

func min(x, y uint) uint {
	if x <= y {
		return x
	}
	return y
}

func FuzzDecodeRequestFragmented(f *testing.F) {

	req := dehex(f, requestBin)

	f.Add(uint(0), uint(0))
	f.Fuzz(func(t *testing.T, a, b uint) {

		// ensure all of req is sent through
		d := NewDecoder()
		got := make([]headerField, 0, 8)
		f := func(name, value string) {
			got = append(got, headerField{name, value})
		}

		r := req
		for _, n := range []uint{a, b, uint(len(req))} {
			nn := min(n, uint(len(r)))
			if err := d.Decode(r[:nn], f); err != nil {
				t.Errorf("decode error: %v", err)
			}
			r = r[nn:]
			if len(r) == 0 {
				break
			}
		}
		if !headersEqual(requestQuack, got) {
			t.Errorf("expected %v, got %v", requestQuack, got)
		}
	})
}

func BenchmarkDecoder(b *testing.B) {
	headers := make([]headerField, 0, 16)
	f := func(name, value string) {
		headers = append(headers, headerField{name, value})
	}

	in := dehex(b, requestBin)
	d := NewDecoder()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Reset()
		headers = headers[:0]
		err := d.Decode(in, f)
		_ = err
	}
}
