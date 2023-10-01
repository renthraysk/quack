package inst

import (
	"github.com/renthraysk/quack/ascii"
	"github.com/renthraysk/quack/huffman"
	"github.com/renthraysk/quack/varint"
)

// AppendSetDynamicTableCapacity https://www.rfc-editor.org/rfc/rfc9204.html#name-set-dynamic-table-capacity
func AppendSetDynamicTableCapacity(p []byte, capacity uint64) []byte {
	const (
		P = 0b0010_0000 // Prefix
		M = 0b0001_1111 // Mask
	)
	return varint.Append(p, P, M, capacity)
}

// AppendInsertWithNameReference https://www.rfc-editor.org/rfc/rfc9204.html#name-insert-with-name-reference
func AppendInsertWithNameReference(p []byte, i uint64, isStatic bool) []byte {
	const (
		P = 0b1000_0000 // Prefix
		T = 0b0100_0000 // Static Table
		M = 0b0011_1111 // Mask
	)
	return varint.Append(p, P|t(isStatic, T), M, i)
}

// AppendInsertWithLiteralName https://www.rfc-editor.org/rfc/rfc9204.html#name-insert-with-literal-name
func AppendInsertWithLiteralName(p []byte, name string) []byte {
	const (
		P = 0b0100_0000 // Prefix
		H = 0b0010_0000 // Huffman Encoded
		M = 0b0001_1111 // Mask
	)

	n := uint64(len(name))
	if h := huffman.EncodeLengthLower(name); h < n {
		p = varint.Append(p, P|H, M, h)
		return huffman.AppendStringLower(p, name)
	}
	p = varint.Append(p, P, M, n)
	return ascii.AppendLower(p, name)
}

// AppendDuplicate https://www.rfc-editor.org/rfc/rfc9204.html#name-duplicate
func AppendDuplicate(p []byte, i uint64) []byte {
	const (
		P = 0b0000_0000 // Prefix
		M = 0b0001_1111 // Mask
	)
	return varint.Append(p, P, M, i)
}
