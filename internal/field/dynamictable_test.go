package field

import (
	"encoding/hex"
	"testing"
)

// https://datatracker.ietf.org/doc/html/rfc9204#name-dynamic-table-insert-evicti
func TestRFC9204_B1(t *testing.T) {

	in := dehex(t, "0000510b2f696e6465782e68746d6c")
	got := make([]header, 0, 2)

	d := &Decoder{}
	err := d.Decode(in, func(name, value string) {
		got = append(got, header{name, value})
	})
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if exp := []header{{":path", "/index.html"}}; !Equal(got, exp) {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func TestRFC9204_DecodingEncoderInstructions(t *testing.T) {

	dt := DT{maxCapacity: 1 << 10}

	// https://datatracker.ietf.org/doc/html/rfc9204#name-dynamic-table-2
	{
		in := dehex(t, "3fbd01c00f7777772e6578616d706c652e636f6dc10c2f73616d706c652f70617468")
		err := dt.DecodeEncoderInstructions(in)
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}
		if dt.capacity != 220 {
			t.Errorf("expected capacity 220, got %d", dt.capacity)
		}
		if dt.size != 106 {
			t.Errorf("expected size 106, got %d", dt.size)
		}
		exp := []header{
			{":authority", "www.example.com"},
			{":path", "/sample/path"}}
		if !Equal(dt.headers, exp) {
			t.Errorf("expected %v, got %v", exp, dt.headers)
		}
	}
	// https://datatracker.ietf.org/doc/html/rfc9204#name-speculative-insert
	{
		in := dehex(t, "4a637573746f6d2d6b65790c637573746f6d2d76616c7565")
		err := dt.DecodeEncoderInstructions(in)
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}
		if dt.size != 160 {
			t.Errorf("expected size 160, got %d", dt.size)
		}
		exp := []header{
			{":authority", "www.example.com"},
			{":path", "/sample/path"},
			{"Custom-Key", "custom-value"},
		}
		if !Equal(dt.headers, exp) {
			t.Errorf("expected %v, got %v", exp, dt.headers)
		}
	}
	// https://datatracker.ietf.org/doc/html/rfc9204#name-duplicate-instruction-strea
	{
		in := dehex(t, "02")
		err := dt.DecodeEncoderInstructions(in)
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}
		if dt.size != 217 {
			t.Errorf("expected size 217, got %d", dt.size)
		}
		exp := []header{
			{":authority", "www.example.com"},
			{":path", "/sample/path"},
			{"Custom-Key", "custom-value"},
			{":authority", "www.example.com"},
		}
		if !Equal(dt.headers, exp) {
			t.Errorf("expected %v, got %v", exp, dt.headers)
		}
	}
	// https://datatracker.ietf.org/doc/html/rfc9204#appendix-B.5
	{
		in := dehex(t, "810d637573746f6d2d76616c756532")
		err := dt.DecodeEncoderInstructions(in)
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}
		if dt.size != 215 {
			t.Errorf("expected size 215, got %d", dt.size)
		}
		exp := []header{
			{":path", "/sample/path"},
			{"Custom-Key", "custom-value"},
			{":authority", "www.example.com"},
			{"Custom-Key", "custom-value2"},
		}
		if !Equal(dt.headers, exp) {
			t.Errorf("expected %v, got %v", exp, dt.headers)
		}
	}
}

func TestSnapshot(t *testing.T) {

	in := DT{}
	in.mu.Lock()
	in.setCapacityLocked(1 << 10)
	in.insertLocked("Server", "proto")
	in.insertLocked("Server", "proto")
	in.insertLocked("Server", "proto2")
	in.insertLocked("A", "B")
	in.mu.Unlock()

	p := in.appendSnapshot(nil)

	out := DT{}
	out.DecodeEncoderInstructions(p)

	if out.size != in.size {
		t.Errorf("expected size %d, got %d", in.size, out.size)
	}
	if out.capacity != in.capacity {
		t.Errorf("expected capacity %d, got %d", in.capacity, out.capacity)
	}
	if !Equal(in.headers, out.headers) {
		t.Errorf("expected headers %v, got %v", in.headers, out.headers)
	}
}

func Equal[T comparable](x, y []T) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func dehex(tb testing.TB, s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		tb.Fatalf("failed to decode hex")
	}
	return b
}
