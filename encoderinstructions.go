package quack

import (
	"errors"
	"math/bits"

	"github.com/renthraysk/quack/ascii"
	"github.com/renthraysk/quack/huffman"
)

func (dt *DT) appendEncoderInstructionLocked(p []byte, name, value string) []byte {
	i, isStatic, m := dt.current.lookup(name, value)
	if m == matchNameValue {
		// @TODO Duplicate?
		return p
	}
	ctrl := headerControl(name)
	if ctrl.neverIndex() {
		// If already have a name match, no point attempting an insert
		// if prevented from inserting a (name, value) pair.
		if m == matchName {
			return p
		}
		value = ""
	}
	if ok := dt.insertLocked(name, value); !ok {
		return p
	}
	// successful insertion into dynamic table, so need an encoder
	// instruction to inform peer
	switch m {
	case matchName:
		p = appendInsertWithNameReference(p, i, isStatic)
	case matchNone:
		p = appendInsertWithLiteralName(p, name)
	}
	return appendStringLiteral(p, value, ctrl.shouldHuffman())
}

// Encoder Instructions
// https://www.rfc-editor.org/rfc/rfc9204.html#name-encoder-instructions
func (dt *DT) appendEncoderInstructions(p []byte, header map[string][]string) ([]byte, *fieldEncoder) {

	dt.mu.Lock()
	defer dt.mu.Unlock()

	for name, values := range header {
		for _, value := range values {
			p = dt.appendEncoderInstructionLocked(p, name, value)
		}
	}
	// Build a fieldEncoder for when peer acks.
	m := make(map[string][]value, len(dt.headers))
	for i, hf := range dt.headers {
		m[hf.Name] = append(m[hf.Name], value{value: hf.Value, index: uint64(i)})
	}

	var reqInsertCount uint64
	var encodedInsertCount uint64

	if reqInsertCount > 0 {
		maxEntries := dt.capacity / 32
		encodedInsertCount = (reqInsertCount % (2 * maxEntries)) + 1
	}

	deltaBase, sign := bits.Sub64(dt.base, reqInsertCount, 0)
	if sign != 0 {
		deltaBase = reqInsertCount - dt.base - 1
		sign = 0x80
	}
	return p, &fieldEncoder{m: m, encodedInsertCount: encodedInsertCount, deltaBase: deltaBase, sign: byte(sign)}

}

func (dt *DT) parseEncoderInstructions(p []byte) error {
	var decodeBuf [256]byte

	dt.mu.Lock()
	defer dt.mu.Unlock()

	for len(p) > 0 {
		switch p[0] >> 5 {
		case 0b000:
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-duplicate
			const M = 0b0001_1111

			i, q, err := readVarint(p, M)
			if err != nil {
				return err
			}
			if i > uint64(len(dt.headers)) {
				return errors.New("over")
			}
			if ok := dt.insertLocked(dt.headers[i].Name, dt.headers[i].Value); !ok {
			}
			p = q

		case 0b001:
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-set-dynamic-table-capacity
			const M = 0b0001_1111

			capacity, q, err := readVarint(p, M)
			if err != nil {
				return err
			}
			if ok := dt.setCapacityLocked(capacity); !ok {
			}
			p = q

		case 0b010, 0b011:
			name, q, err := nameInsertWithLiteralName(p, decodeBuf[:0])
			if err != nil {
				return err
			}
			value, q, err := readStringLiteral(q, decodeBuf[:0])
			if err != nil {
				return err
			}
			if ok := dt.insertLocked(name, value); !ok {
			}
			p = q

		default:
			name, q, err := dt.nameInsertWithNameReference(p)
			if err != nil {
				return err
			}
			value, q, err := readStringLiteral(q, decodeBuf[:0])
			if err != nil {
				return err
			}
			if ok := dt.insertLocked(name, value); !ok {
			}
			p = q
		}
	}
	return nil
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-insert-with-literal-name
func nameInsertWithLiteralName(p, buf []byte) (string, []byte, error) {
	const H = 0b0010_0000
	const M = 0b0001_1111

	n, q, err := readVarint(p, M)
	if err != nil {
		return "", p, err
	}
	if n > uint64(len(q)) {
		return "", p, errors.New("")
	}
	b := q[:n]
	if p[0]&H == H {
		b, err = huffman.Decode(buf[:0], b)
		if err != nil {
			return "", p, err
		}
	}
	if !ascii.IsName3Valid(b) {
	}
	return ascii.ToCanonical(b), q[n:], nil
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-insert-with-name-reference
func (dt *DT) nameInsertWithNameReference(p []byte) (string, []byte, error) {
	const T = 0b0100_0000
	const M = 0b0011_1111

	i, q, err := readVarint(p, M)
	if err != nil {
		return "", p, err
	}
	if p[0]&T == T {
		if i >= uint64(len(staticTable)) {
			return "", p, errors.New("invalid static table index")
		}
		return staticTable[i].Name, q, nil
	}
	if i >= uint64(len(dt.headers)) {
		return "", p, errors.New("invalid dynamic table index")
	}
	return dt.headers[i].Name, q, nil
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

//
