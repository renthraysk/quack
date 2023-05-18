package varint

import "errors"

const (
	maxVarint62    = (1 << 62) - 1
	maxVarint62Len = (62 + 6) / 7
)

var errUnexpectedEnd = errors.New("unexpected end")

// https://datatracker.ietf.org/doc/html/rfc9204#section-4.1.1

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
	if len(q) > maxVarint62Len {
		q = q[:maxVarint62Len]
	}
	for i, b := range q {
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

func Append(p []byte, x uint64, mask, prefix byte) []byte {
	if x < uint64(mask) {
		return append(p, prefix|byte(x))
	}
	p = append(p, prefix|mask)
	x -= uint64(mask)
	for x >= 0x80 {
		p = append(p, 0x80|byte(x))
		x >>= 7
	}
	return append(p, byte(x))
}
