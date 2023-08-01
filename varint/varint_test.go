package varint

import (
	"errors"
	"testing"
)

func TestReadVarint(t *testing.T) {

	bad := [16]byte{}
	memset(bad[:], 0xFF)

	tests := []struct {
		name        string
		in          []byte
		mask        uint8
		expected    uint64
		expectedErr error
	}{
		{"rfc7541_c1_1", []byte{0b01010}, 0x1F, 10, nil},
		{"rfc7541_c1_2", []byte{0b11111, 0b1001_1010, 0b000_1010}, 0x1F, 1337, nil},
		{"rfc7541_c1_3", []byte{0b00101010}, 0xFF, 42, nil},

		{"empty", []byte{}, 0x01, 0, errUnexpectedEnd},
		{"zero-7", []byte{0x00}, 0x7F, 0, nil},
		{"one-1", []byte{0x01, 0x00}, 0x01, 1, nil},

		{"maxint62-1", Append(nil, 0, 1, maxVarint62), 1, maxVarint62, nil},
		{"maxint62-7", Append(nil, 0, 0x7F, maxVarint62), 0x7F, maxVarint62, nil},

		{"overflow-1", Append(nil, 0, 1, maxVarint62+1), 1, 0, errVarintOverflow},
		{"overflow-7", Append(nil, 0, 0x7F, maxVarint62+1), 0x7F, 0, errVarintOverflow},
		{"overflow-8", Append(nil, 0, 0xFF, maxVarint62+1), 0xFF, 0, errVarintOverflow},

		{"short", bad[:9], 0x7F, 0, errUnexpectedEnd},
		{"overflow", bad[:10], 0x7F, 0, errVarintOverflow},
		{"long", bad[:], 0x7F, 0, errVarintOverflow},

		{"eos", []byte{0x7F, 0x80}, 0x7F, 0, errUnexpectedEnd},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, r, err := Read(tt.in, tt.mask)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected error expected %v, got %v", tt.expectedErr, err)
			}
			if x != tt.expected {
				t.Errorf("expected x to be %d, got %d", tt.expected, x)
			}
			if err != nil {
				if len(tt.in) > 0 && &tt.in[0] != &r[0] {
					t.Error("expected remain to be unchanged")
				}
			} else if len(r) != 0 {
				t.Errorf("expected remain to be empty")
			}
		})
	}
}

func memset(p []byte, x byte) []byte {
	for i := range p {
		p[i] = x
	}
	return p
}
