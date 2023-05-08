package quack

import (
	"errors"
	"math/bits"
)

type DT struct {
	headers  []headerField
	base     uint64 // evicted count
	capacity uint64
	size     uint64
}

func New(maxCapacity uint64) DT {
	return DT{
		headers:  make([]headerField, 0, maxCapacity/32),
		capacity: maxCapacity,
	}
}

func (dt *DT) lookup(name, value string) (uint64, match) {
	i := len(dt.headers) - 1
	for i >= 0 && dt.headers[i].Name != name {
		i--
	}
	if i < 0 || i >= len(dt.headers) {
		return 0, matchNone
	}
	if dt.headers[i].Value == value {
		return uint64(i) + dt.base, matchNameValue
	}
	j := i - 1
	for j >= 0 && dt.headers[j].Name != name || dt.headers[j].Value != value {
		j--
	}
	if j < 0 {
		return uint64(i) + dt.base, matchName
	}
	return uint64(j) + dt.base, matchNameValue
}

// evict attempts to evict headerFields until size is less than or equal to
// targetSize. Returns true if was able to ensure the dynamic table size
// is or below targetSize, false otherwise.
func (dt *DT) evict(targetSize uint64) bool {
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
	dt.base = b
	dt.size = size
	dt.headers = append(dt.headers[:0], dt.headers[i:]...)
	return true
}

func (dt *DT) insert(name, value string) bool {
	s := size(name, value)
	if s > dt.capacity {
		return false
	}
	if ok := dt.evict(dt.capacity - s); !ok {
		return false
	}
	dt.size += s
	dt.headers = append(dt.headers, headerField{name, value})
	return true
}

func (dt *DT) maxEntries() uint64 { return dt.capacity / 32 }

// https://www.rfc-editor.org/rfc/rfc9204.html#name-encoded-field-section-prefi
func (dt *DT) appendFieldSectionPrefix(p []byte, reqInsertCount uint64) []byte {
	var insertCount uint64

	// https://www.rfc-editor.org/rfc/rfc9204.html#name-required-insert-count
	if reqInsertCount > 0 {
		insertCount = (reqInsertCount % (2 * dt.maxEntries())) + 1
	}
	p = appendVarint(p, insertCount, 0xFF, 0)

	// https://www.rfc-editor.org/rfc/rfc9204.html#name-base
	var sign byte
	base := dt.base - reqInsertCount
	if dt.base < reqInsertCount {
		base = reqInsertCount - dt.base - 1
		sign = 0b1000_0000
	}
	return appendVarint(p, base, 0b0111_1111, sign)
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
