package quack

import (
	"errors"

	"github.com/renthraysk/quack/internal/field"
	"github.com/renthraysk/quack/internal/inst"
	"github.com/renthraysk/quack/varint"
)

func (dt *DT) appendEncoderInstructionLocked(p []byte, name, value string) []byte {
	i, isStatic, m := dt.current.Lookup(name, value)
	if m == field.MatchNameValue {
		// @TODO Duplicate?
		return p
	}
	ctrl := field.HeaderControl(name)
	if ctrl.NeverIndex() {
		// If already have a name match, no point attempting an insert
		// if prevented from inserting a (name, value) pair.
		if m == field.MatchName {
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
	case field.MatchName:
		p = inst.AppendInsertWithNameReference(p, i, isStatic)
	case field.MatchNone:
		p = inst.AppendInsertWithLiteralName(p, name)
	}
	return inst.AppendStringLiteral(p, value, ctrl.ShouldHuffman())
}

// Encoder Instructions
// https://www.rfc-editor.org/rfc/rfc9204.html#name-encoder-instructions
func (dt *DT) appendEncoderInstructions(p []byte, header map[string][]string) ([]byte, *field.Encoder) {

	dt.mu.Lock()
	defer dt.mu.Unlock()

	for name, values := range header {
		for _, value := range values {
			p = dt.appendEncoderInstructionLocked(p, name, value)
		}
	}
	// Build a fieldEncoder for when peer acks.
	m := make(field.NameValues, len(dt.headers))
	for i, hf := range dt.headers {
		m[hf.Name] = append(m[hf.Name], field.Value{Value: hf.Value, Index: uint64(i)})
	}

	var reqInsertCount uint64

	return p, field.New(m, reqInsertCount, dt.base, dt.capacity)
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

			i, q, err := varint.Read(p, M)
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

			capacity, q, err := varint.Read(p, M)
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
