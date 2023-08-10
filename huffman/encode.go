package huffman

import (
	"time"

	"github.com/renthraysk/quack/ascii"
)

// EncodeLength returns the number of bytes would it would take to huffman
// encode s.
func EncodeLength(s string) uint64 {
	var n uint64
	for i := 0; i < len(s); i++ {
		n += uint64(codeLengths[s[i]])
	}
	return (n + 7) / 8
}

// AppendString appends huffman encoded s to p, returning the result.
func AppendString(p []byte, s string) []byte {
	var (
		x uint64
		n uint
	)
	for i := 0; i < len(s); i++ {
		p, x, n = appendByte(p, x, n, s[i])
	}
	return appendFinal(p, x, n)
}

// EncodeLengthLower returns the number of bytes would it would take to encode
// the ASCII lower cased version of s.
func EncodeLengthLower(s string) uint64 {
	var n uint64
	for i := 0; i < len(s); i++ {
		n += uint64(codeLengths[ascii.ToLower(s[i])])
	}
	return (n + 7) / 8
}

// AppendStringLower appends the huffman encoded ASCII lower case of s to p.
func AppendStringLower(p []byte, s string) []byte {
	var (
		x uint64
		n uint
	)
	for i := 0; i < len(s); i++ {
		p, x, n = appendByte(p, x, n, ascii.ToLower(s[i]))
	}
	return appendFinal(p, x, n)
}

// AppendInt appends the huffman encoded ASCII representation of v to p,
// return the result.
func AppendInt(p []byte, v int64) []byte {
	if v == 0 {
		return appendFinal(p, uint64(codes['0']), uint(codeLengths['0']))
	}
	var x uint64
	var n uint
	var a [16]uint8

	u := uint64(v)
	if v < 0 {
		u = uint64(-v)
		x = uint64(codes['-'])
		n = uint(codeLengths['-'])
	}
	i := len(a)
	for u >= 100 {
		w := u / 100
		i--
		a[i] = uint8(u - w*100)
		u = w
	}
	// u < 100
	switch {
	case u >= 10:
		p, x, n = append00To99(p, x, n, int(u))
	case u > 0:
		p, x, n = appendByte(p, x, n, byte(u+'0'))
	}
	for ; i < len(a); i++ {
		p, x, n = append00To99(p, x, n, int(a[i]))
	}
	return appendFinal(p, x, n)
}

// AppendHttpTime appends the huffman encoding of time.Time t in RFC1123
// format to p returning the result.
func AppendHttpTime(p []byte, t time.Time) []byte {
	u := t.UTC()
	year, month, day := u.Date()
	y := year / 100
	if y >= 100 {
		panic("year overflows 4 digits")
	}

	d := days[u.Weekday()]
	x, n := uint64(d.code), uint(d.length) // "Day,"
	p, x, n = appendByte(p, x, n, ' ')     // "Day, "
	p, x, n = append00To99(p, x, n, day)   // "Day, 01"

	m := months[month-1]
	p, x, n = appendCode(p, x, n, m.code, m.length) // "Day, 01 Mon "

	p, x, n = append00To99(p, x, n, y)
	p, x, n = append00To99(p, x, n, year-(y*100)) // "Day, 01 Mon 1990"
	p, x, n = appendByte(p, x, n, ' ')            // "Day, 01 Mon 1990 "

	hour, minute, second := u.Clock()
	p, x, n = append00To99(p, x, n, hour) // "Day, 01 Mon 1990 HH"
	p, x, n = appendByte(p, x, n, ':')
	p, x, n = append00To99(p, x, n, minute) // "Day, 01 Mon 1990 HH:MM"
	p, x, n = appendByte(p, x, n, ':')
	p, x, n = append00To99(p, x, n, second)           // "Day, 01 Mon 1990 HH:MM:SS"
	p, x, n = appendCode(p, x, n, gmtCode, gmtLength) // "Day, 01 Mon 1990 HH:MM:SS GMT"
	return appendFinal(p, x, n)
}

// appendByte appends the huffman code for byte c into the tuple (p, x, n)
// Ensures returned x has last than 32 valid bits, and n is less than 32.
// Assumes x has less than 32 valid bits, and n to be less than 32.
func appendByte(p []byte, x uint64, n uint, c byte) ([]byte, uint64, uint) {
	// inlines
	return appendCode(p, x, n, codes[c], codeLengths[c])
}

// append00To99 appends the huffman codes for the two digit ASCII
// representation of i into the tuple (p, x, n)
// Assumes x has less than 32 valid bits, and n to be less than 32.
// Ensures returned x has less than 32 valid bits, and n is less than 32.
func append00To99(p []byte, x uint64, n uint, i int) ([]byte, uint64, uint) {
	// inlines
	return appendCode(p, x, n, uint32(codes00To99[i].code), codes00To99[i].length)
}

// appendCode appends the huffman code of given length to (p, x, n)
// Assumes x has less than 32 valid bits, and n to be less than 32
// Ensures returned x has less than 32 valid bits, and n is less than 32.
func appendCode(p []byte, x uint64, n uint, code uint32, length uint8) ([]byte, uint64, uint) {
	// inlines
	x <<= length % 64
	x |= uint64(code)
	n += uint(length)
	if n >= 32 {
		n %= 32
		y := uint32(x >> n)
		p = append(p, byte(y>>24), byte(y>>16), byte(y>>8), byte(y))
	}
	return p, x, n
}

// appendFinal appends remaining n least significant bits of x to p, adding
// padding as necessary, completing a huffman encoding.
// Assumes x has less than 33 valid bits, and n to be less than 33.
func appendFinal(p []byte, x uint64, n uint) []byte {
	const eosPadByte = eosCode >> (eosLen - 8)

	if over := n % 8; over > 0 {
		pad := 8 - over
		x = (x << pad) | (eosPadByte >> over)
		n += pad
	}
	switch n / 8 {
	case 0:
		return p
	case 1:
		return append(p, byte(x))
	case 2:
		return append(p, byte(x>>8), byte(x))
	case 3:
		return append(p, byte(x>>16), byte(x>>8), byte(x))
	}
	return append(p, byte(x>>24), byte(x>>16), byte(x>>8), byte(x))
}
