package inst

import (
	"github.com/renthraysk/quack/ascii"
	"github.com/renthraysk/quack/huffman"
	"github.com/renthraysk/quack/varint"
)

// https://www.rfc-editor.org/rfc/rfc9204.html#name-set-dynamic-table-capacity
func AppendSetDynamicTableCapacity(p []byte, capacity uint64) []byte {
	const (
		P = 0b0010_0000
		M = 0b0001_1111
	)
	return varint.Append(p, capacity, M, P)
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-insert-with-name-reference
func AppendInsertWithNameReference(p []byte, i uint64, isStatic bool) []byte {
	const (
		P = 0b1000_0000
		T = 0b0100_0000
		M = 0b0011_1111
	)
	return varint.Append(p, i, M, P|t(isStatic, T))
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-insert-with-literal-name
func AppendInsertWithLiteralName(p []byte, name string) []byte {
	const (
		P = 0b0100_0000
		H = 0b0010_0000
		M = 0b0001_1111
	)

	n := uint64(len(name))
	if h := huffman.EncodeLengthLower(name); h < n {
		p = varint.Append(p, h, M, P|H)
		return huffman.AppendStringLower(p, name)
	}
	p = varint.Append(p, n, M, P)
	return ascii.AppendLower(p, name)
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-duplicate
func AppendDuplicate(p []byte, i uint64) []byte {
	const (
		P = 0b0000_0000
		M = 0b0001_1111
	)
	return varint.Append(p, i, M, P)
}

//
