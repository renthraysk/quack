package ascii

import "testing"

func TestIsValueValid(t *testing.T) {
	if IsValueValid([]byte{'\v'}) {
		t.Error("\\v expected false, got true")
	}
	if !IsValueValid([]byte{'!'}) {
		t.Error("! expected true, got false")
	}
	if !IsValueValid([]byte{'~'}) {
		t.Error("~ expected true, got false")
	}
	if IsValueValid([]byte{'\x7F'}) {
		t.Error("\\x7F expected false, got true")
	}
	if !IsValueValid([]byte{'\x80'}) {
		t.Error("\\x80 expected true, got false")
	}
	if !IsValueValid([]byte{'\xFF'}) {
		t.Error("\\xFF expected true, got false")
	}
}

func TestAppendLower(t *testing.T) {
	tt := []struct {
		p        string
		s        string
		expected string
	}{
		{"", "ABCDEFGHJKLMNOPQRSTUVWXYZ", "abcdefghjklmnopqrstuvwxyz"},
		{"ABCDEFGHJKLMN", "OPQRSTUVWXYZ", "ABCDEFGHJKLMNopqrstuvwxyz"},
	}
	for _, c := range tt {
		if got := AppendLower([]byte(c.p), c.s); string(got) != c.expected {
			t.Errorf("expected %q, got %q", c.expected, got)
		}
	}
}

func TestToCanonical(t *testing.T) {
	tt := []struct {
		b        string
		expected string
	}{
		{"content-type", "Content-Type"},
		{"cONTENT-tYPE", "Content-Type"},
	}
	for _, c := range tt {
		b := []byte(c.b)
		ToCanonical(b)
		if got := string(b); string(got) != c.expected {
			t.Errorf("expected %q, got %q", c.expected, got)
		}
	}
}
