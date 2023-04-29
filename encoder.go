package quack

import (
	"github.com/renthraysk/quack/ascii"
	"github.com/renthraysk/quack/huffman"
)

// match returned status of a table search
type match uint

const (
	// matchNone No match
	matchNone match = iota
	// matchName Matched name only
	matchName
	// matchNameValue Matched name & value
	matchNameValue
)

// Control controls finer details of how a specific headers should be encoded.
type Control uint8

const (
	// NeverIndex header field should never be put in the dynamic table.
	NeverIndex Control = 1 << iota
	// NeverHuffman never compress the value field.
	NeverHuffman
)

func (c Control) NeverIndex() bool    { return c&NeverIndex != 0 }
func (c Control) ShouldHuffman() bool { return c&NeverHuffman == 0 }

// neverIndex returns the value of yes if the header should not be index, 0
// otherwise.
func (c Control) neverIndex(yes byte) byte {
	if c.NeverIndex() {
		return yes
	}
	return 0
}

// defaultHeaderControls default set of headers that require special encoding
// treatment
var defaultHeaderControls = map[string]Control{
	"authorization":       NeverIndex | NeverHuffman,
	"content-md5":         NeverIndex | NeverHuffman,
	"date":                NeverIndex,
	"etag":                NeverIndex,
	"if-modified-since":   NeverIndex,
	"if-unmodified-since": NeverIndex,
	"last-modified":       NeverIndex,
	"location":            NeverIndex,
	"match":               NeverIndex,
	"range":               NeverIndex,
	"retry-after":         NeverIndex,
	"set-cookie":          NeverIndex,
}

type Encoder struct {
	dt             DT
	headerControls map[string]Control
}

func NewEncoder() *Encoder {
	return &Encoder{headerControls: defaultHeaderControls}
}

func (e *Encoder) encodeHeader(p []byte, header map[string][]string) []byte {
	var lowerBuf [len("access-control-allow-credentials")]byte

	for name, values := range header {
		// @TODO combine validation & AppendLower
		if !ascii.IsNameValid(name) {
			continue
		}
		lower := string(ascii.AppendLower(lowerBuf[:0], name))

		for _, value := range values {
			if ascii.IsValueValid(value) {
				p = e.encodeHeaderField(p, lower, value)
			}
		}
	}
	return p
}

func (e *Encoder) encodeHeaderField(p []byte, name, value string) []byte {
	const (
		// https://www.rfc-editor.org/rfc/rfc9204.html#name-literal-field-line-with-nam

		// P '01' 2-bit Pattern of literal field line with name reference
		P = 0b0100_0000
		// N Never index bit of literal field line with name reference
		N = 0b0010_0000
		// T Static table bit of literal field line with name reference
		T = 0b0001_0000
		// M Mask of literal field line with name reference
		M = 0b0000_1111
	)

	i, m := staticLookup(name, value)
	if m == matchNameValue {
		// https://www.rfc-editor.org/rfc/rfc9204.html#name-indexed-field-line
		return appendVarint(p, i, 0b0011_1111, 0b1100_0000)
	}

	ctrl := e.headerControls[name]
	if false /* @TODO */ {
		switch di, dm := e.dt.lookup(name, value); dm {
		case matchNameValue:
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-indexed-field-line
			indexBits, prefix := byte(0b0011_1111), byte(0b1000_0000)
			if false {
				// https://www.rfc-editor.org/rfc/rfc9204.html#name-indexed-field-line-with-pos
				indexBits, prefix = 0b0000_1111, 0b0001_0000
			}
			return appendVarint(p, di, indexBits, prefix)
		case matchName:
			// Prefer static table name matches over dynamic table name matches.
			if m == matchNone {
				p = appendVarint(p, di, M, P|ctrl.neverIndex(N))
				return appendStringLiteral(p, value, ctrl)
			}
		}
	}
	switch m {
	case matchName:
		p = appendVarint(p, i, M, P|ctrl.neverIndex(N)|T)
	case matchNone:
		p = appendLiteralName(p, name, ctrl)
	}
	return appendStringLiteral(p, value, ctrl)
}

func appendLiteralName(p []byte, name string, ctrl Control) []byte {
	const (
		// https://www.rfc-editor.org/rfc/rfc9204.html#name-literal-field-line-with-lit
		// '001' 3 bit pattern
		P = 0b0010_0000
		// Never index bit
		N = 0b0001_0000
		// Huffman encoded bit
		H = 0b0000_1000
		// Mask
		M = 0b0000_0111
	)
	n := uint64(len(name))
	if h := huffman.EncodeLengthLower(name); h < n {
		p = appendVarint(p, h, M, P|ctrl.neverIndex(N)|H)
		return huffman.AppendStringLower(p, name)
	}
	p = appendVarint(p, n, M, P|ctrl.neverIndex(N))
	return ascii.AppendLower(p, name)
}

// appendStringLiteral appends the QPACK encoded string literal s to p.
func appendStringLiteral(p []byte, s string, ctrl Control) []byte {
	const (
		// layout of the first byte of the length of a string literal
		// H Huffman encoded
		H = 0b1000_0000
		// M Mask
		M = 0b0111_1111
	)

	n := uint64(len(s))
	if n > 2 && ctrl.ShouldHuffman() {
		if h := huffman.EncodeLength(s); h < n {
			p = appendVarint(p, h, M, H)
			return huffman.AppendString(p, s)
		}
	}
	p = appendVarint(p, n, M, 0)
	return append(p, s...)
}
