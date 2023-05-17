package quack

import (
	"errors"
	"math/bits"
	"sync"
	"sync/atomic"
)

type DT struct {
	mu       sync.Mutex
	headers  []headerField
	base     uint64
	capacity uint64
	size     uint64

	// current is the encoder that encodes in manner understood by the peer's
	// decoder
	current *fieldEncoder

	// next is the encoder to switch to when the peer acks insertions into
	// it's own dynamic table
	next *fieldEncoder
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

func (dt *DT) changeEncoder(p *atomic.Pointer[fieldEncoder], knownReceivedCount uint64) error {
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

	encodedInsertCount, q, err := readVarint(p, 0xFF)
	if err != nil {
		return p, 0, 0, err
	}

	const (
		SIGN = 0b1000_0000
		M    = 0b0111_1111
	)

	deltaBase, r, err := readVarint(q, M)
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
	p = appendSetDynamicTableCapacity(p, dt.capacity)
	dt.mu.Unlock()

	n := make(map[string]uint64, len(headers))
	m := make(map[headerField]uint64, len(headers))

	for i, hf := range headers {
		if j, ok := m[hf]; ok {
			p = appendDuplicate(p, j)
			continue
		}
		m[hf] = uint64(i)
		if j, ok := n[hf.Name]; ok {
			p = appendInsertWithNameReference(p, j, false)
		} else {
			n[hf.Name] = uint64(i)
			p = appendInsertWithLiteralName(p, hf.Name)
		}
		p = appendStringLiteral(p, hf.Value, true)
	}
	return p
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
