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
	mu           sync.Mutex
	headers      []Header
	fieldEncoder *Encoder
	evicted      uint64
	size         uint64
	capacity     uint64
	maxCapacity  uint64
}

func (dt *DT) insertCountLocked() uint64 {
	return dt.evicted + uint64(len(dt.headers))
}

func (dt *DT) headerFromRelativePosLocked(rel uint64) (Header, bool) {
	abs := dt.insertCountLocked() - rel - 1
	if abs >= uint64(len(dt.headers)) {
		return Header{}, false
	}
	return dt.headers[abs], true
}

func (dt *DT) setMaxCapacity(maxCapacity uint64) {
	dt.mu.Lock()
	defer dt.mu.Unlock()
	dt.maxCapacity = maxCapacity
	if dt.capacity > maxCapacity {
		dt.setCapacityLocked(maxCapacity)
	}
}

func (dt *DT) setCapacityLocked(capacity uint64) bool {
	if capacity > dt.maxCapacity {
		return false
	}
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
		size -= dt.headers[i].size()
		i++
	}
	if i == 0 {
		return true
	}
	if i >= len(dt.headers) {
		return false
	}
	evicted, c := bits.Add64(dt.evicted, uint64(i), 0)
	if c != 0 {
		return false
	}
	// Eviction can proceed, modify state.
	dt.evicted = evicted
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
	// This addition cannot overflow as dt.size <= dt.capacity - s
	dt.size += s
	dt.headers = append(dt.headers, Header{Name: name, Value: value})
	return true
}

func (dt *DT) appendSnapshot(p []byte) []byte {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	p = inst.AppendSetDynamicTableCapacity(p, dt.capacity)

	n := make(map[string]int, len(dt.headers))
	m := make(map[Header]int, len(dt.headers))

	for i, hf := range dt.headers {
		if j, ok := m[hf]; ok {
			p = inst.AppendDuplicate(p, uint64(j))
			continue
		}
		m[hf] = i
		if j, ok := n[hf.Name]; ok {
			p = inst.AppendInsertWithNameReference(p, uint64(j), false)
		} else {
			n[hf.Name] = i
			p = inst.AppendInsertWithLiteralName(p, hf.Name)
		}
		p = inst.AppendStringLiteral(p, hf.Value, true)
	}
	return p
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-insert-with-literal-name
func decodeNameInsertWithLiteralName(p, buf []byte) (string, []byte, error) {
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
func (dt *DT) decodeNameInsertWithNameReference(p []byte) (string, []byte, error) {
	const T = 0b0100_0000
	const M = 0b0011_1111

	i, q, err := varint.Read(p, M)
	if err != nil {
		return "", p, err
	}
	if p[0]&T != 0 {
		if i >= uint64(len(staticTable)) {
			return "", p, errors.New("invalid static table index")
		}
		return staticTable[i].Name, q, nil
	}
	h, ok := dt.headerFromRelativePosLocked(i)
	if !ok {
		return "", p, errors.New("invalid dynamic table index")
	}
	return h.Name, q, nil
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

	base := dt.insertCountLocked()

	for name, values := range header {
		for _, value := range values {
			p = dt.appendEncoderInstructionLocked(p, name, value)
		}
	}
	// Build a fieldEncoder for when peer acks.
	nv := make(nameValues, len(dt.headers))
	for i, hf := range dt.headers {
		nv[hf.Name] = append(nv[hf.Name], value{value: hf.Value, index: uint64(i)})
	}

	return p, newEncoder(nv, base, dt.insertCountLocked()-base, dt.capacity)
}

func (dt *DT) DecodeEncoderInstructions(p []byte) error {
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
			h, ok := dt.headerFromRelativePosLocked(i)
			if !ok {
				return errors.New("duplicate: non-existant header")
			}
			if ok := dt.insertLocked(h.Name, h.Value); !ok {
				return errors.New("duplicate: failed to insert")
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
			name, q, err := decodeNameInsertWithLiteralName(p, decodeBuf[:0])
			if err != nil {
				return err
			}
			value, q, err := readStringLiteral(q, decodeBuf[:0])
			if err != nil {
				return err
			}
			if ok := dt.insertLocked(name, value); !ok {
				return errors.New("failed to insert header with literal name")
			}
			p = q

		default:
			name, q, err := dt.decodeNameInsertWithNameReference(p)
			if err != nil {
				return err
			}
			value, q, err := readStringLiteral(q, decodeBuf[:0])
			if err != nil {
				return err
			}
			if ok := dt.insertLocked(name, value); !ok {
				return errors.New("failed to insert header with name reference")
			}
			p = q
		}
	}
	return nil
}
