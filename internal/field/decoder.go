package field

import (
	"errors"

	"github.com/renthraysk/quack/ascii"
	"github.com/renthraysk/quack/huffman"
	"github.com/renthraysk/quack/varint"
)

var (
	errUnexpectedEnd         = errors.New("unexpected end")
	errUnexpectedTypeByte    = errors.New("unexpected type byte 0b000X_XXXX")
	errStaticIndexOutOfRange = errors.New("static index out of range")
	errNameInvalid           = errors.New("invalid name")
	errValueInvalid          = errors.New("invalid value")
)

type Decoder struct {
	insertCount uint64
	capacity    uint64
}

// https://datatracker.ietf.org/doc/html/rfc9204#name-encoded-field-section-prefi
func (d *Decoder) readFieldSectionPrefix(p []byte) ([]byte, uint64, uint64, error) {
	var reqInsertCount uint64

	encodedInsertCount, q, err := varint.Read(p, 0xFF)
	if err != nil {
		return p, 0, 0, err
	}
	// https://datatracker.ietf.org/doc/html/rfc9204#name-required-insert-count
	if encodedInsertCount != 0 {
		maxEntries := d.capacity / 32

		fullRange := 2 * maxEntries
		if encodedInsertCount > fullRange {
			return p, 0, 0, errors.New("encodedInsertCount > fullRange")
		}
		maxValue := d.insertCount + maxEntries
		maxWrapped := (maxValue / fullRange) * fullRange
		reqInsertCount = maxWrapped + encodedInsertCount - 1
		if reqInsertCount > maxValue {
			if reqInsertCount <= fullRange {
				return p, 0, 0, errors.New("reqInsertCount <= fullRange")
			}
			reqInsertCount -= fullRange
		}
		if reqInsertCount == 0 {
			return p, 0, 0, errors.New("reqInsertCount of 0 not encoded as 0")
		}
	}

	// https://datatracker.ietf.org/doc/html/rfc9204#name-base
	const (
		S = 0b1000_0000
		M = 0b0111_1111
	)
	deltaBase, r, err := varint.Read(q, M)
	if err != nil {
		return p, 0, 0, err
	}

	base := reqInsertCount + deltaBase
	if q[0]&S != 0 {
		base = reqInsertCount - deltaBase - 1
	}
	return r, reqInsertCount, base, nil
}

// Decode decodes the header fields in p.
func (d *Decoder) Decode(p []byte, accept func(string, string)) error {

	q, reqInsertCount, base, err := d.readFieldSectionPrefix(p)
	if err != nil {
		return err
	}
	_, _ = reqInsertCount, base

	buf := make([]byte, 0, 256) // Huffman decode scratch buffer

	for len(q) > 0 {
		switch (q[0] >> 4) & 0b1111 { // & 0b1111 should be unnecessary
		case 0b0000:
			//  0000_NXXX Literal Field Line with Post-Base Name Reference
			// 	https://www.rfc-editor.org/rfc/rfc9204.html#name-literal-field-line-with-pos
			const NeverIndex = 0b0000_1000

			index, r, err := varint.Read(q, 0b0000_0111)
			if err != nil {
				return err
			}
			value, r, err := readStringLiteral(r, buf)
			if err != nil {
				return err
			}
			name, err := d.baseNameIndex(index)
			if err != nil {
				return err
			}
			q = r
			accept(name, value)

		case 0b0001:
			// 0001_XXXX Indexed Field Line with Post-Base Index
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-indexed-field-line-with-pos
			index, r, err := varint.Read(q, 0b0000_1111)
			if err != nil {
				return err
			}
			name, value, err := d.baseLineIndex(index)
			if err != nil {
				return err
			}
			q = r
			accept(name, value)

		case 0b0010, 0b0011:
			// 001N_HXXX Literal Field Line with Literal Name
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-literal-field-line-with-lit

			name, r, err := readLiteralName(q, buf)
			if err != nil {
				return err
			}
			value, r, err := readStringLiteral(r, buf)
			if err != nil {
				return err
			}
			q = r
			accept(name, value)

		case 0b0100, 0b0110:
			// 01N0_XXXX: Literal Field Line with Name Reference in dynamic table
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-literal-field-line-with-nam

			index, r, err := varint.Read(q, 0b0000_1111)
			if err != nil {
				return err
			}
			name, err := d.nameIndex(index)
			if err != nil {
				return err
			}
			value, r, err := readStringLiteral(r, buf)
			if err != nil {
				return err
			}
			q = r
			accept(name, value)

		case 0b0101, 0b0111:
			// 01N1_XXXX: Literal Field Line with Name Reference in static table
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-literal-field-line-with-nam

			index, r, err := varint.Read(q, 0b0000_1111)
			if err != nil {
				return err
			}
			if index >= uint64(len(staticTable)) {
				return errStaticIndexOutOfRange
			}
			value, r, err := readStringLiteral(r, buf)
			if err != nil {
				return err
			}
			q = r
			accept(staticTable[index].Name, value)

		case 0b1000, 0b1001, 0b1010, 0b1011:
			// 10XX_XXXX Indexed Field Line in dynamic table
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-indexed-field-line
			index, r, err := varint.Read(q, 0b0011_1111)
			if err != nil {
				return err
			}
			name, value, err := d.lineIndex(index)
			if err != nil {
				return err
			}
			q = r
			accept(name, value)

		case 0b1100, 0b1101, 0b1110, 0b1111:
			// 11XX_XXXX Indexed Field Line in static table
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-indexed-field-line
			index, r, err := varint.Read(q, 0b0011_1111)
			if err != nil {
				return err
			}
			if index >= uint64(len(staticTable)) {
				return errStaticIndexOutOfRange
			}
			q = r
			accept(staticTable[index].Name, staticTable[index].Value)
		}
	}
	return nil
}

func (d *Decoder) nameIndex(index uint64) (string, error) {
	return "", errors.New("TODO unsupported")
}

func (d *Decoder) baseNameIndex(index uint64) (string, error) {
	return "", errors.New("TODO unsupported")
}

func (d *Decoder) lineIndex(index uint64) (string, string, error) {
	return "", "", errors.New("TODO unsupported")
}

func (d *Decoder) baseLineIndex(index uint64) (string, string, error) {
	return "", "", errors.New("TODO unsupported")
}

// readLiteralName reads a literal name from p. Will use decodeBuf if the
// string needs to be huffman decoded.
func readLiteralName(p, decodeBuf []byte) (string, []byte, error) {
	const (
		// layout of the first byte of a literal name length
		P = 0b0010_0000
		N = 0b0001_0000
		H = 0b0000_1000
		M = 0b0000_0111
	)
	if len(p) <= 0 {
		return "", p, errUnexpectedEnd
	}
	n, q, err := varint.Read(p, M)
	if err != nil {
		return "", p, err
	}
	if n > uint64(len(q)) {
		return "", p, errUnexpectedEnd
	}
	b := q[:n:n]
	if p[0]&H == H {
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
		H = 0b1000_0000
		M = 0b0111_1111
	)
	return readLiteral(p, decodeBuf, M, H)
}

func readLiteral(p, decodeBuf []byte, m, h uint8) (string, []byte, error) {
	if len(p) <= 0 {
		return "", p, errUnexpectedEnd
	}
	n, q, err := varint.Read(p, m)
	if err != nil {
		return "", p, err
	}
	if n > uint64(len(q)) {
		return "", p, errUnexpectedEnd
	}
	b := q[:n:n]
	if p[0]&h == h {
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
