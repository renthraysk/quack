package field

import (
	"errors"
	"math/bits"
	"sync"
	"sync/atomic"

	"github.com/renthraysk/quack/ascii"
	"github.com/renthraysk/quack/huffman"
	"github.com/renthraysk/quack/internal/inst"
	"github.com/renthraysk/quack/varint"
)

type DT struct {
	mu       sync.Mutex
	headers  []Header
	base     uint64
	capacity uint64
	size     uint64

	// fieldEncoder is the encoder that encodes in manner understood by the peer's
	// decoder
	fieldEncoder *Encoder
}

func (dt *DT) setCapacityLocked(capacity uint64) bool {
	if dt.evictLocked(capacity) {
		dt.capacity = capacity
		return true
	}
	return false
}

func (dt *DT) changeEncoder(p *atomic.Pointer[Encoder], knownReceivedCount uint64) error {
	/* TODO */
	return nil
}

func (dt *DT) changeDecoder(p *atomic.Pointer[Decoder], knownReceivedCount uint64) error {
	/* TODO */
	return nil
}

// evictLocked attempts to evict field.Headers until size is less than or equal to
// targetSize. Returns true if was able to ensure the dynamic table size
// is or below targetSize, false otherwise.
func (dt *DT) evictLocked(targetSize uint64) bool {
	var i int

	size := dt.size
	for i < len(dt.headers) && size > targetSize {
		size -= dt.headers[i].Size()
		i++
	}
	if i == 0 {
		return true
	}
	if i >= len(dt.headers) {
		return false
	}
	b, c := bits.Add64(dt.base, uint64(i), 0)
	if c != 0 {
		return false
	}
	// Eviction can proceed, modify state.
	dt.base = b
	dt.size = size
	dt.headers = append(dt.headers[:0], dt.headers[i:]...)
	return true
}

func (dt *DT) insertLocked(name, value string) bool {
	s := headerSize(name, value)

	if s > dt.capacity {
		return false
	}
	if ok := dt.evictLocked(dt.capacity - s); !ok {
		return false
	}
	dt.size += s
	dt.headers = append(dt.headers, Header{Name: name, Value: value})
	return true
}

// appendSnapshot will append encoder instructions to recreate the current state
// of the dynamic table
func (dt *DT) appendSnapshot(p []byte) []byte {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	p = inst.AppendSetDynamicTableCapacity(p, dt.capacity)

	n := make(map[string]uint64, len(dt.headers))
	m := make(map[Header]uint64, len(dt.headers))

	for i, hf := range dt.headers {
		if j, ok := m[hf]; ok {
			p = inst.AppendDuplicate(p, j)
			continue
		}
		m[hf] = uint64(i)
		if j, ok := n[hf.Name]; ok {
			p = inst.AppendInsertWithNameReference(p, j, false)
		} else {
			n[hf.Name] = uint64(i)
			p = inst.AppendInsertWithLiteralName(p, hf.Name)
		}
		p = inst.AppendStringLiteral(p, hf.Value, true)
	}
	return p
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-insert-with-literal-name
func nameInsertWithLiteralName(p, buf []byte) (string, []byte, error) {
	const H = 0b0010_0000
	const M = 0b0001_1111

	n, q, err := varint.Read(p, M)
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

	i, q, err := varint.Read(p, M)
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

func (dt *DT) appendEncoderInstructionLocked(p []byte, name, value string) []byte {
	i, isStatic, m := dt.fieldEncoder.lookup(name, value)
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
		p = inst.AppendInsertWithNameReference(p, i, isStatic)
	case matchNone:
		p = inst.AppendInsertWithLiteralName(p, name)
	}
	return inst.AppendStringLiteral(p, value, ctrl.shouldHuffman())
}

// Encoder Instructions
// https://www.rfc-editor.org/rfc/rfc9204.html#name-encoder-instructions
func (dt *DT) appendEncoderInstructions(p []byte, header map[string][]string) ([]byte, *Encoder) {

	dt.mu.Lock()
	defer dt.mu.Unlock()

	for name, values := range header {
		for _, value := range values {
			p = dt.appendEncoderInstructionLocked(p, name, value)
		}
	}
	// Build a fieldEncoder for when peer acks.
	m := make(nameValues, len(dt.headers))
	for i, hf := range dt.headers {
		m[hf.Name] = append(m[hf.Name], value{value: hf.Value, index: uint64(i)})
	}

	var reqInsertCount uint64

	return p, newEncoder(m, reqInsertCount, dt.base, dt.capacity)
}

func (dt *DT) ParseEncoderInstructions(p []byte) error {
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
				return errors.New("not existing header")
			}
			if ok := dt.insertLocked(dt.headers[i].Name, dt.headers[i].Value); !ok {
				return errors.New("failed to insert duplicate")
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
				return errors.New("failed to set capacity")
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
				return errors.New("failed to insert header with name referernce")
			}
			p = q
		}
	}
	return nil
}
