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

type Encoder struct {
	dt DT
}

func NewEncoder(capacity uint64) *Encoder {
	return &Encoder{dt: DT{capacity: capacity}}
}

func (e *Encoder) Parse(p []byte) error {
	var streamID, increment uint64
	var err error

	for len(p) > 0 {
		switch p[0] & 0b1100_0000 {
		case 0b0000_0000:
			// insert count increment
			increment, p, err = readVarint(p, 0b0011_1111)
			if err != nil {
				return err
			}
			_ = increment
			// @TODO

		case 0b0100_0000:
			// stream cancellation
			streamID, p, err = readVarint(p, 0b0011_1111)
			if err != nil {
				return err
			}
			_ = streamID
			// @TODO

		case 0b1000_0000, 0b1100_0000:
			// section acknowledgement
			streamID, p, err = readVarint(p, 0b0111_1111)
			if err != nil {
				return err
			}
			_ = streamID
			// @TODO
		}
	}
	return nil
}

// Encoder Instructions
// https://www.rfc-editor.org/rfc/rfc9204.html#name-encoder-instructions

func (e *Encoder) appendEncoderInstructions(p []byte, header map[string][]string) ([]byte, uint64) {

	var reqInsertCount uint64 // @TODO

	for name, values := range header {
		for _, value := range values {
			si, sm := staticLookup(name, value)
			if sm == matchNameValue {
				continue
			}
			di, dm := e.dt.lookup(name, value)
			if dm == matchNameValue {
				continue
			}
			ctrl := e.headerControl(name)
			if ctrl.neverIndex() {
				// If already have a name match, no point attempting an insert
				// if prevented from inserting a (name, value) pair.
				if sm == matchName || dm == matchName {
					continue
				}
				// sm == matchNone && dm == matchNone
				value = ""
			}
			if ok := e.dt.insert(name, value); !ok {
				continue
			}
			// successful insertion into dynamic table, so need an encoder
			// instruction to inform peer
			switch {
			case sm == matchName:
				p = appendInsertWithNameReference(p, si, true)

			case dm == matchName:
				p = appendInsertWithNameReference(p, di, false)

			default:
				p = appendInsertWithLiteralName(p, name)
			}
			p = appendStringLiteral(p, value, ctrl.shouldHuffman())
		}
	}
	return p, reqInsertCount
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-set-dynamic-table-capacity
func appendSetDynamicTableCapacity(p []byte, capacity uint64) []byte {
	const (
		P = 0b0010_0000
		M = 0b0001_1111
	)
	return appendVarint(p, capacity, M, P)
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-insert-with-name-reference
func appendInsertWithNameReference(p []byte, i uint64, isStatic bool) []byte {
	const (
		P = 0b1000_0000
		T = 0b0100_0000
		M = 0b0011_1111
	)
	return appendVarint(p, i, M, P|t(isStatic, T))
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-insert-with-literal-name
func appendInsertWithLiteralName(p []byte, name string) []byte {
	const (
		P = 0b0100_0000
		H = 0b0010_0000
		M = 0b0001_1111
	)

	n := uint64(len(name))
	if h := huffman.EncodeLengthLower(name); h < n {
		p = appendVarint(p, h, M, P|H)
		return huffman.AppendStringLower(p, name)
	}
	p = appendVarint(p, n, M, P)
	return ascii.AppendLower(p, name)
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-duplicate
func appendDuplicate(p []byte, i uint64) []byte {
	const (
		P = 0b0000_0000
		M = 0b0001_1111
	)
	return appendVarint(p, i, M, P)
}

// Field Line Representations
// https://www.rfc-editor.org/rfc/rfc9204.html#name-field-line-representations

func (e *Encoder) appendFieldLines(p []byte, header map[string][]string) []byte {
	for name, values := range header {
		for _, value := range values {
			si, sm := staticLookup(name, value)
			if sm == matchNameValue {
				p = appendIndexedLine(p, si, true)
				continue
			}
			di, dm := e.dt.lookup(name, value)
			if dm == matchNameValue {
				p = appendIndexedLinePostBase(p, di-e.dt.base)
				continue
			}
			ctrl := e.headerControl(name)
			switch {
			case sm == matchName:
				p = appendNamedReference(p, si, ctrl.neverIndex(), true)
			case dm == matchName:
				p = appendNamedReference(p, di, ctrl.neverIndex(), false)
			default:
				p = appendLiteralName(p, name, ctrl.neverIndex())
			}
			p = appendStringLiteral(p, value, ctrl.shouldHuffman())
		}
	}
	return p
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
func appendNamedReference(p []byte, i uint64, neverIndex bool, isStatic bool) []byte {
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

//

// appendStringLiteral appends the QPACK encoded string literal s to p.
func appendStringLiteral(p []byte, s string, shouldHuffman bool) []byte {
	const (
		H = 0b1000_0000
		M = 0b0111_1111
	)

	n := uint64(len(s))
	if n > 2 && shouldHuffman {
		if h := huffman.EncodeLength(s); h < n {
			p = appendVarint(p, h, M, H)
			return huffman.AppendString(p, s)
		}
	}
	p = appendVarint(p, n, M, 0)
	return append(p, s...)
}

func t(b bool, t byte) byte {
	if b {
		return t
	}
	return 0
}
