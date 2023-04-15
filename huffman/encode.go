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

// AppendInt
func AppendInt(p []byte, i int64) []byte {
	if i == 0 {
		return appendFinal(p, uint64(codes['0']), uint(codeLengths['0']))
	}
	var x uint64
	var n uint

	u := uint64(i)
	if i < 0 {
		u = uint64(-i)
		x = uint64(codes['-'])
		n = uint(codeLengths['-'])
	}
	a := make([]uint8, 0, 16)
	for u >= 100 {
		w := u / 100
		a = append(a, uint8(u-w*100))
		u = w
	}
	// u < 100
	switch {
	case u >= 10:
		p, x, n = append00To99(p, x, n, int(u))
	case u > 0:
		p, x, n = appendByte(p, x, n, byte(u+'0'))
	}
	for i := len(a) - 1; i >= 0; i-- {
		p, x, n = append00To99(p, x, n, int(a[i]))
	}
	return appendFinal(p, x, n)
}

// AppendRFC1123Time appends the huffman encoding of time.Time t in RFC1123
// format to p returning the result.
func AppendRFC1123Time(p []byte, t time.Time) []byte {
	u := t.UTC()
	year, month, day := u.Date()
	d := days[u.Weekday()]
	x, n := uint64(d.code), uint(d.length) // "Day, "
	p, x, n = append00To99(p, x, n, day)   // "Day, 01"

	m := months[month-1]
	x <<= m.length % 64
	n += uint(m.length)
	x |= uint64(m.code) // "Day, 01 Mon "

	y := year / 100
	if y >= 100 {
		panic("year overflows 4 digits")
	}
	p, x, n = append00To99(p, x, n, y)
	p, x, n = append00To99(p, x, n, year-(y*100))
	p, x, n = appendByte(p, x, n, ' ')

	hour, minute, second := u.Clock()
	p, x, n = append00To99(p, x, n, hour)
	p, x, n = appendByte(p, x, n, ':')
	p, x, n = append00To99(p, x, n, minute)
	p, x, n = appendByte(p, x, n, ':')
	p, x, n = append00To99(p, x, n, second)
	// " GMT"
	x <<= 27
	x |= 0x0298b46f // "Day, 01 Mon 1990 HH:MM:SS GMT"
	n += 27
	if n >= 32 {
		n %= 32
		y := uint32(x >> n)
		p = append(p, byte(y>>24), byte(y>>16), byte(y>>8), byte(y))
	}
	return appendFinal(p, x, n)
}

var days = [...]struct {
	code   uint64
	length uint8
}{
	{0x1badabe94, 33}, // "Sun, "
	{0xd07abe94, 32},  // "Mon, "
	{0xdf697e94, 32},  // "Tue, "
	{0xe4593e94, 32},  // "Wed, "
	{0x1be7b7e94, 33}, // "Thu, "
	{0xc361be94, 32},  // "Fri, "
	{0x6e1a7e94, 31},  // "Sat, "
}

var months = [...]struct {
	code   uint32
	length uint8
}{
	{0x14ca3a94, 30}, // " Jan "
	{0x14c258d4, 30}, // " Feb "
	{0x14d03b14, 30}, // " Mar "
	{0x1486bb14, 30}, // " Apr "
	{0x29a07e94, 31}, // " May "
	{0x2996da94, 31}, // " Jun "
	{0x2996da14, 31}, // " Jul "
	{0x1486d994, 30}, // " Aug "
	{0x14dc5ad4, 30}, // " Sep "
	{0x0a6a2254, 29}, // " Oct "
	{0x29a4fdd4, 31}, // " Nov "
	{0x0a5f2914, 29}, // " Dec "
}

// appendByte appends the huffman code for byte c into the tuple (p, x, n)
// Ensures returned x has last than 32 valid bits, and n is less than 32.
// Assumes x has less than 32 valid bits, and n to be less than 32.
func appendByte(p []byte, x uint64, n uint, c byte) ([]byte, uint64, uint) {
	// inlines
	s := uint(codeLengths[c])
	x <<= s % 64
	x |= uint64(codes[c])
	n += s
	if n >= 32 {
		n %= 32
		y := uint32(x >> n)
		p = append(p, byte(y>>24), byte(y>>16), byte(y>>8), byte(y))
	}
	return p, x, n
}

// append00To99 appends the huffman codes for the two digit ASCII
// representation of i into the tuple (p, x, n)
// Ensures returned x has last than 32 valid bits, and n is less than 32.
// Assumes x has less than 32 valid bits, and n to be less than 32.
func append00To99(p []byte, x uint64, n uint, i int) ([]byte, uint64, uint) {
	// inlines
	s := uint(codes00To99[i].length)
	x <<= s % 64
	x |= uint64(codes00To99[i].code)
	n += s
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
		y := uint16(x)
		return append(p, byte(y>>8), byte(y))
	case 3:
		y := uint16(x >> 8)
		return append(p, byte(y>>8), byte(y), byte(x))
	}
	y := uint32(x)
	return append(p, byte(y>>24), byte(y>>16), byte(y>>8), byte(y))
}
