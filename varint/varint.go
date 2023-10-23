package varint

import "errors"

const (
	maxVarint62    = (1 << 62) - 1
	maxVarint62Len = (62 + 6) / 7
)

var errUnexpectedEnd = errors.New("quack: unexpected end")
var errVarintOverflow = errors.New("quack: varint overflow")

func Read(p []byte, mask uint8) (uint64, []byte, error) {
	if len(p) <= 0 {
		return 0, p, errUnexpectedEnd
	}
	x, q := uint64(p[0]&mask), p[1:]
	if x < uint64(mask) {
		return x, q, nil
	}
	var s uint

	// MaxVarint62Len (9) bytes (9*7, 63 bits) cannot overflow a uint64, even
	// with x possibly being 0xFF from the byte above.
	for i, b := range q[:min(len(q), maxVarint62Len)] {
		x += uint64(b&0x7F) << s
		s += 7
		if b < 0x80 {
			if x > maxVarint62 {
				return 0, p, errVarintOverflow
			}
			return x, p[1+i+1:], nil
		}
	}
	// either no bytes or all continuation bits were set
	if len(q) < maxVarint62Len {
		// Looks like a truncated varint
		return 0, p, errUnexpectedEnd
	}
	return 0, p, errVarintOverflow
}

func Append(p []byte, prefix, mask byte, x uint64) []byte {
	if x < uint64(mask) {
		return append(p, prefix|byte(x))
	}
	x -= uint64(mask)
	prefix |= mask
	if x <= 0x7F {
		return append(p, prefix, byte(x))
	}
	p = append(p, prefix, byte(x|0x80))
	x >>= 7
	for x > 0x7F {
		p = append(p, byte(x)|0x80)
		x >>= 7
	}
	return append(p, byte(x))
}
