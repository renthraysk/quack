//go:build ignore

package main

/*
	This generates various lookup tables for huffman encoding & decoding.

	Having a direct lookup for short codes of 13 bits or less, means all but 6
	visible ASCII characters and space (VCHAR / SP) are decoded with the fast
	path, with only needing a ~8KB lookup table.

	go run tables_codegen.go > tables.go
*/

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	// Huffman code length limits
	minCodeLength = 5
	maxCodeLength = 30

	eosCode = 0x3fffffff
	eosLen  = 30

	stringWidth = 70
)

type bound struct {
	delta  uint32
	length uint8
	offset uint8
}

type codeIndex struct {
	shortCodes []byte
	longCodes  [256]byte
	bounds     []bound

	//
	maxShortCode int
	delta        uint32
	offset30     uint8
}

type buckets [maxCodeLength + 1][]byte

// next returns the next bit length that has output symbols.
func (b *buckets) next(length int) (int, bool) {
	i := length + 1
	for i < len(b) && len(b[i]) == 0 {
		i++
	}
	return i, i < len(b)
}

func New(maxShortCode int) *codeIndex {
	ci := &codeIndex{
		shortCodes:   make([]byte, 1<<maxShortCode),
		maxShortCode: maxShortCode,
	}

	// bucket syms by their code length
	li := new(buckets)
	for sym, n := range codeLengths {
		li[n] = append(li[n], byte(sym))
	}

	{ // Short code table generation
		i := 0
		for length := minCodeLength; length <= maxShortCode; length++ {
			n := 1 << (maxShortCode - length)
			for _, b := range li[length] {
				memset(ci.shortCodes[i:i+n], b)
				i += n
			}
		}
		ci.shortCodes = ci.shortCodes[:i]
	}

	{ // Long code table generation
		i := 0
		for _, syms := range li {
			i += copy(ci.longCodes[i:], syms)
		}
	}

	{ // Bounds table generation
		offset := 0
		for length, syms := range li {
			if len(syms) == 0 {
				continue
			}

			nextLength, ok := li.next(length)
			if !ok {
				break
			}

			nextCode := codes[li[nextLength][0]] << (32 - nextLength)
			code := codes[syms[0]] << (32 - length)

			ci.bounds = append(ci.bounds, bound{delta: nextCode - code, length: byte(length), offset: byte(offset)})
			offset += len(syms)
		}
		ci.offset30 = byte(offset)
	}
	{
		nextLength, ok := li.next(maxShortCode)
		if !ok {
			panic("couldn't find initial delta")
		}
		ci.delta = codes[li[nextLength][0]] << (32 - nextLength)
	}
	return ci
}

func code(s string) (uint64, uint) {
	var x uint64
	var n uint

	for _, c := range s {
		x <<= codeLengths[c]
		n += uint(codeLengths[c])
		x |= uint64(codes[c])
	}
	return x, n
}

func main() {
	const maxShortCode = 13
	const days = "Sun,Mon,Tue,Wed,Thu,Fri,Sat,"
	const months = " Jan  Feb  Mar  Apr  May  Jun  Jul  Aug  Sep  Oct  Nov  Dec "

	w := os.Stdout

	ci := New(maxShortCode)

	fmt.Fprintln(w, "package huffman")
	fmt.Fprintln(w)

	{ // " GMT" preencoded
		x, n := code(" GMT")
		fmt.Fprintln(w, "const (")
		fmt.Fprintf(w, "\tgmtCode   = %#08x // \" GMT\"\n", x)
		fmt.Fprintf(w, "\tgmtLength = %d\n", n)
		fmt.Fprintln(w, ")")
	}

	fmt.Fprintln(w)
	fmt.Println(`type code[T uint16 | uint32 | uint64] struct {
	_      [0]func()
	code   T
	length uint8
}`)
	fmt.Fprintln(w)

	fmt.Fprintln(w, "var days = [...]code[uint32] {")
	for i := 0; i < len(days); i += 4 {
		x, n := code(days[i : i+4])
		fmt.Fprintf(w, "\t{code: %#08x, length: %d}, // %q\n", x, n, days[i:i+4])
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
	
	fmt.Fprintln(w, "var months = [...]code[uint32] {")
	for i := 0; i < len(months); i += 5 {
		x, n := code(months[i : i+5])
		fmt.Fprintf(w, "\t{code: %#08x, length: %d}, // %q\n", x, n, months[i:i+5])
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)

	fmt.Fprintf(w, `const (
	// decoding constants
	maxShortCode = %d
	delta = %#08x
	offset30 = %d
)`, maxShortCode, ci.delta, ci.offset30)

	fmt.Fprintln(w)
	fmt.Fprintln(w)

	fmt.Fprintln(w, `var bounds = []struct {
	_      [0]func()
	delta  uint32 //
	length uint8  // number of bits for the code
	offset uint8  // offset into longCodes for first sym that has this length
}{`)
	for _, b := range ci.bounds {
		s := "\t{delta: %#08x, length: %d, offset: %d},\n"
		if b.length <= maxShortCode {
			s = "\t// {delta: %#08x, length: %d, offset: %d},\n"
		}
		fmt.Fprintf(w, s, b.delta, b.length, b.offset)
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
	printCodes(w, &codes)
	fmt.Fprintln(w)
	fmt.Fprintln(w, `const codeLengths = "" +`)
	printGoString(w, codeLengths, stringWidth)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "const shortCodes = \"\" +")
	printGoString(w, ci.shortCodes, stringWidth)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "const longCodes = \"\" +")
	printGoString(w, ci.longCodes[:], stringWidth)

	fmt.Fprintln(w, `
var codes00To99 = [100]code[uint16] {`)
	for i := byte('0'); i <= '9'; i++ {
		for j := byte('0'); j <= '9'; j++ {
			fmt.Fprintf(w, "\t{length: %d, code: %#04x}, // %c%c\n",
				codeLengths[i]+codeLengths[j],
				(codes[i]<<codeLengths[j])|codes[j], i, j)
		}
	}
	fmt.Fprintln(w, "}")
}

func isIn(c byte, lo, hi uint64) bool {
	m := lo
	if c >= 64 {
		m = hi
	}
	return 1<<(c%64)&m != 0
}

func printGoString[T ~string | ~[]byte](w io.Writer, s T, width int) {
	const hex = "0123456789abcdef"
	const escChars = (1 << 0x20) - 1 | 1<<'"' | 1<<'\\' | 1<<'\x7F'
	const tabWidth = 4

	if len(s) == 0 {
		return
	}
	var b bytes.Buffer

	ww := width - tabWidth
	for i := 0; i < len(s); {
		b.Reset()
		b.WriteString("\t\"")
		for ; i < len(s) && b.Len() < ww; i++ {
			c := s[i]
			if c > '~' || isIn(c, escChars%(1<<64), escChars>>64) {
				b.WriteByte('\\')
				switch c {
				case '\a':
					c = 'a'
				case '\b':
					c = 'b'
				case '\t':
					c = 't'
				case '\n':
					c = 'n'
				case '\v':
					c = 'v'
				case '\f':
					c = 'f'
				case '\r':
					c = 'r'
				case '"', '\\':
				default:
					b.WriteByte('x')
					b.WriteByte(hex[c>>4])
					c = hex[c&0xF]
				}
			}
			b.WriteByte(c)
		}
		eol := "\" +\n"
		if i >= len(s) {
			eol = "\"\n"
		}
		b.WriteString(eol)
		b.WriteTo(w)
	}
}

func printCodes(w io.Writer, codes *[256]uint32) {
	zeroPad := append(make([]byte, 0, 7+8), "0000000"...)

	io.WriteString(w, "var codes = [256]uint32{\n\t")

	for i, code := range codes {
		io.WriteString(w, "0x")
		v := strconv.AppendUint(zeroPad, uint64(code), 16)
		w.Write(v[len(v)-8:])
		delim := ", "
		if i%8 == 7 {
			delim = ",\n\t"
			if i >= len(codes)-1 {
				delim = ",\n"
			}
		}
		io.WriteString(w, delim)
	}
	io.WriteString(w, "}\n")
}

func memset(p []byte, b byte) {
	for i := range p {
		p[i] = b
	}
}

const codeLengths = "" +
	"\x0d\x17\x1c\x1c\x1c\x1c\x1c\x1c\x1c\x18\x1e\x1c\x1c\x1e\x1c\x1c" +
	"\x1c\x1c\x1c\x1c\x1c\x1c\x1e\x1c\x1c\x1c\x1c\x1c\x1c\x1c\x1c\x1c" +
	"\x06\x0a\x0a\x0c\x0d\x06\x08\x0b\x0a\x0a\x08\x0b\x08\x06\x06\x06" +
	"\x05\x05\x05\x06\x06\x06\x06\x06\x06\x06\x07\x08\x0f\x06\x0c\x0a" +
	"\x0d\x06\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07" +
	"\x07\x07\x07\x07\x07\x07\x07\x07\x08\x07\x08\x0d\x13\x0d\x0e\x06" +
	"\x0f\x05\x06\x05\x06\x05\x06\x06\x06\x05\x07\x07\x06\x06\x06\x05" +
	"\x06\x07\x06\x05\x05\x06\x07\x07\x07\x07\x07\x0f\x0b\x0e\x0d\x1c" +
	"\x14\x16\x14\x14\x16\x16\x16\x17\x16\x17\x17\x17\x17\x17\x18\x17" +
	"\x18\x18\x16\x17\x18\x17\x17\x17\x17\x15\x16\x17\x16\x17\x17\x18" +
	"\x16\x15\x14\x16\x16\x17\x17\x15\x17\x16\x16\x18\x15\x16\x17\x17" +
	"\x15\x15\x16\x15\x17\x16\x17\x17\x14\x16\x16\x16\x17\x16\x16\x17" +
	"\x1a\x1a\x14\x13\x16\x17\x16\x19\x1a\x1a\x1a\x1b\x1b\x1a\x18\x19" +
	"\x13\x15\x1a\x1b\x1b\x1a\x1b\x18\x15\x15\x1a\x1a\x1c\x1b\x1b\x1b" +
	"\x14\x18\x14\x15\x16\x15\x15\x17\x16\x16\x19\x19\x18\x18\x1a\x17" +
	"\x1a\x1b\x1a\x1a\x1b\x1b\x1b\x1b\x1b\x1c\x1b\x1b\x1b\x1b\x1b\x1a"

var codes = [256]uint32{
	0x00001ff8, 0x007fffd8, 0x0fffffe2, 0x0fffffe3, 0x0fffffe4, 0x0fffffe5, 0x0fffffe6, 0x0fffffe7,
	0x0fffffe8, 0x00ffffea, 0x3ffffffc, 0x0fffffe9, 0x0fffffea, 0x3ffffffd, 0x0fffffeb, 0x0fffffec,
	0x0fffffed, 0x0fffffee, 0x0fffffef, 0x0ffffff0, 0x0ffffff1, 0x0ffffff2, 0x3ffffffe, 0x0ffffff3,
	0x0ffffff4, 0x0ffffff5, 0x0ffffff6, 0x0ffffff7, 0x0ffffff8, 0x0ffffff9, 0x0ffffffa, 0x0ffffffb,
	0x00000014, 0x000003f8, 0x000003f9, 0x00000ffa, 0x00001ff9, 0x00000015, 0x000000f8, 0x000007fa,
	0x000003fa, 0x000003fb, 0x000000f9, 0x000007fb, 0x000000fa, 0x00000016, 0x00000017, 0x00000018,
	0x00000000, 0x00000001, 0x00000002, 0x00000019, 0x0000001a, 0x0000001b, 0x0000001c, 0x0000001d,
	0x0000001e, 0x0000001f, 0x0000005c, 0x000000fb, 0x00007ffc, 0x00000020, 0x00000ffb, 0x000003fc,
	0x00001ffa, 0x00000021, 0x0000005d, 0x0000005e, 0x0000005f, 0x00000060, 0x00000061, 0x00000062,
	0x00000063, 0x00000064, 0x00000065, 0x00000066, 0x00000067, 0x00000068, 0x00000069, 0x0000006a,
	0x0000006b, 0x0000006c, 0x0000006d, 0x0000006e, 0x0000006f, 0x00000070, 0x00000071, 0x00000072,
	0x000000fc, 0x00000073, 0x000000fd, 0x00001ffb, 0x0007fff0, 0x00001ffc, 0x00003ffc, 0x00000022,
	0x00007ffd, 0x00000003, 0x00000023, 0x00000004, 0x00000024, 0x00000005, 0x00000025, 0x00000026,
	0x00000027, 0x00000006, 0x00000074, 0x00000075, 0x00000028, 0x00000029, 0x0000002a, 0x00000007,
	0x0000002b, 0x00000076, 0x0000002c, 0x00000008, 0x00000009, 0x0000002d, 0x00000077, 0x00000078,
	0x00000079, 0x0000007a, 0x0000007b, 0x00007ffe, 0x000007fc, 0x00003ffd, 0x00001ffd, 0x0ffffffc,
	0x000fffe6, 0x003fffd2, 0x000fffe7, 0x000fffe8, 0x003fffd3, 0x003fffd4, 0x003fffd5, 0x007fffd9,
	0x003fffd6, 0x007fffda, 0x007fffdb, 0x007fffdc, 0x007fffdd, 0x007fffde, 0x00ffffeb, 0x007fffdf,
	0x00ffffec, 0x00ffffed, 0x003fffd7, 0x007fffe0, 0x00ffffee, 0x007fffe1, 0x007fffe2, 0x007fffe3,
	0x007fffe4, 0x001fffdc, 0x003fffd8, 0x007fffe5, 0x003fffd9, 0x007fffe6, 0x007fffe7, 0x00ffffef,
	0x003fffda, 0x001fffdd, 0x000fffe9, 0x003fffdb, 0x003fffdc, 0x007fffe8, 0x007fffe9, 0x001fffde,
	0x007fffea, 0x003fffdd, 0x003fffde, 0x00fffff0, 0x001fffdf, 0x003fffdf, 0x007fffeb, 0x007fffec,
	0x001fffe0, 0x001fffe1, 0x003fffe0, 0x001fffe2, 0x007fffed, 0x003fffe1, 0x007fffee, 0x007fffef,
	0x000fffea, 0x003fffe2, 0x003fffe3, 0x003fffe4, 0x007ffff0, 0x003fffe5, 0x003fffe6, 0x007ffff1,
	0x03ffffe0, 0x03ffffe1, 0x000fffeb, 0x0007fff1, 0x003fffe7, 0x007ffff2, 0x003fffe8, 0x01ffffec,
	0x03ffffe2, 0x03ffffe3, 0x03ffffe4, 0x07ffffde, 0x07ffffdf, 0x03ffffe5, 0x00fffff1, 0x01ffffed,
	0x0007fff2, 0x001fffe3, 0x03ffffe6, 0x07ffffe0, 0x07ffffe1, 0x03ffffe7, 0x07ffffe2, 0x00fffff2,
	0x001fffe4, 0x001fffe5, 0x03ffffe8, 0x03ffffe9, 0x0ffffffd, 0x07ffffe3, 0x07ffffe4, 0x07ffffe5,
	0x000fffec, 0x00fffff3, 0x000fffed, 0x001fffe6, 0x003fffe9, 0x001fffe7, 0x001fffe8, 0x007ffff3,
	0x003fffea, 0x003fffeb, 0x01ffffee, 0x01ffffef, 0x00fffff4, 0x00fffff5, 0x03ffffea, 0x007ffff4,
	0x03ffffeb, 0x07ffffe6, 0x03ffffec, 0x03ffffed, 0x07ffffe7, 0x07ffffe8, 0x07ffffe9, 0x07ffffea,
	0x07ffffeb, 0x0ffffffe, 0x07ffffec, 0x07ffffed, 0x07ffffee, 0x07ffffef, 0x07fffff0, 0x03ffffee,
}
