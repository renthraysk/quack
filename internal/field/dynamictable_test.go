package field

import (
	"testing"
)

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
	out.ParseEncoderInstructions(p)

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
