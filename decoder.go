package quack

import (
	"errors"

	"github.com/renthraysk/quack/ascii"
	"github.com/renthraysk/quack/huffman"
)

var (
	errUnexpectedEnd         = errors.New("unexpected end")
	errUnexpectedTypeByte    = errors.New("unexpected type byte 0b000X_XXXX")
	errStaticIndexOutOfRange = errors.New("static index out of range")
	errNameInvalid           = errors.New("invalid name")
	errValueInvalid          = errors.New("invalid value")
)

type Decoder struct {
	dt DT
}

func NewDecoder(maxCapacity uint64) *Decoder {
	return &Decoder{dt: DT{capacity: maxCapacity}}
}

// Decode decodes the header fields in p.
func (d *Decoder) Decode(p []byte, accept func(string, string)) error {
	_, p, err := readVarint(p, 0xFF)
	if err != nil {
		return err
	}
	_, p, err = readVarint(p, 0x7F)
	if err != nil {
		return err
	}

	buf := make([]byte, 0, 256) // Huffman decode scratch buffer

	for len(p) > 0 {
		switch (p[0] >> 4) & 0b1111 { // & 0b1111 should be unnecessary
		case 0b0000:
			//  0000_NXXX Literal Field Line with Post-Base Name Reference
			// 	https://www.rfc-editor.org/rfc/rfc9204.html#name-literal-field-line-with-pos
			const NeverIndex = 0b0000_1000

			index, q, err := readVarint(p, 0b0000_0111)
			if err != nil {
				return err
			}
			value, q, err := readStringLiteral(q, buf)
			if err != nil {
				return err
			}
			name, err := d.dt.baseNameIndex(index)
			if err != nil {
				return err
			}
			if p[0]&NeverIndex != NeverIndex {
				// Index
			}
			p = q
			accept(name, value)

		case 0b0001:
			// 0001_XXXX Indexed Field Line with Post-Base Index
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-indexed-field-line-with-pos
			index, q, err := readVarint(p, 0b0000_1111)
			if err != nil {
				return err
			}
			name, value, err := d.dt.baseLineIndex(index)
			if err != nil {
				return err
			}
			p = q
			accept(name, value)

		case 0b0010, 0b0011:
			// 001N_HXXX Literal Field Line with Literal Name
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-literal-field-line-with-lit
			const NeverIndex = 0b0001_0000

			name, q, err := readLiteralName(p, buf)
			if err != nil {
				return err
			}
			value, q, err := readStringLiteral(q, buf)
			if err != nil {
				return err
			}
			if p[0]&NeverIndex != NeverIndex {
				// Index
			}
			p = q
			accept(name, value)

		case 0b0100, 0b0110:
			// 01N0_XXXX: Literal Field Line with Name Reference in dynamic table
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-literal-field-line-with-nam
			const NeverIndex = 0b0010_0000

			index, q, err := readVarint(p, 0b0000_1111)
			if err != nil {
				return err
			}
			name, err := d.dt.nameIndex(index)
			if err != nil {
				return err
			}
			value, q, err := readStringLiteral(q, buf)
			if err != nil {
				return err
			}
			if p[0]&NeverIndex != NeverIndex {
				// Index
			}
			p = q
			accept(name, value)

		case 0b0101, 0b0111:
			// 01N1_XXXX: Literal Field Line with Name Reference in static table
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-literal-field-line-with-nam
			const NeverIndex = 0b0010_0000

			index, q, err := readVarint(p, 0b0000_1111)
			if err != nil {
				return err
			}
			if index >= uint64(len(staticTable)) {
				return errStaticIndexOutOfRange
			}
			value, q, err := readStringLiteral(q, buf)
			if err != nil {
				return err
			}
			if p[0]&NeverIndex != NeverIndex {
				// Index
			}
			p = q
			accept(staticTable[index].Name, value)

		case 0b1000, 0b1001, 0b1010, 0b1011:
			// 10XX_XXXX Indexed Field Line in dynamic table
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-indexed-field-line
			index, q, err := readVarint(p, 0b0011_1111)
			if err != nil {
				return err
			}
			name, value, err := d.dt.lineIndex(index)
			if err != nil {
				return err
			}
			p = q
			accept(name, value)

		case 0b1100, 0b1101, 0b1110, 0b1111:
			// 11XX_XXXX Indexed Field Line in static table
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-indexed-field-line
			index, q, err := readVarint(p, 0b0011_1111)
			if err != nil {
				return err
			}
			if index >= uint64(len(staticTable)) {
				return errStaticIndexOutOfRange
			}
			p = q
			accept(staticTable[index].Name, staticTable[index].Value)
		}
	}
	return nil
}

// readLiteralName reads a literal name from p. Will use decodeBuf if the
// string needs to be huffman decoded.
func readLiteralName(p, decodeBuf []byte) (string, []byte, error) {
	const (
		// layout of the first byte of a literal name length
		Prefix         = 0b0010_0000
		NeverIndex     = 0b0001_0000
		HuffmanEncoded = 0b0000_1000
		NameLenBits    = 0b0000_0111
	)
	if len(p) <= 0 {
		return "", p, errUnexpectedEnd
	}
	n, q, err := readVarint(p, NameLenBits)
	if err != nil {
		return "", p, err
	}
	if n > uint64(len(q)) {
		return "", p, errUnexpectedEnd
	}
	b := q[:n:n]
	if p[0]&HuffmanEncoded == HuffmanEncoded {
		b, err = huffman.Decode(decodeBuf[:0], b)
		if err != nil {
			return "", p, err
		}
	}
	// Don't allocate for obvious garbage.
	if !ascii.IsNameValid(b) {
		return "", p, errNameInvalid
	}
	return ascii.ToCanonical(b), q[n:], nil // Allocation
}

// readStringLiteral reads a string literal from p. Will use decodeBuf if the
// string needs to be huffman decoded.
func readStringLiteral(p, decodeBuf []byte) (string, []byte, error) {
	const (
		// layout of the first byte of a string literal length
		HuffmanEncoded = 0b1000_0000
		StringLenBits  = 0b0111_1111
	)

	if len(p) <= 0 {
		return "", p, errUnexpectedEnd
	}
	n, q, err := readVarint(p, StringLenBits)
	if err != nil {
		return "", p, err
	}
	if n > uint64(len(q)) {
		return "", p, errUnexpectedEnd
	}
	b := q[:n:n]
	if p[0]&HuffmanEncoded == HuffmanEncoded {
		b, err = huffman.Decode(decodeBuf[:0], b)
		if err != nil {
			return "", p, err
		}
	}
	// Don't allocate for obvious garbage.
	if !ascii.IsValueValid(b) {
		return "", p, errValueInvalid
	}
	return string(b), q[n:], nil // Allocation
}
