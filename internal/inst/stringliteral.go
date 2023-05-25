package inst

import (
	"github.com/renthraysk/quack/huffman"
	"github.com/renthraysk/quack/varint"
)

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
