package quack

import (
	"github.com/renthraysk/quack/ascii"
	"github.com/renthraysk/quack/huffman"
)

type value struct {
	value string
	index uint64
}

// fieldEncoder field line encoder, immutable once created
type fieldEncoder struct {
	m                  map[string][]value
	encodedInsertCount uint64
	deltaBase          uint64
	sign               uint8
}

// Field Line Representations
// https://www.rfc-editor.org/rfc/rfc9204.html#name-field-line-representations

// https://www.rfc-editor.org/rfc/rfc9204.html#name-encoded-field-section-prefi
func (fe *fieldEncoder) appendFieldSectionPrefix(p []byte) []byte {
	if fe == nil {
		// Operating with only static table.
		return append(p, 0, 0)
	}
	p = appendVarint(p, fe.encodedInsertCount, 0xFF, 0)
	return appendVarint(p, fe.deltaBase, 0b0111_1111, fe.sign)
}

func (fe *fieldEncoder) lookup(name, value string) (uint64, bool, match) {
	i, m := staticLookup(name, value)
	if fe == nil || m == matchNameValue {
		// Operating with only static table or have the best match already.
		return i, true, m
	}
	values, ok := fe.m[name]
	if ok {
		for _, v := range values {
			if v.value == value {
				return v.index, false, matchNameValue
			}
		}
	}
	switch {
	case m == matchName:
		return i, true, matchName
	case ok && len(values) > 0:
		return values[0].index, false, matchName
	}
	return 0, false, matchNone
}

func (fe *fieldEncoder) appendFieldLines(p []byte, header map[string][]string) []byte {
	for name, values := range header {
		for _, value := range values {
			p = fe.appendFieldLine(p, name, value)
		}
	}
	return p
}

func (fe *fieldEncoder) appendFieldLine(p []byte, name, value string) []byte {
	ctrl := headerControl(name)
	switch i, isStatic, m := fe.lookup(name, value); m {
	case matchNameValue:
		if isStatic {
			return appendIndexedLine(p, i, true)
		}
		return appendIndexedLinePostBase(p, i)

	case matchName:
		p = appendNamedReference(p, i, ctrl.neverIndex(), isStatic)
	case matchNone:
		p = appendLiteralName(p, name, ctrl.neverIndex())
	}
	return appendStringLiteral(p, value, ctrl.shouldHuffman())
}

/* */

type fieldDecoder struct {
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-indexed-field-line
func appendIndexedLine(p []byte, i uint64, isStatic bool) []byte {
	const (
		P = 0b1000_0000
		T = 0b0100_0000
		M = 0b0011_1111
	)
	return appendVarint(p, i, M, P|t(isStatic, T))
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-indexed-field-line-with-pos
func appendIndexedLinePostBase(p []byte, i uint64) []byte {
	const P = 0b0001_0000
	const M = 0b0000_1111

	return appendVarint(p, i, M, P)
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-literal-field-line-with-nam
func appendNamedReference(p []byte, i uint64, neverIndex, isStatic bool) []byte {
	const (
		P = 0b0100_0000
		N = 0b0010_0000
		T = 0b0001_0000
		M = 0b0000_1111
	)
	return appendVarint(p, i, M, P|t(neverIndex, N)|t(isStatic, T))
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-literal-field-line-with-lit
func appendLiteralName(p []byte, name string, neverIndex bool) []byte {
	const (
		P = 0b0010_0000
		N = 0b0001_0000
		H = 0b0000_1000
		M = 0b0000_0111
	)
	n := uint64(len(name))
	if h := huffman.EncodeLengthLower(name); h < n {
		p = appendVarint(p, h, M, P|t(neverIndex, N)|H)
		return huffman.AppendStringLower(p, name)
	}
	p = appendVarint(p, n, M, P|t(neverIndex, N))
	return ascii.AppendLower(p, name)
}
