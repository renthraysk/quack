package inst

import (
	"github.com/renthraysk/quack/ascii"
	"github.com/renthraysk/quack/huffman"
	"github.com/renthraysk/quack/varint"
)

// Field Line Representations
// https://www.rfc-editor.org/rfc/rfc9204.html#name-field-line-representations

// AppendStaticIndexReference https://www.rfc-editor.org/rfc/rfc9204.html#name-indexed-field-line
func AppendStaticIndexReference(p []byte, i uint64) []byte {
	const (
		P = 0b1000_0000 // Prefix
		T = 0b0100_0000 // Static Table
		M = 0b0011_1111 // Mask
	)
	return varint.Append(p, P|T, M, i)
}

// AppendIndexedLinePostBase https://www.rfc-editor.org/rfc/rfc9204.html#name-indexed-field-line-with-pos
func AppendIndexedLinePostBase(p []byte, i uint64) []byte {
	const P = 0b0001_0000 // Prefix
	const M = 0b0000_1111 // Mask

	return varint.Append(p, P, M, i)
}

// AppendNamedReference https://www.rfc-editor.org/rfc/rfc9204.html#name-literal-field-line-with-nam
func AppendNamedReference(p []byte, i uint64, neverIndex, isStatic bool) []byte {
	const (
		P = 0b0100_0000 // Prefix
		N = 0b0010_0000 // Never Index
		T = 0b0001_0000 // Static Table
		M = 0b0000_1111 // Mask
	)
	return varint.Append(p, P|t(neverIndex, N)|t(isStatic, T), M, i)
}

// AppendLiteralName https://www.rfc-editor.org/rfc/rfc9204.html#name-literal-field-line-with-lit
func AppendLiteralName(p []byte, name string, neverIndex bool) []byte {
	const (
		P = 0b0010_0000 // Prefix
		N = 0b0001_0000 // Never Index
		H = 0b0000_1000 // Huffman Encoded
		M = 0b0000_0111 // Mask
	)
	n := uint64(len(name))
	if h := huffman.EncodeLengthLower(name); h < n {
		p = varint.Append(p, P|t(neverIndex, N)|H, M, h)
		return huffman.AppendStringLower(p, name)
	}
	p = varint.Append(p, P|t(neverIndex, N), M, n)
	return ascii.AppendLower(p, name)
}

func t(b bool, t byte) byte {
	if b {
		return t
	}
	return 0
}
