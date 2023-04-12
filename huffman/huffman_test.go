package huffman

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"testing"
	"time"
)

func dehex(tb testing.TB, s string) []byte {
	r, err := hex.DecodeString(s)
	if err != nil {
		tb.Fatalf("dehex: %v", err)
	}
	return r
}

func TestQuick(t *testing.T) {
	encoded := dehex(t, "1c6490b2cd39ba75a29a8f5f6b109b7bf8f3ebdf")
	got, err := Decode(make([]byte, 0, 32), encoded)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(got) != "abcdefghijklmnopqrstuvwxyz" {
		t.Errorf("expected %q, got %q", "abcdefghijklmnopqrstuvwxyz", string(got))
	}
}

func TestEOS(t *testing.T) {
	t.Run("eos", func(t *testing.T) {
		in := binary.BigEndian.AppendUint32(nil, eosCode<<(32-eosLen))
		_, err := Decode(nil, in)
		if err != errEOSEncoded {
			t.Errorf("expected error %v, got %v", errEOSEncoded, err)
		}
	})
	t.Run("offset-eos", func(t *testing.T) {
		in := []byte{0xf8} // 8 bit code for &
		in = binary.BigEndian.AppendUint32(in, eosCode<<(32-eosLen))
		_, err := Decode(nil, in)
		if err != errEOSEncoded {
			t.Errorf("expected error %v, got %v", errEOSEncoded, err)
		}
	})
}

func TestDecode(t *testing.T) {
	tests := []struct {
		in       string
		expected string
	}{
		{"", ""},
		{"1c6490b2cd39ba75a29a8f5f6b109b7bf8f3ebdf", "abcdefghijklmnopqrstuvwxyz"},
		{"1001", "200"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			encoded := dehex(t, tt.in)
			got, err := Decode(nil, encoded)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, string(got))
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	buf := make([]byte, 1024)
	in := dehex(b, "1c6490b2cd39ba75a29a8f5f6b109b7bf8f3ebdf")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Decode(buf[:0], in)
	}
}
func TestEncodeInt(t *testing.T) {

	tests := []struct {
		value   int64
		length  uint64
		encoded []byte
	}{
		{math.MinInt64, 15, []byte{0x59, 0xf1, 0x09, 0x96, 0x5d, 0x10, 0x19, 0x71, 0xe6, 0xda, 0x75, 0xd6, 0xde, 0x03, 0xdf}},
		{-9e6, 6, []byte{0x59, 0xf0, 0x00, 0x00, 0x00, 0x3f}},
		{0, 1, []byte{0x07}},
		{9e6, 5, []byte{0x7c, 0x00, 0x00, 0x00, 0x0f}},
		{math.MaxInt64, 14, []byte{0x7c, 0x42, 0x65, 0x97, 0x44, 0x06, 0x5c, 0x79, 0xb6, 0x9d, 0x75, 0xb7, 0x80, 0xef}},
	}

	for _, q := range tests {
		t.Run(fmt.Sprintf("%d", q.value), func(t *testing.T) {
			if got := AppendInt(nil, q.value); !bytes.Equal(got, q.encoded) {
				t.Errorf("AppendInt: expected %x, got %x", q.encoded, got)
			}
		})
	}
}

func TestEncodeString(t *testing.T) {
	tests := []struct {
		s       string
		length  uint64
		encoded []byte
	}{
		{"", 0, []byte{}},
		{"200", 2, []byte{0x10, 0x01}},
		{"text/html; charset=utf-8", 18, []byte{0x49, 0x7c, 0xa5, 0x89, 0xd3, 0x4d, 0x1f, 0x6a, 0x12, 0x71, 0xd8, 0x82, 0xa6, 0x0b, 0x53, 0x2a, 0xcf, 0x7f}},
	}

	for _, s := range tests {
		t.Run(s.s, func(t *testing.T) {
			if n := EncodeLength(s.s); n != s.length {
				t.Errorf("expected length %d, got %d", s.length, n)
			}
			if got := AppendString(nil, s.s); !bytes.Equal(got, s.encoded) {
				t.Errorf("expected encoding %x, got %x", s.encoded, got)
			}
		})
	}
}

func TestEncodeStringLower(t *testing.T) {

	lower := []byte{
		0x1c, 0x64, 0x90, 0xb2, 0xcd, 0x39, 0xba, 0x75, 0xa2, 0x9a,
		0x8f, 0x5f, 0x6b, 0x10, 0x9b, 0x7b, 0xf8, 0xf3, 0xeb, 0xdf}

	tests := []struct {
		s       string
		length  uint64
		encoded []byte
	}{
		{"", 0, []byte{}},
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZ", 20, lower},
		{"abcdefghijklmnopqrstuvwxyz", 20, lower},
		{"Content-Type", 9, []byte{0x21, 0xea, 0x49, 0x6a, 0x4a, 0xc9, 0xf5, 0x59, 0x7f}},
	}

	for _, s := range tests {
		t.Run(s.s, func(t *testing.T) {
			if n := EncodeLengthLower(s.s); n != s.length {
				t.Errorf("expected %d, got %d", s.length, n)
			}
			if got := AppendStringLower(nil, s.s); !bytes.Equal(got, s.encoded) {
				t.Errorf("expected %x, got %x", s.encoded, got)
			}
		})
	}
}

func TestEncodeRFC3339Time(t *testing.T) {

	times := []struct {
		time    string
		length  uint64
		encoded string
	}{
		{"2006-01-02T15:04:05Z", 15, "1000e2c00ac016f0b7700d5c037fbf"},
		{"2006-01-02T15:04:05+07:00", 19, "1000e2c00ac016f0b7700d5c037fec0edc003f"},
		{"2006-01-02T15:04:05-07:00", 18, "1000e2c00ac016f0b7700d5c036b01db8007"},
	}
	for _, c := range times {
		t.Run(c.time, func(t *testing.T) {
			tm, err := time.Parse(time.RFC3339, c.time)
			if err != nil {
				t.Fatalf("failed to parse %q", c.time)
			}
			expected := dehex(t, c.encoded)
			if got := AppendRFC3339Time(nil, tm); !bytes.Equal(got, expected) {
				t.Errorf("expected %x, got %x", c.encoded, got)
			}
		})
	}
}

func FuzzEncodeDecode(f *testing.F) {
	f.Add("")
	f.Add("z")
	f.Add("\x00")
	f.Fuzz(func(t *testing.T, in string) {
		encoded := AppendString(nil, in)
		decoded, err := Decode(nil, encoded)
		if err != nil {
			t.Errorf("decode error: %v", err)
		}
		if string(decoded) != in {
			t.Errorf("decode failed, expected %q, got %q", in, decoded)
		}
	})
}

func FuzzInt(f *testing.F) {
	f.Add(int64(0))
	f.Add(int64(math.MaxInt64))
	f.Add(int64(math.MinInt64))
	f.Fuzz(func(t *testing.T, in int64) {
		encoded := AppendInt(nil, in)
		decoded, err := Decode(nil, encoded)
		if err != nil {
			t.Errorf("decode error: %v", err)
		}
		got, err := strconv.ParseInt(string(decoded), 10, 64)
		if err != nil {
			t.Errorf("ParseInt error: %v", err)
		}
		if in != got {
			t.Errorf("expected %d, got %d", in, got)
		}
	})
}

func FuzzTime(f *testing.F) {
	f.Add(int64(0))
	f.Add(int64(1000000000))
	f.Fuzz(func(t *testing.T, sec int64) {
		expected := time.Unix(sec, 0) // no nano in RFC3339

		// avoid any panics if 5+ digit year somehow.
		if expected.Year() > 9999 {
			return
		}
		encoded := AppendRFC3339Time(nil, expected)
		decoded, err := Decode(nil, encoded)
		if err != nil {
			t.Errorf("decode error: %v", err)
		}
		got, err := time.Parse(time.RFC3339, string(decoded))
		if err != nil {
			t.Errorf("parse error: %v", err)
		}
		if !expected.Equal(got) {
			t.Errorf("expected %q, got %q",
				expected.Format(time.RFC3339),
				got.Format(time.RFC3339))
		}
	})
}
