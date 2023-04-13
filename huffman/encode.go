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

// AppendRFC3339Time appends the huffman encoding of time.Time t in RFC3339
// format to p returning the result.
func AppendRFC3339Time(p []byte, t time.Time) []byte {
	year, month, day := t.Date()
	y := year / 100
	if y >= 100 {
		panic("year overflows 4 digits")
	}
	x, n := uint64(codes00To99[y].code), uint(codes00To99[y].length)
	p, x, n = append00To99(p, x, n, year-(y*100))
	p, x, n = appendByte(p, x, n, '-')
	p, x, n = append00To99(p, x, n, int(month))
	p, x, n = appendByte(p, x, n, '-')
	p, x, n = append00To99(p, x, n, day)
	p, x, n = appendByte(p, x, n, 'T')
	hour, minute, second := t.Clock()
	p, x, n = append00To99(p, x, n, hour)
	p, x, n = appendByte(p, x, n, ':')
	p, x, n = append00To99(p, x, n, minute)
	p, x, n = appendByte(p, x, n, ':')
	p, x, n = append00To99(p, x, n, second)
	_, offsetSec := t.Zone()
	if offsetSec == 0 {
		p, x, n = appendByte(p, x, n, 'Z')
		return appendFinal(p, x, n)
	}
	s := byte('+')
	offsetMin := offsetSec / 60
	if offsetMin < 0 {
		s = '-'
		offsetMin = -offsetMin
	}
	p, x, n = appendByte(p, x, n, s)
	offsetHour := offsetMin / 60
	p, x, n = append00To99(p, x, n, offsetHour)
	p, x, n = appendByte(p, x, n, ':')
	p, x, n = append00To99(p, x, n, offsetMin-(60*offsetHour))
	return appendFinal(p, x, n)
}

func AppendRFC1123Time(p []byte, t time.Time) []byte {
	const days = "SunMonTueWedThuFriSat"
	const months = "JanFebMarAprMayJunJulAugSepOctNovDec"

	u := t.UTC()
	year, month, day := u.Date()
	dayName := days[u.Weekday()*3:]
	monthName := months[3*(month-1):]
	_ = dayName[2]
	x, n := uint64(codes[dayName[0]]), uint(codeLengths[dayName[0]])
	p, x, n = appendByte(p, x, n, dayName[1])
	p, x, n = appendByte(p, x, n, dayName[2])
	p, x, n = appendByte(p, x, n, ',')
	p, x, n = appendByte(p, x, n, ' ')
	p, x, n = append00To99(p, x, n, day)
	p, x, n = appendByte(p, x, n, ' ')
	_ = monthName[2]
	p, x, n = appendByte(p, x, n, monthName[0])
	p, x, n = appendByte(p, x, n, monthName[1])
	p, x, n = appendByte(p, x, n, monthName[2])
	p, x, n = appendByte(p, x, n, ' ')
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
	p, x, n = appendByte(p, x, n, ' ')
	p, x, n = appendByte(p, x, n, 'G')
	p, x, n = appendByte(p, x, n, 'M')
	p, x, n = appendByte(p, x, n, 'T')
	return appendFinal(p, x, n)
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
