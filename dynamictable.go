package quack

import (
	"errors"
	"math/bits"
	"sync"
	"sync/atomic"

	"github.com/renthraysk/quack/ascii"
	"github.com/renthraysk/quack/huffman"
	"github.com/renthraysk/quack/internal/field"
	"github.com/renthraysk/quack/internal/inst"
	"github.com/renthraysk/quack/varint"
)

type DT struct {
	mu       sync.Mutex
	headers  []headerField
	base     uint64
	capacity uint64
	size     uint64

	// current is the encoder that encodes in manner understood by the peer's
	// decoder
	current *field.Encoder

	// next is the encoder to switch to when the peer acks insertions into
	// it's own dynamic table
	next *field.Encoder
}

func New(maxCapacity uint64) DT {
	return DT{
		headers:  make([]headerField, 0, maxCapacity/32),
		capacity: maxCapacity,
	}
}

func (dt *DT) setCapacityLocked(capacity uint64) bool {
	if dt.evictLocked(capacity) {
		dt.capacity = capacity
		return true
	}
	return false
}

func (dt *DT) changeEncoder(p *atomic.Pointer[field.Encoder], knownReceivedCount uint64) error {
	p.Store(dt.next)
	return nil
}

// evictLocked attempts to evict headerFields until size is less than or equal to
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
	b, c := bits.Add64(dt.base, uint64(i), 0)
	if c != 0 {
		return false
	}
	// Successful eviction, modify state.
	dt.base = b
	dt.size = size
	dt.headers = append(dt.headers[:0], dt.headers[i:]...)
	return true
}

func (dt *DT) insertLocked(name, value string) bool {
	s := size(name, value)

	if s > dt.capacity {
		return false
	}
	if ok := dt.evictLocked(dt.capacity - s); !ok {
		return false
	}
	dt.size += s
	dt.headers = append(dt.headers, headerField{name, value})
	return true
}

// https://datatracker.ietf.org/doc/html/rfc9204#name-encoded-field-section-prefi
func (dt *DT) readFieldSectionPrefix(p []byte) ([]byte, uint64, uint64, error) {
	var reqInsertCount uint64

	encodedInsertCount, q, err := varint.Read(p, 0xFF)
	if err != nil {
		return p, 0, 0, err
	}

	const (
		SIGN = 0b1000_0000
		M    = 0b0111_1111
	)

	deltaBase, r, err := varint.Read(q, M)
	if err != nil {
		return p, 0, 0, err
	}

	// https://datatracker.ietf.org/doc/html/rfc9204#name-required-insert-count
	if encodedInsertCount != 0 {
		fullRange := 2 * (dt.capacity / 32)
		if encodedInsertCount > fullRange {
			return p, 0, 0, errors.New("")
		}
		maxValue := /* @TODO */ uint64(len(dt.headers)) + (dt.capacity / 32)
		maxWrapped := (maxValue / fullRange) * fullRange
		reqInsertCount = maxWrapped + encodedInsertCount - 1
		if reqInsertCount > maxValue {
			if reqInsertCount <= fullRange {
				return p, 0, 0, errors.New("")
			}
			reqInsertCount -= fullRange
		}
		if reqInsertCount == 0 {
			return p, 0, 0, errors.New("")
		}
	}

	// https://datatracker.ietf.org/doc/html/rfc9204#name-base
	base := reqInsertCount + deltaBase
	if q[0]&SIGN != 0 {
		base = reqInsertCount - deltaBase - 1
	}
	return r, reqInsertCount, base, nil
}

// appendSnapshot will append encoder instructions to recreate the current state
// of the dynamic table
func (dt *DT) appendSnapshot(p []byte) []byte {

	dt.mu.Lock()
	headers := dt.headers
	p = inst.AppendSetDynamicTableCapacity(p, dt.capacity)
	dt.mu.Unlock()

	n := make(map[string]uint64, len(headers))
	m := make(map[headerField]uint64, len(headers))

	for i, hf := range headers {
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

func (dt *DT) nameIndex(index uint64) (string, error) {
	return "", errors.New("TODO unsupported")
}

func (dt *DT) baseNameIndex(index uint64) (string, error) {
	return "", errors.New("TODO unsupported")
}

func (dt *DT) lineIndex(index uint64) (string, string, error) {
	return "", "", errors.New("TODO unsupported")
}

func (dt *DT) baseLineIndex(index uint64) (string, string, error) {
	return "", "", errors.New("TODO unsupported")
}
