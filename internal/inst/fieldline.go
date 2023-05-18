package inst

import (
	"github.com/renthraysk/quack/ascii"
	"github.com/renthraysk/quack/huffman"
	"github.com/renthraysk/quack/varint"
)

// Field Line Representations
// https://www.rfc-editor.org/rfc/rfc9204.html#name-field-line-representations

// https://www.rfc-editor.org/rfc/rfc9204.html#name-indexed-field-line
func AppendIndexedLine(p []byte, i uint64, isStatic bool) []byte {
	const (
		P = 0b1000_0000
		T = 0b0100_0000
		M = 0b0011_1111
	)
	return varint.Append(p, i, M, P|t(isStatic, T))
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-indexed-field-line-with-pos
func AppendIndexedLinePostBase(p []byte, i uint64) []byte {
	const P = 0b0001_0000
	const M = 0b0000_1111

	return varint.Append(p, i, M, P)
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-literal-field-line-with-nam
func AppendNamedReference(p []byte, i uint64, neverIndex, isStatic bool) []byte {
	const (
		P = 0b0100_0000
		N = 0b0010_0000
		T = 0b0001_0000
		M = 0b0000_1111
	)
	return varint.Append(p, i, M, P|t(neverIndex, N)|t(isStatic, T))
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-literal-field-line-with-lit
func AppendLiteralName(p []byte, name string, neverIndex bool) []byte {
	const (
		P = 0b0010_0000
		N = 0b0001_0000
		H = 0b0000_1000
		M = 0b0000_0111
	)
	n := uint64(len(name))
	if h := huffman.EncodeLengthLower(name); h < n {
		p = varint.Append(p, h, M, P|t(neverIndex, N)|H)
		return huffman.AppendStringLower(p, name)
	}
	p = varint.Append(p, n, M, P|t(neverIndex, N))
	return ascii.AppendLower(p, name)
}

// AppendStringLiteral appends the QPACK encoded string literal s to p.
func AppendStringLiteral(p []byte, s string, shouldHuffman bool) []byte {
	const (
		H = 0b1000_0000
		M = 0b0111_1111
	)

	n := uint64(len(s))
	if n > 2 && shouldHuffman {
		if h := huffman.EncodeLength(s); h < n {
			p = varint.Append(p, h, M, H)
			return huffman.AppendString(p, s)
		}
	}
	p = varint.Append(p, n, M, 0)
	return append(p, s...)
}

func t(b bool, t byte) byte {
	if b {
		return t
	}
	return 0
}
